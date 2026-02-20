package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strings"

	"github.com/tacokumo/tacoctl/cmd"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var logLevel slog.Level

	switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true,
	}))

	c := cmd.New()
	if err := c.ExecuteContext(ctx); err != nil {
		logger.ErrorContext(ctx, "failed to execute command", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
