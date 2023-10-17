//go:build ( ( darwin && amd64 ) || ( darwin && arm64 ) || ( darwin && 386 ) || ( linux && amd64 ) || ( linux && arm64 ) || ( linux && 386 ) || windows )

package browser

import (
  "fmt"
  "sitebrush/pkg/config"
  webview "github.com/webview/webview_go"
)

func Run(config Config.Settings) {
  fmt.Printf("Browsing http://%s\n", config.LOCALHOST_LISTENER_ADDRESS)
  // Создаем новое окно WebView
  w := webview.New(false)

  //defer w.Destroy()
  defer func() {
    fmt.Println("Browsing finished.")
    w.Destroy()
  }()

  w.SetTitle(config.APP_NAME)
  w.SetSize(800, 600, webview.HintNone)
  //w.Navigate(config.LOCALHOST_LISTENER_ADDRESS)
  w.Navigate(fmt.Sprintf("http://%s?edit=true&%s=%s", config.LOCALHOST_LISTENER_ADDRESS, config.APP_NAME, config.VERSION))
  // Запускаем окно WebView
  w.Run()
}
