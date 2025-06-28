package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/necofuryai/bocchi-the-map/api/application/clients"
	"github.com/necofuryai/bocchi-the-map/api/domain/entities"
	"github.com/necofuryai/bocchi-the-map/api/pkg/auth"
	"github.com/necofuryai/bocchi-the-map/api/pkg/errors"
	"github.com/necofuryai/bocchi-the-map/api/pkg/converters"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userClient    *clients.UserClient
	userConverter *converters.UserConverter
}

// NewUserHandler creates a new user handler
func NewUserHandler(userClient *clients.UserClient) *UserHandler {
	return &UserHandler{
		userClient:    userClient,
		userConverter: converters.NewUserConverter(),
	}
}

// CreateUserInput represents the OAuth user creation/update request (Auth.js compatible)
type CreateUserInput struct {
	Body struct {
		Email          string `json:"email" maxLength:"255" doc:"User email address"`
		DisplayName    string `json:"name" maxLength:"100" doc:"User display name"`
		AvatarURL      string `json:"image,omitempty" doc:"User avatar URL"`
		AuthProvider   string `json:"provider" enum:"google,twitter,x" doc:"OAuth provider (google, twitter, or x)"`
		AuthProviderID string `json:"provider_id" doc:"Provider-specific user ID"`
	}
}

// CreateUserOutput represents the response for user creation
type CreateUserOutput struct {
	Body struct {
		ID          string    `json:"id" doc:"User ID"`
		Email       string    `json:"email" doc:"User email"`
		DisplayName string    `json:"name" doc:"User display name"`
		AvatarURL   string    `json:"image,omitempty" doc:"User avatar URL"`
		CreatedAt   time.Time `json:"created_at" doc:"Creation timestamp"`
		UpdatedAt   time.Time `json:"updated_at" doc:"Last update timestamp"`
	}
}

// GetCurrentUserInput represents the request to get current user
type GetCurrentUserInput struct {
	// Empty - user context comes from authentication
}

// GetCurrentUserOutput represents the response for getting current user
type GetCurrentUserOutput struct {
	Body struct {
		ID          string                 `json:"id" doc:"User ID"`
		Email       string                 `json:"email" doc:"User email"`
		DisplayName string                 `json:"display_name" doc:"User display name"`
		AvatarURL   string                 `json:"avatar_url,omitempty" doc:"User avatar URL"`
		Preferences map[string]interface{} `json:"preferences,omitempty" doc:"User preferences"`
		CreatedAt   time.Time              `json:"created_at" doc:"Creation timestamp"`
		UpdatedAt   time.Time              `json:"updated_at" doc:"Last update timestamp"`
	}
}

// UpdatePreferencesInput represents the request to update user preferences
type UpdatePreferencesInput struct {
	Body struct {
		Preferences map[string]interface{} `json:"preferences" doc:"User preferences to update"`
	}
}

// UpdatePreferencesOutput represents the response for updating preferences
type UpdatePreferencesOutput struct {
	Body struct {
		ID          string                 `json:"id" doc:"User ID"`
		Email       string                 `json:"email" doc:"User email"`
		DisplayName string                 `json:"display_name" doc:"User display name"`
		AvatarURL   string                 `json:"avatar_url,omitempty" doc:"User avatar URL"`
		Preferences map[string]interface{} `json:"preferences,omitempty" doc:"Updated user preferences"`
		UpdatedAt   time.Time              `json:"updated_at" doc:"Last update timestamp"`
	}
}

// RegisterRoutes registers user routes (both standard API and Auth.js compatible)
func (h *UserHandler) RegisterRoutes(api huma.API) {
	// Standard REST API routes
	huma.Register(api, huma.Operation{
		OperationID: "create-user",
		Method:      http.MethodPost,
		Path:        "/api/v1/users",
		Summary:     "Create or update a user",
		Description: "Create a new user or update existing user information",
		Tags:        []string{"Users"},
	}, h.CreateUser)

	huma.Register(api, huma.Operation{
		OperationID: "get-current-user",
		Method:      http.MethodGet,
		Path:        "/api/v1/users/me",
		Summary:     "Get current user",
		Description: "Get the current authenticated user's information",
		Tags:        []string{"Users"},
	}, h.GetCurrentUser)

	huma.Register(api, huma.Operation{
		OperationID: "update-user-preferences",
		Method:      http.MethodPatch,
		Path:        "/api/v1/users/me/preferences",
		Summary:     "Update user preferences",
		Description: "Update the current user's preferences",
		Tags:        []string{"Users"},
	}, h.UpdatePreferences)
}

// CreateUser creates or updates a user (upsert for OAuth flow)
func (h *UserHandler) CreateUser(ctx context.Context, input *CreateUserInput) (*CreateUserOutput, error) {
	// Add operation context for error tracking
	ctx = errors.WithOperation(ctx, "create_user_http")

	// Convert provider string to domain enum using standardized converter
	authProvider, err := h.userConverter.ConvertHTTPAuthProviderToEntity(input.Body.AuthProvider)
	if err != nil {
		return nil, errors.HandleHTTPError(ctx, err, "convert_auth_provider", "invalid auth provider")
	}

	// Check if user already exists
	existingUser, err := h.userClient.GetUserByAuthProvider(ctx, authProvider, input.Body.AuthProviderID)
	if err != nil && !errors.Is(err, errors.ErrTypeNotFound) {
		return nil, errors.HandleHTTPError(ctx, err, "check_existing_user", "failed to check existing user")
	}

	var user *entities.User
	if existingUser != nil {
		// Update existing user
		existingUser.Email = input.Body.Email
		existingUser.DisplayName = input.Body.DisplayName
		existingUser.AvatarURL = input.Body.AvatarURL
		user, err = h.userClient.UpdateUser(ctx, existingUser)
		if err != nil {
			return nil, errors.HandleHTTPError(ctx, err, "update_user", "failed to update user")
		}
	} else {
		// Create new user
		newUser := &entities.User{
			Email:          input.Body.Email,
			DisplayName:    input.Body.DisplayName,
			AvatarURL:      input.Body.AvatarURL,
			AuthProvider:   authProvider,
			AuthProviderID: input.Body.AuthProviderID,
			Preferences: entities.UserPreferences{
				Language: "ja",
				DarkMode: false,
				Timezone: "Asia/Tokyo",
			},
		}
		user, err = h.userClient.CreateUser(ctx, newUser)
		if err != nil {
			return nil, errors.HandleHTTPError(ctx, err, "create_user", "failed to create user")
		}
	}

	// Convert domain entity to HTTP response
	resp := &CreateUserOutput{}
	resp.Body.ID = user.ID
	resp.Body.Email = user.Email
	resp.Body.DisplayName = user.DisplayName
	resp.Body.AvatarURL = user.AvatarURL
	resp.Body.CreatedAt = user.CreatedAt
	resp.Body.UpdatedAt = user.UpdatedAt

	return resp, nil
}

// GetCurrentUser gets the current authenticated user
func (h *UserHandler) GetCurrentUser(ctx context.Context, input *GetCurrentUserInput) (*GetCurrentUserOutput, error) {
	// Add operation context for error tracking
	ctx = errors.WithOperation(ctx, "get_current_user_http")

	// Extract authenticated user ID from context
	userID := errors.GetUserID(ctx)
	if userID == "" {
		return nil, huma.Error401Unauthorized("authentication required to get current user")
	}

	// Get user via gRPC client with authenticated user ID
	user, err := h.userClient.GetCurrentUserFromGRPC(ctx)
	if err != nil {
		return nil, errors.HandleHTTPError(ctx, err, "get_current_user", "failed to get current user")
	}

	// Convert preferences to map format for HTTP response using standardized converter
	preferences := h.userConverter.ConvertEntityPreferencesToHTTP(user.Preferences)

	// Convert domain entity to HTTP response
	resp := &GetCurrentUserOutput{}
	resp.Body.ID = user.ID
	resp.Body.Email = user.Email
	resp.Body.DisplayName = user.DisplayName
	resp.Body.AvatarURL = user.AvatarURL
	resp.Body.Preferences = preferences
	resp.Body.CreatedAt = user.CreatedAt
	resp.Body.UpdatedAt = user.UpdatedAt

	return resp, nil
}

// UpdatePreferences updates user preferences
func (h *UserHandler) UpdatePreferences(ctx context.Context, input *UpdatePreferencesInput) (*UpdatePreferencesOutput, error) {
	// Add operation context for error tracking
	ctx = errors.WithOperation(ctx, "update_user_preferences_http")

	// Extract authenticated user ID from context
	userID := errors.GetUserID(ctx)
	if userID == "" {
		return nil, huma.Error401Unauthorized("authentication required to update preferences")
	}

	// Convert HTTP preferences to domain preferences using standardized converter
	prefs := h.userConverter.ConvertHTTPPreferencesToEntity(input.Body.Preferences)

	// Update preferences via gRPC client
	user, err := h.userClient.UpdateUserPreferencesFromGRPC(ctx, userID, prefs)
	if err != nil {
		return nil, errors.HandleHTTPError(ctx, err, "update_user_preferences", "failed to update preferences")
	}

	// Convert preferences to map format for HTTP response using standardized converter
	preferences := h.userConverter.ConvertEntityPreferencesToHTTP(user.Preferences)

	// Convert domain entity to HTTP response
	resp := &UpdatePreferencesOutput{}
	resp.Body.ID = user.ID
	resp.Body.Email = user.Email
	resp.Body.DisplayName = user.DisplayName
	resp.Body.AvatarURL = user.AvatarURL
	resp.Body.Preferences = preferences
	resp.Body.UpdatedAt = user.UpdatedAt

	return resp, nil
}

// CreateHumaAuthMiddleware creates a reusable Huma-compatible authentication middleware
func CreateHumaAuthMiddleware(authMiddleware *auth.AuthMiddleware) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		// Extract JWT token from request and validate it
		claims, err := authMiddleware.ExtractAndValidateTokenFromContext(ctx)
		if err != nil {
			if strings.Contains(err.Error(), "no token found") {
				ctx.SetStatus(http.StatusUnauthorized)
				ctx.SetBody(map[string]string{"error": "Authentication required - no valid token found"})
			} else if strings.Contains(err.Error(), "authentication service error") {
				ctx.SetStatus(http.StatusInternalServerError)
				ctx.SetBody(map[string]string{"error": "Authentication service error"})
			} else {
				ctx.SetStatus(http.StatusUnauthorized)
				ctx.SetBody(map[string]string{"error": "Invalid token"})
			}
			return
		}

		// Add user context to request
		requestCtx := ctx.Context()
		requestCtx = errors.WithUserID(requestCtx, claims.UserID)
		requestCtx = errors.WithRequestID(requestCtx, ctx.Header("X-Request-ID"))
		ctx.SetContext(requestCtx)

		// Continue to next middleware/handler
		next(ctx)
	}
}

// RegisterRoutesWithAuth registers user routes with authentication middleware for secure endpoints
func (h *UserHandler) RegisterRoutesWithAuth(api huma.API, authMiddleware *auth.AuthMiddleware) {
	// Create Huma-compatible authentication middleware
	authHumaMiddleware := CreateHumaAuthMiddleware(authMiddleware)

	// Public endpoint - user creation/OAuth doesn't require authentication
	huma.Register(api, huma.Operation{
		OperationID: "create-user",
		Method:      http.MethodPost,
		Path:        "/api/v1/users",
		Summary:     "Create or update a user",
		Description: "Create a new user or update existing user information",
		Tags:        []string{"Users"},
	}, h.CreateUser)

	// Protected endpoint - get current user requires authentication
	huma.Register(api, huma.Operation{
		OperationID: "get-current-user",
		Method:      http.MethodGet,
		Path:        "/api/v1/users/me",
		Summary:     "Get current user",
		Description: "Get the current authenticated user's information",
		Tags:        []string{"Users"},
		Middlewares: huma.Middlewares{authHumaMiddleware},
	}, h.GetCurrentUser)

	// Protected endpoint - update preferences requires authentication
	huma.Register(api, huma.Operation{
		OperationID: "update-user-preferences",
		Method:      http.MethodPatch,
		Path:        "/api/v1/users/me/preferences",
		Summary:     "Update user preferences",
		Description: "Update the current user's preferences",
		Tags:        []string{"Users"},
		Middlewares: huma.Middlewares{authHumaMiddleware},
	}, h.UpdatePreferences)
}