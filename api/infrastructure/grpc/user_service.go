package grpc

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/necofuryai/bocchi-the-map/api/domain/entities"
	"github.com/necofuryai/bocchi-the-map/api/infrastructure/database"
	"github.com/google/uuid"
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
type AuthProvider int32

const (
	AuthProviderUnspecified AuthProvider = 0
	AuthProviderGoogle      AuthProvider = 1
	AuthProviderX           AuthProvider = 2
)

type UserPreferences struct {
	Language string // Simplified language handling
	DarkMode bool
	Timezone string
}

type User struct {
	ID             string
	AnonymousID    string
	Email          string
	DisplayName    string
	AvatarURL      string
	AuthProvider   AuthProvider
	AuthProviderID string
	Preferences    *UserPreferences
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type GetCurrentUserRequest struct {
	// Empty request - user context comes from authentication
}

type GetCurrentUserResponse struct {
	User *User
}

type UpdateUserPreferencesRequest struct {
	Preferences *UserPreferences
}

type UpdateUserPreferencesResponse struct {
	User *User
}

// GetCurrentUser retrieves the current authenticated user
func (s *UserService) GetCurrentUser(ctx context.Context, req *GetCurrentUserRequest) (*GetCurrentUserResponse, error) {
	// TODO: Extract user ID from context/authentication
	// For now, return dummy data
	user := &User{
		ID:             "user_123",
		AnonymousID:    "anon_456",
		Email:          "user@example.com",
		DisplayName:    "Example User",
		AvatarURL:      "https://example.com/avatar.jpg",
		AuthProvider:   AuthProviderGoogle,
		AuthProviderID: "google_123",
		Preferences: &UserPreferences{
			Language: "ja",
			DarkMode: false,
			Timezone: "Asia/Tokyo",
		},
		CreatedAt: time.Now().Add(-30 * 24 * time.Hour),
		UpdatedAt: time.Now(),
	}

	return &GetCurrentUserResponse{User: user}, nil
}

// GetUserByAuthProvider retrieves a user by authentication provider and provider ID
func (s *UserService) GetUserByAuthProvider(ctx context.Context, authProvider entities.AuthProvider, providerID string) (*entities.User, error) {
	// Convert domain auth provider to database enum
	var dbAuthProvider database.UsersAuthProvider
	switch authProvider {
	case entities.AuthProviderGoogle:
		dbAuthProvider = database.UsersAuthProviderGoogle
	case entities.AuthProviderX:
		dbAuthProvider = database.UsersAuthProviderX
	default:
		return nil, errors.New("invalid auth provider")
	}

	dbUser, err := s.queries.GetUserByProviderID(ctx, database.GetUserByProviderIDParams{
		AuthProvider:   dbAuthProvider,
		AuthProviderID: providerID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return s.convertDatabaseUserToEntity(dbUser), nil
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	// Generate UUID for new user
	user.ID = uuid.New().String()

	// Convert preferences to JSON
	prefsJSON, err := json.Marshal(user.Preferences)
	if err != nil {
		return nil, err
	}

	// Convert domain auth provider to database enum
	var dbAuthProvider database.UsersAuthProvider
	switch user.AuthProvider {
	case entities.AuthProviderGoogle:
		dbAuthProvider = database.UsersAuthProviderGoogle
	case entities.AuthProviderX:
		dbAuthProvider = database.UsersAuthProviderX
	default:
		return nil, errors.New("invalid auth provider")
	}

	// Convert avatar URL to sql.NullString
	var avatarURL sql.NullString
	if user.AvatarURL != "" {
		avatarURL = sql.NullString{String: user.AvatarURL, Valid: true}
	}

	err = s.queries.CreateUser(ctx, database.CreateUserParams{
		ID:             user.ID,
		Email:          user.Email,
		DisplayName:    user.DisplayName,
		AvatarUrl:      avatarURL,
		AuthProvider:   dbAuthProvider,
		AuthProviderID: user.AuthProviderID,
		Preferences:    prefsJSON,
	})
	if err != nil {
		return nil, err
	}

	// Set timestamps
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	return user, nil
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	// Convert preferences to JSON
	prefsJSON, err := json.Marshal(user.Preferences)
	if err != nil {
		return nil, err
	}

	// Convert domain auth provider to database enum
	var dbAuthProvider database.UsersAuthProvider
	switch user.AuthProvider {
	case entities.AuthProviderGoogle:
		dbAuthProvider = database.UsersAuthProviderGoogle
	case entities.AuthProviderX:
		dbAuthProvider = database.UsersAuthProviderX
	default:
		return nil, errors.New("invalid auth provider")
	}

	// Convert avatar URL to sql.NullString
	var avatarURL sql.NullString
	if user.AvatarURL != "" {
		avatarURL = sql.NullString{String: user.AvatarURL, Valid: true}
	}

	// Use upsert to update user
	err = s.queries.UpsertUser(ctx, database.UpsertUserParams{
		ID:             user.ID,
		Email:          user.Email,
		DisplayName:    user.DisplayName,
		AvatarUrl:      avatarURL,
		AuthProvider:   dbAuthProvider,
		AuthProviderID: user.AuthProviderID,
		Preferences:    prefsJSON,
	})
	if err != nil {
		return nil, err
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
	case database.UsersAuthProviderX:
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
		return nil, status.Error(codes.InvalidArgument, "preferences are required")
	}

	// TODO: Extract user ID from context and update preferences in repository
	// For now, return dummy updated user
	user := &User{
		ID:             "user_123",
		AnonymousID:    "anon_456",
		Email:          "user@example.com",
		DisplayName:    "Example User",
		AvatarURL:      "https://example.com/avatar.jpg",
		AuthProvider:   AuthProviderGoogle,
		AuthProviderID: "google_123",
		Preferences:    req.Preferences,
		CreatedAt:      time.Now().Add(-30 * 24 * time.Hour),
		UpdatedAt:      time.Now(),
	}

	return &UpdateUserPreferencesResponse{User: user}, nil
}