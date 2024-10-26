package api

import (
	"context"
	"fmt"
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

func LogInfo(ctx context.Context, msg any) {
	Logger(ctx).Info(fmt.Sprintf("%s", msg))
}

func LogError(ctx context.Context, msg any) {
	Logger(ctx).Error(fmt.Sprintf("%s", msg))
}

func LogErrorInternal(ctx context.Context, err error) {
	if err != nil && ErrorStatusCode(err) == http.StatusInternalServerError {
		LogError(ctx, err)
	}
}
