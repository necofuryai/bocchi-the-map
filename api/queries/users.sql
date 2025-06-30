-- User management queries for Bocchi The Map API
-- These queries support Auth0 integration and user profile management

-- name: GetUserByID :one
SELECT * FROM users WHERE id = ? LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ? LIMIT 1;

-- name: GetUserByProviderID :one
SELECT * FROM users WHERE auth_provider = ? AND auth_provider_id = ? LIMIT 1;

-- name: CreateUser :exec
INSERT INTO users (
    id, email, display_name, avatar_url, auth_provider, auth_provider_id, 
    preferences, created_at, updated_at
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, NOW(), NOW()
);

-- name: UpsertUser :exec
INSERT INTO users (
    id, email, display_name, avatar_url, auth_provider, auth_provider_id, 
    preferences, created_at, updated_at
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, NOW(), NOW()
)
ON DUPLICATE KEY UPDATE
    email = VALUES(email),
    display_name = VALUES(display_name),
    avatar_url = VALUES(avatar_url),
    preferences = VALUES(preferences),
    updated_at = NOW();

-- name: UpdateUserAvatar :exec
UPDATE users SET avatar_url = ?, updated_at = NOW() WHERE id = ?;

-- name: UpdateUserPreferences :exec
UPDATE users SET preferences = ?, updated_at = NOW() WHERE id = ?;

-- Token blacklist queries for logout and security
-- name: AddToBlacklist :exec
INSERT INTO token_blacklist (jti, user_id, token_type, expires_at, revoked_at, reason) 
VALUES (?, ?, ?, ?, NOW(), ?);

-- name: BlacklistAccessToken :exec
INSERT INTO token_blacklist (jti, user_id, token_type, expires_at, revoked_at, reason) 
VALUES (?, ?, 'access', ?, NOW(), 'logout');

-- name: BlacklistRefreshToken :exec
INSERT INTO token_blacklist (jti, user_id, token_type, expires_at, revoked_at, reason) 
VALUES (?, ?, 'refresh', ?, NOW(), 'logout');

-- name: IsTokenBlacklisted :one
SELECT EXISTS(
    SELECT 1 FROM token_blacklist 
    WHERE jti = ? AND expires_at > NOW()
) as is_blacklisted;

-- name: CleanupExpiredTokens :exec
DELETE FROM token_blacklist WHERE expires_at <= NOW();