package auth

import (
	"context"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/necofuryai/bocchi-the-map/api/infrastructure/database"
	"github.com/stretchr/testify/assert"
)

// Mock TokenBlacklistQuerier for testing
type mockTokenBlacklistQuerier struct {
	isBlacklistedFunc           func(ctx context.Context, jti string) (bool, error)
	blacklistAccessTokenFunc   func(ctx context.Context, arg database.BlacklistAccessTokenParams) error
	blacklistRefreshTokenFunc  func(ctx context.Context, arg database.BlacklistRefreshTokenParams) error
}

func (m *mockTokenBlacklistQuerier) IsTokenBlacklisted(ctx context.Context, jti string) (bool, error) {
	if m.isBlacklistedFunc != nil {
		return m.isBlacklistedFunc(ctx, jti)
	}
	return false, nil
}

func (m *mockTokenBlacklistQuerier) BlacklistAccessToken(ctx context.Context, arg database.BlacklistAccessTokenParams) error {
	if m.blacklistAccessTokenFunc != nil {
		return m.blacklistAccessTokenFunc(ctx, arg)
	}
	return nil
}

func (m *mockTokenBlacklistQuerier) BlacklistRefreshToken(ctx context.Context, arg database.BlacklistRefreshTokenParams) error {
	if m.blacklistRefreshTokenFunc != nil {
		return m.blacklistRefreshTokenFunc(ctx, arg)
	}
	return nil
}

func TestJWTClaims_TokenType(t *testing.T) {
	middleware := NewAuthMiddleware("test-secret", &mockTokenBlacklistQuerier{})

	t.Run("GenerateToken should create access token with correct TokenType", func(t *testing.T) {
		userID := "test-user-id"
		email := "test@example.com"

		tokenString, err := middleware.GenerateToken(userID, email)
		assert.NoError(t, err)
		assert.NotEmpty(t, tokenString)

		// Parse the token and verify claims
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("test-secret"), nil
		})
		assert.NoError(t, err)

		claims, ok := token.Claims.(*JWTClaims)
		assert.True(t, ok)
		assert.Equal(t, "access", claims.TokenType)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, email, claims.Email)
	})

	t.Run("GenerateRefreshToken should create refresh token with correct TokenType", func(t *testing.T) {
		userID := "test-user-id"
		email := "test@example.com"

		tokenString, err := middleware.GenerateRefreshToken(userID, email)
		assert.NoError(t, err)
		assert.NotEmpty(t, tokenString)

		// Parse the token and verify claims
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("test-secret"), nil
		})
		assert.NoError(t, err)

		claims, ok := token.Claims.(*JWTClaims)
		assert.True(t, ok)
		assert.Equal(t, "refresh", claims.TokenType)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, email, claims.Email)
	})
}

func TestBlacklistToken_TokenType(t *testing.T) {
	accessTokenBlacklisted := false
	refreshTokenBlacklisted := false

	mock := &mockTokenBlacklistQuerier{
		blacklistAccessTokenFunc: func(ctx context.Context, arg database.BlacklistAccessTokenParams) error {
			accessTokenBlacklisted = true
			return nil
		},
		blacklistRefreshTokenFunc: func(ctx context.Context, arg database.BlacklistRefreshTokenParams) error {
			refreshTokenBlacklisted = true
			return nil
		},
	}

	middleware := NewAuthMiddleware("test-secret", mock)

	t.Run("BlacklistToken should blacklist access token correctly", func(t *testing.T) {
		userID := "test-user-id"
		email := "test@example.com"

		// Generate access token
		accessToken, err := middleware.GenerateToken(userID, email)
		assert.NoError(t, err)

		// Blacklist it
		err = middleware.BlacklistToken(context.Background(), accessToken, "test")
		assert.NoError(t, err)
		assert.True(t, accessTokenBlacklisted)
		assert.False(t, refreshTokenBlacklisted)
	})

	t.Run("BlacklistToken should blacklist refresh token correctly", func(t *testing.T) {
		// Reset flags
		accessTokenBlacklisted = false
		refreshTokenBlacklisted = false

		userID := "test-user-id"
		email := "test@example.com"

		// Generate refresh token
		refreshToken, err := middleware.GenerateRefreshToken(userID, email)
		assert.NoError(t, err)

		// Blacklist it
		err = middleware.BlacklistToken(context.Background(), refreshToken, "test")
		assert.NoError(t, err)
		assert.False(t, accessTokenBlacklisted)
		assert.True(t, refreshTokenBlacklisted)
	})
}