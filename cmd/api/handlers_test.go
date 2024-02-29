package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Тестирование функции status.
func TestStatus(t *testing.T) {
	// Создание экземпляра приложения для теста.
	app := &application{}

	// Создание HTTP-запроса.
	req := httptest.NewRequest("GET", "/status", nil)

	// Создание записи для записи HTTP-ответа.
	w := httptest.NewRecorder()

	// Вызов функции status с записью ответа.
	app.status(w, req)

	// Проверка кода статуса ответа.
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// TODO: Проверка тела ответа и других ожидаемых результатов.
}

// Тестирование функции protected.
func TestProtected(t *testing.T) {
	// Создание экземпляра приложения для теста.
	app := &application{}

	// Создание HTTP-запроса.
	req := httptest.NewRequest("GET", "/protected", nil)

	// Создание записи для записи HTTP-ответа.
	w := httptest.NewRecorder()

	// Вызов функции protected с записью ответа.
	app.protected(w, req)

	// Проверка кода статуса ответа.
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// TODO: Проверка тела ответа и других ожидаемых результатов.
}
