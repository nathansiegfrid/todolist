package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/nathansiegfrid/todolist/internal/api"
)

// Recoverer recovers from panics and logs the error.
// It should be used after Logger middleware.
func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log with request ID and user ID from context.
				api.Logger(r.Context()).Error(fmt.Sprintf("panic: %s", err), "trace", string(debug.Stack()))
				// Respond with 500 Internal Server Error.
				api.WriteError(w, errors.New("unknown error"))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
