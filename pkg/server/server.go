package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

const gracefulShutdownTimeout = 15 * time.Second

// ListenAndServe starts an HTTP server with graceful shutdown.
func ListenAndServe(port int, router http.Handler) error {
	svr := http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", port),
		Handler: router,
	}

	// Run HTTP server.
	var err error
	signalCtx, signalCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	go func() {
		defer signalCancel()
		err = svr.ListenAndServe()
	}()
	<-signalCtx.Done() // Wait for interrupt/terminate signals.
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("listen and serve: %w", err)
	}

	// Shut down HTTP server gracefully.
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
	defer shutdownCancel()
	if err := svr.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shut down: %w", err)
	}
	return nil
}
