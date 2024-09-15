package todo

import (
	"context"
	"database/sql"
	"errors"

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

func (r *Repository) GetAll(ctx context.Context) ([]*Todo, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, description, completed
		FROM todo`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*Todo
	for rows.Next() {
		todo := &Todo{}
		err := rows.Scan(&todo.ID, &todo.UserID, &todo.Description, &todo.Completed)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func (r *Repository) Get(ctx context.Context, id uuid.UUID) (*Todo, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, description, completed
		FROM todo
		WHERE id = $1`,
		id,
	)

	todo := &Todo{}
	err := row.Scan(&todo.ID, &todo.UserID, &todo.Description, &todo.Completed)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, service.ErrNotFound(id)
		}
		return nil, err
	}
	return todo, nil
}

func (r *Repository) Create(ctx context.Context, t *Todo) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO todo (user_id, description, completed)
		VALUES ($1, $2, $3)`,
		t.UserID, t.Description, t.Completed,
	)
	if err != nil {
		// This check is unnecessary, it's here as template.
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			// 23505 is the PostgreSQL error code for unique_violation.
			return service.ErrConflict("ID", t.ID.String())
		}
		return err
	}
	return nil
}

func (r *Repository) Update(ctx context.Context, id uuid.UUID, t *Todo) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE todo
		SET description = $1, completed = $2
		WHERE id = $3`,
		t.Description, t.Completed, id,
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
	return nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM todo WHERE id = $1", id)
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
	return nil
}
