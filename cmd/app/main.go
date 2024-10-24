package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lmittmann/tint"
	"github.com/nathansiegfrid/todolist-go/config"
	"github.com/nathansiegfrid/todolist-go/middleware"
	"github.com/nathansiegfrid/todolist-go/service/auth"
	"github.com/nathansiegfrid/todolist-go/service/todo"
	"github.com/pressly/goose/v3"
)

func main() {
	// These flags are optional and have no effect on services.
	svcName := flag.String("service", "todolist", "Specifies the service name included in log output.")
	devMode := flag.Bool("development", false, "Output logs in human-readable format instead of JSON.")
	flag.Parse()

	// INIT GLOBAL LOGGER
	var logger *slog.Logger
	if *devMode {
		logger = slog.New(tint.NewHandler(os.Stderr, nil))
	} else {
		logger = slog.New(slog.NewJSONHandler(os.Stderr, nil))
	}
	host, _ := os.Hostname()
	slog.SetDefault(logger.With("service", *svcName, "host", host))

	// LOAD APPLICATION CONFIG
	c, err := config.Load()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	// CONNECT TO DATABASE
	db, err := sql.Open("pgx", fmt.Sprintf(
		// Added single quotes to accomodate empty values.
		// https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING
		"host='%s' port='%d' user='%s' password='%s' dbname='%s' sslmode='%s' sslrootcert='%s'",
		c.PGHost, c.PGPort, c.PGUser, c.PGPassword, c.PGDatabase, c.PGSSLMode, c.PGRootCertLoc,
	))
	if err != nil {
		slog.Error(fmt.Sprintf("error setting up database connection: %s", err))
		return
	}

	// Verify DB connection. If error, retry with exponential backoff.
	slog.Info("verifying database connection...")
	start, sleep, timeout := time.Now(), time.Second, 30*time.Second
	for {
		if err := db.Ping(); err == nil {
			slog.Info("connected to database")
			break
		} else if time.Since(start) > timeout {
			slog.Error(fmt.Sprintf("error verifying database connection: %s", err))
			return // Exit if timeout reached.
		}
		time.Sleep(sleep)
		sleep *= 2
	}

	// RUN DATABASE MIGRATIONS
	// Goose supports out of order migration with "allow missing" option.
	if err := goose.Up(db, "migration", goose.WithAllowMissing()); err != nil {
		slog.Error(fmt.Sprintf("error running database migrations: %s", err))
		return
	}

	// INIT SERVICES
	jwtService := auth.NewJWTService([]byte(c.JWTSecret))
	authHandler := auth.NewHandler(db, jwtService)
	todoHandler := todo.NewHandler(db)

	// ADD SERVICE HANDLERS TO HTTP ROUTER
	router := chi.NewRouter()
	router.Use(middleware.Heartbeat("/ping"))
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Authorization"},
		AllowCredentials: true,
	}))
	router.Use(middleware.RequestID)
	router.Use(middleware.VerifyAuth(jwtService))
	router.Use(middleware.Logger)

	router.Route("/api/v1", func(router chi.Router) {
		// Add public routes.
		router.Handle("/login", authHandler.HandleLoginRoute())
		router.Handle("/register", authHandler.HandleRegisterRoute())

		// Add private routes.
		router.Group(func(router chi.Router) {
			router.Use(middleware.RequireAuth)
			router.Handle("/verify-auth", authHandler.HandleVerifyAuthRoute())
			router.Handle("/todo", todoHandler.HandleTodoRoute())
			router.Handle("/todo/{id}", todoHandler.HandleTodoIDRoute())
		})
	})

	// Change default 404 and 405 handlers.
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found.", http.StatusNotFound)
	})
	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
	})

	// RUN HTTP SERVER
	svr := http.Server{
		Addr:    fmt.Sprintf("%s:%d", c.APIHost, c.APIPort),
		Handler: router,
	}

	go func() {
		slog.Info(fmt.Sprintf("HTTP server listening on %s", svr.Addr))
		if err := svr.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			slog.Error(fmt.Sprintf("error running HTTP server: %s", err))
		}
	}()

	// Wait for interrupt or terminate signal.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// GRACEFUL SHUTDOWN
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	slog.Info("shutting down HTTP server...")
	if err := svr.Shutdown(ctx); err != nil {
		slog.Error(fmt.Sprintf("error shutting down HTTP server: %s", err))
		return
	}
	slog.Info("HTTP server shut down gracefully")
}
