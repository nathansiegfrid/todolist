package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/nathansiegfrid/todolist-go/service"
)

// responseWriter is a wrapper for http.ResponseWriter that captures the written HTTP status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode  int
	wroteHeader bool
}

func (w *responseWriter) WriteHeader(statusCode int) {
	if w.wroteHeader {
		return
	}
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
	w.wroteHeader = true
}

// Logger should be used after RequestID and Authentication middlewares.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add request ID and user ID to logs.
		logger := service.Logger(r.Context())

		// Recover and log panics.
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				logger.Error(fmt.Sprintf("panic: %s", err), "trace", string(debug.Stack()))
			}
		}()

		start := time.Now()
		ww := &responseWriter{ResponseWriter: w} // Use wrapper to get HTTP status code.
		next.ServeHTTP(ww, r)

		logger.Info(
			fmt.Sprintf("response: %d %s", ww.statusCode, http.StatusText(ww.statusCode)),
			"status", ww.statusCode,
			"method", r.Method,
			"path", r.URL.EscapedPath(),
			"duration", time.Since(start),
		)
	})
}
