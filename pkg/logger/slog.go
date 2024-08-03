package logger

import (
	"io"
	"log/slog"
)

func SetupLogger(env string, output io.Writer) *slog.Logger {
	var logger *slog.Logger
	switch env {
	case "local":
		logger = slog.New(slog.NewJSONHandler(output, &slog.HandlerOptions{Level: slog.LevelDebug})).With(slog.String("env", env))
	case "development":
		logger = slog.New(slog.NewJSONHandler(output, &slog.HandlerOptions{Level: slog.LevelDebug})).With(slog.String("env", env))
	case "production":
		logger = slog.New(slog.NewJSONHandler(output, &slog.HandlerOptions{Level: slog.LevelInfo})).With(slog.String("env", env))
	}

	return logger
}

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "Error",
		Value: slog.StringValue(err.Error()),
	}
}
