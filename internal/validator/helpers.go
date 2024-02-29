//Код содержит набор функций для валидации различных типов данных, таких как строки, числа, адреса электронной почты и URL

package validator

import (
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"

	"golang.org/x/exp/constraints"
)

// RgxEmail - регулярное выражение для валидации электронной почты.
var RgxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// NotBlank возвращает true, если значение строки не пусто после удаления пробелов по краям.
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MinRunes возвращает true, если длина строки (в кодовых точках Unicode) не меньше заданного значения n.
func MinRunes(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

// MaxRunes возвращает true, если длина строки (в кодовых точках Unicode) не больше заданного значения n.
func MaxRunes(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// Between возвращает true, если значение находится в пределах заданного диапазона [min, max].
func Between[T constraints.Ordered](value, min, max T) bool {
	return value >= min && value <= max
}

// Matches возвращает true, если строка соответствует заданному регулярному выражению.
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// In возвращает true, если значение находится в списке допустимых значений (safelist).
func In[T comparable](value T, safelist ...T) bool {
	for i := range safelist {
		if value == safelist[i] {
			return true
		}
	}
	return false
}

// AllIn возвращает true, если все значения в списке (values) находятся в списке допустимых значений (safelist).
func AllIn[T comparable](values []T, safelist ...T) bool {
	for i := range values {
		if !In(values[i], safelist...) {
			return false
		}
	}
	return true
}

// NotIn возвращает true, если значение не находится в списке недопустимых значений (blocklist).
func NotIn[T comparable](value T, blocklist ...T) bool {
	for i := range blocklist {
		if value == blocklist[i] {
			return false
		}
	}
	return true
}

// NoDuplicates возвращает true, если в списке нет дубликатов значений.
func NoDuplicates[T comparable](values []T) bool {
	uniqueValues := make(map[T]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(values) == len(uniqueValues)
}

// IsEmail возвращает true, если значение строки представляет собой корректный адрес электронной почты.
func IsEmail(value string) bool {
	if len(value) > 254 {
		return false
	}

	return RgxEmail.MatchString(value)
}

// IsURL возвращает true, если значение строки представляет собой корректный URL.
func IsURL(value string) bool {
	u, err := url.ParseRequestURI(value)
	if err != nil {
		return false
	}

	return u.Scheme != "" && u.Host != ""
}
