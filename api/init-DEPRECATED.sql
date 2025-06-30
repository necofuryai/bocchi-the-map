-- DEPRECATED: This file contains unreliable environment variable logic.
-- Use init-production.sql for production environments
-- Use init-test.sql for test/CI environments

-- Initialize database with UTF-8 support
CREATE DATABASE IF NOT EXISTS bocchi_the_map CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Grant privileges to the production database (restricted to localhost for security)
GRANT ALL PRIVILEGES ON bocchi_the_map.* TO 'bocchi_user'@'localhost';

FLUSH PRIVILEGES;