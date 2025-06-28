// +build integration

package handlers

import (
	"testing"

	"github.com/necofuryai/bocchi-the-map/api/tests/helpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Global test resources shared across handler integration tests
var (
	testDB         *helpers.TestDatabase
	fixtureManager *helpers.FixtureManager
	authHelper     *helpers.AuthHelper
)

func TestHandlers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "HTTP Handlers Integration Test Suite")
}

var _ = BeforeSuite(func() {
	By("Setting up handler integration test environment")
	
	// Ensure test database is available
	helpers.EnsureTestDatabase()
	
	// Initialize test resources
	testDB = helpers.NewTestDatabase()
	fixtureManager = helpers.NewFixtureManager(testDB)
	authHelper = helpers.NewAuthHelper()
	
	By("Handler integration test environment setup completed")
})

var _ = AfterSuite(func() {
	By("Cleaning up handler integration test environment")
	
	if testDB != nil {
		testDB.Close()
	}
	
	By("Handler integration test environment cleanup completed")
})

// BeforeEach hook for test isolation
var _ = BeforeEach(func() {
	By("Preparing clean test data for handler tests")
	
	// Clean database before each test for isolation
	testDB.CleanDatabase()
})

// AfterEach hook for test cleanup
var _ = AfterEach(func() {
	By("Cleaning up handler test data")
	
	// Cleanup is handled by BeforeEach, but can add specific cleanup here
	fixtureManager.CleanupFixtures()
})