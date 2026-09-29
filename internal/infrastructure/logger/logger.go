package logger

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	Level string
	// Dir — папка для файла логов. Пусто — логи только в stdout.
	Dir string
	// Service — имя сервиса, из него получается имя файла: <Dir>/<Service>.log.
	Service string
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

	handler := slog.NewJSONHandler(output(cfg), opts)
	l := slog.New(handler)
	return l
}

// output — stdout, а если задан Dir, то ещё и файл: из него логи забирает агент (Alloy) и отправляет в Loki.
func output(cfg Config) io.Writer {
	if cfg.Dir == "" || cfg.Service == "" {
		return os.Stdout
	}

	// lumberjack сам ротирует файл: при достижении MaxSize переименовывает его в
	// <Service>-<время>.log и начинает новый, старые удаляет. Иначе файл рос бы бесконечно.
	file := &lumberjack.Logger{
		Filename:   filepath.Join(cfg.Dir, cfg.Service+".log"),
		MaxSize:    50, // МБ
		MaxBackups: 3,
		MaxAge:     7, // дней
	}
	return io.MultiWriter(os.Stdout, file)
}
