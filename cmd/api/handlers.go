package main

import (
	"net/http"

	"apiapp/internal/response"
)

// status обрабатывает запрос к эндпоинту /status, возвращая JSON-ответ с текущим статусом "OK".
func (app *application) status(w http.ResponseWriter, r *http.Request) {
	// Создание карты данных для JSON-ответа.
	data := map[string]string{
		"Status": "OK",
	}

	// Генерация JSON-ответа с данными и кодом статуса 200 (OK).
	err := response.JSON(w, http.StatusOK, data)
	if err != nil {
		// Если произошла ошибка при генерации JSON-ответа, вызываем обработчик серверной ошибки.
		app.serverError(w, r, err)
	}
}

// protected обрабатывает запрос к защищенному эндпоинту, возвращая строку "This is a protected handler".
func (app *application) protected(w http.ResponseWriter, r *http.Request) {
	// Просто записываем строку в ответ.
	w.Write([]byte("This is a protected handler"))
}
