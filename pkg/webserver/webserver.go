package webserver

import (
  "os"
  "fmt"
  "net/http"
  "sitebrush/pkg/config"
)


func Run(config Config.Settings) {
  http.HandleFunc("/", func(responseWriter http.ResponseWriter, request *http.Request) {
    handleRequest(config, responseWriter, request)
  })
  err := http.ListenAndServe(config.WEB_LISTENER_ADDRESS, nil)
  if err != nil {
    fmt.Println("Ошибка при запуске веб-сервера:", err)
    os.Exit(0)
  }
}

