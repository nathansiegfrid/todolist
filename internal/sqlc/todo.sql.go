// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: todo.sql

package sqlc

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createTodo = `-- name: CreateTodo :one
INSERT INTO todo (user_id, subject, description, due_date)
VALUES ($1, $2, $3, $4) RETURNING id
`

type CreateTodoParams struct {
	UserID      uuid.NullUUID
	Subject     string
	Description string
	DueDate     sql.NullTime
}

func (q *Queries) CreateTodo(ctx context.Context, arg CreateTodoParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, createTodo,
		arg.UserID,
		arg.Subject,
		arg.Description,
		arg.DueDate,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const deleteTodo = `-- name: DeleteTodo :execrows
DELETE FROM todo WHERE id = $1
`

func (q *Queries) DeleteTodo(ctx context.Context, id uuid.UUID) (int64, error) {
	result, err := q.db.ExecContext(ctx, deleteTodo, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const getAllTodos = `-- name: GetAllTodos :many
SELECT id, user_id, subject, description, priority, due_date, completed, created_at, updated_at FROM todo
`

func (q *Queries) GetAllTodos(ctx context.Context) ([]Todo, error) {
	rows, err := q.db.QueryContext(ctx, getAllTodos)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Todo{}
	for rows.Next() {
		var i Todo
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Subject,
			&i.Description,
			&i.Priority,
			&i.DueDate,
			&i.Completed,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTodoByID = `-- name: GetTodoByID :one
SELECT id, user_id, subject, description, priority, due_date, completed, created_at, updated_at FROM todo WHERE id = $1
`

func (q *Queries) GetTodoByID(ctx context.Context, id uuid.UUID) (Todo, error) {
	row := q.db.QueryRowContext(ctx, getTodoByID, id)
	var i Todo
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Subject,
		&i.Description,
		&i.Priority,
		&i.DueDate,
		&i.Completed,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getTodoByUserID = `-- name: GetTodoByUserID :many
SELECT id, user_id, subject, description, priority, due_date, completed, created_at, updated_at FROM todo WHERE user_id = $1
`

func (q *Queries) GetTodoByUserID(ctx context.Context, userID uuid.NullUUID) ([]Todo, error) {
	rows, err := q.db.QueryContext(ctx, getTodoByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Todo{}
	for rows.Next() {
		var i Todo
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Subject,
			&i.Description,
			&i.Priority,
			&i.DueDate,
			&i.Completed,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateTodo = `-- name: UpdateTodo :execrows
UPDATE todo
SET user_id = $2,
    subject = $3,
    description = $4,
    priority = $5,
    due_date = $6,
    completed = $7,
    updated_at = NOW()
WHERE id = $1
`

type UpdateTodoParams struct {
	ID          uuid.UUID
	UserID      uuid.NullUUID
	Subject     string
	Description string
	Priority    int32
	DueDate     sql.NullTime
	Completed   bool
}

func (q *Queries) UpdateTodo(ctx context.Context, arg UpdateTodoParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, updateTodo,
		arg.ID,
		arg.UserID,
		arg.Subject,
		arg.Description,
		arg.Priority,
		arg.DueDate,
		arg.Completed,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
