package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/nathansiegfrid/todolist-go/config"
	"github.com/nathansiegfrid/todolist-go/service/todo"
	"github.com/pressly/goose/v3"
)

func main() {
	// LOAD APPLICATION CONFIG
	c, err := config.Load()
	if err != nil {
		log.Fatalf("error loading config: %s", err)
	}

	// CONNECT TO DATABASE
	db, err := sql.Open("pgx", fmt.Sprintf(
		// Added single quotes to accomodate empty values.
		// https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING
		"host='%s' port='%d' user='%s' password='%s' dbname='%s' sslmode='%s' sslrootcert='%s'",
		c.PGHost, c.PGPort, c.PGUser, c.PGPassword, c.PGDatabase, c.PGSSLMode, c.PGRootCertLoc,
	))
	if err != nil {
		log.Fatalf("error setting up database connection: %s", err)
	}
	// Verify DB connection. If error, retry with exponential backoff.
	log.Print("verifying database connection...")
	start, sleep, timeout := time.Now(), time.Second, 30*time.Second
	for {
		err := db.Ping()
		if err == nil {
			log.Print("connected to database")
			break
		}
		if time.Since(start) > timeout {
			log.Fatalf("error verifying database connection: %s", err)
			break
		}
		time.Sleep(sleep)
		sleep *= 2
	}

	// RUN DATABASE MIGRATIONS
	// Goose supports out of order migration with "allow missing" option.
	if err := goose.Up(db, "migration", goose.WithAllowMissing()); err != nil {
		log.Fatalf("error running database migrations: %s", err)
	}

	// ADD SERVICE HANDLERS TO HTTP ROUTER
	router := chi.NewRouter()
	router.Route("/api/v1", func(router chi.Router) {
		// Add public routes.
		// TODO: Implement these routes!
		router.Post("/login", func(w http.ResponseWriter, r *http.Request) {})
		router.Post("/register", func(w http.ResponseWriter, r *http.Request) {})

		// Add private routes.
		router.Group(func(router chi.Router) {
			// TODO: Use auth middleware here!
			todoHandler := todo.NewHandler(todo.NewService(todo.NewRepository(db)))
			router.Mount("/todo", todoHandler.HTTPHandler())
		})
	})

	// RUN HTTP SERVER
	svr := http.Server{
		Addr:    fmt.Sprintf("%s:%d", c.APIHost, c.APIPort),
		Handler: router,
	}
	log.Printf("HTTP server listening on %s", svr.Addr)
	if err := svr.ListenAndServe(); err != nil {
		log.Fatalf("error running HTTP server: %s", err)
	}
	// TODO: Implement graceful shutdown!
	log.Fatalf("error shutting down HTTP server: %s", err)
}
