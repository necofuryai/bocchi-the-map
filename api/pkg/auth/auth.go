// Package auth provides JWT authentication and authorization functionality for the Bocchi The Map API.
//
// This package implements Auth0 JWT validation middleware that integrates with the Onion Architecture
// pattern used throughout the application. It provides both Chi router middleware and Huma v2 
// compatible middleware functions.
//
// Features:
// - Auth0 JWT token validation
// - Rate limiting for authentication endpoints
// - Request context user information
// - Permission-based authorization
// - Integration with New Relic monitoring
// - Comprehensive error handling
//
// Example usage:
//
//	// Initialize Auth middleware
//	authConfig := auth.AuthConfig{
//		Auth0Domain:   "your-domain.auth0.com",
//		Auth0Audience: "your-api-audience",
//		JWTSecret:     "your-jwt-secret",
//		Development:   false,
//	}
//	
//	authMiddleware, err := auth.NewAuthMiddlewareWithConfig(authConfig, queries)
//	if err != nil {
//		log.Fatal(err)
//	}
//	
//	// Use with Chi router
//	router.Use(authMiddleware.RequireAuth())
//	
//	// Use with Huma v2
//	api.UseMiddleware(authMiddleware.HumaMiddleware())
//
// The package follows the project's design principles:
// - BDD testing with Ginkgo framework
// - Comprehensive monitoring and observability
// - Protocol Buffers for type-safe contracts
// - Clear dependency boundaries
// - Security and performance focused
package auth

import (
	"context"
	"time"

	"bocchi/api/infrastructure/database"
	"bocchi/api/pkg/config"
	"bocchi/api/pkg/errors"
	"bocchi/api/pkg/logger"
)

// Service provides a high-level interface for authentication operations
type Service struct {
	middleware   *AuthMiddleware
	rateLimiter  *RateLimiter
	validator    *JWTValidator
	queries      *database.Queries
}

// ServiceConfig holds configuration for the authentication service
type ServiceConfig struct {
	Auth0Domain     string
	Auth0Audience   string
	JWTSecret       string
	Development     bool
	RateLimit       int
	RateLimitWindow time.Duration
	SkipPaths       []string
}

// NewService creates a new authentication service with all components
func NewService(config ServiceConfig, queries *database.Queries) (*Service, error) {
	// Create JWT validator
	validator, err := NewJWTValidator(config.Auth0Domain, config.Auth0Audience)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrTypeInternal, "failed to create JWT validator")
	}

	// Create rate limiter
	rateLimiter := NewRateLimiter(config.RateLimit, config.RateLimitWindow)

	// Create auth middleware
	authConfig := AuthConfig{
		Auth0Domain:   config.Auth0Domain,
		Auth0Audience: config.Auth0Audience,
		JWTSecret:     config.JWTSecret,
		Development:   config.Development,
		SkipPaths:     config.SkipPaths,
	}

	middleware, err := NewAuthMiddlewareWithConfig(authConfig, queries)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrTypeInternal, "failed to create auth middleware")
	}

	service := &Service{
		middleware:  middleware,
		rateLimiter: rateLimiter,
		validator:   validator,
		queries:     queries,
	}

	logger.InfoWithFields("Authentication service initialized", map[string]interface{}{
		"auth0_domain":      config.Auth0Domain,
		"auth0_audience":    config.Auth0Audience,
		"development":       config.Development,
		"rate_limit":        config.RateLimit,
		"rate_limit_window": config.RateLimitWindow.String(),
	})

	return service, nil
}

// NewServiceFromConfig creates a new authentication service from application config
func NewServiceFromConfig(cfg *config.Config, queries *database.Queries) (*Service, error) {
	serviceConfig := ServiceConfig{
		Auth0Domain:     cfg.Auth.Auth0Domain,
		Auth0Audience:   cfg.Auth.Auth0Audience,
		JWTSecret:       cfg.Auth.JWTSecret,
		Development:     cfg.App.Environment == "development",
		RateLimit:       5,                // 5 requests per window
		RateLimitWindow: 5 * time.Minute,  // 5 minute window
		SkipPaths: []string{
			"/health",
			"/metrics",
			"/swagger",
			"/docs",
		},
	}

	return NewService(serviceConfig, queries)
}

// GetMiddleware returns the authentication middleware
func (s *Service) GetMiddleware() *AuthMiddleware {
	return s.middleware
}

// GetRateLimiter returns the rate limiter
func (s *Service) GetRateLimiter() *RateLimiter {
	return s.rateLimiter
}

// GetValidator returns the JWT validator
func (s *Service) GetValidator() *JWTValidator {
	return s.validator
}

// ValidateToken validates a JWT token and returns claims
func (s *Service) ValidateToken(ctx context.Context, token string) (*Claims, error) {
	return s.validator.ValidateToken(token)
}

// CheckPermission validates that a user has a specific permission
func (s *Service) CheckPermission(ctx context.Context, permission string) error {
	if !HasPermission(ctx, permission) {
		return errors.Forbidden("permission", permission)
	}
	return nil
}

// GetUserInfo extracts user information from the request context
func (s *Service) GetUserInfo(ctx context.Context) (map[string]interface{}, error) {
	user, ok := GetUserFromContext(ctx)
	if !ok {
		return nil, errors.Unauthorized("user context not found")
	}
	return user, nil
}

// RequireUser ensures a user is authenticated and returns user info
func (s *Service) RequireUser(ctx context.Context) (map[string]interface{}, error) {
	return s.GetUserInfo(ctx)
}

// RequireUserID ensures a user is authenticated and returns the user ID
func (s *Service) RequireUserID(ctx context.Context) (string, error) {
	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		return "", errors.Unauthorized("user ID not found in context")
	}
	return userID, nil
}

// RequireUserEmail ensures a user is authenticated and returns the user email
func (s *Service) RequireUserEmail(ctx context.Context) (string, error) {
	email, ok := GetUserEmailFromContext(ctx)
	if !ok {
		return "", errors.Unauthorized("user email not found in context")
	}
	return email, nil
}

// Stop stops all authentication service components
func (s *Service) Stop() {
	if s.rateLimiter != nil {
		s.rateLimiter.Stop()
	}
	
	logger.Info("Authentication service stopped")
}

// GetStats returns statistics about the authentication service
func (s *Service) GetStats() map[string]interface{} {
	stats := make(map[string]interface{})
	
	if s.rateLimiter != nil {
		stats["rate_limiter"] = s.rateLimiter.GetStats()
	}
	
	return stats
}

// Health checks the health of the authentication service
func (s *Service) Health(ctx context.Context) error {
	// Check if validator is working by validating the configuration
	if s.validator == nil {
		return errors.Internal("JWT validator not initialized")
	}
	
	// Additional health checks can be added here
	// For example, checking Auth0 connectivity, database connectivity, etc.
	
	return nil
}