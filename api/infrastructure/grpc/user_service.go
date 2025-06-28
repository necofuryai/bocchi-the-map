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

// convertDatabaseUserToEntity converts database user model to domain entity
func (s *UserService) convertDatabaseUserToEntity(dbUser database.User) *entities.User {
	var prefs entities.UserPreferences
	if err := json.Unmarshal(dbUser.Preferences, &prefs); err != nil {
		// Set default preferences if JSON unmarshal fails
		prefs = entities.UserPreferences{
			Language: "ja",
			DarkMode: false,
			Timezone: "Asia/Tokyo",
		}
	}

	var authProvider entities.AuthProvider
	switch dbUser.AuthProvider {
	case database.UsersAuthProviderGoogle:
		authProvider = entities.AuthProviderGoogle
	case database.UsersAuthProviderTwitter:
		authProvider = entities.AuthProviderX
	}

	user := &entities.User{
		ID:             dbUser.ID,
		Email:          dbUser.Email,
		DisplayName:    dbUser.DisplayName,
		AuthProvider:   authProvider,
		AuthProviderID: dbUser.AuthProviderID,
		Preferences:    prefs,
		CreatedAt:      dbUser.CreatedAt,
		UpdatedAt:      dbUser.UpdatedAt,
	}

	if dbUser.AnonymousID.Valid {
		user.AnonymousID = dbUser.AnonymousID.String
	}

	if dbUser.AvatarUrl.Valid {
		user.AvatarURL = dbUser.AvatarUrl.String
	}

	return user
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
	if req.Email == "" {
		return nil, errors.InvalidInput("email", "is required").ToGRPCError()
	}
	if req.DisplayName == "" {
		return nil, errors.InvalidInput("display_name", "is required").ToGRPCError()
	}
	if req.AuthProviderID == "" {
		return nil, errors.InvalidInput("auth_provider_id", "is required").ToGRPCError()
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
		return nil, errors.InvalidInput("auth_provider", "unsupported provider type").ToGRPCError()
	}

	// Create default preferences
	defaultPrefs := map[string]interface{}{
		"language":  "ja",
		"dark_mode": false,
		"timezone":  "Asia/Tokyo",
	}
	prefsJSON, err := json.Marshal(defaultPrefs)
	if err != nil {
		return nil, errors.Internal("failed to marshal default preferences").ToGRPCError()
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
		return nil, errors.Database("create user", err).ToGRPCError()
	}

	// Get created user from database to get accurate timestamps
	dbUser, err := s.queries.GetUserByID(ctx, userID)
	if err != nil {
		return nil, errors.Database("get created user", err).ToGRPCError()
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
	if req.ID == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	// TODO: Implement UpdateUser method in database queries
	// For now, return error to indicate this functionality is not yet implemented
	_ = req // Suppress unused parameter warning
	return nil, status.Error(codes.Unimplemented, "user update functionality not yet implemented")
}

// GetUserByAuthProviderGRPC retrieves a user by authentication provider via gRPC interface
func (s *UserService) GetUserByAuthProviderGRPC(ctx context.Context, req *GetUserByAuthProviderRequest) (*GetUserByAuthProviderResponse, error) {
	if req.AuthProviderID == "" {
		return nil, status.Error(codes.InvalidArgument, "auth provider ID is required")
	}

	// Convert gRPC auth provider to database enum
	var dbAuthProvider database.UsersAuthProvider
	switch req.AuthProvider {
	case AuthProviderGoogle:
		dbAuthProvider = database.UsersAuthProviderGoogle
	case AuthProviderX:
		dbAuthProvider = database.UsersAuthProviderTwitter
	default:
		return nil, status.Error(codes.InvalidArgument, "invalid auth provider")
	}

	// Get user from database
	dbUser, err := s.queries.GetUserByProviderID(ctx, database.GetUserByProviderIDParams{
		AuthProvider:   dbAuthProvider,
		AuthProviderID: req.AuthProviderID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "failed to get user")
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
	if req.ID == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	// Add operation context for error tracking
	ctx = errors.WithOperation(ctx, "get_user_by_id_grpc")
	ctx = errors.WithUserID(ctx, req.ID)

	// Get user from database
	dbUser, err := s.queries.GetUserByID(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, errors.HandleDatabaseError(ctx, err, "get_user_by_id")
	}

	// Convert to gRPC response
	user, err := s.grpcConverter.ConvertDatabaseToGRPC(dbUser)
	if err != nil {
		return nil, errors.GRPCInternal(ctx, "failed to convert user data")
	}
	return &GetUserByIDResponse{User: user}, nil
}