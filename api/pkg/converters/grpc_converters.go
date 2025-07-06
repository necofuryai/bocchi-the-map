package converters

import (
	"encoding/json"

	"google.golang.org/protobuf/types/known/timestamppb"

	"bocchi/api/domain/entities"
	"bocchi/api/gen/user/v1"
	"bocchi/api/infrastructure/database"
	"bocchi/api/pkg/errors"
)

// Default user preference values
const (
	DefaultLanguage = "ja"
	DefaultTimezone = "Asia/Tokyo"
	DefaultDarkMode = false
)

// Use Protocol Buffers generated types
type GRPCUser = userv1.User

// GRPCConverter handles conversions between gRPC types and other representations
type GRPCConverter struct{}

// NewGRPCConverter creates a new GRPCConverter instance
func NewGRPCConverter() *GRPCConverter {
	return &GRPCConverter{}
}

// ConvertEntityToGRPC converts domain entity to gRPC type
func (c *GRPCConverter) ConvertEntityToGRPC(entity *entities.User) *GRPCUser {
	// Convert preferences to JSON string
	prefsBytes, _ := json.Marshal(entity.Preferences)
	
	return &GRPCUser{
		Id:             entity.ID,
		Email:          entity.Email,
		DisplayName:    entity.DisplayName,
		AvatarUrl:      entity.AvatarUrl,
		AuthProvider:   c.ConvertEntityAuthProviderToGRPC(entity.AuthProvider),
		AuthProviderId: entity.AuthProviderID,
		Preferences:    string(prefsBytes),
		CreatedAt:      timestamppb.New(entity.CreatedAt),
		UpdatedAt:      timestamppb.New(entity.UpdatedAt),
	}
}

// ConvertGRPCToEntity converts gRPC type to domain entity
func (c *GRPCConverter) ConvertGRPCToEntity(grpcUser *GRPCUser) (*entities.User, error) {
	authProvider, err := c.ConvertGRPCAuthProviderToEntity(grpcUser.AuthProvider)
	if err != nil {
		return nil, err
	}

	// Convert preferences from JSON string
	prefs := c.ConvertGRPCPreferencesToEntity(grpcUser.Preferences)

	return &entities.User{
		ID:             grpcUser.Id,
		Email:          grpcUser.Email,
		DisplayName:    grpcUser.DisplayName,
		AvatarUrl:      grpcUser.AvatarUrl,
		AuthProvider:   authProvider,
		AuthProviderID: grpcUser.AuthProviderId,
		Preferences:    prefs,
		CreatedAt:      grpcUser.CreatedAt.AsTime(),
		UpdatedAt:      grpcUser.UpdatedAt.AsTime(),
	}, nil
}

// ConvertDatabaseToGRPC converts database type directly to gRPC type
func (c *GRPCConverter) ConvertDatabaseToGRPC(dbUser database.User) (*GRPCUser, error) {
	// Convert preferences to JSON string for protobuf
	var preferencesJSON string
	if len(dbUser.Preferences) > 0 {
		preferencesJSON = string(dbUser.Preferences)
	} else {
		// Default preferences
		defaultPrefs := map[string]interface{}{
			"language": DefaultLanguage,
			"darkMode": DefaultDarkMode,
			"timezone": DefaultTimezone,
		}
		prefsBytes, _ := json.Marshal(defaultPrefs)
		preferencesJSON = string(prefsBytes)
	}

	// Get display name
	var displayName string
	if dbUser.Name.Valid {
		displayName = dbUser.Name.String
	}

	// Handle nullable avatar URL
	var avatarURL string
	if dbUser.Picture.Valid {
		avatarURL = dbUser.Picture.String
	}

	return &GRPCUser{
		Id:             dbUser.ID,
		Email:          dbUser.Email,
		DisplayName:    displayName,
		AvatarUrl:      avatarURL,
		AuthProvider:   dbUser.Provider,
		AuthProviderId: dbUser.ProviderID,
		Preferences:    preferencesJSON,
		CreatedAt:      timestamppb.New(dbUser.CreatedAt),
		UpdatedAt:      timestamppb.New(dbUser.UpdatedAt),
	}, nil
}

// ConvertGRPCAuthProviderToEntity converts gRPC auth provider string to domain entity
func (c *GRPCConverter) ConvertGRPCAuthProviderToEntity(provider string) (entities.AuthProvider, error) {
	switch provider {
	case "google":
		return entities.AuthProviderGoogle, nil
	case "twitter", "x":
		return entities.AuthProviderX, nil
	default:
		return "", errors.InvalidInput("auth_provider", "unknown gRPC auth provider").
			WithField("grpc_provider", provider)
	}
}

// ConvertEntityAuthProviderToGRPC converts domain entity auth provider to gRPC string
func (c *GRPCConverter) ConvertEntityAuthProviderToGRPC(provider entities.AuthProvider) string {
	switch provider {
	case entities.AuthProviderGoogle:
		return "google"
	case entities.AuthProviderX:
		return "x"
	default:
		return ""
	}
}

// ConvertGRPCPreferencesToEntity converts gRPC preferences JSON string to domain entity
func (c *GRPCConverter) ConvertGRPCPreferencesToEntity(prefsJSON string) entities.UserPreferences {
	if prefsJSON == "" {
		return entities.UserPreferences{
			Language: DefaultLanguage,
			DarkMode: DefaultDarkMode,
			Timezone: DefaultTimezone,
		}
	}

	var prefs entities.UserPreferences
	if err := json.Unmarshal([]byte(prefsJSON), &prefs); err != nil {
		// Return defaults if parsing fails
		return entities.UserPreferences{
			Language: DefaultLanguage,
			DarkMode: DefaultDarkMode,
			Timezone: DefaultTimezone,
		}
	}

	return prefs
}

// ConvertEntityPreferencesToGRPC converts domain entity preferences to gRPC JSON string
func (c *GRPCConverter) ConvertEntityPreferencesToGRPC(entityPrefs entities.UserPreferences) string {
	prefsBytes, _ := json.Marshal(entityPrefs)
	return string(prefsBytes)
}