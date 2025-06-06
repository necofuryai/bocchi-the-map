package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserService implements the gRPC UserService
type UserService struct {
	// TODO: Add dependencies like repository interfaces
}

// NewUserService creates a new UserService instance
func NewUserService() *UserService {
	return &UserService{}
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