package service

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

// Logger returns the default slog.Logger with additional information from context.Context.
func Logger(ctx context.Context) *slog.Logger {
	rid := RequestIDFromContext(ctx)
	uid := UserIDFromContext(ctx)

	return slog.Default().
		With("request_id", rid).
		With("user_id", uuid.NullUUID{UUID: uid, Valid: uid != uuid.Nil})
}

func LogInfo(ctx context.Context, err error) {
	Logger(ctx).Info(err.Error())
}

func LogError(ctx context.Context, err error) {
	Logger(ctx).Error(err.Error())
}

func LogErrorInternal(ctx context.Context, err error) {
	if ErrorStatusCode(err) == http.StatusInternalServerError {
		LogError(ctx, err)
	}
}
