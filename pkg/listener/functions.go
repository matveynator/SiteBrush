package listener
import (
	"bytes"
	"fmt"
	"log"
	"net"
	"time"
	"strconv"
	"strings"
	"errors"
	"encoding/csv"
	"encoding/xml"

	"sitebrush/pkg/config"
	"sitebrush/pkg/data"
)
func IsValidXML(data []byte) bool {
	return xml.Unmarshal(data, new(interface{})) == nil
}


func parseCSVLine(data []byte, remoteIPAddress string) (rawData Data.RawData, err error){
	if !IsValidXML(data) {
		// CSV data processing
		var csvField []string
		csvField, err = csv.NewReader(bytes.NewReader(data)).Read()
		if err != nil {
			log.Println(err)
			return
		} else {
			//if one packet received (valid):
			numberOfCSVColumns := len(csvField)
			if numberOfCSVColumns == 3 || numberOfCSVColumns == 4 {
				// Prepare antenna position
				var antennaPosition int64
				antennaPosition, err = strconv.ParseInt(string(strings.TrimSpace(csvField[2])), 10, 64)
				if err != nil {
					log.Println("Recived incorrect Antenna position CSV value:", err)
					return
				} else { 
					rawData.DiscoveryUnixTime, err = strconv.ParseInt(string(strings.TrimSpace(csvField[1])), 10, 64)
					if err != nil {
						log.Println("Recived incorrect discovery unix time CSV value:", err)
						return
					} else {
						rawData.TagId = string(strings.TrimSpace(csvField[0]))

						rawData.Antenna = uint(antennaPosition)

						if numberOfCSVColumns == 3 {
							//reader connection without proxy
							rawData.ReaderIP = remoteIPAddress

							//Debug all received data from RFId reader
							log.Printf("TAG=%s, TIME=%d, Reader-IP=%s, Reader-ANT=%d\n", rawData.TagId, rawData.DiscoveryUnixTime, rawData.ReaderIP, rawData.Antenna)

						} else if numberOfCSVColumns == 4 {
							//proxy connection
							if net.ParseIP(string(strings.TrimSpace(csvField[3]))) != nil {
								//sending data with proxy
								rawData.ReaderIP = string(strings.TrimSpace(csvField[3]))
								rawData.ProxyIP = remoteIPAddress
							} else {
								//sending csvField[3] is not an IP address
								rawData.ReaderIP = remoteIPAddress
							}
							//Debug all received data from PROXY
							log.Printf("TAG=%s, TIME=%d, Reader-IP=%s, Reader-Antenna=%d, Proxy-IP=%s\n", rawData.TagId, rawData.DiscoveryUnixTime, rawData.ReaderIP, rawData.Antenna, rawData.ProxyIP)

						}
					}
				}
			} else {
				err = errors.New(fmt.Sprintf("Error: Parser expected 3 or 4 fields in CSV, but received %d fields.", numberOfCSVColumns))
			}
		}
	} else {
		err = errors.New("Error: received data is not a valid CSV - it is XML!")
	}
	return
}

func parseXMLPacket(data []byte, remoteIPAddress string, config Config.Settings)(rawData Data.RawData, err error) {
	if IsValidXML(data) {
		// XML data processing
		var rawXMLData Data.RawXMLData
		err = xml.Unmarshal(data, &rawXMLData)
		if err != nil {
			log.Println("xml.Unmarshal ERROR:", err)
			return
		}

		var loc *time.Location
		loc, err = time.LoadLocation(config.TIME_ZONE)
		if err != nil {
			log.Println(err)
			return
		}
		xmlTimeFormat := `2006/01/02 15:04:05.000`
		var discoveryTime time.Time
		discoveryTime, err = time.ParseInLocation(xmlTimeFormat, rawXMLData.DiscoveryTime, loc)
		if err != nil {
			log.Println("time.ParseInLocation ERROR:", err)
			return
		}
		rawData.DiscoveryUnixTime = discoveryTime.UnixNano()/int64(time.Millisecond)
		rawData.TagId = strings.ReplaceAll(rawXMLData.TagId, " ", "")
		rawData.Antenna = uint(rawXMLData.Antenna)

		if net.ParseIP(rawXMLData.ReaderIP) != nil {
			//connection received from proxy (not from reader).
			rawData.ReaderIP = rawXMLData.ReaderIP
			rawData.ProxyIP = remoteIPAddress

			//Debug all received data from PROXY
			log.Printf("TAG=%s, TIME=%d, Reader-IP=%s, Reader-Antenna=%d, Proxy-IP=%s\n", rawData.TagId, rawData.DiscoveryUnixTime, rawData.ReaderIP, rawData.Antenna, rawData.ProxyIP)
		} else {
			//connection received from reader (not from proxy)
			rawData.ReaderIP = remoteIPAddress
			//Debug all received data from RFID reader
			log.Printf("TAG=%s, TIME=%d, Reader-IP=%s, Reader-ANT=%d\n", rawData.TagId, rawData.DiscoveryUnixTime, rawData.ReaderIP, rawData.Antenna)
		}
	} else {
		err = errors.New("Error: received data is not a valid XML.")
	}
	return
}
