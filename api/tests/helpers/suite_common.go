package helpers

import (
	"fmt"
	
	. "github.com/onsi/ginkgo/v2"
)

// Common test suite resources and setup logic
type CommonTestSuite struct {
	TestDB         *TestDatabase
	FixtureManager *FixtureManager
	AuthHelper     *AuthHelper
}

// NewCommonTestSuite creates and initializes common test resources
func NewCommonTestSuite() *CommonTestSuite {
	By("Setting up common test resources")
	
	// Ensure test database is available
	EnsureTestDatabase()
	
	// Initialize test resources
	testDB, err := NewTestDatabase()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize test database during NewTestDatabase() operation: %v. "+
			"Common causes: "+
			"1) Database connection failure - check TEST_DATABASE_URL environment variable, "+
			"2) Database ping failure - ensure database server is running and accessible, "+
			"3) Missing test database configuration - verify database exists and credentials are correct.", err))
	}
	fixtureManager := NewFixtureManager(testDB)
	authHelper := NewAuthHelper()
	
	return &CommonTestSuite{
		TestDB:         testDB,
		FixtureManager: fixtureManager,
		AuthHelper:     authHelper,
	}
}

// Cleanup performs common cleanup operations
func (suite *CommonTestSuite) Cleanup() {
	By("Cleaning up common test resources")
	
	if suite.TestDB != nil {
		suite.TestDB.Close()
	}
}

// PrepareCleanTestData prepares clean test data for each test
func (suite *CommonTestSuite) PrepareCleanTestData() {
	By("Preparing clean test data")
	
	// Clean database before each test for isolation
	err := suite.TestDB.CleanDatabase()
	if err != nil {
		Fail(fmt.Sprintf("Failed to clean database: %v", err))
		return
	}
}

// CleanupTestData performs test data cleanup
func (suite *CommonTestSuite) CleanupTestData() {
	By("Cleaning up test data")
	
	// Cleanup is handled by PrepareCleanTestData, but can add specific cleanup here
	suite.FixtureManager.CleanupFixtures()
}