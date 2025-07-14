-- name: CreateUser :exec
INSERT INTO users (
    id, email, name
) VALUES (
    ?, ?, ?
);

-- name: GetUserByID :one
SELECT id, email, name, created_at, updated_at FROM users 
WHERE id = ?;

-- name: GetUserByEmail :one
SELECT id, email, name, created_at, updated_at FROM users 
WHERE email = ?;

-- name: UpdateUser :exec
UPDATE users 
SET name = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: UpdateUserEmail :exec
UPDATE users 
SET email = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: DeleteUser :exec
DELETE FROM users 
WHERE id = ?;

-- name: ListUsers :many
SELECT id, email, name, created_at, updated_at FROM users 
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;