package api

import (
	"context"
	"net/http"
	"song-library/pkg/logger"
	"time"
)

// Server представляет HTTP сервер приложения
type Server struct {
	httpServer *http.Server
	logger     *logger.Logger
}

// NewServer создает новый экземпляр сервера
func NewServer(router *Router, port string, logger *logger.Logger) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:           ":" + port,
			Handler:        router.GetEngine(),
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
		logger: logger,
	}
}

// Run запускает HTTP сервер
func (s *Server) Run() error {
	s.logger.Info("Запуск HTTP сервера", "port", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

// Shutdown останавливает HTTP сервер
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Остановка HTTP сервера")
	return s.httpServer.Shutdown(ctx)
}
