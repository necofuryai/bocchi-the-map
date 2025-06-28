-- Drop token blacklist table and cleanup event
DROP EVENT IF EXISTS cleanup_expired_tokens;
DROP TABLE IF EXISTS token_blacklist;