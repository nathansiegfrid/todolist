package api

import (
	"context"
	"log/slog"
	"runtime/debug"

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

func LogError(ctx context.Context, err error) {
	if err == nil {
		return
	}
	if ErrorStatusCode(err) < 500 {
		// Client errors are logged as INFO.
		Logger(ctx).Info(err.Error())
	} else {
		// ERROR level logs, requires investigation/intervention.
		Logger(ctx).Error(err.Error(), "trace", string(debug.Stack()))
	}
}
