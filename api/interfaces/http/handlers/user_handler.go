package handlers

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/necofuryai/bocchi-the-map/api/application/clients"
	"github.com/necofuryai/bocchi-the-map/api/domain/entities"
)

// CreateUserRequest represents the request payload for creating a user
type CreateUserRequest struct {
	Email          string `json:"email" validate:"required,email" example:"user@example.com" doc:"User email address"`
	DisplayName    string `json:"display_name" validate:"required,min=1,max=100" example:"John Doe" doc:"User display name"`
	AvatarURL      string `json:"avatar_url,omitempty" example:"https://example.com/avatar.jpg" doc:"User avatar URL"`
	AuthProvider   string `json:"auth_provider" validate:"required,oneof=google x" example:"google" doc:"Authentication provider (google or x)"`
	AuthProviderID string `json:"auth_provider_id" validate:"required" example:"google_123456789" doc:"Provider-specific user ID"`
}

// CreateUserResponse represents the response for creating a user
type CreateUserResponse struct {
	Body *entities.User `json:"user" doc:"Created user information"`
}

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userClient *clients.UserClient
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(userClient *clients.UserClient) *UserHandler {
	return &UserHandler{
		userClient: userClient,
	}
}

// CreateUserOperation represents the operation metadata
type CreateUserOperation struct {
	Request  CreateUserRequest
	Response CreateUserResponse
}

// CreateUser handles POST /api/users requests to create or update a user
func (h *UserHandler) CreateUser(ctx context.Context, input *CreateUserRequest) (*CreateUserResponse, error) {
	// Convert string auth provider to domain type
	var authProvider entities.AuthProvider
	switch input.AuthProvider {
	case "google":
		authProvider = entities.AuthProviderGoogle
	case "x":
		authProvider = entities.AuthProviderX  
	default:
		return nil, huma.Error400BadRequest("Invalid auth provider. Must be 'google' or 'x'")
	}

	// Check if user already exists by auth provider and provider ID
	existingUser, err := h.userClient.GetUserByAuthProvider(ctx, authProvider, input.AuthProviderID)
	if err != nil && err.Error() != "user not found" {
		return nil, huma.Error500InternalServerError("Failed to check existing user", err)
	}

	var user *entities.User

	if existingUser != nil {
		// Update existing user
		existingUser.Email = input.Email
		existingUser.DisplayName = input.DisplayName
		existingUser.AvatarURL = input.AvatarURL

		user, err = h.userClient.UpdateUser(ctx, existingUser)
		if err != nil {
			return nil, huma.Error500InternalServerError("Failed to update user", err)
		}
	} else {
		// Create new user
		newUser := entities.NewUser(input.Email, input.DisplayName, authProvider, input.AuthProviderID)
		newUser.AvatarURL = input.AvatarURL

		user, err = h.userClient.CreateUser(ctx, newUser)
		if err != nil {
			return nil, huma.Error500InternalServerError("Failed to create user", err)
		}
	}

	return &CreateUserResponse{
		Body: user,
	}, nil
}

// RegisterUserRoutes registers user-related routes with the API
func RegisterUserRoutes(api huma.API, userHandler *UserHandler) {
	huma.Register(api, huma.Operation{
		OperationID: "CreateUser",
		Method:      http.MethodPost,
		Path:        "/api/users",
		Summary:     "Create or update user",
		Description: "Creates a new user or updates an existing user based on authentication provider information. If a user with the same auth_provider and auth_provider_id already exists, it will be updated.",
		Tags:        []string{"users"},
	}, userHandler.CreateUser)
}