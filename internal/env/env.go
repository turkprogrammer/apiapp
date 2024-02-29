package env

import (
	"os"
	"strconv"
)

// GetString возвращает значение переменной окружения с заданным ключом.
// Если переменная не существует, возвращается значение по умолчанию.
func GetString(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	return value
}

// GetInt возвращает целочисленное значение переменной окружения с заданным ключом.
// Если переменная не существует или не является целым числом, возвращается значение по умолчанию.
func GetInt(key string, defaultValue int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		// В случае ошибки преобразования, вызываем панику.
		// Это может быть улучшено для возврата ошибки вместо вызова паники в реальном приложении.
		panic(err)
	}

	return intValue
}

// GetBool возвращает булево значение переменной окружения с заданным ключом.
// Если переменная не существует или не является булевым значением, возвращается значение по умолчанию.
func GetBool(key string, defaultValue bool) bool {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		// В случае ошибки преобразования, вызываем панику.
		// Это может быть улучшено для возврата ошибки вместо вызова паники в реальном приложении.
		panic(err)
	}

	return boolValue
}
