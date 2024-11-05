package postgres

import (
	"context"
	"database/sql"
	"os"

	"github.com/pressly/goose/v3"
)

func Migrate(db *sql.DB, dir string) ([]*goose.MigrationResult, error) {
	p, err := goose.NewProvider("postgres", db, os.DirFS(dir))
	if err != nil {
		return nil, err
	}
	return p.Up(context.Background())
}
