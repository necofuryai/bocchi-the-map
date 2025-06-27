package converters

import (
	"encoding/json"
	"database/sql"
	"time"

	"github.com/necofuryai/bocchi-the-map/api/domain/entities"
	"github.com/necofuryai/bocchi-the-map/api/infrastructure/database"
	"github.com/necofuryai/bocchi-the-map/api/pkg/errors"
)

// UserConverter handles conversions between different user representations
type UserConverter struct{}

// NewUserConverter creates a new UserConverter instance
func NewUserConverter() *UserConverter {
	return &UserConverter{}
}

// ConvertDatabaseToEntity converts database User to domain entity
func (c *UserConverter) ConvertDatabaseToEntity(dbUser database.User) (*entities.User, error) {
	var prefs entities.UserPreferences
	if len(dbUser.Preferences) > 0 {
		if err := json.Unmarshal(dbUser.Preferences, &prefs); err != nil {
			return nil, errors.Wrap(err, errors.ErrTypeInternal, "failed to unmarshal user preferences").
				WithField("user_id", dbUser.ID)
		}
	} else {
		// Default preferences
		prefs = entities.UserPreferences{
			Language: "ja",
			DarkMode: false,
			Timezone: "Asia/Tokyo",
		}
	}

	// Convert database auth provider enum to domain enum
	var authProvider entities.AuthProvider
	switch dbUser.AuthProvider {
	case database.UsersAuthProviderGoogle:
		authProvider = entities.AuthProviderGoogle
	case database.UsersAuthProviderTwitter:
		authProvider = entities.AuthProviderX
	default:
		return nil, errors.InvalidInput("auth_provider", "unknown provider from database").
			WithField("db_provider", string(dbUser.AuthProvider)).
			WithField("user_id", dbUser.ID)
	}

	// Handle nullable avatar URL
	var avatarURL string
	if dbUser.AvatarUrl.Valid {
		avatarURL = dbUser.AvatarUrl.String
	}

	return &entities.User{
		ID:             dbUser.ID,
		AnonymousID:    dbUser.AnonymousID,
		Email:          dbUser.Email,
		DisplayName:    dbUser.DisplayName,
		AvatarURL:      avatarURL,
		AuthProvider:   authProvider,
		AuthProviderID: dbUser.AuthProviderID,
		Preferences:    prefs,
		CreatedAt:      dbUser.CreatedAt,
		UpdatedAt:      dbUser.UpdatedAt,
	}, nil
}

// ConvertEntityToDatabase converts domain entity to database parameters
func (c *UserConverter) ConvertEntityToDatabase(user *entities.User) (database.CreateUserParams, error) {
	// Convert preferences to JSON
	prefsJSON, err := json.Marshal(user.Preferences)
	if err != nil {
		return database.CreateUserParams{}, errors.Wrap(err, errors.ErrTypeInternal, "failed to marshal user preferences").
			WithField("user_id", user.ID)
	}

	// Convert domain auth provider to database enum
	var dbAuthProvider database.UsersAuthProvider
	switch user.AuthProvider {
	case entities.AuthProviderGoogle:
		dbAuthProvider = database.UsersAuthProviderGoogle
	case entities.AuthProviderX:
		dbAuthProvider = database.UsersAuthProviderTwitter
	default:
		return database.CreateUserParams{}, errors.InvalidInput("auth_provider", "unsupported provider").
			WithField("provider", string(user.AuthProvider)).
			WithField("user_id", user.ID)
	}

	// Convert avatar URL to sql.NullString
	var avatarURL sql.NullString
	if user.AvatarURL != "" {
		avatarURL = sql.NullString{String: user.AvatarURL, Valid: true}
	}

	return database.CreateUserParams{
		ID:             user.ID,
		AnonymousID:    user.AnonymousID,
		Email:          user.Email,
		DisplayName:    user.DisplayName,
		AvatarUrl:      avatarURL,
		AuthProvider:   dbAuthProvider,
		AuthProviderID: user.AuthProviderID,
		Preferences:    prefsJSON,
	}, nil
}

// ConvertEntityToUpsertDatabase converts domain entity to database upsert parameters
func (c *UserConverter) ConvertEntityToUpsertDatabase(user *entities.User) (database.UpsertUserParams, error) {
	// Convert preferences to JSON
	prefsJSON, err := json.Marshal(user.Preferences)
	if err != nil {
		return database.UpsertUserParams{}, errors.Wrap(err, errors.ErrTypeInternal, "failed to marshal user preferences").
			WithField("user_id", user.ID)
	}

	// Convert domain auth provider to database enum
	var dbAuthProvider database.UsersAuthProvider
	switch user.AuthProvider {
	case entities.AuthProviderGoogle:
		dbAuthProvider = database.UsersAuthProviderGoogle
	case entities.AuthProviderX:
		dbAuthProvider = database.UsersAuthProviderTwitter
	default:
		return database.UpsertUserParams{}, errors.InvalidInput("auth_provider", "unsupported provider").
			WithField("provider", string(user.AuthProvider)).
			WithField("user_id", user.ID)
	}

	// Convert avatar URL to sql.NullString
	var avatarURL sql.NullString
	if user.AvatarURL != "" {
		avatarURL = sql.NullString{String: user.AvatarURL, Valid: true}
	}

	return database.UpsertUserParams{
		ID:             user.ID,
		Email:          user.Email,
		DisplayName:    user.DisplayName,
		AvatarUrl:      avatarURL,
		AuthProvider:   dbAuthProvider,
		AuthProviderID: user.AuthProviderID,
		Preferences:    prefsJSON,
	}, nil
}

// AuthProviderToDatabase converts domain auth provider to database enum
func (c *UserConverter) AuthProviderToDatabase(provider entities.AuthProvider) (database.UsersAuthProvider, error) {
	switch provider {
	case entities.AuthProviderGoogle:
		return database.UsersAuthProviderGoogle, nil
	case entities.AuthProviderX:
		return database.UsersAuthProviderTwitter, nil
	default:
		return "", errors.InvalidInput("auth_provider", "unsupported provider").
			WithField("provider", string(provider))
	}
}

// AuthProviderFromDatabase converts database auth provider to domain enum
func (c *UserConverter) AuthProviderFromDatabase(provider database.UsersAuthProvider) (entities.AuthProvider, error) {
	switch provider {
	case database.UsersAuthProviderGoogle:
		return entities.AuthProviderGoogle, nil
	case database.UsersAuthProviderTwitter:
		return entities.AuthProviderX, nil
	default:
		return "", errors.InvalidInput("auth_provider", "unknown provider from database").
			WithField("db_provider", string(provider))
	}
}

// ConvertHTTPPreferencesToEntity converts HTTP preferences map to domain entity
func (c *UserConverter) ConvertHTTPPreferencesToEntity(prefsMap map[string]interface{}) entities.UserPreferences {
	prefs := entities.UserPreferences{
		Language: "ja",       // Default
		DarkMode: false,      // Default
		Timezone: "Asia/Tokyo", // Default
	}

	if prefsMap != nil {
		if lang, ok := prefsMap["language"].(string); ok {
			prefs.Language = lang
		}
		if darkMode, ok := prefsMap["dark_mode"].(bool); ok {
			prefs.DarkMode = darkMode
		}
		if timezone, ok := prefsMap["timezone"].(string); ok {
			prefs.Timezone = timezone
		}
	}

	return prefs
}

// ConvertEntityPreferencesToHTTP converts domain entity preferences to HTTP map
func (c *UserConverter) ConvertEntityPreferencesToHTTP(prefs entities.UserPreferences) map[string]interface{} {
	return map[string]interface{}{
		"language":  prefs.Language,
		"dark_mode": prefs.DarkMode,
		"timezone":  prefs.Timezone,
	}
}

// ConvertHTTPAuthProviderToEntity converts HTTP auth provider string to domain entity
func (c *UserConverter) ConvertHTTPAuthProviderToEntity(provider string) (entities.AuthProvider, error) {
	switch provider {
	case "google":
		return entities.AuthProviderGoogle, nil
	case "twitter", "x":
		return entities.AuthProviderX, nil
	default:
		return "", errors.InvalidInput("auth_provider", "unsupported provider").
			WithField("provider", provider)
	}
}