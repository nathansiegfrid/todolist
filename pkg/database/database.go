package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const retryConnectTimeout = 10 * time.Second

func ConnectPostgres(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("error setting up database connection: %w", err)
	}

	// Verify DB connection. If error, retry with exponential backoff.
	start := time.Now()
	sleep := time.Second
	for {
		err := db.Ping()
		if err == nil {
			return db, nil
		}
		if time.Since(start) > retryConnectTimeout {
			return nil, fmt.Errorf("error verifying database connection: %w", err)
		}
		time.Sleep(sleep)
		sleep *= 2
	}
}
