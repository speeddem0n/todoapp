package todo

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(port string, handler http.Handler) error { // Метод запуска сервера
	s.httpServer = &http.Server{
		Addr:           ":" + port, // Server address
		Handler:        handler,    // Handler
		MaxHeaderBytes: 1 << 20,    // 1 MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error { // Метод остановки сервера
	return s.httpServer.Shutdown(ctx)
}
