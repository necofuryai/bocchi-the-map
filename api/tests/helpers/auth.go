package helpers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/necofuryai/bocchi-the-map/api/domain/entities"
	. "github.com/onsi/gomega"
)

// AuthHelper provides utilities for authentication testing
type AuthHelper struct {
	jwtSecret string
}

// NewAuthHelper creates a new authentication helper
func NewAuthHelper() *AuthHelper {
	return &AuthHelper{
		// This JWT secret is hardcoded intentionally for testing purposes only
		// NEVER use this secret in production environments
		jwtSecret: "test-secret-key-for-bdd-testing",
	}
}

// CreateTestUser creates a test user entity with common defaults
func (ah *AuthHelper) CreateTestUser(overrides ...func(*entities.User)) *entities.User {
	user := &entities.User{
		ID:       "test-user-" + generateRandomID(),
		Email:    "test@example.com",
		DisplayName: "Test User",
		AuthProvider: entities.AuthProviderGoogle,
		Preferences: entities.UserPreferences{
			Language: "en",
			DarkMode: false,
			Timezone: "UTC",
		},
	}
	
	// Apply any overrides
	for _, override := range overrides {
		override(user)
	}
	
	return user
}

// CreateValidJWTToken creates a valid JWT token for testing
func (ah *AuthHelper) CreateValidJWTToken(userID string) string {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(ah.jwtSecret))
	Expect(err).NotTo(HaveOccurred(), "Failed to create JWT token")
	
	return tokenString
}

// CreateExpiredJWTToken creates an expired JWT token for testing
func (ah *AuthHelper) CreateExpiredJWTToken(userID string) string {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(-time.Hour).Unix(), // Expired 1 hour ago
		"iat": time.Now().Add(-2 * time.Hour).Unix(),
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(ah.jwtSecret))
	Expect(err).NotTo(HaveOccurred(), "Failed to create expired JWT token")
	
	return tokenString
}

// CreateInvalidJWTToken creates an invalid JWT token for testing
func (ah *AuthHelper) CreateInvalidJWTToken() string {
	return "invalid.jwt.token"
}

// contextKey is a private type for context keys to avoid collisions
type contextKey string

// userIDKey is the context key for user ID
const userIDKey contextKey = "user_id"

// CreateAuthenticatedContext creates a context with user authentication
func (ah *AuthHelper) CreateAuthenticatedContext(userID string) context.Context {
	ctx := context.Background()
	return context.WithValue(ctx, userIDKey, userID)
}

// CreateUnauthenticatedContext creates a context without authentication
func (ah *AuthHelper) CreateUnauthenticatedContext() context.Context {
	return context.Background()
}

// AuthTestData provides common authentication test scenarios
type AuthTestData struct {
	ValidUserID     string
	ValidToken      string
	ExpiredToken    string
	InvalidToken    string
	TestUser        *entities.User
}

// NewAuthTestData creates a complete set of authentication test data
func (ah *AuthHelper) NewAuthTestData() *AuthTestData {
	userID := "test-user-" + generateRandomID()
	testUser := ah.CreateTestUser(func(u *entities.User) {
		u.ID = userID
	})
	
	return &AuthTestData{
		ValidUserID:  userID,
		ValidToken:   ah.CreateValidJWTToken(userID),
		ExpiredToken: ah.CreateExpiredJWTToken(userID),
		InvalidToken: ah.CreateInvalidJWTToken(),
		TestUser:     testUser,
	}
}

var fallbackCounter int64

// generateRandomID generates a random ID for testing
func generateRandomID() string {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp with counter to ensure uniqueness
		counter := atomic.AddInt64(&fallbackCounter, 1)
		return fmt.Sprintf("%s_%d", time.Now().Format("20060102150405.000000"), counter)
	}
	return hex.EncodeToString(bytes)
}