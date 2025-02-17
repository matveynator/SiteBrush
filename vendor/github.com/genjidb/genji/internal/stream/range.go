package stream

import (
	"strings"

	"github.com/genjidb/genji/document"
	"github.com/genjidb/genji/internal/database"
	"github.com/genjidb/genji/internal/environment"
	"github.com/genjidb/genji/internal/expr"
	"github.com/genjidb/genji/internal/stringutil"
)

// Range represents a range to select values after or before
// a given boundary.
type Range struct {
	Min, Max expr.LiteralExprList
	Paths    []document.Path
	// Exclude Min and Max from the results.
	// By default, min and max are inclusive.
	// Exclusive and Exact cannot be set to true at the same time.
	Exclusive bool
	// Used to match an exact value equal to Min.
	// If set to true, Max will be ignored for comparison
	// and for determining the global upper bound.
	Exact bool
}

func (r *Range) Eval(env *environment.Environment) (*database.Range, error) {
	rng := database.Range{
		Exclusive: r.Exclusive,
		Exact:     r.Exact,
	}

	if r.Min != nil {
		min, err := r.Min.Eval(env)
		if err != nil {
			return nil, err
		}

		rng.Min = min.V().(*document.ValueBuffer).Values
	}

	if r.Max != nil {
		max, err := r.Max.Eval(env)
		if err != nil {
			return nil, err
		}
		rng.Max = max.V().(*document.ValueBuffer).Values
	}

	return &rng, nil
}

func (r *Range) String() string {
	format := func(el expr.LiteralExprList) string {
		switch len(el) {
		case 0:
			return "-1"
		case 1:
			return el[0].String()
		default:
			return el.String()
		}
	}

	if r.Exact {
		return stringutil.Sprintf("%v", format(r.Min))
	}

	if r.Exclusive {
		return stringutil.Sprintf("[%v, %v, true]", format(r.Min), format(r.Max))
	}

	return stringutil.Sprintf("[%v, %v]", format(r.Min), format(r.Max))
}

func (r *Range) IsEqual(other *Range) bool {
	if r.Exact != other.Exact {
		return false
	}

	if r.Exclusive != other.Exclusive {
		return false
	}

	if len(r.Min) != len(other.Min) {
		return false
	}

	if len(r.Max) != len(other.Max) {
		return false
	}

	if !r.Min.IsEqual(other.Min) {
		return false
	}

	if !r.Max.IsEqual(other.Max) {
		return false
	}

	return true
}

type Ranges []Range

// Encode each range using the given value encoder.
func (r Ranges) Eval(env *environment.Environment) ([]*database.Range, error) {
	ranges := make([]*database.Range, 0, len(r))

	for i := range r {
		rng, err := r[i].Eval(env)
		if err != nil {
			return nil, err
		}
		if rng != nil {
			ranges = append(ranges, rng)
		}
	}

	return ranges, nil
}

// Append rng to r and return the new slice.
// Duplicate ranges are ignored.
func (r Ranges) Append(rng Range) Ranges {
	// ensure we don't keep duplicate ranges
	isDuplicate := false
	for _, e := range r {
		if e.IsEqual(&rng) {
			isDuplicate = true
			break
		}
	}

	if isDuplicate {
		return r
	}

	return append(r, rng)
}

func (r Ranges) String() string {
	var sb strings.Builder

	for i, rr := range r {
		if i > 0 {
			sb.WriteString(", ")
		}

		sb.WriteString(rr.String())
	}

	return sb.String()
}

// Cost is a best effort function to determine the cost of
// a range lookup.
func (r Ranges) Cost() int {
	var cost int

	for _, rng := range r {
		// if we are looking for an exact value
		// increment by 1
		if rng.Exact {
			cost++
			continue
		}

		// if there are two boundaries, increment by 50
		if len(rng.Min) > 0 && len(rng.Max) > 0 {
			cost += 50
			continue
		}

		// if there is only one boundary, increment by 100
		if len(rng.Min) > 0 || len(rng.Max) > 0 {
			cost += 100
			continue
		}

		// if there are no boundaries, increment by 200
		cost += 200
	}

	return cost
}
