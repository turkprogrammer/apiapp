package main

import (
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// recoverPanic возвращает middleware для восстановления от паники в хендлере.
// Обрабатывает панику, логгирует информацию об ошибке и продолжает выполнение следующего хендлера.
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				// В случае паники логгируется информация об ошибке, и выполнение продолжается.
				app.serverError(w, r, fmt.Errorf("%s", err))
			}
		}()

		// Вызов следующего хендлера в цепочке.
		next.ServeHTTP(w, r)
	})
}

// requireBasicAuthentication возвращает middleware, проверяющее наличие базовой аутентификации в запросе.
// Проверяет корректность учетных данных, в том числе сравнивает хеш пароля.
// В случае ошибки или несоответствия требованиям, вызывает соответствующий хендлер.
func (app *application) requireBasicAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получение учетных данных из заголовков базовой аутентификации.
		username, plaintextPassword, ok := r.BasicAuth()
		if !ok {
			// Если учетные данные отсутствуют, вызывается хендлер требования базовой аутентификации.
			app.basicAuthenticationRequired(w, r)
			return
		}

		// Проверка соответствия имени пользователя требуемому значению.
		if app.config.basicAuth.username != username {
			// В случае несоответствия вызывается хендлер требования базовой аутентификации.
			app.basicAuthenticationRequired(w, r)
			return
		}

		// Сравнение хеша пароля с переданным паролем.
		err := bcrypt.CompareHashAndPassword([]byte(app.config.basicAuth.hashedPassword), []byte(plaintextPassword))
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// В случае несоответствия вызывается хендлер требования базовой аутентификации.
			app.basicAuthenticationRequired(w, r)
			return
		case err != nil:
			// В случае других ошибок вызывается хендлер серверной ошибки.
			app.serverError(w, r, err)
			return
		}

		// Если все проверки успешны, вызывается следующий хендлер в цепочке.
		next.ServeHTTP(w, r)
	})
}
