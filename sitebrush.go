
package main

import (
	"time"
	"sitebrush/pkg/config"
	"sitebrush/pkg/database"
	"sitebrush/pkg/mylog"
  "sitebrush/pkg/webserver"
)

func main() {

	settings := Config.ParseFlags()

	go MyLog.ErrorLogWorker()
  go webserver.Run(settings)
  go database.Run(settings)

	for {
		time.Sleep(10 * time.Second)	
	}

}
