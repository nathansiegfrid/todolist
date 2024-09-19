package service

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type requestIDKey struct{}

const RequestIDHeader = "X-Request-ID"

func RequestIDFromContext(ctx context.Context) (string, bool) {
	reqID, ok := ctx.Value(requestIDKey{}).(string)
	return reqID, ok
}

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get(RequestIDHeader)
		if reqID == "" {
			reqID = uuid.New().String()
			r.Header.Set(RequestIDHeader, reqID)
		}

		ctx := context.WithValue(r.Context(), requestIDKey{}, reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
