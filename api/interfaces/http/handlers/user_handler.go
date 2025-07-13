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



