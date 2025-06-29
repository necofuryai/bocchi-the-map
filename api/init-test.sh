#!/bin/bash

# Test database initialization script
# This script runs in the MySQL Docker container during initialization

set -e

echo "Creating test database and user..."

mysql -u root -p"${MYSQL_ROOT_PASSWORD}" <<-EOSQL
    -- Create test database with UTF-8 support
    CREATE DATABASE IF NOT EXISTS bocchi_test CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
    
    -- Create user with password from environment variable
    CREATE USER IF NOT EXISTS 'bocchi_user'@'localhost' IDENTIFIED BY '${MYSQL_PASSWORD}';
    
    -- Grant privileges to the test database (restricted to localhost for security)
    GRANT SELECT, INSERT, UPDATE, DELETE ON bocchi_test.* TO 'bocchi_user'@'localhost';
    
    FLUSH PRIVILEGES;
EOSQL

echo "Test database initialization completed."