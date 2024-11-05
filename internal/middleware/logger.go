package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/nathansiegfrid/todolist/internal/api"
)

// responseWriter is a wrapper for http.ResponseWriter that
// captures the written HTTP status code and byte size.
type responseWriter struct {
	http.ResponseWriter
	wroteHeader bool
	statusCode  int
	bytesOut    int // Reading the Content-Length header doesn't work.
}

func (w *responseWriter) WriteHeader(statusCode int) {
	if w.wroteHeader {
		return
	}
	w.ResponseWriter.WriteHeader(statusCode)
	w.wroteHeader = true
	w.statusCode = statusCode
}

func (w *responseWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err := w.ResponseWriter.Write(b)
	w.bytesOut += n
	return n, err
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
		api.Logger(r.Context()).Info(
			fmt.Sprintf("API response: %d %s.", ww.statusCode, http.StatusText(ww.statusCode)),
			"status", ww.statusCode,
			"method", r.Method,
			"path", r.URL.EscapedPath(),
			"bytes_in", r.ContentLength,
			"bytes_out", ww.bytesOut,
			"duration", time.Since(start),
		)
	})
}
