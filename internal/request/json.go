//Этот код предоставляет функции для декодирования JSON-тела HTTP-запроса в структуру,
//а также обеспечивает обработку различных ошибок, которые могут возникнуть при этом процессе.

package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// DecodeJSON декодирует JSON-тело HTTP-запроса в структуру, переданную в параметре dst.
// В случае ошибок валидации, ошибка возвращается с соответствующим сообщением.
func DecodeJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	return decodeJSON(w, r, dst, false)
}

// DecodeJSONStrict выполняет строгую декодировку JSON-тела HTTP-запроса в структуру, переданную в параметре dst.
// Строгая декодировка запрещает неизвестные поля в JSON-теле.
func DecodeJSONStrict(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	return decodeJSON(w, r, dst, true)
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst interface{}, disallowUnknownFields bool) error {
	// Ограничиваем размер JSON-тела, чтобы предотвратить атаки на переполнение буфера.
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)

	// Включаем или отключаем строгую проверку неизвестных полей в JSON-теле.
	if disallowUnknownFields {
		dec.DisallowUnknownFields()
	}

	// Декодируем JSON-тело в переданную структуру.
	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		case errors.As(err, &invalidUnmarshalError):
			// Исключение в случае некорректного использования API.
			panic(err)

		default:
			return err
		}
	}

	// Проверяем, что в JSON-теле содержится только одно значение.
	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}
