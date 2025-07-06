package helpers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/necofuryai/bocchi-the-map/api/domain/entities"
	"github.com/necofuryai/bocchi-the-map/api/infrastructure/database"
	"github.com/onsi/gomega"
)

var (
	fixedTestTime = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
)

// FixtureManager manages test data fixtures
type FixtureManager struct {
	db      *TestDatabase
}

// NewFixtureManager creates a new fixture manager
func NewFixtureManager(db *TestDatabase) *FixtureManager {
	return &FixtureManager{
		db: db,
	}
}

// SpotFixture represents a test spot fixture
type SpotFixture struct {
	ID          string
	Name        string
	NameI18n    map[string]string
	Latitude    float64
	Longitude   float64
	Category    string
	Address     string
	AddressI18n map[string]string
	CountryCode string
}

// UserFixture represents a test user fixture
type UserFixture struct {
	ID             string
	Email          string
	DisplayName    string
	AuthProvider   string
	AuthProviderID string
	Preferences    entities.UserPreferences
	Permissions    []string // JWT permissions for authentication testing
	Role           string   // User role for authorization testing
}

// ReviewFixture represents a test review fixture
type ReviewFixture struct {
	ID            string
	SpotID        string
	UserID        string
	Rating        int
	Comment       string
	RatingAspects map[string]int
}

// CreateSpotFixture creates a spot fixture in the database
func (fm *FixtureManager) CreateSpotFixture(ctx context.Context, fixture SpotFixture) *entities.Spot {
	
	// Marshal i18n maps to JSON
	var nameI18nJSON []byte
	if fixture.NameI18n != nil {
		var err error
		nameI18nJSON, err = json.Marshal(fixture.NameI18n)
		gomega.Expect(err).NotTo(gomega.HaveOccurred(), "Failed to marshal name i18n")
	}
	
	var addressI18nJSON []byte
	if fixture.AddressI18n != nil {
		var err error
		addressI18nJSON, err = json.Marshal(fixture.AddressI18n)
		gomega.Expect(err).NotTo(gomega.HaveOccurred(), "Failed to marshal address i18n")
	}
	
	// Create spot in database
	params := database.CreateSpotParams{
		ID:          fixture.ID,
		Name:        fixture.Name,
		NameI18n:    nameI18nJSON,
		Latitude:    fmt.Sprintf("%.8f", fixture.Latitude),
		Longitude:   fmt.Sprintf("%.8f", fixture.Longitude),
		Category:    fixture.Category,
		Address:     fixture.Address,
		AddressI18n: addressI18nJSON,
		CountryCode: fixture.CountryCode,
	}
	
	err := fm.db.Queries.CreateSpot(ctx, params)
	gomega.Expect(err).NotTo(gomega.HaveOccurred(), "Failed to create spot fixture: %v", err)
	
	return &entities.Spot{
		ID:          fixture.ID,
		Name:        fixture.Name,
		NameI18n:    fixture.NameI18n,
		Latitude:    fixture.Latitude,
		Longitude:   fixture.Longitude,
		Category:    fixture.Category,
		Address:     fixture.Address,
		AddressI18n: fixture.AddressI18n,
		CountryCode: fixture.CountryCode,
		CreatedAt:   fixedTestTime,
		UpdatedAt:   fixedTestTime,
	}
}

// CreateUserFixture creates a user fixture in the database
func (fm *FixtureManager) CreateUserFixture(ctx context.Context, fixture UserFixture) *entities.User {
	
	params := database.CreateUserParams{
		ID:            fixture.ID,
		Email:         fixture.Email,
		Name:          sql.NullString{String: fixture.DisplayName, Valid: true},
		Nickname:      sql.NullString{String: fixture.DisplayName, Valid: true},
		Picture:       sql.NullString{},
		Provider:      fixture.AuthProvider,
		ProviderID:    fixture.AuthProviderID,
		EmailVerified: true,
		Preferences:   []byte(`{}`),
	}
	
	err := fm.db.Queries.CreateUser(ctx, params)
	gomega.Expect(err).NotTo(gomega.HaveOccurred(), "Failed to create user fixture: %v", err)
	
	return &entities.User{
		ID:             fixture.ID,
		Email:          fixture.Email,
		DisplayName:    fixture.DisplayName,
		AuthProvider:   entities.AuthProvider(fixture.AuthProvider),
		AuthProviderID: fixture.AuthProviderID,
		Preferences:    fixture.Preferences,
		CreatedAt:      fixedTestTime,
		UpdatedAt:      fixedTestTime,
	}
}

// CreateReviewFixture creates a review fixture in the database
func (fm *FixtureManager) CreateReviewFixture(ctx context.Context, fixture ReviewFixture) *entities.Review {
	
	var ratingAspectsJSON []byte
	if fixture.RatingAspects != nil {
		var err error
		ratingAspectsJSON, err = json.Marshal(fixture.RatingAspects)
		gomega.Expect(err).NotTo(gomega.HaveOccurred(), "Failed to marshal rating aspects")
	}
	
	params := database.CreateReviewParams{
		ID:            fixture.ID,
		SpotID:        fixture.SpotID,
		UserID:        sql.NullString{String: fixture.UserID, Valid: true},
		Rating:        int32(fixture.Rating),
		Comment:       sql.NullString{String: fixture.Comment, Valid: fixture.Comment != ""},
		RatingAspects: ratingAspectsJSON,
	}
	
	err := fm.db.Queries.CreateReview(ctx, params)
	gomega.Expect(err).NotTo(gomega.HaveOccurred(), "Failed to create review fixture: %v", err)
	
	return &entities.Review{
		ID:            fixture.ID,
		SpotID:        fixture.SpotID,
		UserID:        fixture.UserID,
		Rating:        fixture.Rating,
		Comment:       fixture.Comment,
		RatingAspects: fixture.RatingAspects,
		CreatedAt:     fixedTestTime,
		UpdatedAt:     fixedTestTime,
	}
}

// DefaultSpotFixtures returns common spot test data
func (fm *FixtureManager) DefaultSpotFixtures() []SpotFixture {
	return []SpotFixture{
		{
			ID:        "spot-cafe-tokyo",
			Name:      "Solo Cafe Tokyo",
			NameI18n:  map[string]string{"ja": "ソロカフェ東京"},
			Latitude:  35.6762,
			Longitude: 139.6503,
			Category:  "cafe",
			Address:   "1-1-1 Shibuya, Tokyo",
			AddressI18n: map[string]string{
				"ja": "東京都渋谷区1-1-1",
			},
			CountryCode: "JP",
		},
		{
			ID:        "spot-library-osaka",
			Name:      "Quiet Study Library",
			NameI18n:  map[string]string{"ja": "静かな勉強図書館"},
			Latitude:  34.6937,
			Longitude: 135.5023,
			Category:  "library",
			Address:   "2-2-2 Namba, Osaka",
			AddressI18n: map[string]string{
				"ja": "大阪府大阪市2-2-2",
			},
			CountryCode: "JP",
		},
	}
}

// DefaultUserFixtures returns common user test data
func (fm *FixtureManager) DefaultUserFixtures() []UserFixture {
	return []UserFixture{
		{
			ID:             "user-creator-1",
			Email:          "creator1@example.com",
			DisplayName:    "Spot Creator 1",
			AuthProvider:   "google",
			AuthProviderID: "google_123",
			Permissions:    []string{"read:spots", "write:reviews", "edit:profile"},
			Role:           "user",
			Preferences: entities.UserPreferences{
				Language: "ja",
				DarkMode: true,
				Timezone: "Asia/Tokyo",
			},
		},
		{
			ID:             "user-creator-2",
			Email:          "creator2@example.com",
			DisplayName:    "Spot Creator 2",
			AuthProvider:   "google",
			AuthProviderID: "google_456",
			Permissions:    []string{"read:spots", "write:reviews", "edit:profile"},
			Role:           "user",
			Preferences: entities.UserPreferences{
				Language: "en",
				DarkMode: false,
				Timezone: "UTC",
			},
		},
	}
}

// createStandardFixtures creates a standard set of test data and returns the created entities
func (fm *FixtureManager) createStandardFixtures(ctx context.Context) ([]*entities.User, []*entities.Spot) {
	var users []*entities.User
	var spots []*entities.Spot
	
	// Create users first (for foreign key constraints)
	for _, userFixture := range fm.DefaultUserFixtures() {
		user := fm.CreateUserFixture(ctx, userFixture)
		users = append(users, user)
	}
	
	// Create spots
	for _, spotFixture := range fm.DefaultSpotFixtures() {
		spot := fm.CreateSpotFixture(ctx, spotFixture)
		spots = append(spots, spot)
	}
	
	return users, spots
}

// SetupStandardFixtures creates a standard set of test data
func (fm *FixtureManager) SetupStandardFixtures(ctx context.Context) {
	fm.createStandardFixtures(ctx)
}

// SetupStandardFixturesWithReturn creates a standard set of test data and returns the created entities
func (fm *FixtureManager) SetupStandardFixturesWithReturn() ([]*entities.User, []*entities.Spot) {
	return fm.createStandardFixtures(context.Background())
}

// CleanupFixtures removes all fixture data
func (fm *FixtureManager) CleanupFixtures(t *testing.T) error {
	err := fm.db.CleanDatabase()
	if err != nil {
		// Log error but don't panic in cleanup
		t.Logf("Warning: Failed to cleanup fixtures: %v", err)
		return err
	}
	return nil
}

// AuthUserFixtures returns user fixtures with various authentication permissions for testing
func (fm *FixtureManager) AuthUserFixtures() []UserFixture {
	return []UserFixture{
		{
			ID:             "auth-admin-user",
			Email:          "admin@bocchi-map.com",
			DisplayName:    "Admin User",
			AuthProvider:   "auth0",
			AuthProviderID: "auth0_admin_123",
			Permissions:    []string{"read:spots", "write:spots", "delete:spots", "admin:users", "admin:system"},
			Role:           "admin",
			Preferences: entities.UserPreferences{
				Language: "en",
				DarkMode: false,
				Timezone: "UTC",
			},
		},
		{
			ID:             "auth-moderator-user",
			Email:          "moderator@bocchi-map.com",
			DisplayName:    "Moderator User",
			AuthProvider:   "auth0",
			AuthProviderID: "auth0_moderator_456",
			Permissions:    []string{"read:spots", "write:spots", "moderate:reviews", "edit:spots"},
			Role:           "moderator",
			Preferences: entities.UserPreferences{
				Language: "en",
				DarkMode: true,
				Timezone: "UTC",
			},
		},
		{
			ID:             "auth-regular-user",
			Email:          "user@bocchi-map.com",
			DisplayName:    "Regular User",
			AuthProvider:   "auth0",
			AuthProviderID: "auth0_user_789",
			Permissions:    []string{"read:spots", "write:reviews", "edit:profile"},
			Role:           "user",
			Preferences: entities.UserPreferences{
				Language: "ja",
				DarkMode: false,
				Timezone: "Asia/Tokyo",
			},
		},
		{
			ID:             "auth-viewer-user",
			Email:          "viewer@bocchi-map.com",
			DisplayName:    "Viewer User",
			AuthProvider:   "auth0",
			AuthProviderID: "auth0_viewer_101",
			Permissions:    []string{"read:spots"},
			Role:           "viewer",
			Preferences: entities.UserPreferences{
				Language: "en",
				DarkMode: false,
				Timezone: "UTC",
			},
		},
		{
			ID:             "auth-premium-user",
			Email:          "premium@bocchi-map.com",
			DisplayName:    "Premium User",
			AuthProvider:   "auth0",
			AuthProviderID: "auth0_premium_202",
			Permissions:    []string{"read:spots", "write:reviews", "edit:profile", "premium:features", "export:data"},
			Role:           "premium",
			Preferences: entities.UserPreferences{
				Language: "ja",
				DarkMode: true,
				Timezone: "Asia/Tokyo",
			},
		},
		{
			ID:             "auth-restricted-user",
			Email:          "restricted@bocchi-map.com",
			DisplayName:    "Restricted User",
			AuthProvider:   "auth0",
			AuthProviderID: "auth0_restricted_303",
			Permissions:    []string{"read:spots"}, // Very limited permissions
			Role:           "restricted",
			Preferences: entities.UserPreferences{
				Language: "en",
				DarkMode: false,
				Timezone: "UTC",
			},
		},
	}
}

// CreateAuthUserFixture creates an auth user fixture with JWT claims compatibility
func (fm *FixtureManager) CreateAuthUserFixture(ctx context.Context, fixture UserFixture) *entities.User {
	// Create the basic user first
	user := fm.CreateUserFixture(ctx, fixture)
	
	// The permissions and role information is stored in the fixture for JWT testing
	// but not directly in the database user entity (that would be handled by Auth0)
	return user
}

// SetupAuthUserFixtures creates authentication-related user test data
func (fm *FixtureManager) SetupAuthUserFixtures(ctx context.Context) []*entities.User {
	var users []*entities.User
	
	for _, userFixture := range fm.AuthUserFixtures() {
		user := fm.CreateAuthUserFixture(ctx, userFixture)
		users = append(users, user)
	}
	
	return users
}

// GetAdminUserFixture returns an admin user fixture for testing
func (fm *FixtureManager) GetAdminUserFixture() UserFixture {
	fixtures := fm.AuthUserFixtures()
	for _, fixture := range fixtures {
		if fixture.Role == "admin" {
			return fixture
		}
	}
	// Return default admin if not found
	return UserFixture{
		ID:             "default-admin",
		Email:          "admin@example.com",
		DisplayName:    "Default Admin",
		AuthProvider:   "auth0",
		AuthProviderID: "auth0_default_admin",
		Permissions:    []string{"read:spots", "write:spots", "delete:spots", "admin:users", "admin:system"},
		Role:           "admin",
		Preferences: entities.UserPreferences{
			Language: "en",
			DarkMode: false,
			Timezone: "UTC",
		},
	}
}

// GetRegularUserFixture returns a regular user fixture for testing
func (fm *FixtureManager) GetRegularUserFixture() UserFixture {
	fixtures := fm.AuthUserFixtures()
	for _, fixture := range fixtures {
		if fixture.Role == "user" {
			return fixture
		}
	}
	// Return default user if not found
	return UserFixture{
		ID:             "default-user",
		Email:          "user@example.com",
		DisplayName:    "Default User",
		AuthProvider:   "auth0",
		AuthProviderID: "auth0_default_user",
		Permissions:    []string{"read:spots", "write:reviews", "edit:profile"},
		Role:           "user",
		Preferences: entities.UserPreferences{
			Language: "en",
			DarkMode: false,
			Timezone: "UTC",
		},
	}
}

// GetViewerUserFixture returns a viewer user fixture for testing
func (fm *FixtureManager) GetViewerUserFixture() UserFixture {
	fixtures := fm.AuthUserFixtures()
	for _, fixture := range fixtures {
		if fixture.Role == "viewer" {
			return fixture
		}
	}
	// Return default viewer if not found
	return UserFixture{
		ID:             "default-viewer",
		Email:          "viewer@example.com",
		DisplayName:    "Default Viewer",
		AuthProvider:   "auth0",
		AuthProviderID: "auth0_default_viewer",
		Permissions:    []string{"read:spots"},
		Role:           "viewer",
		Preferences: entities.UserPreferences{
			Language: "en",
			DarkMode: false,
			Timezone: "UTC",
		},
	}
}

// CreateUserFixtureWithPermissions creates a user fixture with specific permissions
func (fm *FixtureManager) CreateUserFixtureWithPermissions(ctx context.Context, userID, email, role string, permissions []string) *entities.User {
	fixture := UserFixture{
		ID:             userID,
		Email:          email,
		DisplayName:    fmt.Sprintf("Test %s User", role),
		AuthProvider:   "auth0",
		AuthProviderID: fmt.Sprintf("auth0_%s_%s", role, userID),
		Permissions:    permissions,
		Role:           role,
		Preferences: entities.UserPreferences{
			Language: "en",
			DarkMode: false,
			Timezone: "UTC",
		},
	}
	
	return fm.CreateAuthUserFixture(ctx, fixture)
}

// PermissionTestScenarios returns common permission testing scenarios
func (fm *FixtureManager) PermissionTestScenarios() map[string][]string {
	return map[string][]string{
		"admin": {
			"read:spots", "write:spots", "delete:spots", 
			"read:reviews", "write:reviews", "delete:reviews",
			"admin:users", "admin:system", "admin:analytics",
		},
		"moderator": {
			"read:spots", "write:spots", "edit:spots",
			"read:reviews", "write:reviews", "moderate:reviews",
			"moderate:users",
		},
		"premium": {
			"read:spots", "write:reviews", "edit:profile",
			"premium:features", "export:data", "advanced:search",
		},
		"user": {
			"read:spots", "write:reviews", "edit:profile",
		},
		"viewer": {
			"read:spots",
		},
		"restricted": {
			"read:spots", // Very limited access
		},
	}
}