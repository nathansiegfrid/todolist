package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/nathansiegfrid/todolist-go/config"
	"github.com/nathansiegfrid/todolist-go/service"
	"github.com/nathansiegfrid/todolist-go/service/todo"
	"github.com/pressly/goose/v3"
)

func main() {
	svcName := flag.String("service-name", "todolist", "Specifies the service name included in log output.")
	devMode := flag.Bool("development", false, "Output logs in human-readable format instead of JSON.")
	flag.Parse()

	// INIT GLOBAL LOGGER
	var logger *slog.Logger
	if *devMode {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	} else {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}
	slog.SetDefault(logger.With("service", *svcName))

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
		err := db.Ping()
		if err == nil {
			slog.Info("connected to database")
			break
		}
		if time.Since(start) > timeout {
			slog.Error(fmt.Sprintf("error verifying database connection: %s", err))
			return
		}
		time.Sleep(sleep)
		sleep *= 2
	}

	// RUN DATABASE MIGRATIONS
	// Goose supports out of order migration with "allow missing" option.
	err = goose.Up(db, "migration", goose.WithAllowMissing())
	if err != nil {
		slog.Error(fmt.Sprintf("error running database migrations: %s", err))
		return
	}

	// ADD SERVICE HANDLERS TO HTTP ROUTER
	router := chi.NewRouter()
	router.Use(middleware.Heartbeat("/ping"))
	router.Use(service.RequestIDMiddleware)
	router.Use(service.LoggerMiddleware)
	router.Route("/api/v1", func(router chi.Router) {
		// Add public routes.
		// TODO: Implement these routes!
		router.Post("/login", func(w http.ResponseWriter, r *http.Request) {})
		router.Post("/register", func(w http.ResponseWriter, r *http.Request) {})

		// Add private routes.
		router.Group(func(router chi.Router) {
			// TODO: Use auth middleware here!
			todoHandler := todo.NewHandler(db)
			router.Mount("/todo", todoHandler.HTTPHandler())
		})
	})

	// RUN HTTP SERVER
	svr := http.Server{
		Addr:    fmt.Sprintf("%s:%d", c.APIHost, c.APIPort),
		Handler: router,
	}

	slog.Info(fmt.Sprintf("HTTP server listening on %s", svr.Addr))
	err = svr.ListenAndServe()
	if err != nil {
		slog.Error(fmt.Sprintf("error running HTTP server: %s", err))
		return
	}
	// TODO: Implement graceful shutdown!
}
