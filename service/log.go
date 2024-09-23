package service

import (
	"context"
	"log/slog"
	"net/http"
)

// Logger returns the default slog.Logger with additional information from context.Context.
func Logger(ctx context.Context) *slog.Logger {
	logger := slog.Default()
	if reqID, ok := RequestIDFromContext(ctx); ok {
		logger = logger.With("requestID", reqID)
	}
	if userID, ok := UserIDFromContext(ctx); ok {
		logger = logger.With("userID", userID)
	}
	return logger
}

func LogInfo(ctx context.Context, err error) {
	Logger(ctx).Info(err.Error())
}

func LogError(ctx context.Context, err error) {
	Logger(ctx).Error(err.Error())
}

func LogInternalError(ctx context.Context, err error) {
	if ErrorStatusCode(err) != http.StatusInternalServerError {
		return
	}
	LogError(ctx, err)
}
