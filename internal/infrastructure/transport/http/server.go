package http

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"
)

type Config struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Server struct {
	httpServer *http.Server
	log        *slog.Logger
}

func NewServer(log *slog.Logger, handler http.Handler, cfg Config) *Server {
	srv := &http.Server{
		Addr:         cfg.Addr,
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	return &Server{
		httpServer: srv,
		log:        log,
	}
}

func (s *Server) Start() {
	go func() {
		s.log.Info("http server started", "addr", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.log.Error("http server failed", "error", err)
		}
	}()
}

func (s *Server) Stop(ctx context.Context) error {
	s.log.Info("stopping http server")
	return s.httpServer.Shutdown(ctx)
}
