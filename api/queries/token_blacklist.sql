-- name: AddToBlacklist :exec
INSERT INTO token_blacklist (jti, user_id, token_type, expires_at, reason)
VALUES (?, ?, ?, ?, ?);

-- name: IsTokenBlacklisted :one
SELECT COUNT(*) FROM token_blacklist 
WHERE jti = ? AND expires_at > NOW();

-- name: BlacklistAccessToken :exec
INSERT INTO token_blacklist (jti, user_id, token_type, expires_at, reason)
VALUES (?, ?, 'access', ?, 'logout');

-- name: BlacklistRefreshToken :exec  
INSERT INTO token_blacklist (jti, user_id, token_type, expires_at, reason)
VALUES (?, ?, 'refresh', ?, 'logout');

-- name: CleanupExpiredTokens :exec
DELETE FROM token_blacklist 
WHERE expires_at < NOW() - INTERVAL 24 HOUR;