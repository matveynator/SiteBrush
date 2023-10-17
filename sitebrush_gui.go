// +build darwin linux windows

package main

import (
  "time"
  "sitebrush/pkg/config"
  "sitebrush/pkg/database"
  "sitebrush/pkg/mylog"
  "sitebrush/pkg/webserver"
  "sitebrush/pkg/browser"

)

func main() {

  settings := Config.ParseFlags()

  //run error log daemon
  go MyLog.ErrorLogWorker()
  go database.Run(settings)
  go webserver.Run(settings)

  if ( settings.GUI == true ) {
    browser.Run(settings)
  } else {
    for {
      time.Sleep(10 * time.Second)
    }
  }
}

