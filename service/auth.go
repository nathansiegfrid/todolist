package service

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type userIDKey struct{}

func ContextWithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey{}, userID)
}

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(userIDKey{}).(uuid.UUID)
	return userID, ok
}

func AuthenticationMiddleware(next http.Handler) http.Handler {
	panic("not implemented")
}
