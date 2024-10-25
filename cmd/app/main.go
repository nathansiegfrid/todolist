package main

import (
	"log/slog"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/nathansiegfrid/todolist-go/config"
	"github.com/nathansiegfrid/todolist-go/middleware"
	"github.com/nathansiegfrid/todolist-go/pkg/database"
	"github.com/nathansiegfrid/todolist-go/pkg/logger"
	"github.com/nathansiegfrid/todolist-go/pkg/server"
	"github.com/nathansiegfrid/todolist-go/service"
	"github.com/nathansiegfrid/todolist-go/service/auth"
	"github.com/nathansiegfrid/todolist-go/service/todo"
)

func main() {
	// LOGGER
	useJSONLog := len(os.Args) > 1 && os.Args[1] == "production"
	logger.Setup(useJSONLog)

	// CONFIG
	c, err := config.Load()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	// DATABASE
	db, err := database.Connect(c.PGString())
	if err != nil {
		slog.Error(err.Error())
		return
	}
	if err := database.Migrate(db, "migration"); err != nil {
		slog.Error(err.Error())
		return
	}

	// SERVICE HANDLERS
	jwtService := auth.NewJWTService([]byte(c.JWTSecret))
	authHandler := auth.NewHandler(db, jwtService)
	todoHandler := todo.NewHandler(db)

	// ROUTER
	router := chi.NewRouter()
	router.NotFound(service.NotFound)
	router.MethodNotAllowed(service.MethodNotAllowed)
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
	router.Use(middleware.Recoverer)

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

	// RUN SERVER
	server.ListenAndServe(c.APIAddr(), router)
}
