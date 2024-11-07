-- name: GetAllTodos :many
SELECT * FROM todo;

-- name: GetTodoByID :one
SELECT * FROM todo WHERE id = $1;

-- name: GetTodoByUserID :many
SELECT * FROM todo WHERE user_id = $1;

-- name: CreateTodo :one
INSERT INTO todo (user_id, subject, description, due_date)
VALUES ($1, $2, $3, $4) RETURNING id;

-- name: UpdateTodo :execrows
UPDATE todo
SET user_id = $2,
    subject = $3,
    description = $4,
    priority = $5,
    due_date = $6,
    completed = $7,
    updated_at = NOW()
WHERE id = $1;

-- name: DeleteTodo :execrows
DELETE FROM todo WHERE id = $1;
