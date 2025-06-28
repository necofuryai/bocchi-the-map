package grpc

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/necofuryai/bocchi-the-map/api/domain/entities"
	"github.com/necofuryai/bocchi-the-map/api/infrastructure/database"
	"github.com/necofuryai/bocchi-the-map/api/pkg/errors"
	"github.com/necofuryai/bocchi-the-map/api/pkg/converters"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserService implements the gRPC UserService
type UserService struct {
	queries       *database.Queries
	userConverter *converters.UserConverter
	grpcConverter *converters.GRPCConverter
}

// NewUserService creates a new UserService instance
func NewUserService(db *sql.DB) *UserService {
	return &UserService{
		queries:       database.New(db),
		userConverter: converters.NewUserConverter(),
		grpcConverter: converters.NewGRPCConverter(),
	}
}

// Use converters package types (will be replaced by protobuf generated types)
type (
	AuthProvider     = converters.GRPCAuthProvider
	UserPreferences  = converters.GRPCUserPreferences
	User             = converters.GRPCUser
)

// Re-export constants for compatibility
const (
	AuthProviderUnspecified = converters.GRPCAuthProviderUnspecified
	AuthProviderGoogle      = converters.GRPCAuthProviderGoogle
	AuthProviderX           = converters.GRPCAuthProviderX
)

type GetCurrentUserRequest struct {
	// Empty request - user context comes from authentication
}

type GetCurrentUserResponse struct {
	User *User
}

type GetUserByIDRequest struct {
	ID string
}

type GetUserByIDResponse struct {
	User *User
}

type UpdateUserPreferencesRequest struct {
	Preferences *UserPreferences
}

type UpdateUserPreferencesResponse struct {
	User *User
}

// Additional gRPC request/response types for user operations
type CreateUserRequest struct {
	Email          string
	DisplayName    string
	AvatarURL      string
	AuthProvider   AuthProvider
	AuthProviderID string
}

type CreateUserResponse struct {
	User *User
}

type UpdateUserRequest struct {
	ID             string
	Email          string
	DisplayName    string
	AvatarURL      string
	Preferences    *UserPreferences
}

type UpdateUserResponse struct {
	User *User
}

type GetUserByAuthProviderRequest struct {
	AuthProvider   AuthProvider
	AuthProviderID string
}

type GetUserByAuthProviderResponse struct {
	User *User
}

// GetCurrentUser retrieves the current authenticated user
func (s *UserService) GetCurrentUser(ctx context.Context, req *GetCurrentUserRequest) (*GetCurrentUserResponse, error) {
	// Extract user ID from authentication context
	userID := errors.GetUserID(ctx)
	if userID == "" {
		return nil, errors.GRPCUnauthenticated(ctx, "user not authenticated")
	}

	// Add operation context for error tracking
	ctx = errors.WithOperation(ctx, "get_current_user")
	ctx = errors.WithUserID(ctx, userID)

	// Get user from database
	dbUser, err := s.queries.GetUserByID(ctx, userID)
	if err != nil {
		return nil, errors.HandleDatabaseError(ctx, err, "get_user_by_id")
	}

	// Convert database user to gRPC response
	user, err := s.grpcConverter.ConvertDatabaseToGRPC(dbUser)
	if err != nil {
		return nil, errors.HandleGRPCError(ctx, err, "convert_database_to_grpc", "failed to convert user data")
	}
	return &GetCurrentUserResponse{User: user}, nil
}

// GetUserByAuthProvider retrieves a user by authentication provider and provider ID
func (s *UserService) GetUserByAuthProvider(ctx context.Context, authProvider entities.AuthProvider, providerID string) (*entities.User, error) {
	// Add operation context for error tracking
	ctx = errors.WithOperation(ctx, "get_user_by_auth_provider")

	// Convert domain auth provider to database enum using standardized converter
	dbAuthProvider, err := s.userConverter.AuthProviderToDatabase(authProvider)
	if err != nil {
		return nil, err
	}

	dbUser, err := s.queries.GetUserByProviderID(ctx, database.GetUserByProviderIDParams{
		AuthProvider:   dbAuthProvider,
		AuthProviderID: providerID,
	})
	if err != nil {
		return nil, errors.HandleDatabaseError(ctx, err, "get_user_by_provider_id")
	}

	// Convert database user to entity using standardized converter
	entity, err := s.userConverter.ConvertDatabaseToEntity(dbUser)
	if err != nil {
		return nil, errors.HandleGRPCError(ctx, err, "convert_database_to_entity", "failed to convert user data")
	}
	return entity, nil
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	// Add operation context for error tracking
	ctx = errors.WithOperation(ctx, "create_user")

	// Generate UUID for new user
	user.ID = uuid.New().String()

	// Convert entity to database parameters using standardized converter
	createParams, err := s.userConverter.ConvertEntityToDatabase(user)
	if err != nil {
		return nil, errors.HandleGRPCError(ctx, err, "convert_entity_to_database", "failed to prepare user data")
	}

	err = s.queries.CreateUser(ctx, createParams)
	if err != nil {
		return nil, errors.HandleDatabaseError(ctx, err, "create_user")
	}

	// Set timestamps
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	return user, nil
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	// Add operation context for error tracking
	ctx = errors.WithOperation(ctx, "update_user")
	ctx = errors.WithUserID(ctx, user.ID)

	// Convert entity to database upsert parameters using standardized converter
	upsertParams, err := s.userConverter.ConvertEntityToUpsertDatabase(user)
	if err != nil {
		return nil, errors.HandleGRPCError(ctx, err, "convert_entity_to_upsert_database", "failed to prepare user update data")
	}

	// Use upsert to update user
	err = s.queries.UpsertUser(ctx, upsertParams)
	if err != nil {
		return nil, errors.HandleDatabaseError(ctx, err, "upsert_user")
	}

	// Update timestamp
	user.UpdatedAt = time.Now()

	return user, nil
}



// UpdateUserPreferences updates user preferences
func (s *UserService) UpdateUserPreferences(ctx context.Context, req *UpdateUserPreferencesRequest) (*UpdateUserPreferencesResponse, error) {
	if req.Preferences == nil {
		return nil, errors.GRPCInvalidArgument(ctx, "preferences", "preferences are required")
	}

	// Extract user ID from authentication context
	userID := errors.GetUserID(ctx)
	if userID == "" {
		return nil, errors.GRPCUnauthenticated(ctx, "user not authenticated")
	}

	// Convert gRPC preferences to JSON
	prefsJSON, err := json.Marshal(map[string]interface{}{
		"language":  req.Preferences.Language,
		"dark_mode": req.Preferences.DarkMode,
		"timezone":  req.Preferences.Timezone,
	})
	if err != nil {
		return nil, errors.GRPCInternal(ctx, "failed to marshal preferences")
	}

	// Update preferences in database
	err = s.queries.UpdateUserPreferences(ctx, database.UpdateUserPreferencesParams{
		ID:          userID,
		Preferences: prefsJSON,
	})
	if err != nil {
		return nil, errors.GRPCInternal(ctx, "failed to update preferences")
	}

	// Get updated user from database
	dbUser, err := s.queries.GetUserByID(ctx, userID)
	if err != nil {
		return nil, errors.GRPCInternal(ctx, "failed to get updated user")
	}

	// Convert database user to gRPC response
	user, err := s.grpcConverter.ConvertDatabaseToGRPC(dbUser)
	if err != nil {
		return nil, errors.GRPCInternal(ctx, "failed to convert user data")
	}
	return &UpdateUserPreferencesResponse{User: user}, nil
}

// CreateUserGRPC creates a new user via gRPC interface
func (s *UserService) CreateUserGRPC(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
	// Add operation context for error tracking
	ctx = errors.WithOperation(ctx, "create_user_grpc")

	if req.Email == "" {
		return nil, errors.GRPCInvalidArgument(ctx, "email", "is required")
	}
	if req.DisplayName == "" {
		return nil, errors.GRPCInvalidArgument(ctx, "display_name", "is required")
	}
	if req.AuthProviderID == "" {
		return nil, errors.GRPCInvalidArgument(ctx, "auth_provider_id", "is required")
	}

	// Generate UUID for new user
	userID := uuid.New().String()

	// Convert gRPC auth provider to database enum
	var dbAuthProvider database.UsersAuthProvider
	switch req.AuthProvider {
	case AuthProviderGoogle:
		dbAuthProvider = database.UsersAuthProviderGoogle
	case AuthProviderX:
		dbAuthProvider = database.UsersAuthProviderTwitter
	default:
		return nil, errors.GRPCInvalidArgument(ctx, "auth_provider", "unsupported provider type")
	}

	// Create default preferences
	defaultPrefs := map[string]interface{}{
		"language":  "ja",
		"dark_mode": false,
		"timezone":  "Asia/Tokyo",
	}
	prefsJSON, err := json.Marshal(defaultPrefs)
	if err != nil {
		return nil, errors.GRPCInternal(ctx, "failed to marshal default preferences")
	}

	// Convert avatar URL to sql.NullString
	var avatarURL sql.NullString
	if req.AvatarURL != "" {
		avatarURL = sql.NullString{String: req.AvatarURL, Valid: true}
	}

	// Create user in database
	err = s.queries.CreateUser(ctx, database.CreateUserParams{
		ID:             userID,
		Email:          req.Email,
		DisplayName:    req.DisplayName,
		AvatarUrl:      avatarURL,
		AuthProvider:   dbAuthProvider,
		AuthProviderID: req.AuthProviderID,
		Preferences:    prefsJSON,
	})
	if err != nil {
		return nil, errors.HandleDatabaseError(ctx, err, "create_user")
	}

	// Get created user from database to get accurate timestamps
	dbUser, err := s.queries.GetUserByID(ctx, userID)
	if err != nil {
		return nil, errors.HandleDatabaseError(ctx, err, "get_created_user")
	}

	// Convert to gRPC response
	user, err := s.grpcConverter.ConvertDatabaseToGRPC(dbUser)
	if err != nil {
		return nil, errors.GRPCInternal(ctx, "failed to convert user data")
	}
	return &CreateUserResponse{User: user}, nil
}

// UpdateUserGRPC updates an existing user via gRPC interface
func (s *UserService) UpdateUserGRPC(ctx context.Context, req *UpdateUserRequest) (*UpdateUserResponse, error) {
	// Add operation context for error tracking
	ctx = errors.WithOperation(ctx, "update_user_grpc")
	ctx = errors.WithUserID(ctx, req.ID)

	if req.ID == "" {
		return nil, errors.GRPCInvalidArgument(ctx, "id", "is required")
	}
	if req.Email == "" {
		return nil, errors.GRPCInvalidArgument(ctx, "email", "is required")
	}

	// Get existing user to verify it exists and get auth provider info
	existingUser, err := s.queries.GetUserByID(ctx, req.ID)
	if err != nil {
		return nil, errors.HandleDatabaseError(ctx, err, "get_user_by_id")
	}

	// Convert avatar URL to sql.NullString
	var avatarURL sql.NullString
	if req.AvatarURL != "" {
		avatarURL = sql.NullString{String: req.AvatarURL, Valid: true}
	}

	// Prepare preferences JSON
	var prefsJSON json.RawMessage
	if req.Preferences != nil {
		prefsMap := map[string]interface{}{
			"language":  req.Preferences.Language,
			"dark_mode": req.Preferences.DarkMode,
			"timezone":  req.Preferences.Timezone,
		}
		prefsJSON, err = json.Marshal(prefsMap)
		if err != nil {
			return nil, errors.GRPCInternal(ctx, "failed to marshal preferences")
		}
	} else {
		// Keep existing preferences if not provided
		prefsJSON = existingUser.Preferences
	}

	// Use UpsertUser to update the user
	err = s.queries.UpsertUser(ctx, database.UpsertUserParams{
		ID:             req.ID,
		Email:          req.Email,
		DisplayName:    req.DisplayName,
		AvatarUrl:      avatarURL,
		AuthProvider:   existingUser.AuthProvider,   // Keep existing auth provider
		AuthProviderID: existingUser.AuthProviderID, // Keep existing auth provider ID
		Preferences:    prefsJSON,
	})
	if err != nil {
		return nil, errors.HandleDatabaseError(ctx, err, "upsert_user")
	}

	// Get updated user from database
	dbUser, err := s.queries.GetUserByID(ctx, req.ID)
	if err != nil {
		return nil, errors.HandleDatabaseError(ctx, err, "get_updated_user")
	}

	// Convert to gRPC response
	user, err := s.grpcConverter.ConvertDatabaseToGRPC(dbUser)
	if err != nil {
		return nil, errors.GRPCInternal(ctx, "failed to convert user data")
	}
	return &UpdateUserResponse{User: user}, nil
}

// GetUserByAuthProviderGRPC retrieves a user by authentication provider via gRPC interface
func (s *UserService) GetUserByAuthProviderGRPC(ctx context.Context, req *GetUserByAuthProviderRequest) (*GetUserByAuthProviderResponse, error) {
	// Add operation context for error tracking
	ctx = errors.WithOperation(ctx, "get_user_by_auth_provider_grpc")

	if req.AuthProviderID == "" {
		return nil, errors.GRPCInvalidArgument(ctx, "auth_provider_id", "is required")
	}

	// Convert gRPC auth provider to database enum
	var dbAuthProvider database.UsersAuthProvider
	switch req.AuthProvider {
	case AuthProviderGoogle:
		dbAuthProvider = database.UsersAuthProviderGoogle
	case AuthProviderX:
		dbAuthProvider = database.UsersAuthProviderTwitter
	default:
		return nil, errors.GRPCInvalidArgument(ctx, "auth_provider", "invalid provider type")
	}

	// Get user from database
	dbUser, err := s.queries.GetUserByProviderID(ctx, database.GetUserByProviderIDParams{
		AuthProvider:   dbAuthProvider,
		AuthProviderID: req.AuthProviderID,
	})
	if err != nil {
		return nil, errors.HandleDatabaseError(ctx, err, "get_user_by_provider_id")
	}

	// Convert to gRPC response
	user, err := s.grpcConverter.ConvertDatabaseToGRPC(dbUser)
	if err != nil {
		return nil, errors.GRPCInternal(ctx, "failed to convert user data")
	}
	return &GetUserByAuthProviderResponse{User: user}, nil
}

// GetUserByID retrieves a user by ID via gRPC interface
func (s *UserService) GetUserByID(ctx context.Context, req *GetUserByIDRequest) (*GetUserByIDResponse, error) {
	// Add operation context for error tracking
	ctx = errors.WithOperation(ctx, "get_user_by_id_grpc")
	ctx = errors.WithUserID(ctx, req.ID)

	if req.ID == "" {
		return nil, errors.GRPCInvalidArgument(ctx, "id", "is required")
	}

	// Get user from database
	dbUser, err := s.queries.GetUserByID(ctx, req.ID)
	if err != nil {
		return nil, errors.HandleDatabaseError(ctx, err, "get_user_by_id")
	}

	// Convert to gRPC response
	user, err := s.grpcConverter.ConvertDatabaseToGRPC(dbUser)
	if err != nil {
		return nil, errors.GRPCInternal(ctx, "failed to convert user data")
	}
	return &GetUserByIDResponse{User: user}, nil
}