package postgres

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const retryConnectTimeout = 3 * time.Second

func Connect(connStr string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}

	// Try to ping database until timeout.
	start := time.Now()
	sleep := time.Second
	for {
		err := db.Ping()
		if err == nil {
			return db, nil
		}
		if time.Since(start) > retryConnectTimeout {
			return nil, fmt.Errorf("ping: %w", err)
		}
		// Exponential backoff.
		time.Sleep(sleep)
		sleep *= 2
	}
}
