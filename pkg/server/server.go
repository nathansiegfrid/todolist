package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func ListenAndServe(addr string, router http.Handler) {
	svr := http.Server{Addr: addr, Handler: router}

	go func() {
		// Run HTTP server.
		slog.Info(fmt.Sprintf("HTTP server listening on %s", svr.Addr))
		if err := svr.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			slog.Error(fmt.Sprintf("error running HTTP server: %s", err))
		}
	}()

	// Wait for interrupt or terminate signal.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Set timeout for graceful shutdown.
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	slog.Info("shutting down HTTP server...")
	if err := svr.Shutdown(ctx); err != nil {
		slog.Error(fmt.Sprintf("error shutting down HTTP server: %s", err))
	}
}
