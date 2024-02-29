package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Константы для настройки тайм-аутов и периода завершения сервера.
const (
	defaultIdleTimeout    = time.Minute      // Время ожидания бездействия до закрытия соединения.
	defaultReadTimeout    = 5 * time.Second  // Время ожидания чтения данных из запроса.
	defaultWriteTimeout   = 10 * time.Second // Время ожидания записи данных в ответ.
	defaultShutdownPeriod = 30 * time.Second // Период завершения сервера при получении сигнала завершения.
)

// serveHTTP запускает HTTP-сервер с настройками, определенными в конфигурации приложения.
// Использует http.Server для управления сервером и Gorilla Mux для обработки маршрутов.
// Производит Graceful Shutdown сервера при получении сигнала завершения.
func (app *application) serveHTTP() error {
	// Создание нового HTTP-сервера с использованием настроек из конфигурации приложения.
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.httpPort),
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelWarn),
		IdleTimeout:  defaultIdleTimeout,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}

	// Канал для передачи ошибки завершения сервера.
	shutdownErrorChan := make(chan error)

	// Горутина для обработки сигналов завершения (SIGINT, SIGTERM).
	go func() {
		quitChan := make(chan os.Signal, 1)
		signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM)
		<-quitChan

		// Создание контекста с таймаутом для Graceful Shutdown.
		ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownPeriod)
		defer cancel()

		// Вызов Shutdown для Graceful Shutdown сервера и передача результата в канал.
		shutdownErrorChan <- srv.Shutdown(ctx)
	}()

	// Логгирование информации о запуске сервера.
	app.logger.Info("starting server", slog.Group("server", "addr", srv.Addr))

	// Запуск сервера для обработки входящих HTTP-запросов.
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// Ожидание завершения Graceful Shutdown и проверка ошибки завершения.
	err = <-shutdownErrorChan
	if err != nil {
		return err
	}

	// Логгирование информации о завершении работы сервера.
	app.logger.Info("stopped server", slog.Group("server", "addr", srv.Addr))

	// Ожидание завершения всех фоновых задач перед завершением.
	app.wg.Wait()

	return nil
}
