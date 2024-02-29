// Пакет main является входной точкой приложения.
package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"
	"sync"

	"apiapp/internal/env"
	"apiapp/internal/version"

	"github.com/lmittmann/tint"
)

// Функция main инициализирует логгер, запускает приложение и обрабатывает ошибки.
func main() {
	// Инициализация логгера с цветным выводом и уровнем отладки.
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug}))

	// Запуск приложения и обработка возможных ошибок.
	err := run(logger)
	if err != nil {
		// В случае ошибки логгирование вместе с трассировкой стека и завершение с ненулевым кодом статуса.
		trace := string(debug.Stack())
		logger.Error(err.Error(), "trace", trace)
		os.Exit(1)
	}
}

// Структура config содержит конфигурационные параметры для приложения.
type config struct {
	baseURL   string
	httpPort  int
	basicAuth struct {
		username       string
		hashedPassword string
	}
}

// Структура application инкапсулирует состояние приложения, включая конфигурацию, логгер и wait group.
type application struct {
	config config
	logger *slog.Logger
	wg     sync.WaitGroup
}

// Функция run инициализирует конфигурацию, парсит флаги командной строки и запускает HTTP-сервер.
func run(logger *slog.Logger) error {
	// Инициализация конфигурации с значениями по умолчанию или значениями из переменных окружения.
	var cfg config
	cfg.baseURL = env.GetString("BASE_URL", "http://localhost:4444")
	cfg.httpPort = env.GetInt("HTTP_PORT", 4444)
	cfg.basicAuth.username = env.GetString("BASIC_AUTH_USERNAME", "admin")
	cfg.basicAuth.hashedPassword = env.GetString("BASIC_AUTH_HASHED_PASSWORD", "$2a$10$jRb2qniNcoCyQM23T59RfeEQUbgdAXfR6S0scynmKfJa5Gj3arGJa")

	// Парсинг флагов командной строки, включая флаг для отображения версии.
	showVersion := flag.Bool("version", false, "отобразить версию и завершить программу")
	flag.Parse()

	// Если установлен флаг версии, вывод версии и завершение без запуска приложения.
	if *showVersion {
		fmt.Printf("версия: %s\n", version.Get())
		return nil
	}

	// Создание экземпляра приложения с сконфигурированными значениями и логгером.
	app := &application{
		config: cfg,
		logger: logger,
	}

	// Запуск обслуживания HTTP-запросов и обработка возможных ошибок.
	return app.serveHTTP()
}
