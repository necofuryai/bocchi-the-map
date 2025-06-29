-- Initialize database with UTF-8 support
CREATE DATABASE IF NOT EXISTS bocchi_the_map CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Grant privileges to the production database
GRANT ALL PRIVILEGES ON bocchi_the_map.* TO 'bocchi_user'@'%';

-- Test database setup (only for CI/test environments)
-- Set environment variable CI=true or TEST_ENV=true to enable test database
SET @is_test_env = IF(
    (@ci := COALESCE(
        NULLIF(TRIM(@ci), ''),
        NULLIF(TRIM(COALESCE(@@global.init_connect, '')), ''),
        'false'
    )) = 'true' OR 
    (@test_env := COALESCE(
        NULLIF(TRIM(@test_env), ''),
        'false'
    )) = 'true',
    1, 
    0
);

-- Create test database and grant privileges only in test environments
SET @sql = IF(@is_test_env = 1,
    'CREATE DATABASE IF NOT EXISTS bocchi_test CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci; GRANT ALL PRIVILEGES ON bocchi_test.* TO ''bocchi_user''@''%'';',
    'SELECT "Skipping test database setup in production environment" AS message;'
);

PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

FLUSH PRIVILEGES;