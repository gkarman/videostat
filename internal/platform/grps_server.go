package platform

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/gkarman/demo/internal/config"
	"github.com/gkarman/demo/internal/infrastructure/dispatcher"
	grpc2 "github.com/gkarman/demo/internal/infrastructure/transport/grpc"
	"github.com/gkarman/demo/internal/infrastructure/transport/grpc/interceptor"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

func NewGRPCServer(ctx context.Context, log *slog.Logger, db *pgxpool.Pool, cfg *config.Config, d *dispatcher.Dispatcher) (*grpc2.Server, error) {
	grpcConf := grpc2.Config{
		Addr: cfg.ServerGRPC.Addr,
	}
	grpcServer, err := grpc2.NewServer(
		ctx,
		log,
		grpcConf,
		// Спан на каждый gRPC-вызов, с продолжением трейса из метаданных запроса.
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			interceptor.Recovery(),
			interceptor.Logger(log),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("create gRPC server with interceptors: %w", err)
	}
	grpc2.RegisterServices(grpcServer, log, db, d)

	return grpcServer, nil
}
