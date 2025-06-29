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

func init() {
	testSuite = helpers.SetupGinkgoHooks("handler integration")
}