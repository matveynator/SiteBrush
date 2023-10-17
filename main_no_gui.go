// +build android aix dragonfly freebsd illumos js netbsd openbsd plan9 solaris zos



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
	go database.Run(settings)
  go webserver.Run(settings)

	for {
		time.Sleep(10 * time.Second)	
	}

}
