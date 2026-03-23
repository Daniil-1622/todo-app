package todo_app

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port, // адрес, например :8080
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,          // максимальный размер заголовка = 1 MB
		ReadTimeout:    10 * time.Second, // таймаут на чтение = 10 секунд
		WriteTimeout:   10 * time.Second, // таймаут на запись = 10 секунд
	}

	return s.httpServer.ListenAndServe() // Запускает сервер — он начинает слушать входящие HTTP запросы на указанном порту.
}

// Метод Shutdown, принимает в аргументы context - который задает deadline
func (s *Server) Shutdown(ctx context.Context) error {
	/*
		Shutdown - стандартный метод Go, который:

		1.Перестаёт принимать новые соединения
		2.Ждёт завершения активных запросов (не обрывает их резко)
		3.Если контекст истекает раньше — возвращает ошибку
	*/
	return s.httpServer.Shutdown(ctx)
}
