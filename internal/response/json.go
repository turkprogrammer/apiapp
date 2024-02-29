//Этот код предоставляет функции для отправки JSON-ответов с указанным статус-кодом, данными и заголовками.

package response

import (
	"encoding/json"
	"net/http"
)

// JSON отправляет JSON-ответ с указанным статус-кодом и данными.
// Если произойдет ошибка при маршалинге данных, функция возвращает эту ошибку.
func JSON(w http.ResponseWriter, status int, data interface{}) error {
	return JSONWithHeaders(w, status, data, nil)
}

// JSONWithHeaders отправляет JSON-ответ с указанным статус-кодом, данными и заголовками.
// Если произойдет ошибка при маршалинге данных или при установке заголовков, функция возвращает эту ошибку.
func JSONWithHeaders(w http.ResponseWriter, status int, data interface{}, headers http.Header) error {
	// Маршалинг данных в формат JSON с отступами для читаемости.
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	// Добавление символа новой строки в конец JSON для улучшения читаемости.
	js = append(js, '\n')

	// Установка переданных заголовков ответа.
	for key, value := range headers {
		w.Header()[key] = value
	}

	// Установка заголовка Content-Type как "application/json".
	w.Header().Set("Content-Type", "application/json")

	// Установка статус-кода ответа.
	w.WriteHeader(status)

	// Запись JSON-данных в тело ответа.
	w.Write(js)

	return nil
}
