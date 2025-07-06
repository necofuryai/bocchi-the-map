package grpc

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/necofuryai/bocchi-the-map/api/infrastructure/database"
	"github.com/necofuryai/bocchi-the-map/api/pkg/errors"
	"github.com/necofuryai/bocchi-the-map/api/pkg/logger"
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

// Temporary structs until protobuf generates them
type User struct {
	ID             string            `json:"id"`
	Email          string            `json:"email"`
	DisplayName    string            `json:"display_name"`
	AvatarUrl      string            `json:"avatar_url,omitempty"`
	AuthProvider   string            `json:"auth_provider"`
	AuthProviderID string            `json:"auth_provider_id"`
	Preferences    map[string]interface{} `json:"preferences,omitempty"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

type GetUserRequest struct {
	ID string `json:"id"`
}

type GetUserResponse struct {
	User *User `json:"user"`
}

type GetUserByEmailRequest struct {
	Email string `json:"email"`
}

type GetUserByEmailResponse struct {
	User *User `json:"user"`
}

type CreateUserRequest struct {
	Email          string                 `json:"email"`
	DisplayName    string                 `json:"display_name"`
	AvatarUrl      string                 `json:"avatar_url,omitempty"`
	AuthProvider   string                 `json:"auth_provider"`
	AuthProviderID string                 `json:"auth_provider_id"`
	Preferences    map[string]interface{} `json:"preferences,omitempty"`
}

type CreateUserResponse struct {
	User *User `json:"user"`
}

type UpdateUserRequest struct {
	ID          string                 `json:"id"`
	DisplayName string                 `json:"display_name,omitempty"`
	AvatarUrl   string                 `json:"avatar_url,omitempty"`
	Preferences map[string]interface{} `json:"preferences,omitempty"`
}

type UpdateUserResponse struct {
	User *User `json:"user"`
}

type DeleteUserRequest struct {
	ID string `json:"id"`
}

type DeleteUserResponse struct {
	Success bool `json:"success"`
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error) {
	if req.ID == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	// Get user from database
	dbUser, err := s.queries.GetUserByID(ctx, req.ID)
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
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	// Get user from database
	dbUser, err := s.queries.GetUserByEmail(ctx, req.Email)
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
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if req.DisplayName == "" {
		return nil, status.Error(codes.InvalidArgument, "display name is required")
	}
	if req.AuthProvider == "" {
		return nil, status.Error(codes.InvalidArgument, "auth provider is required")
	}
	if req.AuthProviderID == "" {
		return nil, status.Error(codes.InvalidArgument, "auth provider ID is required")
	}

	// Check if user already exists
	_, err := s.queries.GetUserByEmail(ctx, req.Email)
	if err == nil {
		return nil, status.Error(codes.AlreadyExists, "user with this email already exists")
	} else if err != sql.ErrNoRows {
		logger.ErrorWithContext(ctx, "Failed to check existing user", err)
		return nil, status.Error(codes.Internal, "failed to check existing user")
	}

	// Generate UUID for new user
	userID := uuid.New().String()

	// Convert preferences to JSON
	preferencesJSON, err := json.Marshal(req.Preferences)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to marshal preferences")
	}
	if req.Preferences == nil {
		preferencesJSON = []byte("{}")
	}

	// Convert avatar URL to nullable string
	var avatarUrl sql.NullString
	if req.AvatarUrl != "" {
		avatarUrl = sql.NullString{String: req.AvatarUrl, Valid: true}
	}

	// Validate auth provider
	var authProvider string
	switch req.AuthProvider {
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
		Email:         req.Email,
		Name:          sql.NullString{String: req.DisplayName, Valid: true},
		Nickname:      sql.NullString{String: req.DisplayName, Valid: true},
		Picture:       avatarUrl,
		Provider:      authProvider,
		ProviderID:    req.AuthProviderID,
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

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(ctx context.Context, req *UpdateUserRequest) (*UpdateUserResponse, error) {
	if req.ID == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	// Get current user to verify it exists
	_, err := s.queries.GetUserByID(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		logger.ErrorWithContext(ctx, "Failed to get user for update", err)
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	// Update avatar if provided
	if req.AvatarUrl != "" {
		var avatarUrl sql.NullString
		if req.AvatarUrl != "" {
			avatarUrl = sql.NullString{String: req.AvatarUrl, Valid: true}
		}
		
		err = s.queries.UpdateUserAvatar(ctx, database.UpdateUserAvatarParams{
			ID:      req.ID,
			Picture: avatarUrl,
		})
		if err != nil {
			logger.ErrorWithContext(ctx, "Failed to update user avatar", err)
			return nil, status.Error(codes.Internal, "failed to update user avatar")
		}
	}

	// Update preferences if provided
	if req.Preferences != nil {
		preferencesJSON, err := json.Marshal(req.Preferences)
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to marshal preferences")
		}

		err = s.queries.UpdateUserPreferences(ctx, database.UpdateUserPreferencesParams{
			ID:          req.ID,
			Preferences: preferencesJSON,
		})
		if err != nil {
			logger.ErrorWithContext(ctx, "Failed to update user preferences", err)
			return nil, status.Error(codes.Internal, "failed to update user preferences")
		}
	}

	// Retrieve the updated user
	dbUser, err := s.queries.GetUserByID(ctx, req.ID)
	if err != nil {
		logger.ErrorWithContext(ctx, "Failed to retrieve updated user", err)
		return nil, status.Error(codes.Internal, "failed to retrieve updated user")
	}

	// Convert database user to gRPC response
	user := s.convertDatabaseUserToGRPC(dbUser)
	return &UpdateUserResponse{User: user}, nil
}

// DeleteUser deletes a user (hard deletion)
func (s *UserService) DeleteUser(ctx context.Context, req *DeleteUserRequest) (*DeleteUserResponse, error) {
	if req.ID == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	// Extract authenticated user ID from context
	authUserID := errors.GetUserID(ctx)
	if authUserID == "" {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	// Check if user is trying to delete themselves or has admin permissions
	if authUserID != req.ID {
		// TODO: Add admin permission check here
		return nil, status.Error(codes.PermissionDenied, "insufficient permissions to delete user")
	}

	// Check if user exists
	_, err := s.queries.GetUserByID(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		logger.ErrorWithContext(ctx, "Failed to get user for deletion", err)
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	// Delete user (CASCADE will handle related reviews)
	err = s.queries.DeleteUser(ctx, req.ID)
	if err != nil {
		logger.ErrorWithContext(ctx, "Failed to delete user", err)
		return nil, status.Error(codes.Internal, "failed to delete user")
	}

	logger.InfoWithFields("User deleted successfully", map[string]interface{}{
		"user_id":      req.ID,
		"auth_user_id": authUserID,
	})

	return &DeleteUserResponse{Success: true}, nil
}

// convertDatabaseUserToGRPC converts database user model to gRPC user struct
func (s *UserService) convertDatabaseUserToGRPC(dbUser database.User) *User {
	// Parse preferences JSON
	var preferences map[string]interface{}
	if len(dbUser.Preferences) > 0 {
		if err := json.Unmarshal(dbUser.Preferences, &preferences); err != nil {
			logger.Error("Failed to unmarshal user preferences", err)
			preferences = make(map[string]interface{})
		}
	}

	return &User{
		ID:             dbUser.ID,
		Email:          dbUser.Email,
		DisplayName:    dbUser.Name.String,
		AvatarUrl:      dbUser.Picture.String,
		AuthProvider:   dbUser.Provider,
		AuthProviderID: dbUser.ProviderID,
		Preferences:    preferences,
		CreatedAt:      dbUser.CreatedAt,
		UpdatedAt:      dbUser.UpdatedAt,
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