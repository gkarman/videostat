package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gkarman/demo/internal/app"
	"github.com/gkarman/demo/internal/platform"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	// Отправить накопленные спаны перед выходом (os.Exit ниже defer не выполнит — на ошибке старта это не важно).
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		platform.ShutdownTracing(shutdownCtx)
	}()

	worker, err := app.NewWorkerCore(ctx)
	if err != nil {
		slog.Error("worker build failed", "error", err)
		os.Exit(1)
	}

	if err := worker.Run(ctx); err != nil {
		slog.Error("worker failed", "error", err)
		os.Exit(1)
	}
}
