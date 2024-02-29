package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// routes возвращает HTTP-обработчик, представляющий конфигурацию маршрутов приложения.
// Использует библиотеку Gorilla Mux для управления маршрутами и их обработчиками.
func (app *application) routes() http.Handler {
	// Создание нового маршрутизатора с использованием Gorilla Mux.
	mux := mux.NewRouter()

	// Установка обработчиков для случаев, когда маршрут не найден или используется неподдерживаемый HTTP-метод.
	mux.NotFoundHandler = http.HandlerFunc(app.notFound)
	mux.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowed)

	// Использование middleware для восстановления от паник в обработчиках.
	mux.Use(app.recoverPanic)

	// Установка обработчика для маршрута "/status" с методом GET.
	mux.HandleFunc("/status", app.status).Methods("GET")

	// Создание подмаршрута для защищенных ресурсов, используя middleware для базовой аутентификации.
	protectedRoutes := mux.NewRoute().Subrouter()
	protectedRoutes.Use(app.requireBasicAuthentication)
	// Установка обработчика для защищенного маршрута "/basic-auth-protected" с методом GET.
	protectedRoutes.HandleFunc("/basic-auth-protected", app.protected).Methods("GET")

	// Возврат маршрутизатора как HTTP-обработчика.
	return mux
}
