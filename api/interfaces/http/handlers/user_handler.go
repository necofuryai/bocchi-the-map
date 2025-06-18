package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/necofuryai/bocchi-the-map/api/infrastructure/database"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	queries *database.Queries
}

// NewUserHandler creates a new user handler
func NewUserHandler(queries *database.Queries) *UserHandler {
	return &UserHandler{
		queries: queries,
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
	// Convert provider string to enum
	var authProvider database.UsersAuthProvider
	switch input.Body.AuthProvider {
	case "google":
		authProvider = database.UsersAuthProviderGoogle
	case "twitter":
		authProvider = database.UsersAuthProviderTwitter
	case "x":
		authProvider = database.UsersAuthProviderX
	default:
		return nil, huma.Error400BadRequest("invalid auth provider")
	}

	// Check if user already exists to avoid ID conflicts and preserve preferences
	existingUser, err := h.queries.GetUserByProviderID(ctx, database.GetUserByProviderIDParams{
		AuthProvider:   authProvider,
		AuthProviderID: input.Body.AuthProviderID,
	})
	
	var userID string
	var prefsJSON []byte
	
	if err == sql.ErrNoRows {
		// New user - generate UUID and use default preferences
		userID = uuid.New().String()
		defaultPrefs := map[string]interface{}{
			"language":  "ja",
			"dark_mode": false,
			"timezone":  "Asia/Tokyo",
		}
		prefsJSON, err = json.Marshal(defaultPrefs)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to marshal preferences")
		}
	} else if err != nil {
		return nil, huma.Error500InternalServerError("failed to check existing user")
	} else {
		// Existing user - preserve ID and existing preferences
		userID = existingUser.ID
		prefsJSON = existingUser.Preferences
	}

	// Convert avatar URL to nullable string
	var avatarURL sql.NullString
	if input.Body.AvatarURL != "" {
		avatarURL = sql.NullString{String: input.Body.AvatarURL, Valid: true}
	}

	// Upsert user in database
	err = h.queries.UpsertUser(ctx, database.UpsertUserParams{
		ID:             userID,
		Email:          input.Body.Email,
		DisplayName:    input.Body.DisplayName,
		AvatarUrl:      avatarURL,
		AuthProvider:   authProvider,
		AuthProviderID: input.Body.AuthProviderID,
		Preferences:    prefsJSON,
	})
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to create/update user")
	}

	// Retrieve the created/updated user to get accurate timestamps
	user, err := h.queries.GetUserByProviderID(ctx, database.GetUserByProviderIDParams{
		AuthProvider:   authProvider,
		AuthProviderID: input.Body.AuthProviderID,
	})
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to retrieve user")
	}

	// Convert response
	resp := &CreateUserOutput{}
	resp.Body.ID = user.ID
	resp.Body.Email = user.Email
	resp.Body.DisplayName = user.DisplayName
	if user.AvatarUrl.Valid {
		resp.Body.AvatarURL = user.AvatarUrl.String
	}
	resp.Body.CreatedAt = user.CreatedAt
	resp.Body.UpdatedAt = user.UpdatedAt

	return resp, nil
}

// GetCurrentUser gets the current authenticated user
func (h *UserHandler) GetCurrentUser(ctx context.Context, input *GetCurrentUserInput) (*GetCurrentUserOutput, error) {
	// TODO: CRITICAL - Implement authentication context extraction
	// This endpoint is currently INSECURE and should not be used in production
	// Required implementation:
	// 1. Add authentication middleware to extract user ID from JWT/session
	// 2. Add user ID to request context in middleware
	// 3. Extract user ID from context here using: userID, ok := ctx.Value("user_id").(string)
	// 4. Return 401 Unauthorized if user ID is not found in context
	
	// SECURITY WARNING: This hardcoded user ID is for development only
	userID := "user_123"
	if userID == "" {
		return nil, huma.Error401Unauthorized("user not authenticated")
	}
	user, err := h.queries.GetUserByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, huma.Error404NotFound("user not found")
		}
		return nil, huma.Error500InternalServerError("failed to get user")
	}

	// Parse preferences JSON
	var preferences map[string]interface{}
	if len(user.Preferences) > 0 {
		if err := json.Unmarshal(user.Preferences, &preferences); err != nil {
			// If preferences are malformed, use defaults
			preferences = map[string]interface{}{
				"language":  "ja",
				"dark_mode": false,
				"timezone":  "Asia/Tokyo",
			}
		}
	}

	// Convert response
	resp := &GetCurrentUserOutput{}
	resp.Body.ID = user.ID
	resp.Body.Email = user.Email
	resp.Body.DisplayName = user.DisplayName
	if user.AvatarUrl.Valid {
		resp.Body.AvatarURL = user.AvatarUrl.String
	}
	resp.Body.Preferences = preferences
	resp.Body.CreatedAt = user.CreatedAt
	resp.Body.UpdatedAt = user.UpdatedAt

	return resp, nil
}

// UpdatePreferences updates user preferences
func (h *UserHandler) UpdatePreferences(ctx context.Context, input *UpdatePreferencesInput) (*UpdatePreferencesOutput, error) {
	// TODO: CRITICAL - Implement authentication context extraction
	// This endpoint is currently INSECURE and should not be used in production
	// Required implementation:
	// 1. Add authentication middleware to extract user ID from JWT/session
	// 2. Add user ID to request context in middleware
	// 3. Extract user ID from context here using: userID, ok := ctx.Value("user_id").(string)
	// 4. Return 401 Unauthorized if user ID is not found in context
	
	// SECURITY WARNING: This hardcoded user ID is for development only
	userID := "user_123"
	if userID == "" {
		return nil, huma.Error401Unauthorized("user not authenticated")
	}

	// Convert preferences to JSON
	prefsJSON, err := json.Marshal(input.Body.Preferences)
	if err != nil {
		return nil, huma.Error400BadRequest("invalid preferences format")
	}

	// Update preferences in database
	err = h.queries.UpdateUserPreferences(ctx, database.UpdateUserPreferencesParams{
		ID:          userID,
		Preferences: prefsJSON,
	})
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to update preferences")
	}

	// Retrieve updated user
	user, err := h.queries.GetUserByID(ctx, userID)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to retrieve updated user")
	}

	// Parse updated preferences
	var preferences map[string]interface{}
	if len(user.Preferences) > 0 {
		json.Unmarshal(user.Preferences, &preferences)
	}

	// Convert response
	resp := &UpdatePreferencesOutput{}
	resp.Body.ID = user.ID
	resp.Body.Email = user.Email
	resp.Body.DisplayName = user.DisplayName
	if user.AvatarUrl.Valid {
		resp.Body.AvatarURL = user.AvatarUrl.String
	}
	resp.Body.Preferences = preferences
	resp.Body.UpdatedAt = user.UpdatedAt

	return resp, nil
}