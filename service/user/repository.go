package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nathansiegfrid/todolist-go/service"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) GetAll(ctx context.Context, filter *UserFilter) ([]*User, error) {
	// Translate filter into WHERE conditions and args.
	where, args, argIndex := []string{"TRUE"}, []any{}, 1
	if v := filter.ID; v != nil {
		where = append(where, fmt.Sprintf("id = $%d", argIndex))
		args = append(args, *v)
		argIndex += 1
	}
	if v := filter.Email; v != nil {
		where = append(where, fmt.Sprintf("email = $%d", argIndex))
		args = append(args, *v)
		argIndex += 1
	}

	var limit, offset string
	if filter.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d ", filter.Limit)
	}
	if filter.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d ", filter.Offset)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, email, password_hash, created_at, updated_at
		FROM "user"
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY email ASC`+
		limit+offset,
		args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		u := &User{}
		err := rows.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *Repository) Get(ctx context.Context, id uuid.UUID) (*User, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, email, password_hash, created_at, updated_at
		FROM "user"
		WHERE id = $1`,
		id,
	)

	u := &User{}
	err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, service.ErrNotFound(id)
		}
		return nil, err
	}
	return u, nil
}

func (r *Repository) Create(ctx context.Context, u *User) error {
	u.ID = uuid.New()
	u.CreatedAt = time.Now()
	u.UpdatedAt = u.CreatedAt

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO "user" (id, email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)`,
		u.ID, u.Email, u.PasswordHash, u.CreatedAt, u.UpdatedAt,
	)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			// 23505 is the PostgreSQL error code for unique_violation.
			return service.ErrConflict("Email", u.Email)
		}
		return err
	}
	return nil
}

func (r *Repository) Update(ctx context.Context, id uuid.UUID, update *UserUpdate) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// FOR UPDATE will lock selected row, which prevents new writes and locks to the same row
	// before current Tx is done.
	row := tx.QueryRowContext(ctx, `
		SELECT id, email, password_hash, created_at, updated_at
		FROM "user"
		WHERE id = $1
		FOR UPDATE`,
		id,
	)

	u := &User{}
	err = row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return service.ErrNotFound(id)
		}
		return err
	}

	// Check if resource is owned by user.
	if u.ID != service.UserIDFromContext(ctx) {
		return service.ErrPermission()
	}

	// Apply patch update.
	var updated bool
	if v := update.Email; v != nil {
		u.Email = *v
		updated = true
	}
	if v := update.Password; v != nil {
		u.SetNewPassword(*v)
		updated = true
	}
	if updated {
		u.UpdatedAt = time.Now()
	}

	result, err := tx.ExecContext(ctx, `
		UPDATE "user"
		SET email = $1, password_hash = $2, updated_at = $3
		WHERE id = $4`,
		u.Email, u.PasswordHash, u.UpdatedAt, id,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return service.ErrNotFound(id)
	}

	return tx.Commit()
}
