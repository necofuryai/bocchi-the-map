package clients

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	grpcSvc "github.com/necofuryai/bocchi-the-map/api/infrastructure/grpc"
)

// UserClient wraps gRPC client calls for user operations
type UserClient struct {
	service *grpcSvc.UserService
	conn    *grpc.ClientConn
}

// NewUserClient creates a new user client
func NewUserClient(serviceAddr string) (*UserClient, error) {
	// For internal communication in monolith, we can use direct service calls
	// In a true microservice setup, this would connect to remote gRPC service
	if serviceAddr == "internal" {
		return &UserClient{
			service: grpcSvc.NewUserService(),
		}, nil
	}

	// For external gRPC service connection
	// TODO: Use TLS credentials in production
	var creds credentials.TransportCredentials
	if os.Getenv("GRPC_INSECURE") == "true" {
		creds = insecure.NewCredentials()
	} else {
		creds = credentials.NewTLS(&tls.Config{
			MinVersion: tls.VersionTLS13,
		})
	}
	conn, err := grpc.NewClient(serviceAddr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}

	return &UserClient{
		conn: conn,
		// TODO: Use generated gRPC client when protobuf is available
		service: grpcSvc.NewUserService(),
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