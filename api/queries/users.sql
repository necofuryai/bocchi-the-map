-- name: CreateUser :exec
INSERT INTO users (
    id, email, display_name, avatar_url, auth_provider, auth_provider_id, preferences
) VALUES (
    ?, ?, ?, ?, ?, ?, ?
);

-- name: GetUserByProviderID :one
SELECT * FROM users 
WHERE auth_provider = ? AND auth_provider_id = ?;

-- name: GetUserByEmail :one
SELECT * FROM users 
WHERE email = ?;

-- name: GetUserByID :one
SELECT * FROM users 
WHERE id = ?;

-- name: UpdateUserPreferences :exec
UPDATE users 
SET preferences = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: UpdateUserAvatar :exec
UPDATE users 
SET avatar_url = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: UpsertUser :exec
INSERT INTO users (
    id, email, display_name, avatar_url, auth_provider, auth_provider_id, preferences
) VALUES (
    ?, ?, ?, ?, ?, ?, ?
) ON DUPLICATE KEY UPDATE
    display_name = VALUES(display_name),
    avatar_url = VALUES(avatar_url),
    updated_at = CURRENT_TIMESTAMP;