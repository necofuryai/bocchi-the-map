package converters

import (
	"encoding/json"
	"time"

	"github.com/necofuryai/bocchi-the-map/api/domain/entities"
	"github.com/necofuryai/bocchi-the-map/api/infrastructure/database"
	"github.com/necofuryai/bocchi-the-map/api/pkg/errors"
)

// Default user preference values
const (
	DefaultLanguage = "ja"
	DefaultTimezone = "Asia/Tokyo"
	DefaultDarkMode = false
)

// Temporary gRPC structs until protobuf generates them
// These will be replaced once protobuf code generation is working
type GRPCAuthProvider int32

const (
	GRPCAuthProviderUnspecified GRPCAuthProvider = 0
	GRPCAuthProviderGoogle      GRPCAuthProvider = 1
	GRPCAuthProviderX           GRPCAuthProvider = 2
)

type GRPCUserPreferences struct {
	Language string
	DarkMode bool
	Timezone string
}

type GRPCUser struct {
	ID             string
	AnonymousID    string
	Email          string
	DisplayName    string
	AvatarURL      string
	AuthProvider   GRPCAuthProvider
	AuthProviderID string
	Preferences    *GRPCUserPreferences
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// GRPCConverter handles conversions between gRPC types and other representations
type GRPCConverter struct{}

// NewGRPCConverter creates a new GRPCConverter instance
func NewGRPCConverter() *GRPCConverter {
	return &GRPCConverter{}
}

// ConvertEntityToGRPC converts domain entity to gRPC type
func (c *GRPCConverter) ConvertEntityToGRPC(entity *entities.User) *GRPCUser {
	var grpcAuthProvider GRPCAuthProvider
	switch entity.AuthProvider {
	case entities.AuthProviderGoogle:
		grpcAuthProvider = GRPCAuthProviderGoogle
	case entities.AuthProviderX:
		grpcAuthProvider = GRPCAuthProviderX
	default:
		grpcAuthProvider = GRPCAuthProviderUnspecified
	}

	return &GRPCUser{
		ID:             entity.ID,
		AnonymousID:    entity.AnonymousID,
		Email:          entity.Email,
		DisplayName:    entity.DisplayName,
		AvatarURL:      entity.AvatarURL,
		AuthProvider:   grpcAuthProvider,
		AuthProviderID: entity.AuthProviderID,
		Preferences: &GRPCUserPreferences{
			Language: entity.Preferences.Language,
			DarkMode: entity.Preferences.DarkMode,
			Timezone: entity.Preferences.Timezone,
		},
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

// ConvertGRPCToEntity converts gRPC type to domain entity
func (c *GRPCConverter) ConvertGRPCToEntity(grpcUser *GRPCUser) (*entities.User, error) {
	var authProvider entities.AuthProvider
	switch grpcUser.AuthProvider {
	case GRPCAuthProviderGoogle:
		authProvider = entities.AuthProviderGoogle
	case GRPCAuthProviderX:
		authProvider = entities.AuthProviderX
	default:
		return nil, errors.InvalidInput("auth_provider", "unknown gRPC auth provider").
			WithField("grpc_provider", int32(grpcUser.AuthProvider)).
			WithField("user_id", grpcUser.ID)
	}

	var prefs entities.UserPreferences
	if grpcUser.Preferences != nil {
		prefs = entities.UserPreferences{
			Language: grpcUser.Preferences.Language,
			DarkMode: grpcUser.Preferences.DarkMode,
			Timezone: grpcUser.Preferences.Timezone,
		}
	} else {
		// Default preferences
		prefs = entities.UserPreferences{
			Language: DefaultLanguage,
			DarkMode: DefaultDarkMode,
			Timezone: DefaultTimezone,
		}
	}

	return &entities.User{
		ID:             grpcUser.ID,
		AnonymousID:    grpcUser.AnonymousID,
		Email:          grpcUser.Email,
		DisplayName:    grpcUser.DisplayName,
		AvatarURL:      grpcUser.AvatarURL,
		AuthProvider:   authProvider,
		AuthProviderID: grpcUser.AuthProviderID,
		Preferences:    prefs,
		CreatedAt:      grpcUser.CreatedAt,
		UpdatedAt:      grpcUser.UpdatedAt,
	}, nil
}

// ConvertDatabaseToGRPC converts database type directly to gRPC type
func (c *GRPCConverter) ConvertDatabaseToGRPC(dbUser database.User) (*GRPCUser, error) {
	var prefs GRPCUserPreferences
	if len(dbUser.Preferences) > 0 {
		var entityPrefs entities.UserPreferences
		if err := json.Unmarshal(dbUser.Preferences, &entityPrefs); err != nil {
			return nil, errors.Wrap(err, errors.ErrTypeInternal, "failed to unmarshal user preferences").
				WithField("user_id", dbUser.ID)
		}
		prefs = GRPCUserPreferences{
			Language: entityPrefs.Language,
			DarkMode: entityPrefs.DarkMode,
			Timezone: entityPrefs.Timezone,
		}
	} else {
		// Default preferences
		prefs = GRPCUserPreferences{
			Language: DefaultLanguage,
			DarkMode: DefaultDarkMode,
			Timezone: DefaultTimezone,
		}
	}

	// Convert database auth provider to gRPC enum
	var grpcAuthProvider GRPCAuthProvider
	switch dbUser.AuthProvider {
	case database.UsersAuthProviderGoogle:
		grpcAuthProvider = GRPCAuthProviderGoogle
	case database.UsersAuthProviderTwitter:
		grpcAuthProvider = GRPCAuthProviderX
	default:
		return nil, errors.InvalidInput("auth_provider", "unknown database auth provider").
			WithField("db_provider", string(dbUser.AuthProvider)).
			WithField("user_id", dbUser.ID)
	}

	// Handle nullable avatar URL
	var avatarURL string
	if dbUser.AvatarUrl.Valid {
		avatarURL = dbUser.AvatarUrl.String
	}

	// Handle nullable anonymous ID
	var anonymousID string
	if dbUser.AnonymousID.Valid {
		anonymousID = dbUser.AnonymousID.String
	}

	return &GRPCUser{
		ID:             dbUser.ID,
		AnonymousID:    anonymousID,
		Email:          dbUser.Email,
		DisplayName:    dbUser.DisplayName,
		AvatarURL:      avatarURL,
		AuthProvider:   grpcAuthProvider,
		AuthProviderID: dbUser.AuthProviderID,
		Preferences:    &prefs,
		CreatedAt:      dbUser.CreatedAt,
		UpdatedAt:      dbUser.UpdatedAt,
	}, nil
}

// ConvertGRPCAuthProviderToEntity converts gRPC auth provider to domain entity
func (c *GRPCConverter) ConvertGRPCAuthProviderToEntity(provider GRPCAuthProvider) (entities.AuthProvider, error) {
	switch provider {
	case GRPCAuthProviderGoogle:
		return entities.AuthProviderGoogle, nil
	case GRPCAuthProviderX:
		return entities.AuthProviderX, nil
	default:
		return "", errors.InvalidInput("auth_provider", "unknown gRPC auth provider").
			WithField("grpc_provider", int32(provider))
	}
}

// ConvertEntityAuthProviderToGRPC converts domain entity auth provider to gRPC
func (c *GRPCConverter) ConvertEntityAuthProviderToGRPC(provider entities.AuthProvider) GRPCAuthProvider {
	switch provider {
	case entities.AuthProviderGoogle:
		return GRPCAuthProviderGoogle
	case entities.AuthProviderX:
		return GRPCAuthProviderX
	default:
		return GRPCAuthProviderUnspecified
	}
}

// ConvertGRPCPreferencesToEntity converts gRPC preferences to domain entity
func (c *GRPCConverter) ConvertGRPCPreferencesToEntity(grpcPrefs *GRPCUserPreferences) entities.UserPreferences {
	if grpcPrefs == nil {
		return entities.UserPreferences{
			Language: DefaultLanguage,
			DarkMode: DefaultDarkMode,
			Timezone: DefaultTimezone,
		}
	}

	return entities.UserPreferences{
		Language: grpcPrefs.Language,
		DarkMode: grpcPrefs.DarkMode,
		Timezone: grpcPrefs.Timezone,
	}
}

// ConvertEntityPreferencesToGRPC converts domain entity preferences to gRPC
func (c *GRPCConverter) ConvertEntityPreferencesToGRPC(entityPrefs entities.UserPreferences) *GRPCUserPreferences {
	return &GRPCUserPreferences{
		Language: entityPrefs.Language,
		DarkMode: entityPrefs.DarkMode,
		Timezone: entityPrefs.Timezone,
	}
}