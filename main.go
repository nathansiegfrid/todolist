package main

import (
	"log/slog"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/nathansiegfrid/todolist/internal/api"
	"github.com/nathansiegfrid/todolist/internal/api/auth"
	"github.com/nathansiegfrid/todolist/internal/api/todo"
	"github.com/nathansiegfrid/todolist/internal/middleware"
	"github.com/nathansiegfrid/todolist/pkg/config"
	"github.com/nathansiegfrid/todolist/pkg/database"
	"github.com/nathansiegfrid/todolist/pkg/logger"
	"github.com/nathansiegfrid/todolist/pkg/server"
)

func main() {
	// LOGGER
	useJSONLog := len(os.Args) > 1 && os.Args[1] == "production"
	logger.SetDefaultSlog(os.Stderr, useJSONLog)

	// CONFIG
	c, err := config.Load()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	// DATABASE
	db, err := database.Connect(c.PostgresDSN())
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer db.Close()
	if err := database.Migrate(db, "migrations"); err != nil {
		slog.Error(err.Error())
		return
	}

	// SERVICE HANDLERS
	jwtService := auth.NewJWTService([]byte(c.JWTSecret))
	authHandler := auth.NewHandler(db, jwtService)
	todoHandler := todo.NewHandler(db)

	// ROUTER
	router := chi.NewRouter()
	router.NotFound(api.NotFound)
	router.MethodNotAllowed(api.MethodNotAllowed)
	router.Use(middleware.Heartbeat("/ping"))
	router.Use(middleware.CORSAllowOrigins("http://localhost:3000", "http://localhost:5173"))
	router.Use(middleware.RequestID)
	router.Use(middleware.VerifyAuth(jwtService))
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
	server.ListenAndServe(c.APIAddr(), router)
}
