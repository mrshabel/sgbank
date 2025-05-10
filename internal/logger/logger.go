package logger

import (
	"log/slog"
	"os"

	"github.com/mrshabel/sgbank/internal/config"
)

func New(env config.ENV) *slog.Logger {
	opts := slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelInfo,
	}

	// override log options in development
	if env == config.DEV {
		opts.Level = slog.LevelDebug
		opts.AddSource = true
	}

	// write logs to the handler
	logger := slog.New(slog.NewTextHandler(os.Stdout, &opts))
	// override application default logger
	slog.SetDefault(logger)
	return logger
}
