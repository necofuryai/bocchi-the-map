-- Add composite index for token blacklist performance optimization
-- This index optimizes the IsTokenBlacklisted query that filters by jti and expires_at
ALTER TABLE token_blacklist
  ADD INDEX idx_token_blacklist_jti_expires (jti, expires_at)
  ALGORITHM=INPLACE, LOCK=NONE;