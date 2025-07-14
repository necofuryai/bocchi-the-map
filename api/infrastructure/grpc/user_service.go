package grpc

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"bocchi/api/gen/user/v1"
	"bocchi/api/infrastructure/database"
	"bocchi/api/pkg/logger"
)

// UserService implements the gRPC UserService
type UserService struct {
	queries *database.Queries
}

// NewUserService creates a new UserService instance
func NewUserService(db *sql.DB) *UserService {
	return &UserService{
		queries: database.New(db),
	}
}

// Use Protocol Buffers generated types
type (
	User                   = userv1.User
	GetUserRequest         = userv1.GetUserRequest
	GetUserResponse        = userv1.GetUserResponse
	GetUserByEmailRequest  = userv1.GetUserByEmailRequest
	GetUserByEmailResponse = userv1.GetUserByEmailResponse
	CreateUserRequest      = userv1.CreateUserRequest
	CreateUserResponse     = userv1.CreateUserResponse
)

// GetUser retrieves a user by ID
func (s *UserService) GetUser(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	// Get user from database
	dbUser, err := s.queries.GetUserByID(ctx, req.GetId())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		logger.ErrorWithContext(ctx, "Failed to get user by ID", err)
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	// Convert database user to gRPC response
	user := s.convertDatabaseUserToGRPC(dbUser)
	return &GetUserResponse{User: user}, nil
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(ctx context.Context, req *GetUserByEmailRequest) (*GetUserByEmailResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	// Get user from database
	dbUser, err := s.queries.GetUserByEmail(ctx, req.GetEmail())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		logger.ErrorWithContext(ctx, "Failed to get user by email", err)
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	// Convert database user to gRPC response
	user := s.convertDatabaseUserToGRPC(dbUser)
	return &GetUserByEmailResponse{User: user}, nil
}

// CreateUser creates a new user (MVP version)
func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
	// Validate request
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	// Check if user already exists
	_, err := s.queries.GetUserByEmail(ctx, req.GetEmail())
	if err == nil {
		return nil, status.Error(codes.AlreadyExists, "user with this email already exists")
	} else if err != sql.ErrNoRows {
		logger.ErrorWithContext(ctx, "Failed to check existing user", err)
		return nil, status.Error(codes.Internal, "failed to check existing user")
	}

	// Generate UUID for new user
	userID := uuid.New().String()

	// Create user in database
	err = s.queries.CreateUser(ctx, database.CreateUserParams{
		ID:    userID,
		Email: req.GetEmail(),
		Name:  req.GetName(),
	})
	if err != nil {
		logger.ErrorWithContext(ctx, "Failed to create user", err)
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	// Retrieve the created user to get accurate timestamps
	dbUser, err := s.queries.GetUserByID(ctx, userID)
	if err != nil {
		logger.ErrorWithContext(ctx, "Failed to retrieve created user", err)
		return nil, status.Error(codes.Internal, "failed to retrieve created user")
	}

	// Convert database user to gRPC response
	user := s.convertDatabaseUserToGRPC(dbUser)
	return &CreateUserResponse{User: user}, nil
}



// convertDatabaseUserToGRPC converts database user model to gRPC user struct (MVP version)
func (s *UserService) convertDatabaseUserToGRPC(dbUser database.User) *User {
	return &User{
		Id:        dbUser.ID,
		Email:     dbUser.Email,
		Name:      dbUser.Name,
		CreatedAt: timestamppb.New(dbUser.CreatedAt),
		UpdatedAt: timestamppb.New(dbUser.UpdatedAt),
	}
}

