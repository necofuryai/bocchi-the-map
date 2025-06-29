//go:build integration

package handlers

import (
	"testing"

	"github.com/necofuryai/bocchi-the-map/api/tests/helpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Global test suite for handler integration tests
var testSuite *helpers.CommonTestSuite

func TestHandlers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "HTTP Handlers Integration Test Suite")
}

var _ = BeforeSuite(func() {
	By("Setting up handler integration test environment")
	
	testSuite = helpers.NewCommonTestSuite()
	
	By("Handler integration test environment setup completed")
})

var _ = AfterSuite(func() {
	By("Cleaning up handler integration test environment")
	
	if testSuite != nil {
		testSuite.Cleanup()
	}
	
	By("Handler integration test environment cleanup completed")
})

// BeforeEach hook for test isolation
var _ = BeforeEach(func() {
	testSuite.PrepareCleanTestData()
})

// AfterEach hook for test cleanup
var _ = AfterEach(func() {
	testSuite.CleanupTestData()
})