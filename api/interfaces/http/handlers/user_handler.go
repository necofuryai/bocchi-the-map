package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"

	"bocchi/api/application/clients"
	"bocchi/api/pkg/auth"
	grpcSvc "bocchi/api/infrastructure/grpc"
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

// GetUserOutput represents the response for getting a user
type GetUserOutput struct {
	Body struct {
		ID             string                 `json:"id" doc:"User ID"`
		Email          string                 `json:"email" doc:"User email"`
		DisplayName    string                 `json:"display_name" doc:"User display name"`
		AvatarUrl      string                 `json:"avatar_url,omitempty" doc:"User avatar URL"`
		AuthProvider   string                 `json:"auth_provider" doc:"Authentication provider"`
		AuthProviderID string                 `json:"auth_provider_id" doc:"Authentication provider ID"`
		Preferences    map[string]interface{} `json:"preferences,omitempty" doc:"User preferences"`
		CreatedAt      time.Time              `json:"created_at" doc:"Creation timestamp"`
		UpdatedAt      time.Time              `json:"updated_at" doc:"Last update timestamp"`
	}
}

// GetCurrentUserInput represents the request to get current user info
type GetCurrentUserInput struct{}

// GetCurrentUserOutput represents the response for getting current user
type GetCurrentUserOutput struct {
	Body struct {
		ID             string                 `json:"id" doc:"User ID"`
		Email          string                 `json:"email" doc:"User email"`
		DisplayName    string                 `json:"display_name" doc:"User display name"`
		AvatarUrl      string                 `json:"avatar_url,omitempty" doc:"User avatar URL"`
		AuthProvider   string                 `json:"auth_provider" doc:"Authentication provider"`
		AuthProviderID string                 `json:"auth_provider_id" doc:"Authentication provider ID"`
		Preferences    map[string]interface{} `json:"preferences,omitempty" doc:"User preferences"`
		CreatedAt      time.Time              `json:"created_at" doc:"Creation timestamp"`
		UpdatedAt      time.Time              `json:"updated_at" doc:"Last update timestamp"`
	}
}

// UpdateUserInput represents the request to update a user
type UpdateUserInput struct {
	ID string `path:"id" doc:"User ID"`
	Body struct {
		DisplayName string                 `json:"display_name,omitempty" minLength:"1" maxLength:"255" doc:"User display name"`
		AvatarUrl   string                 `json:"avatar_url,omitempty" maxLength:"500" doc:"User avatar URL"`
		Preferences map[string]interface{} `json:"preferences,omitempty" doc:"User preferences"`
	}
}

// UpdateUserOutput represents the response for updating a user
type UpdateUserOutput struct {
	Body struct {
		ID             string                 `json:"id" doc:"User ID"`
		Email          string                 `json:"email" doc:"User email"`
		DisplayName    string                 `json:"display_name" doc:"User display name"`
		AvatarUrl      string                 `json:"avatar_url,omitempty" doc:"User avatar URL"`
		AuthProvider   string                 `json:"auth_provider" doc:"Authentication provider"`
		AuthProviderID string                 `json:"auth_provider_id" doc:"Authentication provider ID"`
		Preferences    map[string]interface{} `json:"preferences,omitempty" doc:"User preferences"`
		CreatedAt      time.Time              `json:"created_at" doc:"Creation timestamp"`
		UpdatedAt      time.Time              `json:"updated_at" doc:"Last update timestamp"`
	}
}

// UpdateCurrentUserInput represents the request to update current user
type UpdateCurrentUserInput struct {
	Body struct {
		DisplayName string                 `json:"display_name,omitempty" minLength:"1" maxLength:"255" doc:"User display name"`
		AvatarUrl   string                 `json:"avatar_url,omitempty" maxLength:"500" doc:"User avatar URL"`
		Preferences map[string]interface{} `json:"preferences,omitempty" doc:"User preferences"`
	}
}

// UpdateCurrentUserOutput represents the response for updating current user
type UpdateCurrentUserOutput struct {
	Body struct {
		ID             string                 `json:"id" doc:"User ID"`
		Email          string                 `json:"email" doc:"User email"`
		DisplayName    string                 `json:"display_name" doc:"User display name"`
		AvatarUrl      string                 `json:"avatar_url,omitempty" doc:"User avatar URL"`
		AuthProvider   string                 `json:"auth_provider" doc:"Authentication provider"`
		AuthProviderID string                 `json:"auth_provider_id" doc:"Authentication provider ID"`
		Preferences    map[string]interface{} `json:"preferences,omitempty" doc:"User preferences"`
		CreatedAt      time.Time              `json:"created_at" doc:"Creation timestamp"`
		UpdatedAt      time.Time              `json:"updated_at" doc:"Last update timestamp"`
	}
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
	grpcResp, err := h.userClient.GetUser(ctx, &grpcSvc.GetUserRequest{
		ID: input.ID,
	})
	if err != nil {
		return nil, grpcToHTTPError(err, "failed to get user")
	}

	// Convert gRPC response to HTTP response (public info only)
	resp := &GetUserOutput{}
	resp.Body.ID = grpcResp.User.ID
	resp.Body.DisplayName = grpcResp.User.DisplayName
	resp.Body.AvatarUrl = grpcResp.User.AvatarUrl
	resp.Body.CreatedAt = grpcResp.User.CreatedAt

	// Don't expose sensitive information in public endpoint
	// Email, AuthProvider, AuthProviderID, and Preferences are private

	return resp, nil
}

// GetCurrentUser gets the current authenticated user
func (h *UserHandler) GetCurrentUser(ctx context.Context, input *GetCurrentUserInput) (*GetCurrentUserOutput, error) {
	// Extract user ID from authentication context
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		return nil, huma.Error401Unauthorized("authentication required")
	}

	// Call gRPC service
	grpcResp, err := h.userClient.GetUser(ctx, &grpcSvc.GetUserRequest{
		ID: userID,
	})
	if err != nil {
		return nil, grpcToHTTPError(err, "failed to get current user")
	}

	// Convert gRPC response to HTTP response (full info for authenticated user)
	resp := &GetCurrentUserOutput{}
	resp.Body.ID = grpcResp.User.ID
	resp.Body.Email = grpcResp.User.Email
	resp.Body.DisplayName = grpcResp.User.DisplayName
	resp.Body.AvatarUrl = grpcResp.User.AvatarUrl
	resp.Body.AuthProvider = grpcResp.User.AuthProvider
	resp.Body.AuthProviderID = grpcResp.User.AuthProviderID
	resp.Body.Preferences = grpcResp.User.Preferences
	resp.Body.CreatedAt = grpcResp.User.CreatedAt
	resp.Body.UpdatedAt = grpcResp.User.UpdatedAt

	return resp, nil
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
	grpcResp, err := h.userClient.UpdateUser(ctx, &grpcSvc.UpdateUserRequest{
		ID:          input.ID,
		DisplayName: input.Body.DisplayName,
		AvatarUrl:   input.Body.AvatarUrl,
		Preferences: input.Body.Preferences,
	})
	if err != nil {
		return nil, grpcToHTTPError(err, "failed to update user")
	}

	// Convert gRPC response to HTTP response
	resp := &UpdateUserOutput{}
	resp.Body.ID = grpcResp.User.ID
	resp.Body.Email = grpcResp.User.Email
	resp.Body.DisplayName = grpcResp.User.DisplayName
	resp.Body.AvatarUrl = grpcResp.User.AvatarUrl
	resp.Body.AuthProvider = grpcResp.User.AuthProvider
	resp.Body.AuthProviderID = grpcResp.User.AuthProviderID
	resp.Body.Preferences = grpcResp.User.Preferences
	resp.Body.CreatedAt = grpcResp.User.CreatedAt
	resp.Body.UpdatedAt = grpcResp.User.UpdatedAt

	return resp, nil
}

// UpdateCurrentUser updates the current authenticated user
func (h *UserHandler) UpdateCurrentUser(ctx context.Context, input *UpdateCurrentUserInput) (*UpdateCurrentUserOutput, error) {
	// Extract user ID from authentication context
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		return nil, huma.Error401Unauthorized("authentication required")
	}

	// Call gRPC service
	grpcResp, err := h.userClient.UpdateUser(ctx, &grpcSvc.UpdateUserRequest{
		ID:          userID,
		DisplayName: input.Body.DisplayName,
		AvatarUrl:   input.Body.AvatarUrl,
		Preferences: input.Body.Preferences,
	})
	if err != nil {
		return nil, grpcToHTTPError(err, "failed to update current user")
	}

	// Convert gRPC response to HTTP response
	resp := &UpdateCurrentUserOutput{}
	resp.Body.ID = grpcResp.User.ID
	resp.Body.Email = grpcResp.User.Email
	resp.Body.DisplayName = grpcResp.User.DisplayName
	resp.Body.AvatarUrl = grpcResp.User.AvatarUrl
	resp.Body.AuthProvider = grpcResp.User.AuthProvider
	resp.Body.AuthProviderID = grpcResp.User.AuthProviderID
	resp.Body.Preferences = grpcResp.User.Preferences
	resp.Body.CreatedAt = grpcResp.User.CreatedAt
	resp.Body.UpdatedAt = grpcResp.User.UpdatedAt

	return resp, nil
}

// DeleteCurrentUser deletes the current authenticated user
func (h *UserHandler) DeleteCurrentUser(ctx context.Context, input *DeleteCurrentUserInput) (*DeleteCurrentUserOutput, error) {
	// Extract user ID from authentication context
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		return nil, huma.Error401Unauthorized("authentication required")
	}

	// Call gRPC service to delete the user
	_, err := h.userClient.DeleteUser(ctx, &grpcSvc.DeleteUserRequest{
		ID: userID,
	})
	if err != nil {
		return nil, grpcToHTTPError(err, "failed to delete user")
	}

	// Return empty response for 204 No Content
	return &DeleteCurrentUserOutput{}, nil
}