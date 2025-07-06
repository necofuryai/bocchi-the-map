package entities

import (
	"time"
)

// AuthProvider represents authentication provider types
type AuthProvider string

const (
	AuthProviderGoogle  AuthProvider = "google"
	AuthProviderTwitter AuthProvider = "twitter"
	AuthProviderX       AuthProvider = "x"
	AuthProviderAuth0   AuthProvider = "auth0"
)

// UserPreferences represents user-specific preferences
type UserPreferences struct {
	Language string `json:"language"`
	DarkMode bool   `json:"dark_mode"`
	Timezone string `json:"timezone"`
}

// User represents a user in the domain
type User struct {
	ID             string           `json:"id"`
	Email          string           `json:"email"`
	DisplayName    string           `json:"display_name"`
	AvatarUrl      string           `json:"avatar_url,omitempty"`
	AuthProvider   AuthProvider     `json:"auth_provider"`
	AuthProviderID string           `json:"auth_provider_id"`
	Preferences    UserPreferences  `json:"preferences"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
}

// NewUser creates a new User instance
func NewUser(email, displayName string, authProvider AuthProvider, authProviderID string) *User {
	now := time.Now()
	return &User{
		Email:          email,
		DisplayName:    displayName,
		AuthProvider:   authProvider,
		AuthProviderID: authProviderID,
		Preferences: UserPreferences{
			Language: "en",
			DarkMode: false,
			Timezone: "UTC",
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

// UpdateDisplayName updates the user's display name
func (u *User) UpdateDisplayName(name string) {
	u.DisplayName = name
	u.UpdatedAt = time.Now()
}

// UpdateAvatarUrl updates the user's avatar URL
func (u *User) UpdateAvatarUrl(url string) {
	u.AvatarUrl = url
	u.UpdatedAt = time.Now()
}