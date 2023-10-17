package main

import (
	"time"
	"sitebrush/pkg/config"
	"sitebrush/pkg/listener"
	"sitebrush/pkg/database"
	"sitebrush/pkg/mylog"
)

func main() {

	settings := Config.ParseFlags()

	//run error log daemon
	go MyLog.ErrorLogWorker()
	go Database.Run(settings)
	go listener.Run(settings)

	for {
		time.Sleep(10 * time.Second)	
	}

}
