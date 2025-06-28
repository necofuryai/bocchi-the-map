-- Add composite index for token blacklist performance optimization
-- This index optimizes the IsTokenBlacklisted query that filters by jti and expires_at
CREATE INDEX idx_token_blacklist_jti_expires ON token_blacklist(jti, expires_at);