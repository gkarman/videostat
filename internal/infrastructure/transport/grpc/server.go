package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

type Config struct {
	Addr string
}

type Server struct {
	grpcServer *grpc.Server
	lis        net.Listener
	log        *slog.Logger
}

func (s *Server) Registrar() grpc.ServiceRegistrar {
	return s.grpcServer
}

func NewServer(log *slog.Logger, cfg Config, opts ...grpc.ServerOption) (*Server, error) {
	lis, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return nil, fmt.Errorf("listen: %w", err)
	}

	s := grpc.NewServer(opts...)

	return &Server{
		grpcServer: s,
		lis:        lis,
		log:        log,
	}, nil
}

func (s *Server) Start() {
	go func() {
		s.log.Info("gRPC server started", "addr", s.lis.Addr())
		if err := s.grpcServer.Serve(s.lis); err != nil {
			s.log.Error("gRPC server failed", "error", err)
		}
	}()
}

func (s *Server) Stop(ctx context.Context) error {
	s.log.Info("stopping gRPC server")
	done := make(chan struct{})
	go func() {
		s.grpcServer.GracefulStop()
		close(done)
	}()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		s.grpcServer.Stop()
		return ctx.Err()
	}
}
