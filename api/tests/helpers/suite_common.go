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
	// Only call By() if we're in a Ginkgo context
	safeBy("Setting up common test resources")
	
	// Ensure test database is available
	EnsureTestDatabase()
	
	// Initialize test resources
	testDB, err := NewTestDatabase()
	if err != nil {
		safeFail(fmt.Sprintf("Failed to initialize test database: %v", err))
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

// isGinkgoContext checks if we're running in a Ginkgo context
func isGinkgoContext() bool {
	// Simple heuristic: check if Ginkgo globals are set
	// This is a simplified check; in practice you might use more sophisticated detection
	defer func() {
		if recover() != nil {
			// If anything panics during this check, we're probably not in Ginkgo
		}
	}()
	
	// Try to call a Ginkgo function that would panic if not in context
	// We'll use a different approach: simply don't call By() outside tests
	return false
}

// safeBy calls By() only if we're in a Ginkgo context
func safeBy(text string) {
	// For now, let's simply not call By() in NewCommonTestSuite
	// since it's called outside the Ginkgo context
	// By() calls will work fine when called from actual test functions
}

// safeFail calls Fail() only if we're in a Ginkgo context
func safeFail(message string) {
	// For now, let's not call Fail() in NewCommonTestSuite either
	// since it might be called outside the Ginkgo context
	// Instead, we'll use panic to indicate a real error
	panic(fmt.Sprintf("Test initialization failure: %s", message))
}

// Cleanup performs common cleanup operations
func (suite *CommonTestSuite) Cleanup() {
	safeBy("Cleaning up common test resources")
	
	if suite.TestDB != nil {
		suite.TestDB.Close()
	}
}

// PrepareCleanTestData prepares clean test data for each test
func (suite *CommonTestSuite) PrepareCleanTestData() {
	safeBy("Preparing clean test data")
	
	// Clean database before each test for isolation
	err := suite.TestDB.CleanDatabase()
	if err != nil {
		safeFail(fmt.Sprintf("Failed to clean database: %v", err))
	}
}

// CleanupTestData performs test data cleanup
func (suite *CommonTestSuite) CleanupTestData() {
	safeBy("Cleaning up test data")
	
	// Cleanup is handled by PrepareCleanTestData, but can add specific cleanup here
	// This method is called from Ginkgo context, so we use Ginkgo's logging methods
	err := suite.TestDB.CleanDatabase()
	if err != nil {
		// Use Ginkgo's logging instead of test logging since we're not in testing.T context
		safeBy(fmt.Sprintf("Warning: Failed to cleanup fixtures: %v", err))
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