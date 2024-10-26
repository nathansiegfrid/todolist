package database

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func Migrate(db *sql.DB, path string) error {
	if err := goose.Up(db, path); err != nil {
		return fmt.Errorf("error running migrations: %w", err)
	}
	return nil
}
