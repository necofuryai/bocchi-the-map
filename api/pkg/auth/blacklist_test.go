package auth

import (
	"context"
	"testing"
	"time"
)

// Test helper functions for context manipulation
func TestClaimsWithJTI(t *testing.T) {
	// Test Claims struct with JTI
	claims := &Claims{
		ID:      "test-jti-123",
		Subject: "test-user-456",
	}

	if claims.ID != "test-jti-123" {
		t.Errorf("Expected JTI 'test-jti-123', got '%s'", claims.ID)
	}

	if claims.Subject != "test-user-456" {
		t.Errorf("Expected Subject 'test-user-456', got '%s'", claims.Subject)
	}
}

func TestClaimsWithoutJTI(t *testing.T) {
	// Test Claims struct without JTI
	claims := &Claims{
		Subject: "test-user-789",
	}

	if claims.ID != "" {
		t.Errorf("Expected empty JTI, got '%s'", claims.ID)
	}

	if claims.Subject != "test-user-789" {
		t.Errorf("Expected Subject 'test-user-789', got '%s'", claims.Subject)
	}
}

func TestJTIExtractionFromContext(t *testing.T) {
	ctx := context.Background()

	// Test with no JTI in context
	jti, hasJTI := GetJTIFromContext(ctx)
	if hasJTI || jti != "" {
		t.Errorf("Expected no JTI in empty context, got: %s", jti)
	}

	// Test with JTI in context
	testJTI := "test-jti-123"
	ctxWithJTI := context.WithValue(ctx, "jti", testJTI)
	
	jti, hasJTI = GetJTIFromContext(ctxWithJTI)
	if !hasJTI || jti != testJTI {
		t.Errorf("Expected JTI %s, got: %s (hasJTI: %v)", testJTI, jti, hasJTI)
	}
}

func TestTokenExpirationFromContext(t *testing.T) {
	ctx := context.Background()

	// Test with no expiration in context
	exp, hasExp := GetTokenExpirationFromContext(ctx)
	if hasExp || !exp.IsZero() {
		t.Errorf("Expected no expiration in empty context, got: %v", exp)
	}

	// Test with expiration in context
	testExpiration := time.Now().Add(1 * time.Hour)
	ctxWithExp := context.WithValue(ctx, "token_expires_at", testExpiration)
	
	exp, hasExp = GetTokenExpirationFromContext(ctxWithExp)
	if !hasExp || !exp.Equal(testExpiration) {
		t.Errorf("Expected expiration %v, got: %v (hasExp: %v)", testExpiration, exp, hasExp)
	}
}