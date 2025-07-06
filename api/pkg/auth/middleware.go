package auth

import (
	"context"
	stdErrors "errors"
	"net/http"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-chi/chi/v5"

	"github.com/necofuryai/bocchi-the-map/api/infrastructure/database"
	"github.com/necofuryai/bocchi-the-map/api/pkg/errors"
	"github.com/necofuryai/bocchi-the-map/api/pkg/logger"
	"github.com/necofuryai/bocchi-the-map/api/pkg/monitoring"
)

// AuthMiddleware provides JWT authentication middleware for the API
type AuthMiddleware struct {
	validator    *JWTValidator
	queries      *database.Queries
	jwtSecret    string
	skipPaths    map[string]bool
	development  bool
}

// AuthConfig holds configuration for the authentication middleware
type AuthConfig struct {
	Auth0Domain   string
	Auth0Audience string
	JWTSecret     string
	Development   bool
	SkipPaths     []string
}

// NewAuthMiddleware creates a new authentication middleware instance
func NewAuthMiddleware(jwtSecret string, queries *database.Queries) *AuthMiddleware {
	// For now, we'll create a placeholder that will be properly configured
	// when the full Auth0 configuration is available
	return &AuthMiddleware{
		queries:   queries,
		jwtSecret: jwtSecret,
		skipPaths: make(map[string]bool),
	}
}

// NewAuthMiddlewareWithConfig creates a new authentication middleware with full configuration
func NewAuthMiddlewareWithConfig(config AuthConfig, queries *database.Queries) (*AuthMiddleware, error) {
	// Initialize JWT validator with Auth0 configuration
	validator, err := NewJWTValidator(config.Auth0Domain, config.Auth0Audience)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrTypeInternal, "failed to create JWT validator")
	}

	// Create skip paths map for performance
	skipPaths := make(map[string]bool)
	defaultSkipPaths := []string{
		"/health",
		"/metrics",
		"/debug",
		"/favicon.ico",
	}
	
	// Add default skip paths
	for _, path := range defaultSkipPaths {
		skipPaths[path] = true
	}
	
	// Add custom skip paths
	for _, path := range config.SkipPaths {
		skipPaths[path] = true
	}

	middleware := &AuthMiddleware{
		validator:   validator,
		queries:     queries,
		jwtSecret:   config.JWTSecret,
		skipPaths:   skipPaths,
		development: config.Development,
	}

	logger.InfoWithFields("Auth middleware initialized", map[string]interface{}{
		"auth0_domain":   config.Auth0Domain,
		"auth0_audience": config.Auth0Audience,
		"development":    config.Development,
		"skip_paths":     len(skipPaths),
	})

	return middleware, nil
}

// RequireAuth is a Chi middleware that requires valid JWT authentication
func (m *AuthMiddleware) RequireAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if path should be skipped
			if m.shouldSkipPath(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			// Add request tracing for monitoring
			ctx := monitoring.StartTrace(r.Context(), "auth.validate_token")
			defer monitoring.EndTrace(ctx)

			// Validate token
			claims, err := m.validateRequest(r)
			if err != nil {
				m.handleAuthError(w, r, err)
				return
			}

			// Add user context to request
			userCtx := m.validator.GetUserContext(r.Context(), claims)
			
			// Add JWT ID (JTI) to context for logout functionality
			if claims.ID != "" {
				userCtx = context.WithValue(userCtx, "jti", claims.ID)
			}
			
			// Add token expiration to context
			if claims.ExpiresAt > 0 {
				userCtx = context.WithValue(userCtx, "token_expires_at", time.Unix(claims.ExpiresAt, 0))
			}
			
			// Add user info for monitoring
			monitoring.AddUserContext(userCtx, claims.Subject, claims.Email)

			// Continue with authenticated request
			next.ServeHTTP(w, r.WithContext(userCtx))
		})
	}
}

// HumaMiddleware returns a Huma v2 compatible middleware function
// This creates a wrapper that integrates Auth0 authentication with Huma v2 operations
func (m *AuthMiddleware) HumaMiddleware() func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		// This is a placeholder for Huma v2 specific middleware
		// In practice, it's better to use the Chi middleware with humachi adapter
		next(ctx)
	}
}

// CreateProtectedOperation creates a Huma operation that requires authentication
// This is a helper function to easily create protected endpoints
func (m *AuthMiddleware) CreateProtectedOperation(operation huma.Operation) huma.Operation {
	// Add security requirement to the operation
	if operation.Security == nil {
		operation.Security = []map[string][]string{}
	}
	
	// Add bearer token security requirement
	operation.Security = append(operation.Security, map[string][]string{
		"bearerAuth": {},
	})
	
	return operation
}

// OptionalAuth is a Chi middleware that optionally validates JWT authentication
func (m *AuthMiddleware) OptionalAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Always allow the request to continue, but add user context if token is valid
			claims, err := m.validateRequest(r)
			if err != nil {
				// Log the error but don't block the request
				logger.InfoWithFields("Optional auth failed", map[string]interface{}{
					"error": err.Error(),
					"path":  r.URL.Path,
				})
				next.ServeHTTP(w, r)
				return
			}

			// Add user context if validation succeeded
			userCtx := m.validator.GetUserContext(r.Context(), claims)
			
			// Add JWT ID (JTI) to context for logout functionality
			if claims.ID != "" {
				userCtx = context.WithValue(userCtx, "jti", claims.ID)
			}
			
			// Add token expiration to context
			if claims.ExpiresAt > 0 {
				userCtx = context.WithValue(userCtx, "token_expires_at", time.Unix(claims.ExpiresAt, 0))
			}
			
			monitoring.AddUserContext(userCtx, claims.Subject, claims.Email)
			
			next.ServeHTTP(w, r.WithContext(userCtx))
		})
	}
}

// validateRequest validates the JWT token from the request
func (m *AuthMiddleware) validateRequest(r *http.Request) (*Claims, error) {
	if m.validator == nil {
		return nil, errors.Internal("JWT validator not initialized")
	}

	// Extract and validate token
	claims, err := m.validator.ValidateTokenFromRequest(r)
	if err != nil {
		return nil, err
	}

	// Additional validation can be added here (e.g., token blacklist check)
	if m.queries != nil {
		if err := m.checkTokenBlacklist(r.Context(), claims); err != nil {
			return nil, err
		}
	}

	return claims, nil
}

// checkTokenBlacklist checks if the token is blacklisted
func (m *AuthMiddleware) checkTokenBlacklist(ctx context.Context, claims *Claims) error {
	// Check if JWT ID (JTI) is available
	if claims.ID == "" {
		logger.InfoWithFields("Token missing JTI", map[string]interface{}{
			"subject": claims.Subject,
		})
		// Allow tokens without JTI for now (Auth0 might not always include JTI)
		return nil
	}

	// Check if token is blacklisted
	isBlacklisted, err := m.queries.IsTokenBlacklisted(ctx, claims.ID)
	if err != nil {
		logger.ErrorWithFields("Failed to check token blacklist", err, map[string]interface{}{
			"jti": claims.ID,
		})
		// Don't block authentication on database errors
		return nil
	}

	if isBlacklisted {
		logger.InfoWithFields("Token has been revoked", map[string]interface{}{
			"jti":     claims.ID,
			"subject": claims.Subject,
		})
		return errors.Unauthorized("token has been revoked")
	}

	return nil
}

// shouldSkipPath checks if the given path should skip authentication
func (m *AuthMiddleware) shouldSkipPath(path string) bool {
	// Check exact match
	if m.skipPaths[path] {
		return true
	}

	// Check for path prefixes that should be skipped
	skipPrefixes := []string{
		"/health",
		"/metrics",
		"/debug",
		"/swagger",
		"/docs",
	}

	for _, prefix := range skipPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}

	return false
}

// handleAuthError handles authentication errors for Chi middleware
func (m *AuthMiddleware) handleAuthError(w http.ResponseWriter, r *http.Request, err error) {
	// Extract domain error information
	var domainErr *errors.DomainError
	if !stdErrors.As(err, &domainErr) {
		domainErr = errors.Unauthorized("authentication failed")
	}

	// Log the authentication failure
	logger.InfoWithFields("Authentication failed", map[string]interface{}{
		"error":  err.Error(),
		"path":   r.URL.Path,
		"method": r.Method,
		"ip":     r.RemoteAddr,
	})

	// Add monitoring metrics
	monitoring.RecordAuthFailure(r.Context(), string(domainErr.Type))

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("WWW-Authenticate", "Bearer")
	
	// Set status code
	statusCode := domainErr.ToHTTPStatus()
	w.WriteHeader(statusCode)

	// Write JSON response (simplified for this implementation)
	w.Write([]byte(`{"error":{"type":"` + string(domainErr.Type) + `","message":"` + domainErr.Message + `"}}`))
}

// convertToHumaError converts domain errors to Huma errors
func (m *AuthMiddleware) convertToHumaError(err error) *huma.ErrorModel {
	var domainErr *errors.DomainError
	if !stdErrors.As(err, &domainErr) {
		domainErr = errors.Unauthorized("authentication failed")
	}

	// Log the authentication failure
	logger.InfoWithFields("Huma authentication failed", map[string]interface{}{
		"error": err.Error(),
		"type":  domainErr.Type,
	})

	return &huma.ErrorModel{
		Type:   string(domainErr.Type),
		Title:  "Authentication Failed",
		Status: domainErr.ToHTTPStatus(),
		Detail: domainErr.Message,
	}
}

// RegisterAuthRoutes registers authentication-related routes (placeholder)
func (m *AuthMiddleware) RegisterAuthRoutes(router chi.Router) {
	// TODO: Implement authentication routes like login, logout, refresh token
	// This would include routes for OAuth2 flow with Auth0
}

// Logout handles user logout by invalidating tokens
func (m *AuthMiddleware) Logout(ctx context.Context, tokenJTI string, expiresAt time.Time) error {
	if tokenJTI == "" {
		return errors.InvalidInput("tokenJTI", "JWT ID is required for logout")
	}

	// Add token to blacklist
	err := m.queries.BlacklistAccessToken(ctx, database.BlacklistAccessTokenParams{
		Jti:       tokenJTI,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		logger.ErrorWithFields("Failed to blacklist token", err, map[string]interface{}{
			"jti": tokenJTI,
		})
		return errors.Wrap(err, errors.ErrTypeInternal, "failed to invalidate token")
	}

	logger.InfoWithFields("Token blacklisted successfully", map[string]interface{}{
		"jti":        tokenJTI,
		"expires_at": expiresAt,
	})

	return nil
}

// RefreshToken handles token refresh (placeholder)
func (m *AuthMiddleware) RefreshToken(ctx context.Context, refreshToken string) (*Claims, error) {
	// TODO: Implement token refresh logic
	// This would involve validating the refresh token and issuing a new access token
	return nil, errors.Internal("token refresh not implemented")
}

// GetValidator returns the JWT validator instance
func (m *AuthMiddleware) GetValidator() *JWTValidator {
	return m.validator
}