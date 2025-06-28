package clients

import (
	"context"
	"database/sql"
	"fmt"

	"google.golang.org/grpc"

	"github.com/necofuryai/bocchi-the-map/api/domain/entities"
	grpcSvc "github.com/necofuryai/bocchi-the-map/api/infrastructure/grpc"
	"github.com/necofuryai/bocchi-the-map/api/pkg/converters"
)

// UserClient wraps gRPC client calls for user operations
type UserClient struct {
	service   *grpcSvc.UserService
	conn      *grpc.ClientConn
	converter *converters.UserConverter
}

// NewUserClient creates a new user client
func NewUserClient(serviceAddr string, db *sql.DB) (*UserClient, error) {
	// For internal communication in monolith, we can use direct service calls
	// In a true microservice setup, this would connect to remote gRPC service
	if serviceAddr == "internal" {
		return &UserClient{
			service:   grpcSvc.NewUserService(db),
			converter: converters.NewUserConverter(),
		}, nil
	}

	// TODO: Implement external gRPC service connection when protobuf client is ready
	// For now, return error for external services to avoid silent failures
	return nil, fmt.Errorf("external gRPC service not implemented yet: %s", serviceAddr)
}

// Close closes the gRPC connection
func (c *UserClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// GetCurrentUser retrieves current user via gRPC
func (c *UserClient) GetCurrentUser(ctx context.Context, req *grpcSvc.GetCurrentUserRequest) (*grpcSvc.GetCurrentUserResponse, error) {
	return c.service.GetCurrentUser(ctx, req)
}

// UpdateUserPreferences updates user preferences via gRPC
func (c *UserClient) UpdateUserPreferences(ctx context.Context, req *grpcSvc.UpdateUserPreferencesRequest) (*grpcSvc.UpdateUserPreferencesResponse, error) {
	return c.service.UpdateUserPreferences(ctx, req)
}

// convertAuthProviderToGRPC converts domain AuthProvider to gRPC AuthProvider
func (c *UserClient) convertAuthProviderToGRPC(authProvider entities.AuthProvider) (grpcSvc.AuthProvider, error) {
	switch authProvider {
	case entities.AuthProviderGoogle:
		return grpcSvc.AuthProviderGoogle, nil
	case entities.AuthProviderX:
		return grpcSvc.AuthProviderX, nil
	default:
		return grpcSvc.AuthProviderUnspecified, fmt.Errorf("invalid auth provider")
	}
}

// GetUserByAuthProvider retrieves a user by authentication provider and provider ID
func (c *UserClient) GetUserByAuthProvider(ctx context.Context, authProvider entities.AuthProvider, providerID string) (*entities.User, error) {
	// Convert domain auth provider to gRPC enum
	grpcAuthProvider, err := c.convertAuthProviderToGRPC(authProvider)
	if err != nil {
		return nil, err
	}

	// Call gRPC service method
	resp, err := c.service.GetUserByAuthProviderGRPC(ctx, &grpcSvc.GetUserByAuthProviderRequest{
		AuthProvider:   grpcAuthProvider,
		AuthProviderID: providerID,
	})
	if err != nil {
		return nil, err
	}

	// Convert gRPC response to domain entity
	return c.convertGRPCUserToEntity(resp.User), nil
}

// CreateUser creates a new user
func (c *UserClient) CreateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	// Convert domain auth provider to gRPC enum
	grpcAuthProvider, err := c.convertAuthProviderToGRPC(user.AuthProvider)
	if err != nil {
		return nil, err
	}

	// Call gRPC service method
	resp, err := c.service.CreateUserGRPC(ctx, &grpcSvc.CreateUserRequest{
		Email:          user.Email,
		DisplayName:    user.DisplayName,
		AvatarURL:      user.AvatarURL,
		AuthProvider:   grpcAuthProvider,
		AuthProviderID: user.AuthProviderID,
	})
	if err != nil {
		return nil, err
	}

	// Convert gRPC response to domain entity
	return c.convertGRPCUserToEntity(resp.User), nil
}

// UpdateUser updates an existing user
func (c *UserClient) UpdateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	// Convert domain preferences to gRPC preferences
	grpcPrefs := &grpcSvc.UserPreferences{
		Language: user.Preferences.Language,
		DarkMode: user.Preferences.DarkMode,
		Timezone: user.Preferences.Timezone,
	}

	// Call gRPC service method
	resp, err := c.service.UpdateUserGRPC(ctx, &grpcSvc.UpdateUserRequest{
		ID:          user.ID,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		AvatarURL:   user.AvatarURL,
		Preferences: grpcPrefs,
	})
	if err != nil {
		return nil, err
	}

	// Convert gRPC response to domain entity
	return c.convertGRPCUserToEntity(resp.User), nil
}

// convertGRPCUserToEntity converts gRPC User struct to domain entity
func (c *UserClient) convertGRPCUserToEntity(grpcUser *grpcSvc.User) *entities.User {
	if grpcUser == nil {
		return nil
	}

	// Convert gRPC auth provider to domain enum
	var authProvider entities.AuthProvider
	switch grpcUser.AuthProvider {
	case grpcSvc.AuthProviderGoogle:
		authProvider = entities.AuthProviderGoogle
	case grpcSvc.AuthProviderX:
		authProvider = entities.AuthProviderX
	default:
		authProvider = entities.AuthProviderUnspecified
	}

	// Convert gRPC preferences to domain preferences
	var prefs entities.UserPreferences
	if grpcUser.Preferences != nil {
		prefs = entities.UserPreferences{
			Language: grpcUser.Preferences.Language,
			DarkMode: grpcUser.Preferences.DarkMode,
			Timezone: grpcUser.Preferences.Timezone,
		}
	} else {
		// Default preferences
		prefs = entities.UserPreferences{
			Language: "ja",
			DarkMode: false,
			Timezone: "Asia/Tokyo",
		}
	}

	return &entities.User{
		ID:             grpcUser.ID,
		AnonymousID:    grpcUser.AnonymousID,
		Email:          grpcUser.Email,
		DisplayName:    grpcUser.DisplayName,
		AvatarURL:      grpcUser.AvatarURL,
		AuthProvider:   authProvider,
		AuthProviderID: grpcUser.AuthProviderID,
		Preferences:    prefs,
		CreatedAt:      grpcUser.CreatedAt,
		UpdatedAt:      grpcUser.UpdatedAt,
	}
}

// GetCurrentUserFromGRPC gets the current user via gRPC service  
func (c *UserClient) GetCurrentUserFromGRPC(ctx context.Context) (*entities.User, error) {
	// Call gRPC service method
	resp, err := c.service.GetCurrentUser(ctx, &grpcSvc.GetCurrentUserRequest{})
	if err != nil {
		return nil, err
	}

	// Convert gRPC response to domain entity
	return c.convertGRPCUserToEntity(resp.User), nil
}

// UpdateUserPreferencesFromGRPC updates user preferences via gRPC service
func (c *UserClient) UpdateUserPreferencesFromGRPC(ctx context.Context, prefs entities.UserPreferences) (*entities.User, error) {
	// Convert domain preferences to gRPC preferences
	grpcPrefs := &grpcSvc.UserPreferences{
		Language: prefs.Language,
		DarkMode: prefs.DarkMode,
		Timezone: prefs.Timezone,
	}

	// Call gRPC service method
	resp, err := c.service.UpdateUserPreferences(ctx, &grpcSvc.UpdateUserPreferencesRequest{
		Preferences: grpcPrefs,
	})
	if err != nil {
		return nil, err
	}

	// Convert gRPC response to domain entity
	return c.convertGRPCUserToEntity(resp.User), nil
}

// GetConverter returns the user converter instance
func (c *UserClient) GetConverter() *converters.UserConverter {
	return c.converter
}

// GetUserByID retrieves a user by ID
func (c *UserClient) GetUserByID(ctx context.Context, userID string) (*entities.User, error) {
	// Call gRPC service method
	resp, err := c.service.GetUserByID(ctx, &grpcSvc.GetUserByIDRequest{
		ID: userID,
	})
	if err != nil {
		return nil, err
	}

	// Convert gRPC response to domain entity
	return c.convertGRPCUserToEntity(resp.User), nil
}