package entities

import (
	"time"
)

// User represents an authenticated user
type User struct {
	ID             string          `json:"id"`
	AnonymousID    string          `json:"anonymous_id"`
	Email          string          `json:"email"`
	DisplayName    string          `json:"display_name"`
	AvatarURL      string          `json:"avatar_url,omitempty"`
	AuthProvider   AuthProvider    `json:"auth_provider"`
	AuthProviderID string          `json:"auth_provider_id"`
	Preferences    UserPreferences `json:"preferences"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

// AuthProvider represents the authentication provider
type AuthProvider string

const (
	AuthProviderGoogle AuthProvider = "google"
	AuthProviderX      AuthProvider = "x"
)

// UserPreferences represents user preferences
type UserPreferences struct {
	Language string `json:"language"`
	DarkMode bool   `json:"dark_mode"`
	Timezone string `json:"timezone"`
}

// NewUser creates a new User instance
func NewUser(email, displayName string, provider AuthProvider, providerID string) *User {
	now := time.Now()
	return &User{
		Email:          email,
		DisplayName:    displayName,
		AuthProvider:   provider,
		AuthProviderID: providerID,
		Preferences: UserPreferences{
			Language: "ja",
			DarkMode: false,
			Timezone: "Asia/Tokyo",
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// UpdatePreferences updates user preferences
func (u *User) UpdatePreferences(prefs UserPreferences) {
	u.Preferences = prefs
	u.UpdatedAt = time.Now()
}

// SetAnonymousID sets the anonymous ID for the user
func (u *User) SetAnonymousID(id string) {
	u.AnonymousID = id
	u.UpdatedAt = time.Now()
}