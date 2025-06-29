-- Drop composite index for token blacklist
DROP INDEX IF EXISTS idx_token_blacklist_jti_expires ON token_blacklist;