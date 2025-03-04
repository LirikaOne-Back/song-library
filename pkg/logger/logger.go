package logger

import (
	"context"
	"log/slog"
	"os"
)

// Logger - обертка над slog.Logger
type Logger struct {
	*slog.Logger
}

// NewLogger создает и настраивает новый экземпляр логгера
func NewLogger(level string) *Logger {
	var logLevel slog.Level
	switch level {
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

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	logger := slog.New(handler)
	return &Logger{logger}
}

// WithContext добавляет контекст к логгеру
func (l *Logger) WithContext(ctx context.Context) *slog.Logger {
	return l.Logger.With("requestID", ctx.Value("requestID"))
}
