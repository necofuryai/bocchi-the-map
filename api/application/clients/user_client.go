package clients

import (
	"context"
	"database/sql"
	"fmt"

	"google.golang.org/grpc"

	grpcSvc "bocchi/api/infrastructure/grpc"
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

// GetUser retrieves a user by ID via gRPC
func (c *UserClient) GetUser(ctx context.Context, req *grpcSvc.GetUserRequest) (*grpcSvc.GetUserResponse, error) {
	return c.service.GetUser(ctx, req)
}

// GetUserByEmail retrieves a user by email via gRPC
func (c *UserClient) GetUserByEmail(ctx context.Context, req *grpcSvc.GetUserByEmailRequest) (*grpcSvc.GetUserByEmailResponse, error) {
	return c.service.GetUserByEmail(ctx, req)
}

// CreateUser creates a new user via gRPC
func (c *UserClient) CreateUser(ctx context.Context, req *grpcSvc.CreateUserRequest) (*grpcSvc.CreateUserResponse, error) {
	return c.service.CreateUser(ctx, req)
}

// UpdateUser updates a user via gRPC
func (c *UserClient) UpdateUser(ctx context.Context, req *grpcSvc.UpdateUserRequest) (*grpcSvc.UpdateUserResponse, error) {
	return c.service.UpdateUser(ctx, req)
}

// DeleteUser deletes a user via gRPC
func (c *UserClient) DeleteUser(ctx context.Context, req *grpcSvc.DeleteUserRequest) (*grpcSvc.DeleteUserResponse, error) {
	return c.service.DeleteUser(ctx, req)
}