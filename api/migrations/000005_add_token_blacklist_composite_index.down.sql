-- Drop composite index for token blacklist
ALTER TABLE token_blacklist DROP INDEX IF EXISTS idx_token_blacklist_jti_expires;