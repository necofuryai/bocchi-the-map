-- Production database initialization
-- Initialize database with UTF-8 support
CREATE DATABASE IF NOT EXISTS bocchi_the_map CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Create database user with secure password from environment variable
-- Restrict to localhost and specific application operations only
-- Note: This SQL should be executed with password injection at runtime
-- Example: mysql < init-production.sql -e "SET @DB_PASSWORD = '$DB_PASSWORD';"
SET @password = IFNULL(@DB_PASSWORD, '');
SET @sql = CONCAT('CREATE USER IF NOT EXISTS \'bocchi_user\'@\'localhost\' IDENTIFIED BY \'', @password, '\'');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Grant limited privileges to the production database user
GRANT SELECT, INSERT, UPDATE, DELETE ON bocchi_the_map.* TO 'bocchi_user'@'localhost';

FLUSH PRIVILEGES;