package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func ListenAndServe(addr string, router http.Handler) {
	svr := http.Server{Addr: addr, Handler: router}

	// Create a context that listens for interrupt/terminate signals.
	signalCtx, signalCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer signalCancel()

	go func() {
		// Run HTTP server.
		slog.Info(fmt.Sprintf("HTTP server listening on %s", svr.Addr))
		if err := svr.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			slog.Error(fmt.Sprintf("error running HTTP server: %s", err))
		}
		signalCancel() // Stop listening for signals.
	}()

	// Wait for interrupt/terminate signals.
	<-signalCtx.Done()

	// Set timeout for graceful shutdown.
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer shutdownCancel()

	slog.Info("shutting down HTTP server...")
	if err := svr.Shutdown(shutdownCtx); err != nil {
		slog.Error(fmt.Sprintf("error shutting down HTTP server: %s", err))
	}
}
