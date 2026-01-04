package log

import (
	"log/slog"
	"os"

	"github.com/go-chi/httplog/v3"
)

var (
	logger = slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				ReplaceAttr: httplog.SchemaECS.Concise(true).ReplaceAttr,
			},
		),
	).With(
		slog.String("app", "Watergun"),
		slog.String("version", "v0.0.1"),
	)
)

func Logger() *slog.Logger {
	return logger
}

func Info(msg string, args ...any) {
	logger.Info(msg, args...)
}

func Error(msg string, args ...any) {
	logger.Error(msg, args...)
}
