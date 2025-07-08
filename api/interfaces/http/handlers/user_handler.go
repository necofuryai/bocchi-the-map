package handlers

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"bocchi/api/application/clients"
	"bocchi/api/pkg/auth"
	userv1 "bocchi/api/gen/user/v1"
)


// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userClient *clients.UserClient
}

// NewUserHandler creates a new user handler
func NewUserHandler(userClient *clients.UserClient) *UserHandler {
	if userClient == nil {
		panic("userClient cannot be nil")
	}
	return &UserHandler{
		userClient: userClient,
	}
}

// GetUserInput represents the request to get a user
type GetUserInput struct {
	ID string `path:"id" doc:"User ID"`
}

// GetUserOutput represents the response for getting a user (using protobuf User type)
type GetUserOutput struct {
	Body *userv1.User `json:"user" doc:"User data"`
}

// GetCurrentUserInput represents the request to get current user info
type GetCurrentUserInput struct{}

// GetCurrentUserOutput represents the response for getting current user (using protobuf User type)
type GetCurrentUserOutput struct {
	Body *userv1.User `json:"user" doc:"User data"`
}

// UpdateUserInput represents the request to update a user
type UpdateUserInput struct {
	ID string `path:"id" doc:"User ID"`
	Body struct {
		DisplayName string `json:"display_name,omitempty" minLength:"1" maxLength:"255" doc:"User display name"`
		AvatarUrl   string `json:"avatar_url,omitempty" maxLength:"500" doc:"User avatar URL"`
		Preferences string `json:"preferences,omitempty" doc:"User preferences as JSON string"`
	}
}

// UpdateUserOutput represents the response for updating a user (using protobuf User type)
type UpdateUserOutput struct {
	Body *userv1.User `json:"user" doc:"User data"`
}

// UpdateCurrentUserInput represents the request to update current user
type UpdateCurrentUserInput struct {
	Body struct {
		DisplayName string `json:"display_name,omitempty" minLength:"1" maxLength:"255" doc:"User display name"`
		AvatarUrl   string `json:"avatar_url,omitempty" maxLength:"500" doc:"User avatar URL"`
		Preferences string `json:"preferences,omitempty" doc:"User preferences as JSON string"`
	}
}

// UpdateCurrentUserOutput represents the response for updating current user (using protobuf User type)
type UpdateCurrentUserOutput struct {
	Body *userv1.User `json:"user" doc:"User data"`
}

// DeleteCurrentUserInput represents the request to delete current user
type DeleteCurrentUserInput struct{}

// DeleteCurrentUserOutput represents the response for deleting current user
type DeleteCurrentUserOutput struct{}  // Empty response body for 204 No Content

// RegisterRoutes registers user routes (without authentication)
func (h *UserHandler) RegisterRoutes(api huma.API) {
	// Get user by ID (public endpoint)
	huma.Register(api, huma.Operation{
		OperationID: "get-user",
		Method:      http.MethodGet,
		Path:        "/api/v1/users/{id}",
		Summary:     "Get a user",
		Description: "Get public user information by ID",
		Tags:        []string{"Users"},
	}, h.GetUser)
}

// RegisterRoutesWithAuth registers user routes with authentication middleware
func (h *UserHandler) RegisterRoutesWithAuth(api huma.API, authMiddleware *auth.AuthMiddleware) {
	// Register public routes first
	h.RegisterRoutes(api)

	// Get current user info (requires authentication)
	huma.Register(api, huma.Operation{
		OperationID: "get-current-user",
		Method:      http.MethodGet,
		Path:        "/api/v1/users/me",
		Summary:     "Get current user",
		Description: "Get current authenticated user information",
		Tags:        []string{"Users"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, h.GetCurrentUser)

	// Update user by ID (admin or self only)
	huma.Register(api, huma.Operation{
		OperationID: "update-user",
		Method:      http.MethodPut,
		Path:        "/api/v1/users/{id}",
		Summary:     "Update a user",
		Description: "Update user information (self or admin only)",
		Tags:        []string{"Users"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, h.UpdateUser)

	// Update current user (requires authentication)
	huma.Register(api, huma.Operation{
		OperationID: "update-current-user",
		Method:      http.MethodPut,
		Path:        "/api/v1/users/me",
		Summary:     "Update current user",
		Description: "Update current authenticated user information",
		Tags:        []string{"Users"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, h.UpdateCurrentUser)

	// Delete current user (requires authentication)
	huma.Register(api, huma.Operation{
		OperationID: "delete-current-user",
		Method:      http.MethodDelete,
		Path:        "/api/v1/users/me",
		Summary:     "Delete current user",
		Description: "Delete current authenticated user account",
		Tags:        []string{"Users"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, h.DeleteCurrentUser)
}

// GetUser gets a user by ID (public info only)
func (h *UserHandler) GetUser(ctx context.Context, input *GetUserInput) (*GetUserOutput, error) {
	// Call gRPC service
	grpcResp, err := h.userClient.GetUser(ctx, &userv1.GetUserRequest{
		Id: input.ID,
	})
	if err != nil {
		return nil, grpcToHTTPError(err, "failed to get user")
	}

	// Create public user info (remove sensitive data)
	publicUser := &userv1.User{
		Id:          grpcResp.User.Id,
		DisplayName: grpcResp.User.DisplayName,
		AvatarUrl:   grpcResp.User.AvatarUrl,
		CreatedAt:   grpcResp.User.CreatedAt,
	}

	return &GetUserOutput{Body: publicUser}, nil
}

// GetCurrentUser gets the current authenticated user
func (h *UserHandler) GetCurrentUser(ctx context.Context, input *GetCurrentUserInput) (*GetCurrentUserOutput, error) {
	// Extract user ID from authentication context
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		return nil, huma.Error401Unauthorized("authentication required")
	}

	// Call gRPC service
	grpcResp, err := h.userClient.GetUser(ctx, &userv1.GetUserRequest{
		Id: userID,
	})
	if err != nil {
		return nil, grpcToHTTPError(err, "failed to get current user")
	}

	// Return full user info for authenticated user (protobuf User)
	return &GetCurrentUserOutput{Body: grpcResp.User}, nil
}

// UpdateUser updates a user by ID (admin or self only)
func (h *UserHandler) UpdateUser(ctx context.Context, input *UpdateUserInput) (*UpdateUserOutput, error) {
	// Extract user ID from authentication context
	authUserID, ok := ctx.Value("user_id").(string)
	if !ok || authUserID == "" {
		return nil, huma.Error401Unauthorized("authentication required")
	}

	// Check if user is updating themselves or has admin permissions
	if authUserID != input.ID {
		// TODO: Add admin permission check here
		return nil, huma.Error403Forbidden("insufficient permissions to update this user")
	}

	// Call gRPC service
	grpcResp, err := h.userClient.UpdateUser(ctx, &userv1.UpdateUserRequest{
		Id:          input.ID,
		DisplayName: input.Body.DisplayName,
		AvatarUrl:   input.Body.AvatarUrl,
		Preferences: input.Body.Preferences,
	})
	if err != nil {
		return nil, grpcToHTTPError(err, "failed to update user")
	}

	// Return updated user (protobuf User)
	return &UpdateUserOutput{Body: grpcResp.User}, nil
}

// UpdateCurrentUser updates the current authenticated user
func (h *UserHandler) UpdateCurrentUser(ctx context.Context, input *UpdateCurrentUserInput) (*UpdateCurrentUserOutput, error) {
	// Extract user ID from authentication context
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		return nil, huma.Error401Unauthorized("authentication required")
	}

	// Call gRPC service
	grpcResp, err := h.userClient.UpdateUser(ctx, &userv1.UpdateUserRequest{
		Id:          userID,
		DisplayName: input.Body.DisplayName,
		AvatarUrl:   input.Body.AvatarUrl,
		Preferences: input.Body.Preferences,
	})
	if err != nil {
		return nil, grpcToHTTPError(err, "failed to update current user")
	}

	// Return updated user (protobuf User)
	return &UpdateCurrentUserOutput{Body: grpcResp.User}, nil
}

// DeleteCurrentUser deletes the current authenticated user
func (h *UserHandler) DeleteCurrentUser(ctx context.Context, input *DeleteCurrentUserInput) (*DeleteCurrentUserOutput, error) {
	// Extract user ID from authentication context
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		return nil, huma.Error401Unauthorized("authentication required")
	}

	// Call gRPC service to delete the user
	_, err := h.userClient.DeleteUser(ctx, &userv1.DeleteUserRequest{
		Id: userID,
	})
	if err != nil {
		return nil, grpcToHTTPError(err, "failed to delete user")
	}

	// Return empty response for 204 No Content
	return &DeleteCurrentUserOutput{}, nil
}