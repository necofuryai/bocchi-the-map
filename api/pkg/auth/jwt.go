package auth

import (
	"context"
	"crypto/rsa"
	"net/http"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt/v5"

	"github.com/necofuryai/bocchi-the-map/api/pkg/errors"
	"github.com/necofuryai/bocchi-the-map/api/pkg/logger"
)

// Claims represents the JWT claims structure for Auth0 tokens
type Claims struct {
	Audience  []string `json:"aud,omitempty"`
	ExpiresAt int64    `json:"exp,omitempty"`
	ID        string   `json:"jti,omitempty"`
	IssuedAt  int64    `json:"iat,omitempty"`
	Issuer    string   `json:"iss,omitempty"`
	NotBefore int64    `json:"nbf,omitempty"`
	Subject   string   `json:"sub,omitempty"`
	
	// Auth0 specific claims
	Email         string   `json:"email,omitempty"`
	EmailVerified bool     `json:"email_verified,omitempty"`
	Name          string   `json:"name,omitempty"`
	Nickname      string   `json:"nickname,omitempty"`
	Picture       string   `json:"picture,omitempty"`
	Scope         string   `json:"scope,omitempty"`
	Permissions   []string `json:"permissions,omitempty"`
	
	jwt.RegisteredClaims
}

// JWTValidator handles Auth0 JWT token validation
type JWTValidator struct {
	provider   *oidc.Provider
	verifier   *oidc.IDTokenVerifier
	audience   string
	domain     string
	publicKeys map[string]*rsa.PublicKey
}

// NewJWTValidator creates a new JWT validator for Auth0 tokens
func NewJWTValidator(domain, audience string) (*JWTValidator, error) {
	if domain == "" {
		return nil, errors.InvalidInput("domain", "Auth0 domain is required")
	}
	if audience == "" {
		return nil, errors.InvalidInput("audience", "Auth0 audience is required")
	}

	// Ensure domain has https:// prefix
	issuerURL := domain
	if !strings.HasPrefix(domain, "https://") {
		issuerURL = "https://" + domain + "/"
	} else if !strings.HasSuffix(domain, "/") {
		issuerURL = domain + "/"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Initialize OIDC provider
	provider, err := oidc.NewProvider(ctx, issuerURL)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrTypeExternalAPI, "failed to initialize Auth0 OIDC provider")
	}

	// Configure the ID token verifier
	oidcConfig := &oidc.Config{
		ClientID: audience,
	}
	verifier := provider.Verifier(oidcConfig)

	validator := &JWTValidator{
		provider:   provider,
		verifier:   verifier,
		audience:   audience,
		domain:     domain,
		publicKeys: make(map[string]*rsa.PublicKey),
	}

	logger.InfoWithFields("Auth0 JWT validator initialized successfully", map[string]interface{}{
		"domain":   domain,
		"audience": audience,
		"issuer":   issuerURL,
	})

	return validator, nil
}

// ValidateToken validates an Auth0 JWT token and returns the claims
func (v *JWTValidator) ValidateToken(tokenString string) (*Claims, error) {
	if tokenString == "" {
		return nil, errors.Unauthorized("token is required")
	}

	// Remove "Bearer " prefix if present
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	tokenString = strings.TrimSpace(tokenString)

	if tokenString == "" {
		return nil, errors.Unauthorized("token is empty after processing")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Parse and verify the token
	token, err := v.verifier.Verify(ctx, tokenString)
	if err != nil {
		logger.InfoWithFields("Token verification failed", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, errors.Wrap(err, errors.ErrTypeUnauthorized, "invalid or expired token")
	}

	// Extract claims
	var claims Claims
	if err := token.Claims(&claims); err != nil {
		logger.Error("Failed to extract claims from token", err)
		return nil, errors.Wrap(err, errors.ErrTypeInternal, "failed to extract token claims")
	}

	// Validate audience
	validAudience := false
	for _, aud := range claims.Audience {
		if aud == v.audience {
			validAudience = true
			break
		}
	}
	
	if !validAudience {
		logger.InfoWithFields("Token audience validation failed", map[string]interface{}{
			"expected":  v.audience,
			"audiences": claims.Audience,
		})
		return nil, errors.Unauthorized("token audience is invalid")
	}

	// Validate issuer
	expectedIssuer := v.domain
	if !strings.HasPrefix(v.domain, "https://") {
		expectedIssuer = "https://" + v.domain + "/"
	} else if !strings.HasSuffix(v.domain, "/") {
		expectedIssuer = v.domain + "/"
	}

	if claims.Issuer != expectedIssuer {
		logger.InfoWithFields("Token issuer validation failed", map[string]interface{}{
			"expected": expectedIssuer,
			"actual":   claims.Issuer,
		})
		return nil, errors.Unauthorized("token issuer is invalid")
	}

	// Validate expiration
	if claims.ExpiresAt > 0 && time.Unix(claims.ExpiresAt, 0).Before(time.Now()) {
		return nil, errors.Unauthorized("token has expired")
	}

	// Validate not before
	if claims.NotBefore > 0 && time.Unix(claims.NotBefore, 0).After(time.Now()) {
		return nil, errors.Unauthorized("token is not yet valid")
	}

	logger.InfoWithFields("Token validated successfully", map[string]interface{}{
		"subject": claims.Subject,
		"email":   claims.Email,
	})

	return &claims, nil
}

// ExtractTokenFromRequest extracts the JWT token from HTTP request headers
func (v *JWTValidator) ExtractTokenFromRequest(r *http.Request) (string, error) {
	// Check Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			return parts[1], nil
		}
		return "", errors.Unauthorized("authorization header format must be 'Bearer {token}'")
	}

	// Check query parameter as fallback
	token := r.URL.Query().Get("token")
	if token != "" {
		return token, nil
	}

	return "", errors.Unauthorized("no authorization token found")
}

// ValidateTokenFromRequest validates a JWT token from HTTP request
func (v *JWTValidator) ValidateTokenFromRequest(r *http.Request) (*Claims, error) {
	token, err := v.ExtractTokenFromRequest(r)
	if err != nil {
		return nil, err
	}

	return v.ValidateToken(token)
}

// GetUserContext creates a context with user information from claims
func (v *JWTValidator) GetUserContext(ctx context.Context, claims *Claims) context.Context {
	userInfo := map[string]interface{}{
		"user_id":        claims.Subject,
		"email":          claims.Email,
		"email_verified": claims.EmailVerified,
		"name":           claims.Name,
		"nickname":       claims.Nickname,
		"picture":        claims.Picture,
		"permissions":    claims.Permissions,
	}

	// Set both "user" object and individual context keys for compatibility
	ctx = context.WithValue(ctx, "user", userInfo)
	ctx = context.WithValue(ctx, "user_id", claims.Subject)
	ctx = context.WithValue(ctx, "email", claims.Email)
	
	return ctx
}

// GetUserFromContext extracts user information from context
func GetUserFromContext(ctx context.Context) (map[string]interface{}, bool) {
	user, ok := ctx.Value("user").(map[string]interface{})
	return user, ok
}

// GetUserIDFromContext extracts user ID from context
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	user, ok := GetUserFromContext(ctx)
	if !ok {
		return "", false
	}
	
	userID, ok := user["user_id"].(string)
	return userID, ok
}

// GetUserEmailFromContext extracts user email from context
func GetUserEmailFromContext(ctx context.Context) (string, bool) {
	user, ok := GetUserFromContext(ctx)
	if !ok {
		return "", false
	}
	
	email, ok := user["email"].(string)
	return email, ok
}

// HasPermission checks if the user has a specific permission
func HasPermission(ctx context.Context, permission string) bool {
	user, ok := GetUserFromContext(ctx)
	if !ok {
		return false
	}
	
	permissions, ok := user["permissions"].([]string)
	if !ok {
		return false
	}
	
	for _, p := range permissions {
		if p == permission {
			return true
		}
	}
	
	return false
}