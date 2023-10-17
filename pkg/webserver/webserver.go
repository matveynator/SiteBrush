package webserver

import (
  "fmt"
  "net/http"

  "sitebrush/pkg/config"
)

func webhandler(w http.ResponseWriter, r *http.Request) {
  // Получение всех параметров строки запроса
  values := r.URL.Query()

  // Проход по всем параметрам и вывод их
  for key, value := range values {
    // Поскольку значение параметра возвращает слайс, вы можете использовать value[0] для получения первого значения
    fmt.Fprintf(w, "Key: %s, Value: %s\n", key, value[0])
  }
}

func Run(config Config.Settings ) {
  http.HandleFunc("/", webhandler)
  err := http.ListenAndServe(config.WEB_LISTENER_ADDRESS, nil)
  if err != nil {
    fmt.Println("Ошибка при запуске веб-сервера:", err)
  }
}

