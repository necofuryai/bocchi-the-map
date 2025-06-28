// +build e2e

package e2e

import (
	"testing"

	"github.com/necofuryai/bocchi-the-map/api/tests/helpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Global test resources shared across all BDD tests
var (
	testDB      *helpers.TestDatabase
	fixtureManager *helpers.FixtureManager
	authHelper  *helpers.AuthHelper
)

func TestE2EScenarios(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bocchi The Map API E2E Test Suite")
}

var _ = BeforeSuite(func() {
	By("Setting up test environment")
	
	// Ensure test database is available
	helpers.EnsureTestDatabase()
	
	// Initialize test resources
	testDB = helpers.NewTestDatabase()
	fixtureManager = helpers.NewFixtureManager(testDB)
	authHelper = helpers.NewAuthHelper()
	
	By("Test environment setup completed")
})

var _ = AfterSuite(func() {
	By("Cleaning up test environment")
	
	if testDB != nil {
		testDB.Close()
	}
	
	By("Test environment cleanup completed")
})

// BeforeEach hook for test isolation
var _ = BeforeEach(func() {
	By("Preparing clean test data")
	
	// Clean database before each test for isolation
	testDB.CleanDatabase()
})

// AfterEach hook for test cleanup
var _ = AfterEach(func() {
	By("Cleaning up test data")
	
	// Cleanup is handled by BeforeEach, but can add specific cleanup here
	fixtureManager.CleanupFixtures()
})