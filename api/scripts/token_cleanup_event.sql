-- Token blacklist cleanup event (run manually or via cron)
-- This creates a MySQL event to automatically cleanup expired tokens

DELIMITER ;;

CREATE EVENT IF NOT EXISTS cleanup_expired_tokens
ON SCHEDULE EVERY 1 HOUR
STARTS CURRENT_TIMESTAMP
DO
BEGIN
    DELETE FROM token_blacklist 
    WHERE expires_at < NOW() - INTERVAL 24 HOUR;
END;;

DELIMITER ;

-- To enable/disable the event:
-- SET GLOBAL event_scheduler = ON;
-- SET GLOBAL event_scheduler = OFF;

-- To manually run cleanup:
-- DELETE FROM token_blacklist WHERE expires_at < NOW() - INTERVAL 24 HOUR;