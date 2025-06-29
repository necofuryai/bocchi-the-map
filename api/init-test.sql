-- Test database initialization
-- Create test database with UTF-8 support
CREATE DATABASE IF NOT EXISTS bocchi_test CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Grant privileges to the test database (restricted to localhost for security)
CREATE USER IF NOT EXISTS 'bocchi_user'@'localhost' IDENTIFIED BY '${MYSQL_USER_PASSWORD}';
GRANT SELECT, INSERT, UPDATE, DELETE ON bocchi_test.* TO 'bocchi_user'@'localhost';

FLUSH PRIVILEGES;