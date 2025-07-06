package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"

	"github.com/necofuryai/bocchi-the-map/api/application/clients"
	"github.com/necofuryai/bocchi-the-map/api/pkg/auth"
	"github.com/necofuryai/bocchi-the-map/api/pkg/logger"
	"github.com/necofuryai/bocchi-the-map/api/pkg/monitoring"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	authMiddleware *auth.AuthMiddleware
	userClient     *clients.UserClient
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authMiddleware *auth.AuthMiddleware, userClient *clients.UserClient) *AuthHandler {
	if authMiddleware == nil {
		panic("authMiddleware cannot be nil")
	}
	if userClient == nil {
		panic("userClient cannot be nil")
	}
	return &AuthHandler{
		authMiddleware: authMiddleware,
		userClient:     userClient,
	}
}

// AuthStatusInput represents the request to check auth status
type AuthStatusInput struct{}

// AuthStatusOutput represents the response for auth status check
type AuthStatusOutput struct {
	Body struct {
		Authenticated bool                   `json:"authenticated" doc:"Whether the user is authenticated"`
		User          *AuthUserInfo          `json:"user,omitempty" doc:"User information if authenticated"`
		Permissions   []string               `json:"permissions,omitempty" doc:"User permissions"`
		TokenInfo     *TokenInfo             `json:"token_info,omitempty" doc:"Token information"`
		Timestamp     time.Time              `json:"timestamp" doc:"Response timestamp"`
	}
}

// AuthUserInfo represents authenticated user information
type AuthUserInfo struct {
	ID             string                 `json:"id" doc:"User ID"`
	Email          string                 `json:"email" doc:"User email"`
	DisplayName    string                 `json:"display_name" doc:"User display name"`
	AvatarUrl      string                 `json:"avatar_url,omitempty" doc:"User avatar URL"`
	AuthProvider   string                 `json:"auth_provider" doc:"Authentication provider"`
	Verified       bool                   `json:"verified" doc:"Whether email is verified"`
	Preferences    map[string]interface{} `json:"preferences,omitempty" doc:"User preferences"`
}

// TokenInfo represents JWT token information
type TokenInfo struct {
	Subject   string    `json:"subject" doc:"Token subject (user ID)"`
	Issuer    string    `json:"issuer" doc:"Token issuer"`
	Audience  []string  `json:"audience" doc:"Token audience"`
	IssuedAt  time.Time `json:"issued_at" doc:"Token issued at"`
	ExpiresAt time.Time `json:"expires_at" doc:"Token expires at"`
	Scopes    []string  `json:"scopes,omitempty" doc:"Token scopes"`
}

// ValidateTokenInput represents the request to validate a token
type ValidateTokenInput struct {
	Body struct {
		Token string `json:"token" minLength:"1" doc:"JWT token to validate"`
	}
}

// ValidateTokenOutput represents the response for token validation
type ValidateTokenOutput struct {
	Body struct {
		Valid       bool       `json:"valid" doc:"Whether the token is valid"`
		Claims      *TokenInfo `json:"claims,omitempty" doc:"Token claims if valid"`
		Error       string     `json:"error,omitempty" doc:"Error message if invalid"`
		Timestamp   time.Time  `json:"timestamp" doc:"Response timestamp"`
	}
}

// LogoutInput represents the request to logout
type LogoutInput struct{}

// LogoutOutput represents the response for logout
type LogoutOutput struct {
	Body struct {
		Success   bool      `json:"success" doc:"Whether logout was successful"`
		Message   string    `json:"message" doc:"Logout message"`
		Timestamp time.Time `json:"timestamp" doc:"Response timestamp"`
	}
}

// AuthStatsInput represents the request to get auth statistics
type AuthStatsInput struct{}

// AuthStatsOutput represents the response for auth statistics
type AuthStatsOutput struct {
	Body struct {
		ServiceStats  map[string]interface{} `json:"service_stats" doc:"Authentication service statistics"`
		HealthStatus  string                 `json:"health_status" doc:"Authentication service health status"`
		Timestamp     time.Time              `json:"timestamp" doc:"Response timestamp"`
	}
}

// RegisterRoutes registers authentication routes (public endpoints)
func (h *AuthHandler) RegisterRoutes(api huma.API) {
	// Validate token endpoint (public, for client-side validation)
	huma.Register(api, huma.Operation{
		OperationID: "validate-token",
		Method:      http.MethodPost,
		Path:        "/api/v1/auth/validate",
		Summary:     "Validate JWT token",
		Description: "Validate a JWT token and return claims information",
		Tags:        []string{"Authentication"},
	}, h.ValidateToken)

	// Authentication status endpoint (can be called without auth)
	huma.Register(api, huma.Operation{
		OperationID: "auth-status",
		Method:      http.MethodGet,
		Path:        "/api/v1/auth/status",
		Summary:     "Get authentication status",
		Description: "Check if the current request is authenticated and get user info",
		Tags:        []string{"Authentication"},
	}, h.GetAuthStatus)
}

// RegisterRoutesWithRateLimit registers authentication routes with rate limiting
func (h *AuthHandler) RegisterRoutesWithRateLimit(api huma.API, rateLimiter *auth.RateLimiter) {
	// Register public routes first
	h.RegisterRoutes(api)

	// Logout endpoint (requires authentication)
	huma.Register(api, huma.Operation{
		OperationID: "logout",
		Method:      http.MethodPost,
		Path:        "/api/v1/auth/logout",
		Summary:     "Logout user",
		Description: "Logout the current user and invalidate tokens",
		Tags:        []string{"Authentication"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, h.Logout)

	// Auth statistics endpoint (admin only)
	huma.Register(api, huma.Operation{
		OperationID: "auth-stats",
		Method:      http.MethodGet,
		Path:        "/api/v1/auth/stats",
		Summary:     "Get authentication statistics",
		Description: "Get authentication service statistics (admin only)",
		Tags:        []string{"Authentication"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, h.GetAuthStats)
}

// GetAuthStatus checks the authentication status of the current request
func (h *AuthHandler) GetAuthStatus(ctx context.Context, input *AuthStatusInput) (*AuthStatusOutput, error) {
	resp := &AuthStatusOutput{}
	resp.Body.Timestamp = time.Now()

	// Try to get user from context (might not be authenticated)
	userInfo, hasUser := auth.GetUserFromContext(ctx)
	if !hasUser {
		// Not authenticated
		resp.Body.Authenticated = false
		return resp, nil
	}

	// User is authenticated
	resp.Body.Authenticated = true

	// Extract user information
	userID, _ := auth.GetUserIDFromContext(ctx)
	email, _ := auth.GetUserEmailFromContext(ctx)

	// Get detailed user information from database
	authUserInfo := &AuthUserInfo{
		ID:    userID,
		Email: email,
	}

	// Extract additional info from context if available
	if displayName, ok := userInfo["name"].(string); ok {
		authUserInfo.DisplayName = displayName
	}
	if picture, ok := userInfo["picture"].(string); ok {
		authUserInfo.AvatarUrl = picture
	}
	if verified, ok := userInfo["email_verified"].(bool); ok {
		authUserInfo.Verified = verified
	}

	// Extract token information from context (simplified approach)
	// Since claims aren't directly stored in userInfo, we'll create basic token info
	if userID != "" {
		tokenInfo := &TokenInfo{
			Subject:   userID,
			Issuer:    "auth0", // Placeholder
			Audience:  []string{"bocchi-the-map-api"},
			IssuedAt:  time.Now().Add(-time.Hour), // Placeholder
			ExpiresAt: time.Now().Add(time.Hour),   // Placeholder
		}
		
		// Extract scopes if available
		if scopes, ok := userInfo["scope"].(string); ok && scopes != "" {
			// Parse space-separated scopes
			tokenInfo.Scopes = []string{scopes} // Simplified for now
		}
		
		resp.Body.TokenInfo = tokenInfo
	}

	// Extract permissions (placeholder for future implementation)
	resp.Body.Permissions = []string{} // TODO: Implement role-based permissions

	resp.Body.User = authUserInfo

	// Add monitoring context
	monitoring.AddUserContext(ctx, userID, email)

	return resp, nil
}

// ValidateToken validates a JWT token
func (h *AuthHandler) ValidateToken(ctx context.Context, input *ValidateTokenInput) (*ValidateTokenOutput, error) {
	resp := &ValidateTokenOutput{}
	resp.Body.Timestamp = time.Now()

	// Use the middleware's validator directly
	// TODO: Refactor to use the auth service properly
	validator := h.authMiddleware.GetValidator()
	if validator == nil {
		resp.Body.Valid = false
		resp.Body.Error = "authentication service not available"
		return resp, nil
	}

	claims, err := validator.ValidateToken(input.Body.Token)
	if err != nil {
		resp.Body.Valid = false
		resp.Body.Error = err.Error()
		
		// Log validation failure for monitoring
		logger.InfoWithFields("Token validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		
		return resp, nil
	}

	// Token is valid
	resp.Body.Valid = true
	resp.Body.Claims = &TokenInfo{
		Subject:   claims.Subject,
		Issuer:    claims.Issuer,
		Audience:  claims.Audience,
		IssuedAt:  time.Unix(claims.IssuedAt, 0),
		ExpiresAt: time.Unix(claims.ExpiresAt, 0),
	}

	// Extract scopes if available
	if claims.Scope != "" {
		resp.Body.Claims.Scopes = []string{claims.Scope} // Simplified for now
	}

	logger.InfoWithFields("Token validation successful", map[string]interface{}{
		"subject": claims.Subject,
		"issuer":  claims.Issuer,
	})

	return resp, nil
}

// Logout logs out the current user
func (h *AuthHandler) Logout(ctx context.Context, input *LogoutInput) (*LogoutOutput, error) {
	resp := &LogoutOutput{}
	resp.Body.Timestamp = time.Now()

	// Extract user ID from authentication context
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		return nil, huma.Error401Unauthorized("authentication required")
	}

	// Extract JWT ID (JTI) from context
	jti, hasJTI := auth.GetJTIFromContext(ctx)
	if !hasJTI || jti == "" {
		logger.InfoWithFields("User logout without JTI", map[string]interface{}{
			"user_id": userID,
			"note":    "Token does not contain JTI - proceeding with client-side logout only",
		})
		resp.Body.Success = true
		resp.Body.Message = "Logout successful. Please remove the token from client storage."
		return resp, nil
	}

	// Extract token expiration from context
	expiresAt, hasExpiration := auth.GetTokenExpirationFromContext(ctx)
	if !hasExpiration {
		expiresAt = time.Now().Add(24 * time.Hour) // Default fallback
	}

	// Blacklist the token
	err := h.authMiddleware.Logout(ctx, jti, expiresAt)
	if err != nil {
		logger.ErrorWithFields("Failed to blacklist token during logout", err, map[string]interface{}{
			"user_id": userID,
			"jti":     jti,
		})
		// Don't fail logout completely, but inform the client
		resp.Body.Success = true
		resp.Body.Message = "Logout processed. Please remove the token from client storage. Note: token may still be valid for a short time."
		return resp, nil
	}

	logger.InfoWithFields("User logout successful", map[string]interface{}{
		"user_id": userID,
		"jti":     jti,
	})

	resp.Body.Success = true
	resp.Body.Message = "Logout successful. Token has been invalidated."

	// Add monitoring context
	monitoring.AddUserContext(ctx, userID, "")
	
	return resp, nil
}

// GetAuthStats returns authentication service statistics
func (h *AuthHandler) GetAuthStats(ctx context.Context, input *AuthStatsInput) (*AuthStatsOutput, error) {
	// Extract user ID from authentication context
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		return nil, huma.Error401Unauthorized("authentication required")
	}

	// TODO: Check if user has admin permissions
	// For now, we'll allow any authenticated user to see basic stats

	resp := &AuthStatsOutput{}
	resp.Body.Timestamp = time.Now()

	// TODO: Get stats from the auth service
	// For now, return placeholder data
	resp.Body.ServiceStats = map[string]interface{}{
		"rate_limiter": map[string]interface{}{
			"active": true,
			"stats":  "Rate limiter statistics would go here",
		},
		"jwt_validator": map[string]interface{}{
			"active": true,
			"stats":  "JWT validator statistics would go here",
		},
	}

	resp.Body.HealthStatus = "healthy"

	logger.InfoWithFields("Auth stats requested", map[string]interface{}{
		"user_id": userID,
	})

	return resp, nil
}

// Helper method to get validator from middleware (temporary)
func (h *AuthHandler) getValidator() *auth.JWTValidator {
	// This is a temporary workaround
	// In a proper implementation, we should inject the auth service
	return nil // TODO: Fix this when refactoring
}