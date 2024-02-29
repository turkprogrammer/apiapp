package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"

	"apiapp/internal/response"
	"apiapp/internal/validator"
)

// reportServerError логгирует подробную информацию о серверной ошибке, включая атрибуты запроса и стек вызовов.
func (app *application) reportServerError(r *http.Request, err error) {
	var (
		message = err.Error()
		method  = r.Method
		url     = r.URL.String()
		trace   = string(debug.Stack())
	)

	// Создание структурированной записи лога с деталями запроса и стеком вызовов.
	requestAttrs := slog.Group("request", "method", method, "url", url)
	app.logger.Error(message, requestAttrs, "trace", trace)
}

// errorMessage генерирует JSON-ответ с сообщением об ошибке и необязательными заголовками, логгирует ошибку
// и устанавливает соответствующий HTTP-статус код.
func (app *application) errorMessage(w http.ResponseWriter, r *http.Request, status int, message string, headers http.Header) {
	// Заглавная буква первого символа сообщения об ошибке.
	message = strings.ToUpper(message[:1]) + message[1:]

	// Генерация JSON-ответа с сообщением об ошибке и заголовками.
	err := response.JSONWithHeaders(w, status, map[string]string{"Error": message}, headers)
	if err != nil {
		// Если произошла ошибка при генерации JSON-ответа, логгирование ошибки и установка HTTP-статуса во внутреннюю ошибку сервера.
		app.reportServerError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// serverError обрабатывает внутренние ошибки сервера, регистрируя ошибку и предоставляя общее сообщение об ошибке в ответе.
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	// Регистрация ошибки сервера.
	app.reportServerError(r, err)

	// Предоставление общего сообщения об ошибке в ответе.
	message := "Сервер столкнулся с проблемой и не может обработать ваш запрос"
	app.errorMessage(w, r, http.StatusInternalServerError, message, nil)
}

// notFound обрабатывает запросы к несуществующим ресурсам, предоставляя ответ 404 Not Found.
func (app *application) notFound(w http.ResponseWriter, r *http.Request) {
	message := "Запрашиваемый ресурс не найден"
	app.errorMessage(w, r, http.StatusNotFound, message, nil)
}

// methodNotAllowed обрабатывает запросы с неподдерживаемыми HTTP-методами, предоставляя ответ 405 Method Not Allowed.
func (app *application) methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	// Генерация сообщения о неподдерживаемом HTTP-методе.
	message := fmt.Sprintf("Метод %s не поддерживается для данного ресурса", r.Method)
	app.errorMessage(w, r, http.StatusMethodNotAllowed, message, nil)
}

// badRequest обрабатывает запросы с некорректным синтаксисом или недопустимыми параметрами, предоставляя ответ 400 Bad Request.
func (app *application) badRequest(w http.ResponseWriter, r *http.Request, err error) {
	// Генерация ответа с сообщением об ошибке.
	app.errorMessage(w, r, http.StatusBadRequest, err.Error(), nil)
}

// failedValidation обрабатывает запросы с ошибками валидации, предоставляя ответ 422 Unprocessable Entity.
func (app *application) failedValidation(w http.ResponseWriter, r *http.Request, v validator.Validator) {
	// Генерация JSON-ответа с ошибками валидации.
	err := response.JSON(w, http.StatusUnprocessableEntity, v)
	if err != nil {
		// Если произошла ошибка при генерации JSON-ответа, логгирование ошибки и установка HTTP-статуса во внутреннюю ошибку сервера.
		app.serverError(w, r, err)
	}
}

// basicAuthenticationRequired обрабатывает запросы, требующие базовой аутентификации, но не содержащие действительных учетных данных.
// Предоставляет ответ 401 Unauthorized с необходимыми заголовками для базовой аутентификации.
func (app *application) basicAuthenticationRequired(w http.ResponseWriter, r *http.Request) {
	// Установка заголовка WWW-Authenticate для базовой аутентификации.
	headers := make(http.Header)
	headers.Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)

	// Генерация ответа с ошибкой доступа и соответствующими заголовками.
	message := "Необходима аутентификация для доступа к ресурсу"
	app.errorMessage(w, r, http.StatusUnauthorized, message, headers)
}
