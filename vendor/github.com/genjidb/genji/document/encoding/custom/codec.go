package custom

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/genjidb/genji/document"
	"github.com/genjidb/genji/internal/binarysort"
	"github.com/genjidb/genji/internal/errors"
	"github.com/genjidb/genji/types"
)

// A Codec is a custom implementation of an encoding.Codec.
type Codec struct{}

// NewCodec creates a custom codec.
func NewCodec() *Codec {
	return &Codec{}
}

// Encoder encodes Genji documents and values
// in MessagePack.
type Encoder struct {
	w io.Writer
}

// NewEncoder creates an Encoder that writes in the given writer.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		w: w,
	}
}

// EncodeDocument encodes d.
func (e *Encoder) EncodeDocument(d types.Document) error {
	data, err := EncodeDocument(d)
	if err != nil {
		return err
	}

	_, err = e.w.Write(data)
	return err
}

// Close does nothing.
func (e *Encoder) Close() {}

// EncodeDocument takes a document and encodes it using the encoding.Format type.
func EncodeDocument(d types.Document) ([]byte, error) {
	if ec, ok := d.(*EncodedDocument); ok {
		return ec.data, nil
	}

	var format Format

	var offset uint64
	var dataList [][]byte

	err := d.Iterate(func(f string, v types.Value) error {
		data, err := EncodeValue(v)
		if err != nil {
			return err
		}

		format.Header.FieldHeaders = append(format.Header.FieldHeaders, FieldHeader{
			NameSize:   uint64(len(f)),
			NameString: f,
			Type:       uint64(v.Type()),
			Size:       uint64(len(data)),
			Offset:     offset,
		})

		offset += uint64(len(data))
		dataList = append(dataList, data)
		return nil
	})
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	_, err = format.Header.WriteTo(&buf)
	if err != nil {
		return nil, err
	}

	buf.Grow(format.Header.BodySize())

	for _, data := range dataList {
		_, err = buf.Write(data)
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

// DecodeDocument takes a byte slice and returns a lazily decoded document.
// If buf is malformed, an error will be returned when calling one of the document method.
func DecodeDocument(buf []byte) types.Document {
	return &EncodedDocument{buf}
}

// EncodeValue encodes any value to a binary representation.
func EncodeValue(v types.Value) ([]byte, error) {
	switch v.Type() {
	case types.DocumentValue:
		return EncodeDocument(v.V().(types.Document))
	case types.ArrayValue:
		return EncodeArray(v.V().(types.Array))
	case types.BlobValue:
		return v.V().([]byte), nil
	case types.TextValue:
		return []byte(v.V().(string)), nil
	case types.BoolValue:
		return binarysort.AppendBool(nil, v.V().(bool)), nil
	case types.IntegerValue:
		return encodeInt64(v.V().(int64)), nil
	case types.DoubleValue:
		return binarysort.AppendFloat64(nil, v.V().(float64)), nil
	case types.NullValue:
		return nil, nil
	}

	return nil, errors.New("unknown type")
}

func encodeInt64(x int64) []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutVarint(buf, x)
	return buf[:n]
}

// An EncodedDocument implements the types.Document interface on top of an encoded representation of a
// document.
// It is useful to avoid decoding the entire document when only a few fields are needed.
type EncodedDocument struct {
	data []byte
}

// GetByField decodes the selected field.
func (e *EncodedDocument) GetByField(field string) (types.Value, error) {
	return decodeValueFromDocument(e.data, field)
}

func (e *EncodedDocument) Reset(data []byte) { e.data = data }

// Iterate decodes each fields one by one and passes them to fn until the end of the document
// or until fn returns an error.
func (e *EncodedDocument) Iterate(fn func(field string, value types.Value) error) error {
	var format Format
	err := format.Decode(e.data)
	if err != nil {
		return err
	}

	for _, fh := range format.Header.FieldHeaders {
		v, err := DecodeValue(types.ValueType(fh.Type), format.Body[fh.Offset:fh.Offset+fh.Size])
		if err != nil {
			return err
		}

		err = fn(string(fh.Name), v)
		if err != nil {
			return err
		}
	}

	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (e *EncodedDocument) MarshalJSON() ([]byte, error) {
	return document.MarshalJSON(e)
}

// An EncodedArray implements the types.Array interface on top of an encoded representation of an
// array.
// It is useful to avoid decoding the entire array when only a few values are needed.
type EncodedArray []byte

// Iterate goes through all the values of the array and calls the given function by passing each one of them.
// If the given function returns an error, the iteration stops.
func (e EncodedArray) Iterate(fn func(i int, value types.Value) error) error {
	var format Format
	err := format.Decode(e)
	if err != nil {
		return err
	}

	for _, fh := range format.Header.FieldHeaders {
		v, err := DecodeValue(types.ValueType(fh.Type), format.Body[fh.Offset:fh.Offset+fh.Size])
		if err != nil {
			return err
		}

		i, _ := binary.Varint(fh.Name)
		err = fn(int(i), v)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetByIndex returns a value by index of the array.
func (e EncodedArray) GetByIndex(i int) (types.Value, error) {
	v, err := decodeValueFromDocument(e, string(encodeInt64(int64(i))))
	if errors.Is(err, document.ErrFieldNotFound) {
		return v, document.ErrValueNotFound
	}

	return v, err
}

// MarshalJSON implements the json.Marshaler interface.
func (e EncodedArray) MarshalJSON() ([]byte, error) {
	return document.MarshalJSONArray(e)
}

func decodeValueFromDocument(data []byte, field string) (types.Value, error) {
	hsize, n := binary.Uvarint(data)
	if n <= 0 {
		return nil, errors.New("cannot decode data")
	}

	hdata := data[n : n+int(hsize)]
	body := data[n+len(hdata):]

	// skip number of fields
	_, n = binary.Uvarint(hdata)
	if n <= 0 {
		return nil, errors.New("cannot decode data")
	}
	hdata = hdata[n:]

	var fh FieldHeader
	for len(hdata) > 0 {
		n, err := fh.Decode(hdata)
		if err != nil {
			return nil, err
		}
		hdata = hdata[n:]

		if field == string(fh.Name) {
			return DecodeValue(types.ValueType(fh.Type), body[fh.Offset:fh.Offset+fh.Size])
		}
	}

	return nil, document.ErrFieldNotFound
}

// EncodeArray encodes a into its binary representation.
func EncodeArray(a types.Array) ([]byte, error) {
	var format Format

	var offset uint64
	var dataList [][]byte

	err := a.Iterate(func(i int, v types.Value) error {
		data, err := EncodeValue(v)
		if err != nil {
			return err
		}

		index := encodeInt64(int64(i))

		format.Header.FieldHeaders = append(format.Header.FieldHeaders, FieldHeader{
			NameSize: uint64(len(index)),
			Name:     index,
			Type:     uint64(v.Type()),
			Size:     uint64(len(data)),
			Offset:   offset,
		})

		offset += uint64(len(data))
		dataList = append(dataList, data)
		return nil
	})
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	_, err = format.Header.WriteTo(&buf)
	if err != nil {
		return nil, err
	}

	buf.Grow(format.Header.BodySize())

	for _, data := range dataList {
		_, err = buf.Write(data)
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

// DecodeArray takes a byte slice and returns a lazily decoded array.
// If buf is malformed, an error will be returned when calling one of the array method.
func DecodeArray(buf []byte) types.Array {
	return EncodedArray(buf)
}

// DecodeValue takes some encoded data and decodes it to the target type t.
func DecodeValue(t types.ValueType, data []byte) (types.Value, error) {
	switch t {
	case types.DocumentValue:
		return types.NewDocumentValue(&EncodedDocument{data}), nil
	case types.ArrayValue:
		return types.NewArrayValue(EncodedArray(data)), nil
	case types.BlobValue:
		return types.NewBlobValue(data), nil
	case types.TextValue:
		return types.NewTextValue(string(data)), nil
	case types.BoolValue:
		x, err := binarysort.DecodeBool(data)
		if err != nil {
			return nil, err
		}
		return types.NewBoolValue(x), nil
	case types.IntegerValue:
		x, _ := binary.Varint(data)
		return types.NewIntegerValue(x), nil
	case types.DoubleValue:
		x, err := binarysort.DecodeFloat64(data)
		if err != nil {
			return nil, err
		}
		return types.NewDoubleValue(x), nil
	case types.NullValue:
		return types.NewNullValue(), nil
	}

	return nil, errors.New("unknown type")
}
