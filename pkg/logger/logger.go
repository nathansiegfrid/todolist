package logger

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

// Setup configures the global standard logger.
// This will affect all log output from `log` and `slog` package.
func Setup(json bool) {
	var handler slog.Handler
	if json {
		handler = slog.NewJSONHandler(os.Stderr, nil)
	} else {
		// Colored human-readable output.
		handler = tint.NewHandler(os.Stderr, nil)
	}
	slog.SetDefault(slog.New(handler))
}
