package todo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/nathansiegfrid/todolist-go/service"
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
		SELECT id, user_id, description, completed
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

func (r *Repository) Create(ctx context.Context, todo *Todo) error {
	todo.ID = uuid.New()
	todo.UserID = uuid.NullUUID{UUID: service.UserIDFromContext(ctx), Valid: true}

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO todo (id, user_id, description, completed)
		VALUES ($1, $2, $3, $4)`,
		todo.ID, todo.UserID, todo.Description, todo.Completed,
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
		SELECT id, user_id, description, completed
		FROM todo
		WHERE id = $1
		FOR UPDATE`,
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

func updateTodo(ctx context.Context, tx *sql.Tx, id uuid.UUID, update *TodoUpdate) error {
	todo, err := getTodoForUpdate(ctx, tx, id)
	if err != nil {
		return err
	}

	// Check if resource is owned by user.
	if todo.UserID.UUID != service.UserIDFromContext(ctx) {
		return service.ErrPermission()
	}

	if v := update.Description; v != nil {
		todo.Description = *v
	}
	if v := update.Completed; v != nil {
		todo.Completed = *v
	}

	result, err := tx.ExecContext(ctx, `
		UPDATE todo
		SET description = $1, completed = $2
		WHERE id = $3`,
		todo.Description, todo.Completed, id,
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

func deleteTodo(ctx context.Context, tx *sql.Tx, id uuid.UUID) error {
	todo, err := getTodoForUpdate(ctx, tx, id)
	if err != nil {
		return err
	}

	// Check if resource is owned by user.
	if todo.UserID.UUID != service.UserIDFromContext(ctx) {
		return service.ErrPermission()
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
		return service.ErrNotFound(id)
	}
	return nil
}
