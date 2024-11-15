package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/nathansiegfrid/todolist/pkg/logger"
	"github.com/nathansiegfrid/todolist/pkg/response"
)

// Recoverer recovers from panics and logs the error.
// It should be used after Logger middleware.
func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error(
					r.Context(),
					fmt.Sprintf("Panic: %s.", err),
					"category", "internal_error",
					"trace", string(debug.Stack()),
				)
				response.WriteError(w, response.ErrorResponse{StatusCode: http.StatusInternalServerError})
			}
		}()
		next.ServeHTTP(w, r)
	})
}
