package logger

import (
	"log/slog"
	"os"
)

func InitLogger(conf Logger) *slog.Logger {
	lvl := parseLevel(conf.Level)

	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
	})

	return slog.New(h)
}
