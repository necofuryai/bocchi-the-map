-- Add composite index for token blacklist performance optimization
-- This index optimizes the IsTokenBlacklisted query that filters by jti and expires_at
-- Note: Production version with ALGORITHM=INPLACE, LOCK=NONE is in migrations/production/
ALTER TABLE token_blacklist
  ADD INDEX idx_token_blacklist_jti_expires (jti, expires_at);