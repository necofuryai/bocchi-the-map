package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/necofuryai/bocchi-the-map/api/infrastructure/database"
	"github.com/necofuryai/bocchi-the-map/api/pkg/errors"
)

// TokenBlacklistQuerier interface for token blacklist operations
type TokenBlacklistQuerier interface {
	IsTokenBlacklisted(ctx context.Context, jti string) (bool, error)
	BlacklistAccessToken(ctx context.Context, arg database.BlacklistAccessTokenParams) error
	BlacklistRefreshToken(ctx context.Context, arg database.BlacklistRefreshTokenParams) error
}

// AuthMiddleware validates JWT tokens and extracts user context
type AuthMiddleware struct {
	jwtSecret []byte
	queries   TokenBlacklistQuerier
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(jwtSecret string, queries TokenBlacklistQuerier) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: []byte(jwtSecret),
		queries:   queries,
	}
}

// JWTClaims represents the claims in the JWT token
type JWTClaims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	TokenType string `json:"token_type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

// extractAndValidateToken extracts JWT token from request and validates it
func (am *AuthMiddleware) extractAndValidateToken(r *http.Request) (*JWTClaims, error) {
	// Extract token from Authorization header or cookie
	var tokenString string
	
	// First try Authorization header (Bearer token)
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			tokenString = parts[1]
		}
	}
	
	// If no Bearer token, try cookie
	if tokenString == "" {
		if cookie, err := r.Cookie("bocchi_access_token"); err == nil {
			tokenString = cookie.Value
		}
	}
	
	// If still no token found, return error
	if tokenString == "" {
		return nil, fmt.Errorf("no token found")
	}

	// Parse and validate JWT token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return am.jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Extract claims
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Check if token is blacklisted (if blacklist querier is available)
	if am.queries != nil && claims.ID != "" {
		ctx := r.Context()
		isBlacklisted, err := am.queries.IsTokenBlacklisted(ctx, claims.ID)
		if err != nil {
			return nil, fmt.Errorf("authentication service error: %w", err)
		}
		if isBlacklisted {
			return nil, fmt.Errorf("token has been revoked")
		}
	}

	return claims, nil
}

// ExtractAndValidateTokenFromContext extracts JWT token from Huma context and validates it
func (am *AuthMiddleware) ExtractAndValidateTokenFromContext(ctx huma.Context) (*JWTClaims, error) {
	// Extract token from Authorization header or cookie
	var tokenString string
	
	// First try Authorization header (Bearer token)
	authHeader := ctx.Header("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			tokenString = parts[1]
		}
	}
	
	// If no Bearer token, try cookie
	if tokenString == "" {
		if cookie := ctx.Header("Cookie"); cookie != "" {
			// Parse cookies manually since Huma doesn't provide direct cookie access
			cookies := strings.Split(cookie, ";")
			for _, c := range cookies {
				parts := strings.SplitN(strings.TrimSpace(c), "=", 2)
				if len(parts) == 2 && parts[0] == "bocchi_access_token" {
					tokenString = parts[1]
					break
				}
			}
		}
	}
	
	// If still no token found, return error
	if tokenString == "" {
		return nil, fmt.Errorf("no token found")
	}

	// Parse and validate JWT token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return am.jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Extract claims
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Check if token is blacklisted (if blacklist querier is available)
	if am.queries != nil && claims.ID != "" {
		requestCtx := ctx.Context()
		isBlacklisted, err := am.queries.IsTokenBlacklisted(requestCtx, claims.ID)
		if err != nil {
			return nil, fmt.Errorf("authentication service error: %w", err)
		}
		if isBlacklisted {
			return nil, fmt.Errorf("token has been revoked")
		}
	}

	return claims, nil
}

// Middleware returns the HTTP middleware function
func (am *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := am.extractAndValidateToken(r)
		if err != nil {
			if strings.Contains(err.Error(), "no token found") {
				http.Error(w, "Authentication required - no valid token found", http.StatusUnauthorized)
			} else if strings.Contains(err.Error(), "authentication service error") {
				http.Error(w, "Authentication service error", http.StatusInternalServerError)
			} else {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
			}
			return
		}

		// Add user context to request
		ctx := r.Context()
		ctx = errors.WithUserID(ctx, claims.UserID)
		ctx = errors.WithRequestID(ctx, r.Header.Get("X-Request-ID"))

		// Call next handler with user context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalMiddleware allows requests without authentication but adds user context if present
func (am *AuthMiddleware) OptionalMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := am.extractAndValidateToken(r)
		if err != nil {
			// Invalid token or no token, continue without user context
			next.ServeHTTP(w, r)
			return
		}

		// Add user context to request
		ctx := r.Context()
		ctx = errors.WithUserID(ctx, claims.UserID)
		ctx = errors.WithRequestID(ctx, r.Header.Get("X-Request-ID"))

		// Call next handler with user context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// generateTokenWithExpiration is a helper method that creates JWT tokens with specified expiration and type
func (am *AuthMiddleware) generateTokenWithExpiration(userID, email string, expiration time.Duration, tokenType string) (string, error) {
	// Generate unique JWT ID for token tracking and revocation
	jti := uuid.New().String()
	
	// Create token claims with user information
	claims := &JWTClaims{
		UserID:    userID,
		Email:     email,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "bocchi-the-map-api",
			Subject:   userID,
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret
	tokenString, err := token.SignedString(am.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign %s JWT token: %w", tokenType, err)
	}

	return tokenString, nil
}

// GenerateToken generates a JWT token for authenticated users
func (am *AuthMiddleware) GenerateToken(userID, email string) (string, error) {
	return am.generateTokenWithExpiration(userID, email, 24*time.Hour, "access")
}

// GenerateRefreshToken generates a longer-lived refresh token for token renewal
func (am *AuthMiddleware) GenerateRefreshToken(userID, email string) (string, error) {
	return am.generateTokenWithExpiration(userID, email, 7*24*time.Hour, "refresh")
}

// BlacklistToken adds a token to the blacklist for revocation
func (am *AuthMiddleware) BlacklistToken(ctx context.Context, tokenString string, reason string) error {
	if am.queries == nil {
		return fmt.Errorf("token blacklist not configured")
	}

	// Parse token to extract claims
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return am.jwtSecret, nil
	})
	if err != nil {
		return fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || claims.ID == "" {
		return fmt.Errorf("invalid token claims or missing JWT ID")
	}

	// Determine token type and add to blacklist
	expiresAt := claims.ExpiresAt.Time
	if claims.TokenType == "refresh" {
		// Refresh token
		return am.queries.BlacklistRefreshToken(ctx, database.BlacklistRefreshTokenParams{
			Jti:       claims.ID,
			UserID:    claims.UserID,
			ExpiresAt: expiresAt,
		})
	} else {
		// Access token (default for backwards compatibility)
		return am.queries.BlacklistAccessToken(ctx, database.BlacklistAccessTokenParams{
			Jti:       claims.ID,
			UserID:    claims.UserID,
			ExpiresAt: expiresAt,
		})
	}
}

// ValidateToken validates a JWT token and returns the claims
func (am *AuthMiddleware) ValidateToken(tokenString string) (*JWTClaims, error) {
	// Parse and validate JWT token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return am.jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT token: %w", err)
	}

	// Extract claims
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}