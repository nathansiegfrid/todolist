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

const gracefulShutdownTimeout = 10 * time.Second

// ListenAndServe starts an HTTP server and with graceful shutdown.
func ListenAndServe(port int, router http.Handler) error {
	svr := http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", port),
		Handler: router,
	}

	// Create a context that listens for interrupt/terminate signals.
	signalCtx, signalCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer signalCancel()

	// Run HTTP server.
	var err error
	go func() {
		err = svr.ListenAndServe()
		signalCancel() // Stop listening for signals.
	}()

	// Wait for interrupt/terminate signals.
	<-signalCtx.Done()
	// Return any error that occurred during ListenAndServe.
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("listen and serve: %w", err)
	}

	// Start graceful shutdown.
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
	defer shutdownCancel()
	if err := svr.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shut down: %w", err)
	}
	return nil
}
