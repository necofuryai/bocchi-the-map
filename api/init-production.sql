-- Production database initialization
-- Initialize database with UTF-8 support
CREATE DATABASE IF NOT EXISTS bocchi_the_map CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Grant limited privileges to the production database user
-- Restrict to localhost and specific application operations only
GRANT SELECT, INSERT, UPDATE, DELETE ON bocchi_the_map.* TO 'bocchi_user'@'localhost';

FLUSH PRIVILEGES;