-- Initialize database with UTF-8 support
CREATE DATABASE IF NOT EXISTS bocchi_the_map CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE bocchi_the_map;

-- Grant privileges to the user
GRANT ALL PRIVILEGES ON bocchi_the_map.* TO 'bocchi_user'@'%';
FLUSH PRIVILEGES;