package service

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

func RequestIDFromContext(ctx context.Context) (string, bool) {
	reqID, ok := ctx.Value(requestIDContextKey).(string)
	return reqID, ok
}

func ContextWithUserID(ctx context.Context, userID uuid.NullUUID) context.Context {
	return context.WithValue(ctx, userIDContextKey, userID)
}

func UserIDFromContext(ctx context.Context) (uuid.NullUUID, bool) {
	userID, ok := ctx.Value(userIDContextKey).(uuid.UUID)
	return uuid.NullUUID{UUID: userID, Valid: ok}, ok
}