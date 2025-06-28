-- Token blacklist cleanup event (run manually or via cron)
-- This creates a MySQL event to automatically cleanup expired tokens

DELIMITER ;;

CREATE EVENT IF NOT EXISTS cleanup_expired_tokens
ON SCHEDULE EVERY 1 HOUR
STARTS CURRENT_TIMESTAMP
DO
BEGIN
    DELETE FROM token_blacklist 
    WHERE expires_at < NOW() - INTERVAL 24 HOUR
    LIMIT 1000;
END;;

DELIMITER ;

-- To enable/disable the event:
-- SET GLOBAL event_scheduler = ON;
-- SET GLOBAL event_scheduler = OFF;
--
-- PRODUCTION NOTE: The event scheduler MUST be enabled (ON) in production
-- environments to ensure automatic token cleanup runs correctly. This is
-- critical for system reliability and preventing token blacklist table
-- from growing indefinitely, which could impact performance and storage.

-- To manually run cleanup:
-- DELETE FROM token_blacklist WHERE expires_at < NOW() - INTERVAL 24 HOUR LIMIT 1000;