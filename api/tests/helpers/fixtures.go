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
	
	// Create spot in database
	params := database.CreateSpotParams{
		ID:          fixture.ID,
		Name:        fixture.Name,
		NameI18n:    fixture.NameI18n,
		Latitude:    fmt.Sprintf("%.8f", fixture.Latitude),
		Longitude:   fmt.Sprintf("%.8f", fixture.Longitude),
		Category:    fixture.Category,
		Address:     fixture.Address,
		AddressI18n: fixture.AddressI18n,
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
		ID:             fixture.ID,
		Email:          fixture.Email,
		DisplayName:    fixture.DisplayName,
		AuthProvider:   database.UsersAuthProvider(fixture.AuthProvider),
		AuthProviderID: fixture.AuthProviderID,
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
		UserID:        fixture.UserID,
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