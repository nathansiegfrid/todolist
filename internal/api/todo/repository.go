package todo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nathansiegfrid/todolist/internal/api"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) GetAll(ctx context.Context, filter *TodoFilter) ([]*Todo, error) {
	// Translate filter into WHERE conditions and args.
	where, args, argIndex := []string{"TRUE"}, []any{}, 1
	if v := filter.ID; v != nil {
		where = append(where, fmt.Sprintf("id = $%d", argIndex))
		args = append(args, *v)
		argIndex++
	}
	if v := filter.UserID; v != nil {
		where = append(where, fmt.Sprintf("user_id = $%d", argIndex))
		args = append(args, *v)
		argIndex++
	}
	if v := filter.Priority; v != nil {
		where = append(where, fmt.Sprintf("priority = $%d", argIndex))
		args = append(args, *v)
		argIndex++
	}
	if v := filter.DueDate; v != nil {
		if !v.Valid {
			where = append(where, "due_date IS NULL")
		} else {
			where = append(where, fmt.Sprintf("due_date::date = $%d::date", argIndex))
			args = append(args, *v)
			argIndex++
		}
	}
	if v := filter.Completed; v != nil {
		where = append(where, fmt.Sprintf("completed = $%d", argIndex))
		args = append(args, *v)
		argIndex++
	}

	var limit, offset string
	if filter.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d ", filter.Limit)
	}
	if filter.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d ", filter.Offset)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, subject, description, priority, due_date, completed, created_at, updated_at
		FROM todo
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY description ASC`+
		limit+offset,
		args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*Todo
	for rows.Next() {
		todo := &Todo{}
		err := rows.Scan(
			&todo.ID,
			&todo.UserID,
			&todo.Subject,
			&todo.Description,
			&todo.Priority,
			&todo.DueDate,
			&todo.Completed,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func (r *Repository) Get(ctx context.Context, id uuid.UUID) (*Todo, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, subject, description, priority, due_date, completed, created_at, updated_at
		FROM todo
		WHERE id = $1`,
		id,
	)

	todo := &Todo{}
	err := row.Scan(
		&todo.ID,
		&todo.UserID,
		&todo.Subject,
		&todo.Description,
		&todo.Priority,
		&todo.DueDate,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, api.ErrIDNotFound(id)
		}
		return nil, err
	}
	return todo, nil
}

func (r *Repository) Create(ctx context.Context, todo *Todo) error {
	todo.ID = uuid.New()
	todo.UserID = api.UserIDFromContext(ctx)
	todo.CreatedAt = time.Now()
	todo.UpdatedAt = todo.CreatedAt

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO todo (id, user_id, subject, description, priority, due_date, completed, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		todo.ID,
		todo.UserID,
		todo.Subject,
		todo.Description,
		todo.Priority,
		todo.DueDate,
		todo.Completed,
		todo.CreatedAt,
		todo.UpdatedAt,
	)
	return err
}

func (r *Repository) Update(ctx context.Context, id uuid.UUID, update *TodoUpdate) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = updateTodo(ctx, tx, id, update)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = deleteTodo(ctx, tx, id)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func getTodoForUpdate(ctx context.Context, tx *sql.Tx, id uuid.UUID) (*Todo, error) {
	// FOR UPDATE will lock selected row, which prevents new writes and locks to the same row
	// before current Tx is done.
	row := tx.QueryRowContext(ctx, `
		SELECT id, user_id, subject, description, priority, due_date, completed, created_at, updated_at
		FROM todo
		WHERE id = $1
		FOR UPDATE`,
		id,
	)

	todo := &Todo{}
	err := row.Scan(
		&todo.ID,
		&todo.UserID,
		&todo.Subject,
		&todo.Description,
		&todo.Priority,
		&todo.DueDate,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, api.ErrIDNotFound(id)
		}
		return nil, err
	}
	return todo, nil
}

func updateTodo(ctx context.Context, tx *sql.Tx, id uuid.UUID, update *TodoUpdate) error {
	todo, err := getTodoForUpdate(ctx, tx, id)
	if err != nil {
		return err
	}

	// Check if resource is owned by user.
	if todo.UserID != api.UserIDFromContext(ctx) {
		return api.ErrPermission()
	}

	todo.Subject = update.Subject.ValueOr(todo.Subject)
	todo.Description = update.Description.ValueOr(todo.Description)
	todo.Priority = update.Priority.ValueOr(todo.Priority)
	todo.DueDate = update.DueDate.ValueOr(todo.DueDate)
	todo.Completed = update.Completed.ValueOr(todo.Completed)
	todo.UpdatedAt = time.Now()

	result, err := tx.ExecContext(ctx, `
		UPDATE todo
		SET subject = $2, description = $3, priority = $4, due_date = $5, completed = $6, updated_at = $7
		WHERE id = $1`,
		id,
		todo.Subject,
		todo.Description,
		todo.Priority,
		todo.DueDate,
		todo.Completed,
		todo.UpdatedAt,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return api.ErrIDNotFound(id)
	}
	return nil
}

func deleteTodo(ctx context.Context, tx *sql.Tx, id uuid.UUID) error {
	todo, err := getTodoForUpdate(ctx, tx, id)
	if err != nil {
		return err
	}

	// Check if resource is owned by user.
	if todo.UserID != api.UserIDFromContext(ctx) {
		return api.ErrPermission()
	}

	result, err := tx.ExecContext(ctx, "DELETE FROM todo WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return api.ErrIDNotFound(id)
	}
	return nil
}
