package interceptor

import (
	"context"
	"log/slog"
	"time"

	appLogger "github.com/gkarman/demo/internal/infrastructure/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func Logger(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		reqLog := log.With("transport", "grpc", "method", info.FullMethod)
		ctx = appLogger.WithLogger(ctx, reqLog)
		startedAt := time.Now()

		resp, err := handler(ctx, req)
		duration := time.Since(startedAt)
		code := status.Code(err)

		if err != nil {
			reqLog.Error("grpc request failed", "duration", duration, "grpc_code", code.String(), "error", err)
			return nil, err
		}

		reqLog.Info("grpc request handled", "duration", duration, "grpc_code", code.String())
		return resp, nil
	}
}
