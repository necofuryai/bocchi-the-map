package grpc

import (
	"context"
	"database/sql"
	"encoding/json"

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

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
	// Validate request
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if req.GetDisplayName() == "" {
		return nil, status.Error(codes.InvalidArgument, "display name is required")
	}
	if req.GetAuthProvider() == "" {
		return nil, status.Error(codes.InvalidArgument, "auth provider is required")
	}
	if req.GetAuthProviderId() == "" {
		return nil, status.Error(codes.InvalidArgument, "auth provider ID is required")
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

	// Convert preferences to JSON (for proto string field)
	var preferencesJSON []byte
	if req.GetPreferences() != "" {
		// Validate that it's valid JSON
		var temp interface{}
		if err := json.Unmarshal([]byte(req.GetPreferences()), &temp); err != nil {
			return nil, status.Error(codes.InvalidArgument, "preferences must be valid JSON")
		}
		preferencesJSON = []byte(req.GetPreferences())
	} else {
		preferencesJSON = []byte("{}")
	}

	// Convert avatar URL to nullable string
	var avatarUrl sql.NullString
	if req.GetAvatarUrl() != "" {
		avatarUrl = sql.NullString{String: req.GetAvatarUrl(), Valid: true}
	}

	// Validate auth provider
	var authProvider string
	switch req.GetAuthProvider() {
	case "google":
		authProvider = "google"
	case "twitter":
		authProvider = "twitter"
	case "x":
		authProvider = "x"
	default:
		return nil, status.Error(codes.InvalidArgument, "invalid auth provider")
	}

	// Create user in database
	err = s.queries.CreateUser(ctx, database.CreateUserParams{
		ID:            userID,
		Email:         req.GetEmail(),
		Name:          sql.NullString{String: req.GetDisplayName(), Valid: true},
		Nickname:      sql.NullString{String: req.GetDisplayName(), Valid: true},
		Picture:       avatarUrl,
		Provider:      authProvider,
		ProviderID:    req.GetAuthProviderId(),
		EmailVerified: true,
		Preferences:   preferencesJSON,
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



// convertDatabaseUserToGRPC converts database user model to gRPC user struct
func (s *UserService) convertDatabaseUserToGRPC(dbUser database.User) *User {
	// Convert preferences to JSON string for proto
	var preferencesStr string
	if len(dbUser.Preferences) > 0 {
		// Validate JSON and pretty-print if needed
		var temp interface{}
		if err := json.Unmarshal(dbUser.Preferences, &temp); err != nil {
			logger.Error("Failed to unmarshal user preferences", err)
			preferencesStr = "{}"
		} else {
			preferencesStr = string(dbUser.Preferences)
		}
	} else {
		preferencesStr = "{}"
	}

	return &User{
		Id:             dbUser.ID,
		Email:          dbUser.Email,
		DisplayName:    dbUser.Name.String,
		AvatarUrl:      dbUser.Picture.String,
		AuthProvider:   dbUser.Provider,
		AuthProviderId: dbUser.ProviderID,
		Preferences:    preferencesStr,
		CreatedAt:      timestamppb.New(dbUser.CreatedAt),
		UpdatedAt:      timestamppb.New(dbUser.UpdatedAt),
	}
}

// UpsertUserFromAuth creates or updates a user from authentication provider
func (s *UserService) UpsertUserFromAuth(ctx context.Context, email, displayName, avatarUrl, authProvider, authProviderID string) (*User, error) {
	// Generate UUID for new user (will be ignored if user exists)
	userID := uuid.New().String()

	// Convert preferences to JSON
	defaultPreferences := map[string]interface{}{
		"language": "en",
		"theme":    "auto",
	}
	preferencesJSON, _ := json.Marshal(defaultPreferences)

	// Convert avatar URL to nullable string
	var avatarUrlNullable sql.NullString
	if avatarUrl != "" {
		avatarUrlNullable = sql.NullString{String: avatarUrl, Valid: true}
	}

	// Validate auth provider
	var dbAuthProvider string
	switch authProvider {
	case "google":
		dbAuthProvider = "google"
	case "twitter":
		dbAuthProvider = "twitter"
	case "x":
		dbAuthProvider = "x"
	default:
		return nil, status.Error(codes.InvalidArgument, "invalid auth provider")
	}

	// Upsert user in database
	err := s.queries.UpsertUser(ctx, database.UpsertUserParams{
		ID:            userID,
		Email:         email,
		Name:          sql.NullString{String: displayName, Valid: true},
		Nickname:      sql.NullString{String: displayName, Valid: true},
		Picture:       avatarUrlNullable,
		Provider:      dbAuthProvider,
		ProviderID:    authProviderID,
		EmailVerified: true,
		Preferences:   preferencesJSON,
	})
	if err != nil {
		logger.ErrorWithContext(ctx, "Failed to upsert user", err)
		return nil, status.Error(codes.Internal, "failed to upsert user")
	}

	// Retrieve the upserted user
	dbUser, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		logger.ErrorWithContext(ctx, "Failed to retrieve upserted user", err)
		return nil, status.Error(codes.Internal, "failed to retrieve upserted user")
	}

	// Convert database user to gRPC response
	return s.convertDatabaseUserToGRPC(dbUser), nil
}