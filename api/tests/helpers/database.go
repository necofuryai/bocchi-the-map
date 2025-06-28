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
func (td *TestDatabase) CleanDatabase() {
	ctx := context.Background()
	
	// Clean up in reverse order of dependencies
	tables := []string{
		"reviews",
		"spots", 
		"users",
		"token_blacklist",
	}
	
	for _, table := range tables {
		_, err := td.DB.ExecContext(ctx, fmt.Sprintf("DELETE FROM `%s`", table))
		if err != nil {
			// Log warning but don't fail - table might not exist in test
			GinkgoWriter.Printf("Warning: Failed to clean table %s: %v\n", table, err)
		}
	}
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
			tx.Rollback()
			panic(r)
		}
	}()
	
	fn(tx)
	
	err = tx.Rollback() // Always rollback for test isolation
	Expect(err).NotTo(HaveOccurred(), "Failed to rollback transaction")
}

// getTestDSN returns the test database connection string
func getTestDSN() string {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		// Default test database configuration
		dsn = "root:password@tcp(localhost:3306)/bocchi_test?parseTime=true&multiStatements=true"
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