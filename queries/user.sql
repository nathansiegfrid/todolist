-- name: GetAllUsers :many
SELECT * FROM "user";

-- name: GetUserByID :one
SELECT * FROM "user" WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM "user" WHERE email = LOWER($1);

-- name: CreateUser :one
INSERT INTO "user" (email, password_hash)
VALUES (LOWER($1), $2) RETURNING id;

-- name: UpdateUser :execrows
UPDATE "user"
SET email = LOWER($2),
    password_hash = $3,
    updated_at = NOW()
WHERE id = $1;

-- name: DeleteUser :execrows
DELETE FROM "user" WHERE id = $1;
