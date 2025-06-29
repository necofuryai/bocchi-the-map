//go:build e2e

package e2e

import (
	"testing"

	"github.com/necofuryai/bocchi-the-map/api/tests/helpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Global test suite for E2E tests
var testSuite *helpers.CommonTestSuite

func TestE2EScenarios(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bocchi The Map API E2E Test Suite")
}

var _ = BeforeSuite(func() {
	By("Setting up E2E test environment")
	
	testSuite = helpers.NewCommonTestSuite()
	
	By("E2E test environment setup completed")
})

var _ = AfterSuite(func() {
	By("Cleaning up E2E test environment")
	
	if testSuite != nil {
		testSuite.Cleanup()
	}
	
	By("E2E test environment cleanup completed")
})

// BeforeEach hook for test isolation
var _ = BeforeEach(func() {
	testSuite.PrepareCleanTestData()
})

// AfterEach hook for test cleanup
var _ = AfterEach(func() {
	testSuite.CleanupTestData()
})