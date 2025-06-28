-- Initialize database with UTF-8 support
CREATE DATABASE IF NOT EXISTS bocchi_the_map CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE DATABASE IF NOT EXISTS bocchi_test CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Grant privileges to the user for both databases
GRANT ALL PRIVILEGES ON bocchi_the_map.* TO 'bocchi_user'@'%';
GRANT ALL PRIVILEGES ON bocchi_test.* TO 'bocchi_user'@'%';
GRANT ALL PRIVILEGES ON bocchi_test.* TO 'root'@'%';
FLUSH PRIVILEGES;