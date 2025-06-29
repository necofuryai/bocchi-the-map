#!/bin/bash

# Test database initialization script
# This script runs in the MySQL Docker container during initialization

set -e

echo "Creating test database and user..."

# Check if required environment variables are set
if [ -z "${MYSQL_ROOT_PASSWORD}" ]; then
    echo "ERROR: MYSQL_ROOT_PASSWORD environment variable is not set" >&2
    exit 1
fi

if [ -z "${MYSQL_PASSWORD}" ]; then
    echo "ERROR: MYSQL_PASSWORD environment variable is not set" >&2
    exit 1
fi

# Execute MySQL commands with explicit error handling
mysql -u root -p"${MYSQL_ROOT_PASSWORD}" <<-EOSQL || {
    echo "ERROR: Failed to execute MySQL commands for test database initialization" >&2
    echo "This could be due to:" >&2
    echo "  - Incorrect MYSQL_ROOT_PASSWORD" >&2
    echo "  - MySQL server not ready" >&2
    echo "  - Database connection issues" >&2
    exit 1
}
    -- Create test database with UTF-8 support
    CREATE DATABASE IF NOT EXISTS bocchi_test CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
    
    -- Create user with password from environment variable
    CREATE USER IF NOT EXISTS 'bocchi_user'@'localhost' IDENTIFIED BY '${MYSQL_PASSWORD}';
    
    -- Grant privileges to the test database (restricted to localhost for security)
    GRANT SELECT, INSERT, UPDATE, DELETE ON bocchi_test.* TO 'bocchi_user'@'localhost';
    
    FLUSH PRIVILEGES;
EOSQL

echo "Test database initialization completed."