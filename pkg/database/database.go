package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("error setting up database connection: %w", err)
	}

	// Verify DB connection. If error, retry with exponential backoff.
	start := time.Now()
	sleep := time.Second
	for {
		if err := db.Ping(); err == nil {
			return db, nil
		} else if time.Since(start) > 20*time.Second {
			return nil, fmt.Errorf("error verifying database connection: %w", err)
		}
		time.Sleep(sleep)
		sleep *= 2
	}
}