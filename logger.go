package watergun

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
		slog.String("app", "watergun"),
		slog.String("version", "v0.1.0"),
	)
)

func Logger() *slog.Logger {
	return logger
}
