package logger

import (
	"io"
	"log/slog"

	"github.com/lmittmann/tint"
)

// SetDefaultSlog configures the global standard logger.
// This will affect all log output from `log` and `slog` package.
func SetDefaultSlog(w io.Writer, json bool) {
	var handler slog.Handler
	if json {
		handler = slog.NewJSONHandler(w, nil)
	} else {
		// Colored human-readable output.
		handler = tint.NewHandler(w, nil)
	}
	slogger := slog.New(handler)
	slog.SetDefault(slogger)
}
