package listener

import (
	"io"
	"log"
	"net"
	"time"

	"sitebrush/pkg/config"
	"sitebrush/pkg/data"
	"sitebrush/pkg/database"

)
//initial function
func init() {
	//set microsecond resolution for logging:
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}

func processConnection(connection net.Conn, config Config.Settings) {
	defer connection.Close()
	var tempDelay time.Duration // how long to sleep on accept failure
	var rawData Data.RawData

	// Read connection in lap
	for {
		buf := make([]byte, 1024)
		size, err := connection.Read(buf)
		if err != nil {
			if err == io.EOF {
				//log.Println("conn.Read(buf) error:", err)
				//log.Println("Message EOF detected - closing LAN connection.")
				break
			}

			networkError, ok := err.(*net.OpError) 
			if ok && networkError.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay = 2 * tempDelay
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				log.Printf("Accept error: %v; retrying in %v", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}

			break
		}
		tempDelay = 0

		data := buf[:size]

		remoteIPAddress := connection.RemoteAddr().(*net.TCPAddr).IP.String()

		//various data formats processing (text csv, xml) start:
		if !IsValidXML(data) {
			// CSV data processing
			rawData, err = parseCSVLine(data, remoteIPAddress)
			if err != nil {
				log.Println(err)
			}
		} else {
			// XML data processing
			rawData, err = parseXMLPacket(data, remoteIPAddress, config)
			if err != nil {
				log.Println(err)
			}
		}

		//create a database task:
		Database.DatabaseSaveRawTask <- rawData

	}
}

// Start data collector from RFID readers.
func Run(config Config.Settings) {
	// Start listener
	collector, err := net.Listen("tcp", config.DAEMON_LISTENER_ADDRESS)
	if err != nil {
		log.Panicln("Error: collector can't start. ", err)
	}
	defer collector.Close()


	// Listen new connections
	for {
		connection, err := collector.Accept()
		if err != nil {
			if err != io.EOF {
				log.Panicln("tcp connection error:", err)
			}
		}

		go processConnection(connection, config)
	}
}

