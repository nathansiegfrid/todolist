package logger

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/nathansiegfrid/todolist/pkg/request"
)

// Logger returns the default slog.Logger with additional information from context.Context.
func Logger(ctx context.Context) *slog.Logger {
	logger := slog.Default()
	if v := request.RequestIDFromContext(ctx); v != "" {
		logger = logger.With("request_id", v)
	}
	if v := request.UserIDFromContext(ctx); v != uuid.Nil {
		logger = logger.With("user_id", v)
	}
	return logger
}

func Info(ctx context.Context, msg string, args ...any) {
	Logger(ctx).Info(msg, args...)
}

func Error(ctx context.Context, msg string, args ...any) {
	Logger(ctx).Error(msg, args...)
}

// func Fatal(ctx context.Context, msg string, args ...any) {
// 	Logger(ctx).Error(msg, args...)
// 	os.Exit(1)
// }
