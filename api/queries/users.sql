-- User management queries for Bocchi The Map API
-- These queries support Auth0 integration and user profile management

-- name: GetUserByID :one
SELECT * FROM users WHERE id = ? LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ? LIMIT 1;

-- name: GetUserByProviderID :one
SELECT * FROM users WHERE provider = ? AND provider_id = ? LIMIT 1;

-- name: CreateUser :exec
INSERT INTO users (
    id, email, name, nickname, picture, provider, provider_id, 
    email_verified, preferences, created_at, updated_at
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW()
);

-- name: UpsertUser :exec
INSERT INTO users (
    id, email, name, nickname, picture, provider, provider_id, 
    email_verified, preferences, created_at, updated_at
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW()
)
ON DUPLICATE KEY UPDATE
    email = VALUES(email),
    name = VALUES(name),
    nickname = VALUES(nickname),
    picture = VALUES(picture),
    email_verified = VALUES(email_verified),
    preferences = VALUES(preferences),
    updated_at = NOW();

-- name: UpdateUserAvatar :exec
UPDATE users SET picture = ?, updated_at = NOW() WHERE id = ?;

-- name: UpdateUserPreferences :exec
UPDATE users SET preferences = ?, updated_at = NOW() WHERE id = ?;

-- Token blacklist queries for logout and security
-- name: AddToBlacklist :exec
INSERT INTO token_blacklist (jti, token_type, expires_at) 
VALUES (?, ?, ?);

-- name: BlacklistAccessToken :exec
INSERT INTO token_blacklist (jti, token_type, expires_at) 
VALUES (?, 'access', ?);

-- name: BlacklistRefreshToken :exec
INSERT INTO token_blacklist (jti, token_type, expires_at) 
VALUES (?, 'refresh', ?);

-- name: IsTokenBlacklisted :one
SELECT EXISTS(
    SELECT 1 FROM token_blacklist 
    WHERE jti = ? AND expires_at > NOW()
) as is_blacklisted;

-- name: CleanupExpiredTokens :exec
DELETE FROM token_blacklist WHERE expires_at <= NOW();

-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?;