package logger

import (
	"log/slog"
	"os"
)

func newLogger(cfg LogConfig) *slog.Logger {
	level := parseLevel(cfg.Level)

	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	return slog.New(h)
}

func InitLogger(project string) (*slog.Logger, error) {
	cfg, err := load(project)
	if err != nil {
		return nil, err
	}

	l := newLogger(cfg)

	return l, nil
}
