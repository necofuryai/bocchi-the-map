package clients

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/necofuryai/bocchi-the-map/api/domain/entities"
	grpcSvc "github.com/necofuryai/bocchi-the-map/api/infrastructure/grpc"
)

// UserClient wraps gRPC client calls for user operations
type UserClient struct {
	service *grpcSvc.UserService
	conn    *grpc.ClientConn
}

// NewUserClient creates a new user client
func NewUserClient(serviceAddr string, db *sql.DB) (*UserClient, error) {
	// For internal communication in monolith, we can use direct service calls
	// In a true microservice setup, this would connect to remote gRPC service
	if serviceAddr == "internal" {
		return &UserClient{
			service: grpcSvc.NewUserService(db),
		}, nil
	}

	// TODO: Implement external gRPC service connection when protobuf client is ready
	// For now, skip connection since we're using local service anyway
	return &UserClient{
		conn: nil, // No connection needed for local service
		service: grpcSvc.NewUserService(db),
	}, nil
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
	// For now, call the gRPC service method
	// In the future, this would call a proper gRPC method
	return c.service.GetUserByAuthProvider(ctx, authProvider, providerID)
}

// CreateUser creates a new user
func (c *UserClient) CreateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	// For now, call the gRPC service method
	// In the future, this would call a proper gRPC method
	return c.service.CreateUser(ctx, user)
}

// UpdateUser updates an existing user
func (c *UserClient) UpdateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	// For now, call the gRPC service method
	// In the future, this would call a proper gRPC method
	return c.service.UpdateUser(ctx, user)
}