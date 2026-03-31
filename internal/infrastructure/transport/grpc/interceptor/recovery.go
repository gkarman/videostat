package interceptor

import (
	"context"
	"fmt"
	"runtime/debug"

	appLogger "github.com/gkarman/demo/internal/infrastructure/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const maxStackLogBytes = 8 * 1024

func Recovery() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				log := appLogger.FromContext(ctx)
				stack, truncated := limitedStack()
				log.Error(
					"grpc panic recovered",
					"method", info.FullMethod,
					"panic", fmt.Sprint(r),
					"stack", stack,
					"stack_truncated", truncated,
				)
				err = status.Error(codes.Internal, "internal error")
			}
		}()

		return handler(ctx, req)
	}
}

func limitedStack() (string, bool) {
	stack := debug.Stack()
	if len(stack) <= maxStackLogBytes {
		return string(stack), false
	}
	return string(stack[:maxStackLogBytes]), true
}
