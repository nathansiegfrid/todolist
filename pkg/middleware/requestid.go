package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/nathansiegfrid/todolist/pkg/request"
)

const requestIDHeader = "X-Request-ID"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get(requestIDHeader)
		if reqID == "" {
			reqID = uuid.New().String()
			r.Header.Set(requestIDHeader, reqID)
		}
		ctx := request.ContextWithRequestID(r.Context(), reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
