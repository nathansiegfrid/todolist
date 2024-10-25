package middleware

import (
	"fmt"
	"net/http"
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

// Logger logs HTTP responses.
// It should be used after RequestID and VerifyAuth middlewares.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Use wrapper to get HTTP status code.
		ww := &responseWriter{ResponseWriter: w}
		start := time.Now()
		next.ServeHTTP(ww, r)

		// Log with request ID and user ID from context.
		service.Logger(r.Context()).Info(
			fmt.Sprintf("response: %d %s", ww.statusCode, http.StatusText(ww.statusCode)),
			"status", ww.statusCode,
			"method", r.Method,
			"path", r.URL.EscapedPath(),
			"duration", time.Since(start),
		)
	})
}
