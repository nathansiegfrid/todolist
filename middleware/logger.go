package middleware

import (
	"errors"
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

// Logger should be used after RequestID and VerifyAuth middlewares.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip logging for OPTIONS requests.
		if r.Method == "OPTIONS" {
			next.ServeHTTP(w, r)
			return
		}

		// Logger adds request ID and user ID from context.
		logger := service.Logger(r.Context())
		// Use wrapper to get HTTP status code.
		ww := &responseWriter{ResponseWriter: w}
		start := time.Now()

		defer func() {
			// Recover panics.
			if err := recover(); err != nil {
				// Log error and stack trace.
				logger.Error(fmt.Sprintf("panic: %s", err), "trace", string(debug.Stack()))
				// Respond with 500 Internal Server Error.
				service.WriteError(ww, errors.New("unknown error"))
			}

			logger.Info(
				fmt.Sprintf("response: %d %s", ww.statusCode, http.StatusText(ww.statusCode)),
				"status", ww.statusCode,
				"method", r.Method,
				"path", r.URL.EscapedPath(),
				"duration", time.Since(start),
			)
		}()

		next.ServeHTTP(ww, r)
	})
}
