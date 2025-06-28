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
	service      *grpcSvc.UserService
	conn         *grpc.ClientConn
	converter    *converters.UserConverter
	grpcConverter *converters.GRPCConverter
}

// NewUserClient creates a new user client
func NewUserClient(serviceAddr string, db *sql.DB) (*UserClient, error) {
	// For internal communication in monolith, we can use direct service calls
	// In a true microservice setup, this would connect to remote gRPC service
	if serviceAddr == "internal" {
		return &UserClient{
			service:       grpcSvc.NewUserService(db),
			converter:     converters.NewUserConverter(),
			grpcConverter: converters.NewGRPCConverter(),
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


// GetUserByAuthProvider retrieves a user by authentication provider and provider ID
func (c *UserClient) GetUserByAuthProvider(ctx context.Context, authProvider entities.AuthProvider, providerID string) (*entities.User, error) {
	// Convert domain auth provider to gRPC enum using converter
	grpcAuthProvider := c.grpcConverter.ConvertEntityAuthProviderToGRPC(authProvider)

	// Call gRPC service method
	resp, err := c.service.GetUserByAuthProviderGRPC(ctx, &grpcSvc.GetUserByAuthProviderRequest{
		AuthProvider:   grpcAuthProvider,
		AuthProviderID: providerID,
	})
	if err != nil {
		return nil, err
	}

	// Convert gRPC response to domain entity
	return c.convertGRPCUserToEntity(resp.User)
}

// CreateUser creates a new user
func (c *UserClient) CreateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	// Convert domain auth provider to gRPC enum using converter
	grpcAuthProvider := c.grpcConverter.ConvertEntityAuthProviderToGRPC(user.AuthProvider)

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
	return c.convertGRPCUserToEntity(resp.User)
}

// UpdateUser updates an existing user
func (c *UserClient) UpdateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	// Convert domain preferences to gRPC preferences using converter
	grpcPrefs := c.grpcConverter.ConvertEntityPreferencesToGRPC(user.Preferences)

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
	return c.convertGRPCUserToEntity(resp.User)
}

// convertGRPCUserToEntity converts gRPC User struct to domain entity
func (c *UserClient) convertGRPCUserToEntity(grpcUser *grpcSvc.User) (*entities.User, error) {
	if grpcUser == nil {
		return nil, nil
	}

	// Convert gRPC auth provider to domain enum using converter
	authProvider, err := c.grpcConverter.ConvertGRPCAuthProviderToEntity(grpcUser.AuthProvider)
	if err != nil {
		return nil, err
	}

	// Convert gRPC preferences to domain preferences using converter
	prefs := c.grpcConverter.ConvertGRPCPreferencesToEntity(grpcUser.Preferences)

	return &entities.User{
		ID:             grpcUser.ID,
		Email:          grpcUser.Email,
		DisplayName:    grpcUser.DisplayName,
		AvatarURL:      grpcUser.AvatarURL,
		AuthProvider:   authProvider,
		AuthProviderID: grpcUser.AuthProviderID,
		Preferences:    prefs,
		CreatedAt:      grpcUser.CreatedAt,
		UpdatedAt:      grpcUser.UpdatedAt,
	}, nil
}

// GetCurrentUserFromGRPC gets the current user via gRPC service  
func (c *UserClient) GetCurrentUserFromGRPC(ctx context.Context) (*entities.User, error) {
	// Call gRPC service method
	resp, err := c.service.GetCurrentUser(ctx, &grpcSvc.GetCurrentUserRequest{})
	if err != nil {
		return nil, err
	}

	// Convert gRPC response to domain entity
	return c.convertGRPCUserToEntity(resp.User)
}

// UpdateUserPreferencesFromGRPC updates user preferences via gRPC service
func (c *UserClient) UpdateUserPreferencesFromGRPC(ctx context.Context, userID string, prefs entities.UserPreferences) (*entities.User, error) {
	// Convert domain preferences to gRPC preferences using converter
	grpcPrefs := c.grpcConverter.ConvertEntityPreferencesToGRPC(prefs)

	// Call gRPC service method
	resp, err := c.service.UpdateUserPreferences(ctx, &grpcSvc.UpdateUserPreferencesRequest{
		UserID:      userID,
		Preferences: grpcPrefs,
	})
	if err != nil {
		return nil, err
	}

	// Convert gRPC response to domain entity
	return c.convertGRPCUserToEntity(resp.User)
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
	return c.convertGRPCUserToEntity(resp.User)
}