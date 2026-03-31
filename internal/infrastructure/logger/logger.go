package logger

import (
	"log/slog"
	"os"
)

type Config struct {
	Level string
}

func New(cfg Config) *slog.Logger {
	level := slog.LevelInfo

	if cfg.Level != "" {
		if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
			level = slog.LevelInfo
		}
	}

	opts := &slog.HandlerOptions{
		AddSource:   false,
		Level:       level,
		ReplaceAttr: nil,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	l := slog.New(handler)
	return l
}
