package helpers

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/necofuryai/bocchi-the-map/api/infrastructure/database"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// TestDatabase manages database connections and cleanup for BDD tests
type TestDatabase struct {
	DB      *sql.DB
	Queries *database.Queries
}

// NewTestDatabase creates a new test database connection with proper isolation
func NewTestDatabase() (*TestDatabase, error) {
	dsn := getTestDSN()
	
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to test database: %w", err)
	}
	
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping test database: %w", err)
	}
	
	// Set connection pool settings for tests
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(5 * time.Minute)
	
	return &TestDatabase{
		DB:      db,
		Queries: database.New(db),
	}, nil
}

// CleanDatabase removes all test data while preserving schema
func (td *TestDatabase) CleanDatabase() error {
	ctx := context.Background()
	
	// Define allowed tables for cleanup to prevent SQL injection
	allowedTables := map[string]bool{
		"reviews":         true,
		"spots":           true,
		"users":           true,
		"token_blacklist": true,
	}
	
	// Clean up in reverse order of dependencies
	tables := []string{
		"reviews",
		"spots", 
		"users",
		"token_blacklist",
	}
	
	var errors []error
	
	for _, table := range tables {
		// Validate table name against whitelist
		if !allowedTables[table] {
			return fmt.Errorf("table '%s' is not allowed for cleanup operations", table)
		}
		
		// Use backticks for table names - safe since validated against whitelist
		_, err := td.DB.ExecContext(ctx, fmt.Sprintf("DELETE FROM `%s`", table))
		if err != nil {
			GinkgoWriter.Printf("Warning: Failed to clean table %s: %v\n", table, err)
			errors = append(errors, fmt.Errorf("failed to clean table %s: %w", table, err))
		}
	}
	
	if len(errors) > 0 {
		return fmt.Errorf("database cleanup failed with %d errors: %v", len(errors), errors)
	}
	
	return nil
}

// Close closes the database connection
func (td *TestDatabase) Close() {
	if td.DB != nil {
		td.DB.Close()
	}
}

// WithTransaction executes a function within a database transaction for isolation
func (td *TestDatabase) WithTransaction(fn func(*sql.Tx)) {
	tx, err := td.DB.Begin()
	Expect(err).NotTo(HaveOccurred(), "Failed to begin transaction")
	
	defer func() {
		if r := recover(); r != nil {
			if err := tx.Rollback(); err != nil {
				GinkgoWriter.Printf("Failed to rollback transaction after panic: %v\n", err)
			}
			panic(r)
		}
	}()
	
	fn(tx)
	
	err = tx.Rollback() // Always rollback for test isolation
	Expect(err).NotTo(HaveOccurred(), "Failed to rollback transaction")
}

// getTestDSN returns the test database connection string
// Requires TEST_DATABASE_URL environment variable to be set for security.
// Example: TEST_DATABASE_URL="user:password@tcp(localhost:3306)/bocchi_test?parseTime=true&multiStatements=true"
func getTestDSN() string {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		Skip("TEST_DATABASE_URL environment variable must be set for test database configuration. " +
			"Example: TEST_DATABASE_URL=\"user:password@tcp(localhost:3306)/bocchi_test?parseTime=true&multiStatements=true\"")
	}
	return dsn
}

// EnsureTestDatabase ensures test database exists and is ready
func EnsureTestDatabase() {
	dsn := getTestDSN()
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		Skip(fmt.Sprintf("Test database not available: %v", err))
	}
	defer db.Close()
	
	if err := db.Ping(); err != nil {
		Skip(fmt.Sprintf("Cannot connect to test database: %v", err))
	}
}