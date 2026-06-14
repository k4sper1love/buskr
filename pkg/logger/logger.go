package logger

import (
	"log/slog"
	"os"
)

// Setup configures the global slog logger based on the environment.
func Setup(env string) {
	var h slog.Handler
	if env == "local" {
		h = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})
	} else {
		h = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})
	}
	slog.SetDefault(slog.New(h))
}
