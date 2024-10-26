package api

import (
	"context"

	"github.com/google/uuid"
)

type contextKey int

const (
	requestIDContextKey contextKey = iota
	userIDContextKey
)

func ContextWithRequestID(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, requestIDContextKey, reqID)
}

func RequestIDFromContext(ctx context.Context) string {
	reqID, _ := ctx.Value(requestIDContextKey).(string)
	return reqID
}

func ContextWithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDContextKey, userID)
}

func UserIDFromContext(ctx context.Context) uuid.UUID {
	userID, _ := ctx.Value(userIDContextKey).(uuid.UUID)
	return userID
}
