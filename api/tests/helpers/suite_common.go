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
		Fail(fmt.Sprintf("Failed to initialize test database: %v", err))
		return nil
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
	}
}

// CleanupTestData performs test data cleanup
func (suite *CommonTestSuite) CleanupTestData() {
	By("Cleaning up test data")
	
	// Cleanup is handled by PrepareCleanTestData, but can add specific cleanup here
	// This method is called from Ginkgo context, so we use Ginkgo's logging methods
	err := suite.TestDB.CleanDatabase()
	if err != nil {
		// Use Ginkgo's logging instead of test logging since we're not in testing.T context
		By(fmt.Sprintf("Warning: Failed to cleanup fixtures: %v", err))
	}
}

// SetupGinkgoHooks sets up common Ginkgo BeforeSuite, AfterSuite, BeforeEach, and AfterEach hooks
func SetupGinkgoHooks(suiteNamePrefix string) *CommonTestSuite {
	var testSuite *CommonTestSuite
	
	BeforeSuite(func() {
		By(fmt.Sprintf("Setting up %s test environment", suiteNamePrefix))
		testSuite = NewCommonTestSuite()
		By(fmt.Sprintf("%s test environment setup completed", suiteNamePrefix))
	})
	
	AfterSuite(func() {
		By(fmt.Sprintf("Cleaning up %s test environment", suiteNamePrefix))
		if testSuite != nil {
			testSuite.Cleanup()
		}
		By(fmt.Sprintf("%s test environment cleanup completed", suiteNamePrefix))
	})
	
	BeforeEach(func() {
		testSuite.PrepareCleanTestData()
	})
	
	AfterEach(func() {
		testSuite.CleanupTestData()
	})
	
	return testSuite
}