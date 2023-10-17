package listener

import (
	"fmt"
	"testing"
	"chicha/pkg/data"
	"chicha/pkg/config"
)

func TestParseCSVLine(t *testing.T) {

	testData := []struct {
		csvInput []byte
		csvOutput Data.RawData
	}{
		{
			csvInput: []byte("100008020200000100000189, 1622570553397, 3"),
			csvOutput: Data.RawData{ "100008020200000100000189", 1622570553397, "8.8.8.8", 3, "" },
		},
		{
			csvInput: []byte("TONIBOU-001, 1622570553397, 3, 9.9.9.9"),
			csvOutput: Data.RawData{"TONIBOU-001", 1622570553397, "9.9.9.9", 3, "8.8.8.8" },
		},
	}

	for _,testCase := range testData {
		//check1:
		testResult, _ := parseCSVLine(testCase.csvInput, "8.8.8.8")
		fmt.Println("CSV testResult.TagId=", testResult.TagId)
		if testResult != testCase.csvOutput  {
			t.Errorf("Incorrect result. Expect %v, got %v", testCase.csvOutput, testResult)
		}
	}
}

func TestParseXMLPacket(t *testing.T) {

	config := Config.ParseFlags()

	testData := []struct {
		xmlInput []byte
		xmlOutput Data.RawData
	}{
		{
			xmlInput: []byte("<Alien-RFID-Tag><TagID>1000 0802 0200 0001 0000 0796</TagID><DiscoveryTime>2021/05/16 12:00:34.730</DiscoveryTime><LastSeenTime>2021/05/16 12:00:34.730</LastSeenTime><Antenna>2</Antenna><ReadCount>1</ReadCount><Protocol>2</Protocol></Alien-RFID-Tag>"),
			xmlOutput: Data.RawData{ "100008020200000100000796", 1621166434730, "8.8.8.8", 2, "" },
		},
		{
			xmlInput: []byte("<Alien-RFID-Tag><TagID>TONIBOU-001</TagID><DiscoveryTime>2021/05/16 12:00:34.823</DiscoveryTime><LastSeenTime>2021/05/16 12:00:34.823</LastSeenTime><Antenna>3</Antenna><ReaderIP>9.9.9.9</ReaderIP><ReadCount>1</ReadCount><Protocol>2</Protocol></Alien-RFID-Tag>"),
			xmlOutput: Data.RawData{"TONIBOU-001", 1621166434823, "9.9.9.9", 3, "8.8.8.8" },
		},
	}

	for _,testCase := range testData {
		testResult, _ := parseXMLPacket(testCase.xmlInput, "8.8.8.8", config)
		fmt.Println("XML testResult.TagId=", testResult.TagId)
		if testResult != testCase.xmlOutput  {
			t.Errorf("Incorrect result. Expect %v, got %v", testCase.xmlOutput, testResult)
		}
	}
}
