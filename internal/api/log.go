package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

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

	var res ErrorResponse
	if errors.As(err, &res) && res.StatusCode < 500 {
		Logger(ctx).Info("Client error: " + res.Message)
	} else {
		Logger(ctx).Error(fmt.Sprintf("Internal error: %s.", err))
	}
}
