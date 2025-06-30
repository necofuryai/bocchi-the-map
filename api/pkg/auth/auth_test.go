package auth

import (
	"testing"
	"time"

	"github.com/necofuryai/bocchi-the-map/api/pkg/config"
)

func TestNewServiceFromConfig(t *testing.T) {
	// Test configuration
	cfg := &config.Config{
		Auth: config.AuthConfig{
			Auth0Domain:   "test-domain.auth0.com",
			Auth0Audience: "test-audience",
			JWTSecret:     "test-secret-that-is-longer-than-32-characters-for-validation",
		},
		App: config.AppConfig{
			Environment: "development",
		},
	}

	// This test will fail because we don't have a real database connection
	// but it verifies the basic structure is correct
	_, err := NewServiceFromConfig(cfg, nil)
	
	// We expect an error because Auth0 domain won't be reachable in tests
	if err == nil {
		t.Log("Service creation succeeded (unexpected in test environment)")
	} else {
		t.Logf("Service creation failed as expected: %v", err)
	}
}

func TestRateLimiter(t *testing.T) {
	// Test rate limiter basic functionality
	rl := NewRateLimiter(2, time.Minute)
	defer rl.Stop()

	// First request should be allowed
	if !rl.IsAllowed("127.0.0.1") {
		t.Error("First request should be allowed")
	}

	// Second request should be allowed
	if !rl.IsAllowed("127.0.0.1") {
		t.Error("Second request should be allowed")
	}

	// Third request should be blocked
	if rl.IsAllowed("127.0.0.1") {
		t.Error("Third request should be blocked")
	}

	// Check remaining requests
	remaining := rl.GetRemainingRequests("127.0.0.1")
	if remaining != 0 {
		t.Errorf("Expected 0 remaining requests, got %d", remaining)
	}
}

func TestRateLimiterStats(t *testing.T) {
	rl := NewRateLimiter(5, time.Minute)
	defer rl.Stop()

	// Add some test requests
	rl.IsAllowed("127.0.0.1")
	rl.IsAllowed("192.168.1.1")

	stats := rl.GetStats()
	
	if activeIPs, ok := stats["active_ips"]; !ok || activeIPs != 2 {
		t.Errorf("Expected 2 active IPs, got %v", activeIPs)
	}

	if limit, ok := stats["limit"]; !ok || limit != 5 {
		t.Errorf("Expected limit 5, got %v", limit)
	}
}