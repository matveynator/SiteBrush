package webserver

import (
  "fmt"
  "net/http"
  "os"
  "strings"
)

// isFileExist проверяет, существует ли файл в локальной директории.
func isFileExist(fileName string) bool {
  _, err := os.Stat(fileName)
  return !os.IsNotExist(err)
}

// isUserLoggedIn проверяет, залогинен ли пользователь.
func isUserLoggedIn(request *http.Request) bool {
  // Здесь должна быть реализация проверки
  return true
}

// loginFunction обрабатывает запросы на авторизацию пользователя.
func loginFunction(responseWriter http.ResponseWriter, request *http.Request) {
  // Здесь реализация функции логина
  fmt.Fprint(responseWriter, "Login Function")
}

// editFunction обрабатывает запросы на редактирование файла.
func editFunction(responseWriter http.ResponseWriter, request *http.Request, fileName string) {}

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

// handleRequest обрабатывает входящие HTTP-запросы и перенаправляет их к соответствующим функциям обработки.
func handleRequest(responseWriter http.ResponseWriter, request *http.Request) {
  fileName := request.URL.Path[1:] // Получение имени файла из URL

  switch {
  case isFileExist(fileName) && !strings.HasSuffix(request.URL.RawQuery, "?"):
    http.ServeFile(responseWriter, request, fileName)

  case strings.HasSuffix(request.URL.RawQuery, "?login"):
    loginFunction(responseWriter, request)

  case isFileExist(fileName) && strings.HasSuffix(request.URL.RawQuery, "?edit") && isUserLoggedIn(request):
    editFunction(responseWriter, request, fileName)

  case isFileExist(fileName) && strings.HasSuffix(request.URL.RawQuery, "?delete") && isUserLoggedIn(request):
    deleteRevisionFunction(responseWriter, request, fileName)

  case isFileExist(fileName) && strings.HasSuffix(request.URL.RawQuery, "?revisions") && isUserLoggedIn(request):
    showRevisionsFunction(responseWriter, request, fileName)

  case isFileExist(fileName) && strings.HasSuffix(request.URL.RawQuery, "?subpages") && isUserLoggedIn(request):
    showSubpagesFunction(responseWriter, request, fileName)

  case isFileExist(fileName) && strings.HasSuffix(request.URL.RawQuery, "?properties") && isUserLoggedIn(request):
    editPropertiesFunction(responseWriter, request, fileName)

  case strings.HasSuffix(request.URL.RawQuery, "?freeze") && isUserLoggedIn(request):
    freezeSiteFunction(responseWriter, request)

  case strings.HasSuffix(request.URL.RawQuery, "?unfreeze") && isUserLoggedIn(request):
    unfreezeSiteFunction(responseWriter, request)

  case strings.HasSuffix(request.URL.RawQuery, "?backup") && isUserLoggedIn(request):
    backupSiteFunction(responseWriter, request)

  case strings.HasSuffix(request.URL.RawQuery, "?profile") && isUserLoggedIn(request):
    showProfileFunction(responseWriter, request)

  case strings.HasSuffix(request.URL.RawQuery, "?logout") && isUserLoggedIn(request):
    logoutFunction(responseWriter, request)

  default:
    http.Error(responseWriter, "Not Found", http.StatusNotFound)
  }
}

