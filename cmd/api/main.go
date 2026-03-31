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

	a, err := app.NewApi(ctx)
	if err != nil {
		slog.Error("application build failed", "error", err)
		os.Exit(1)
	}

	if err := a.Run(ctx); err != nil {
		slog.Error("application failed", "error", err)
		os.Exit(1)
	}
}
