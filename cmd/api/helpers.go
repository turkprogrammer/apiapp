package main

import (
	"fmt"
	"net/http"
)

// backgroundTask запускает фоновую задачу в виде горутины, ожидая её завершения.
// Принимает объект запроса `r` и функцию `fn`, которая выполняется в фоновой горутине.
// Использует sync.WaitGroup, чтобы отслеживать выполнение фоновой задачи.
func (app *application) backgroundTask(r *http.Request, fn func() error) {
	// Увеличение счетчика ожидаемых горутин в sync.WaitGroup.
	app.wg.Add(1)

	// Запуск новой горутины.
	go func() {
		// Уменьшение счетчика при завершении горутины (независимо от результата выполнения).
		defer app.wg.Done()

		// Защита от паники в горутине.
		defer func() {
			err := recover()
			if err != nil {
				// В случае паники логгируется информация об ошибке.
				app.reportServerError(r, fmt.Errorf("%s", err))
			}
		}()

		// Выполнение функции в фоновой горутине и обработка возможных ошибок.
		err := fn()
		if err != nil {
			// В случае ошибки логгируется информация об ошибке.
			app.reportServerError(r, err)
		}
	}()
}
