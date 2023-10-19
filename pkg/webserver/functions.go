package webserver

import (
  "fmt"
  "strings"
  "net/http"
  "os"
  "sitebrush/pkg/config"
)

// checkFileExist проверяет, существует ли файл в локальной директории.
func checkFileExist(fileName string) bool {
  _, err := os.Stat(fileName)
  return !os.IsNotExist(err)
}

// checkUserLoggedIn проверяет, залогинен ли пользователь.
func checkUserLoggedIn(request *http.Request) bool {
  // Здесь должна быть реализация проверки
  return true
}

// loginFunction обрабатывает запросы на авторизацию пользователя.
func loginFunction(responseWriter http.ResponseWriter, request *http.Request) {
  // Здесь реализация функции логина
  fmt.Fprint(responseWriter, "Login Function")
}

// editFunction обрабатывает запросы на редактирование файла.
func editFunction(responseWriter http.ResponseWriter, request *http.Request, fileName string) {
  if checkUserLoggedIn(request) {
    fmt.Fprint(responseWriter, "Edit Function")
  } else {
    fmt.Fprint(responseWriter, "Not authorized to edit this page")
  }
}

// deleteRevisionFunction обрабатывает запросы на удаление последней ревизии файла.
func deleteRevisionFunction(responseWriter http.ResponseWriter, request *http.Request, fileName string) {}

// showRevisionsFunction обрабатывает запросы на отображение всех ревизий файла.
func showRevisionsFunction(responseWriter http.ResponseWriter, request *http.Request, fileName string) {}

// showSubpagesFunction обрабатывает запросы на отображение иерархического дерева файлов.
func showSubpagesFunction(responseWriter http.ResponseWriter, request *http.Request, fileName string) {}

// editPropertiesFunction обрабатывает запросы на редактирование свойств файла.
func editPropertiesFunction(responseWriter http.ResponseWriter, request *http.Request, fileName string) {}

// freezeSiteFunction обрабатывает запросы на заморозку сайта.
func freezeSiteFunction(responseWriter http.ResponseWriter, request *http.Request) {}

// unfreezeSiteFunction обрабатывает запросы на разморозку сайта.
func unfreezeSiteFunction(responseWriter http.ResponseWriter, request *http.Request) {}

// backupSiteFunction обрабатывает запросы на создание резервной копии сайта.
func backupSiteFunction(responseWriter http.ResponseWriter, request *http.Request) {}

// showProfileFunction обрабатывает запросы на отображение свойств учетной записи пользователя.
func showProfileFunction(responseWriter http.ResponseWriter, request *http.Request) {}

// logoutFunction обрабатывает запросы на выход из учетной записи.
func logoutFunction(responseWriter http.ResponseWriter, request *http.Request) {}


func handleRequest(config Config.Settings, responseWriter http.ResponseWriter, request *http.Request) {

  fileName := request.URL.Path // Получение имени файла из URL

  if fileName == "" {
    fileName = "/"
  }

  // Проверяем, заканчивается ли fileName на косую черту
  isDirectory := strings.HasSuffix(fileName, "/")

  if isDirectory {
    // Это каталог; вы можете обработать его соответственно
    // Например, вы можете предоставить индексный файл внутри каталога.
    if fileName == "/" {
      fileName = fmt.Sprintf("%s%s%s", config.WEB_FILE_PATH, fileName, config.WEB_INDEX_FILE)
    } else {
      fileName = fmt.Sprintf("%s%s/%s", config.WEB_FILE_PATH, fileName, config.WEB_INDEX_FILE)
    }
  } else {
    // Это файл; вы можете обработать его по-другому, если нужно
    fileName = fmt.Sprintf("%s%s", config.WEB_FILE_PATH, fileName)
  }

  queryParam := request.URL.RawQuery

  switch {
  case checkFileExist(fileName) && queryParam == "":
    http.ServeFile(responseWriter, request, fileName)

  case queryParam == "login":
    loginFunction(responseWriter, request)

  case queryParam == "edit":
    editFunction(responseWriter, request, fileName)

  case checkFileExist(fileName) && queryParam == "delete":
    deleteRevisionFunction(responseWriter, request, fileName)

  case checkFileExist(fileName) && queryParam == "revisions":
    showRevisionsFunction(responseWriter, request, fileName)

  case checkFileExist(fileName) && queryParam == "subpages":
    showSubpagesFunction(responseWriter, request, fileName)

  case checkFileExist(fileName) && queryParam == "properties":
    editPropertiesFunction(responseWriter, request, fileName)

  case queryParam == "freeze" && checkUserLoggedIn(request):
    freezeSiteFunction(responseWriter, request)

  case queryParam == "unfreeze" && checkUserLoggedIn(request):
    unfreezeSiteFunction(responseWriter, request)

  case queryParam == "backup" && checkUserLoggedIn(request):
    backupSiteFunction(responseWriter, request)

  case queryParam == "profile" && checkUserLoggedIn(request):
    showProfileFunction(responseWriter, request)

  case queryParam == "logout" && checkUserLoggedIn(request):
    logoutFunction(responseWriter, request)

  default:
    http.Error(responseWriter, "Not Found", http.StatusNotFound)
  }
}

