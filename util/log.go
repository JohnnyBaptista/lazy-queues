// Package util in log.go
package util

import (
	"fmt"
	"log/slog"
	"os"
)

var Log *slog.Logger

func init() {
	env := os.Getenv("APP_ENV")

	var level slog.Level
	fmt.Println()
	if env == "production" {
		level = slog.LevelInfo
	} else {
		level = slog.LevelDebug
	}

	var handler slog.Handler

	if env == "production" {
		// JSON estruturado para produção (fácil de ingerir no Datadog, Loki, etc.)
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})
	} else {
		// Texto legível para desenvolvimento
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})
	}

	Log = slog.New(handler)
}
