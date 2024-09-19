package service

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"
)

// Logger returns the default slog.Logger with additional information from context.Context.
func Logger(ctx context.Context) *slog.Logger {
	logger := slog.Default()
	if reqID, ok := RequestIDFromContext(ctx); ok {
		logger = logger.With("requestID", reqID)
	}
	if userID, ok := UserIDFromContext(ctx); ok {
		logger = logger.With("userID", userID)
	}
	return logger
}

func LogError(ctx context.Context, err error) {
	Logger(ctx).Error(err.Error())
}

func LogInternalError(ctx context.Context, err error) {
	if StatusCode(err) != http.StatusInternalServerError {
		return
	}
	LogError(ctx, err)
}

// responseWriter is a wrapper for http.ResponseWriter.
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

// LoggerMiddleware should be used after RequestIDMiddleware and AuthenticationMiddleware.
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := Logger(r.Context())

		// Recover and log panics.
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				logger.Error(fmt.Sprintf("panic: %s", err), "trace", string(debug.Stack()))
			}
		}()

		start := time.Now()
		ww := &responseWriter{ResponseWriter: w}
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
