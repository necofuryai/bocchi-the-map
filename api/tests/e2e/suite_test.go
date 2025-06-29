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

func init() {
	testSuite = helpers.SetupGinkgoHooks("E2E")
}