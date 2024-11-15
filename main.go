package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/nathansiegfrid/todolist/internal/auth"
	"github.com/nathansiegfrid/todolist/internal/todo"
	"github.com/nathansiegfrid/todolist/pkg/config"
	"github.com/nathansiegfrid/todolist/pkg/handler"
	"github.com/nathansiegfrid/todolist/pkg/logger"
	"github.com/nathansiegfrid/todolist/pkg/middleware"
	"github.com/nathansiegfrid/todolist/pkg/postgres"
	"github.com/nathansiegfrid/todolist/pkg/server"
	"github.com/nathansiegfrid/todolist/pkg/token"
)

func main() {
	// LOGGER
	// Use "go run main.go production" to enable JSON logging.
	useJSONLog := len(os.Args) > 1 && os.Args[1] == "production"
	logger.SetDefaultSlog(os.Stderr, useJSONLog)

	// CONFIG
	env := config.NewEnvLoader()
	var (
		serverPort  = env.OptionalInt("SERVER_PORT", 8080)
		postgresURL = env.MandatoryString("POSTGRES_URL")
		jwtSecret   = env.MandatoryString("JWT_SECRET")
	)
	if err := env.Validate(); err != nil {
		slog.Error(fmt.Sprintf("Config error: %s.", err))
		return
	}

	// DATABASE
	db, err := postgres.Connect(postgresURL)
	if err != nil {
		slog.Error(fmt.Sprintf("Database connection error: %s.", err))
		return
	}
	defer db.Close()

	// SCHEMA MIGRATION
	results, err := postgres.Migrate(db, "migrations")
	if err != nil {
		slog.Error(fmt.Sprintf("Schema migration error: %s.", err))
		return
	}
	for _, r := range results {
		slog.Info(fmt.Sprintf("Applied schema migration %s.", r.Source.Path))
	}

	// SERVICE HANDLERS
	jwtAuth := token.NewJWTAuth([]byte(jwtSecret))
	authHandler := auth.NewHandler(db, jwtAuth)
	todoHandler := todo.NewHandler(db)

	// ROUTER
	router := chi.NewRouter()
	router.NotFound(handler.NotFound)
	router.MethodNotAllowed(handler.MethodNotAllowed)
	router.Use(middleware.Heartbeat("/ping"))
	router.Use(middleware.CORSAllowOrigins("http://localhost:3000", "http://localhost:5173"))
	router.Use(middleware.RequestID)
	router.Use(middleware.VerifyAuth(jwtAuth))
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/v1", func(router chi.Router) {
		// Add public routes.
		router.Handle("/login", authHandler.HandleLoginRoute())
		router.Handle("/register", authHandler.HandleRegisterRoute())

		// Add private routes.
		router.Group(func(router chi.Router) {
			router.Use(middleware.RequireAuth)
			router.Handle("/verify-auth", authHandler.HandleVerifyAuthRoute())
			router.Handle("/todos", todoHandler.HandleTodosRoute())
			router.Handle("/todos/{id}", todoHandler.HandleTodosIDRoute())
		})
	})

	// RUN SERVER
	slog.Info(fmt.Sprintf("Listening on port %d.", serverPort))
	if err := server.ListenAndServe(serverPort, router); err != nil {
		slog.Error(fmt.Sprintf("HTTP server error: %s.", err))
		return
	}
}
