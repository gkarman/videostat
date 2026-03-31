package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gkarman/demo/internal/app"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	worker, err := app.NewWorkerCron(ctx)
	if err != nil {
		slog.Error("worker cron build failed", "error", err)
		os.Exit(1)
	}

	if err := worker.Run(ctx); err != nil {
		slog.Error("worker cron failed", "error", err)
		os.Exit(1)
	}
}
