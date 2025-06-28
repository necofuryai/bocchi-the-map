package handlers

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/necofuryai/bocchi-the-map/api/application/clients"
	"github.com/necofuryai/bocchi-the-map/api/pkg/auth"
	"github.com/necofuryai/bocchi-the-map/api/pkg/errors"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	authMiddleware *auth.AuthMiddleware
	userClient     *clients.UserClient
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(authMiddleware *auth.AuthMiddleware, userClient *clients.UserClient) *AuthHandler {
	return &AuthHandler{
		authMiddleware: authMiddleware,
		userClient:     userClient,
	}
}

// GenerateTokenInput represents the request to generate an API token
type GenerateTokenInput struct {
	Body struct {
		Email          string `json:"email" maxLength:"255" doc:"User email address"`
		AuthProvider   string `json:"provider" enum:"google,twitter,x" doc:"OAuth provider (google, twitter, or x)"`
		AuthProviderID string `json:"provider_id" doc:"Provider-specific user ID"`
		SessionToken   string `json:"session_token,omitempty" doc:"Auth.js session token for verification"`
	}
}

// GenerateTokenOutput represents the response for token generation
type GenerateTokenOutput struct {
	SetCookie []string `header:"Set-Cookie" doc:"Authentication cookies"`
	Body struct {
		Message   string    `json:"message" doc:"Success message"`
		ExpiresIn int       `json:"expires_in" example:"86400" doc:"Token expiration time in seconds"`
		ExpiresAt time.Time `json:"expires_at" doc:"Token expiration timestamp"`
	}
}

// RefreshTokenInput represents the request to refresh a token
type RefreshTokenInput struct {
	RefreshToken string `cookie:"bocchi_refresh_token" doc:"Refresh token from cookie"`
	Body struct {
		RefreshToken string `json:"refresh_token,omitempty" doc:"Refresh token (optional if cookie is present)"`
	}
}

// RefreshTokenOutput represents the response for token refresh
type RefreshTokenOutput struct {
	SetCookie []string `header:"Set-Cookie" doc:"Authentication cookies"`
	Body struct {
		Message   string    `json:"message" doc:"Success message"`
		ExpiresIn int       `json:"expires_in" example:"86400" doc:"Token expiration time in seconds"`
		ExpiresAt time.Time `json:"expires_at" doc:"Token expiration timestamp"`
	}
}

// LogoutInput represents the request to logout
type LogoutInput struct {
	AccessToken  string `cookie:"bocchi_access_token" doc:"Access token from cookie"`
	RefreshToken string `cookie:"bocchi_refresh_token" doc:"Refresh token from cookie"`
}

// LogoutOutput represents the response for logout
type LogoutOutput struct {
	SetCookie []string `header:"Set-Cookie" doc:"Clear authentication cookies"`
	Body struct {
		Message string `json:"message" doc:"Logout success message"`
	}
}

// RegisterRoutes registers authentication routes
func (h *AuthHandler) RegisterRoutes(api huma.API) {
	// Token generation endpoint (public - called after Auth.js authentication)
	huma.Register(api, huma.Operation{
		OperationID: "generate-api-token",
		Method:      http.MethodPost,
		Path:        "/api/v1/auth/token",
		Summary:     "Generate API access token",
		Description: "Generate JWT access token for API calls using Auth.js session information",
		Tags:        []string{"Authentication"},
	}, h.GenerateToken)

	// Token refresh endpoint (public - uses refresh token)
	huma.Register(api, huma.Operation{
		OperationID: "refresh-api-token",
		Method:      http.MethodPost,
		Path:        "/api/v1/auth/refresh",
		Summary:     "Refresh API access token",
		Description: "Refresh JWT access token using refresh token",
		Tags:        []string{"Authentication"},
	}, h.RefreshToken)

	// Logout endpoint (public - clears cookies)
	huma.Register(api, huma.Operation{
		OperationID: "logout",
		Method:      http.MethodPost,
		Path:        "/api/v1/auth/logout",
		Summary:     "Logout user",
		Description: "Clear authentication cookies and logout user",
		Tags:        []string{"Authentication"},
	}, h.Logout)
}

// CreateHumaRateLimitMiddleware creates a reusable Huma-compatible rate limit middleware
func CreateHumaRateLimitMiddleware(rateLimiter *auth.RateLimiter) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		// Get client IP from request
		clientIP := ctx.RemoteAddr()
		if xff := ctx.Header("X-Forwarded-For"); xff != "" {
			clientIP = xff
		} else if xri := ctx.Header("X-Real-IP"); xri != "" {
			clientIP = xri
		}
		
		// Check rate limit
		if !rateLimiter.Allow(clientIP) {
			// Set rate limit headers
			ctx.SetHeader("X-RateLimit-Limit", strconv.Itoa(rateLimiter.GetLimit()))
			ctx.SetHeader("X-RateLimit-Window", strconv.Itoa(rateLimiter.GetWindow()))
			ctx.SetHeader("Retry-After", strconv.Itoa(rateLimiter.GetWindow()))
			
			// Return rate limit error
			ctx.SetStatus(http.StatusTooManyRequests)
			return
		}
		
		// Continue to next middleware/handler
		next(ctx)
	}
}

// RegisterRoutesWithRateLimit registers authentication routes with rate limiting protection
func (h *AuthHandler) RegisterRoutesWithRateLimit(api huma.API, rateLimiter *auth.RateLimiter) {
	// Create Huma-compatible rate limit middleware
	rateLimitHumaMiddleware := CreateHumaRateLimitMiddleware(rateLimiter)

	// Token generation endpoint with rate limiting
	huma.Register(api, huma.Operation{
		OperationID: "generate-api-token",
		Method:      http.MethodPost,
		Path:        "/api/v1/auth/token",
		Summary:     "Generate API access token",
		Description: "Generate JWT access token for API calls using Auth.js session information",
		Tags:        []string{"Authentication"},
		Middlewares: huma.Middlewares{rateLimitHumaMiddleware},
	}, h.GenerateToken)

	// Token refresh endpoint with rate limiting
	huma.Register(api, huma.Operation{
		OperationID: "refresh-api-token",
		Method:      http.MethodPost,
		Path:        "/api/v1/auth/refresh",
		Summary:     "Refresh API access token", 
		Description: "Refresh JWT access token using refresh token",
		Tags:        []string{"Authentication"},
		Middlewares: huma.Middlewares{rateLimitHumaMiddleware},
	}, h.RefreshToken)

	// Logout endpoint (no rate limiting needed for logout)
	huma.Register(api, huma.Operation{
		OperationID: "logout",
		Method:      http.MethodPost,
		Path:        "/api/v1/auth/logout",
		Summary:     "Logout user",
		Description: "Clear authentication cookies and logout user",
		Tags:        []string{"Authentication"},
	}, h.Logout)
}

// GenerateToken generates JWT tokens for API access after Auth.js authentication
func (h *AuthHandler) GenerateToken(ctx context.Context, input *GenerateTokenInput) (*GenerateTokenOutput, error) {
	// Add operation context for error tracking
	ctx = errors.WithOperation(ctx, "generate_api_token")

	// Validate required fields
	if input.Body.Email == "" {
		return nil, huma.Error400BadRequest("email is required")
	}
	if input.Body.AuthProvider == "" {
		return nil, huma.Error400BadRequest("auth provider is required")
	}
	if input.Body.AuthProviderID == "" {
		return nil, huma.Error400BadRequest("provider ID is required")
	}

	// Rate limiting protection - prevent token generation abuse

	// Convert provider string to domain enum for user lookup
	userConverter := h.userClient.GetConverter()
	authProvider, err := userConverter.ConvertHTTPAuthProviderToEntity(input.Body.AuthProvider)
	if err != nil {
		return nil, errors.HandleHTTPError(ctx, err, "convert_auth_provider", "invalid auth provider")
	}

	// Get user by auth provider to ensure they exist
	user, err := h.userClient.GetUserByAuthProvider(ctx, authProvider, input.Body.AuthProviderID)
	if err != nil {
		if errors.Is(err, errors.ErrTypeNotFound) {
			// Enhanced error response for better client handling
			return nil, huma.Error404NotFound("user not found - please complete OAuth authentication first")
		}
		// Log error details for monitoring but return generic message
		return nil, errors.HandleHTTPError(ctx, err, "get_user", "authentication failed")
	}

	// Verify email matches (security check to prevent account hijacking)
	if user.Email != input.Body.Email {
		// Log security violation for monitoring
		return nil, huma.Error403Forbidden("authentication failed - invalid credentials")
	}

	// Generate access token (24 hours)
	accessToken, err := h.authMiddleware.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, errors.HandleHTTPError(ctx, err, "generate_access_token", "failed to generate access token")
	}

	// Generate refresh token (7 days)
	refreshToken, err := h.authMiddleware.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		return nil, errors.HandleHTTPError(ctx, err, "generate_refresh_token", "failed to generate refresh token")
	}

	// Calculate expiration time
	expiresAt := time.Now().Add(24 * time.Hour)
	expiresIn := int(24 * time.Hour / time.Second)

	// Create secure httpOnly cookies
	cookies := h.createSecureCookies(accessToken, refreshToken, expiresAt)

	// Return success response without tokens
	resp := &GenerateTokenOutput{}
	resp.SetCookie = cookies
	resp.Body.Message = "Authentication successful"
	resp.Body.ExpiresIn = expiresIn
	resp.Body.ExpiresAt = expiresAt

	return resp, nil
}

// RefreshToken refreshes JWT access tokens using refresh token
func (h *AuthHandler) RefreshToken(ctx context.Context, input *RefreshTokenInput) (*RefreshTokenOutput, error) {
	// Add operation context for error tracking
	ctx = errors.WithOperation(ctx, "refresh_api_token")

	// Get refresh token from cookie (via Huma input)
	refreshToken := input.RefreshToken
	if refreshToken == "" && input.Body.RefreshToken != "" {
		refreshToken = input.Body.RefreshToken
	}
	
	if refreshToken == "" {
		return nil, huma.Error401Unauthorized("refresh token not found")
	}

	// Validate and parse refresh token
	claims, err := h.authMiddleware.ValidateToken(refreshToken)
	if err != nil {
		return nil, huma.Error401Unauthorized("invalid refresh token")
	}

	// Get user to ensure they still exist
	user, err := h.userClient.GetUserByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, errors.ErrTypeNotFound) {
			return nil, huma.Error404NotFound("user not found")
		}
		return nil, errors.HandleHTTPError(ctx, err, "get_user", "failed to get user")
	}

	// Generate new access token (24 hours)
	accessToken, err := h.authMiddleware.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, errors.HandleHTTPError(ctx, err, "generate_access_token", "failed to generate access token")
	}

	// Generate new refresh token (7 days)
	newRefreshToken, err := h.authMiddleware.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		return nil, errors.HandleHTTPError(ctx, err, "generate_refresh_token", "failed to generate refresh token")
	}

	// Calculate expiration time
	expiresAt := time.Now().Add(24 * time.Hour)
	expiresIn := int(24 * time.Hour / time.Second)

	// Create secure httpOnly cookies
	cookies := h.createSecureCookies(accessToken, newRefreshToken, expiresAt)

	// Return success response without tokens
	resp := &RefreshTokenOutput{}
	resp.SetCookie = cookies
	resp.Body.Message = "Token refreshed successfully"
	resp.Body.ExpiresIn = expiresIn
	resp.Body.ExpiresAt = expiresAt

	return resp, nil
}

// Logout clears authentication cookies and blacklists tokens
func (h *AuthHandler) Logout(ctx context.Context, input *LogoutInput) (*LogoutOutput, error) {
	// Add operation context for error tracking
	ctx = errors.WithOperation(ctx, "logout")

	// Try to blacklist current tokens before clearing cookies
	// Blacklist access token if present
	if input.AccessToken != "" {
		if err := h.authMiddleware.BlacklistToken(ctx, input.AccessToken, "logout"); err != nil {
			errors.LogError(ctx, err, "blacklist_access_token_on_logout")
		}
	}
	
	// Blacklist refresh token if present  
	if input.RefreshToken != "" {
		if err := h.authMiddleware.BlacklistToken(ctx, input.RefreshToken, "logout"); err != nil {
			errors.LogError(ctx, err, "blacklist_refresh_token_on_logout")
		}
	}

	// Clear authentication cookies
	clearCookies := h.createClearCookies()

	// Return success response
	resp := &LogoutOutput{}
	resp.SetCookie = clearCookies
	resp.Body.Message = "Logout successful"

	return resp, nil
}

// createSecureCookies creates httpOnly cookies with security settings
func (h *AuthHandler) createSecureCookies(accessToken, refreshToken string, expiresAt time.Time) []string {
	// Determine if we're in production (use Secure flag)
	isProduction := os.Getenv("ENVIRONMENT") == "production"
	
	// Get domain from environment or use default
	domain := os.Getenv("COOKIE_DOMAIN")
	if domain == "" {
		domain = "localhost"
	}

	// Create access token cookie (shorter expiration)
	accessCookie := &http.Cookie{
		Name:     "bocchi_access_token",
		Value:    accessToken,
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   isProduction,
		SameSite: http.SameSiteStrictMode,
		Domain:   domain,
		Path:     "/",
	}

	// Create refresh token cookie (longer expiration)
	refreshCookie := &http.Cookie{
		Name:     "bocchi_refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour), // 7 days
		HttpOnly: true,
		Secure:   isProduction,
		SameSite: http.SameSiteStrictMode,
		Domain:   domain,
		Path:     "/",
	}

	return []string{
		accessCookie.String(),
		refreshCookie.String(),
	}
}

// createClearCookies creates cookies that clear authentication cookies
func (h *AuthHandler) createClearCookies() []string {
	// Get domain from environment or use default
	domain := os.Getenv("COOKIE_DOMAIN")
	if domain == "" {
		domain = "localhost"
	}

	// Clear access token cookie
	accessCookie := &http.Cookie{
		Name:     "bocchi_access_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   os.Getenv("ENVIRONMENT") == "production",
		SameSite: http.SameSiteStrictMode,
		Domain:   domain,
		Path:     "/",
		MaxAge:   -1,
	}

	// Clear refresh token cookie
	refreshCookie := &http.Cookie{
		Name:     "bocchi_refresh_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   os.Getenv("ENVIRONMENT") == "production",
		SameSite: http.SameSiteStrictMode,
		Domain:   domain,
		Path:     "/",
		MaxAge:   -1,
	}

	return []string{
		accessCookie.String(),
		refreshCookie.String(),
	}
}

