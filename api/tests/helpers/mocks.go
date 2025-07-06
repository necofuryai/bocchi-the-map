package helpers

import (
	"context"
	stdErrors "errors"
	"fmt"
	"math"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/necofuryai/bocchi-the-map/api/domain/entities"
	"github.com/necofuryai/bocchi-the-map/api/pkg/auth"
	"github.com/necofuryai/bocchi-the-map/api/pkg/errors"
)

const (
	// DefaultMockAuthProviderID is the default auth provider ID used in mock scenarios
	DefaultMockAuthProviderID = "mock_123"
)

// buildAuthKey constructs authentication key from provider and provider ID
func buildAuthKey(provider entities.AuthProvider, providerID string) string {
	return string(provider) + ":" + providerID
}

// distanceInKm calculates the distance between two coordinates in kilometers using the Haversine formula
func distanceInKm(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadiusKm = 6371

	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
		math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c
}

// MockSpotRepository provides a mock implementation for testing
type MockSpotRepository struct {
	mu                    sync.RWMutex
	spots                 map[string]*entities.Spot
	createSpotFunc        func(ctx context.Context, spot *entities.Spot) error
	getSpotFunc           func(ctx context.Context, id string) (*entities.Spot, error)
	getByCoordinatesFunc  func(ctx context.Context, lat, lng, radiusKm float64) ([]*entities.Spot, error)
	listSpotsFunc         func(ctx context.Context, offset, limit int) ([]*entities.Spot, int, error)
	searchSpotsFunc       func(ctx context.Context, query string, lang string, offset, limit int) ([]*entities.Spot, int, error)
	updateSpotFunc        func(ctx context.Context, spot *entities.Spot) error
	deleteSpotFunc        func(ctx context.Context, id string) error
	updateRatingFunc      func(ctx context.Context, spotID string, averageRating float64, reviewCount int) error
}

// NewMockSpotRepository creates a new mock spot repository
func NewMockSpotRepository() *MockSpotRepository {
	return &MockSpotRepository{
		spots: make(map[string]*entities.Spot),
	}
}

// Implement the SpotRepository interface
func (m *MockSpotRepository) Create(ctx context.Context, spot *entities.Spot) error {
	m.mu.RLock()
	createFunc := m.createSpotFunc
	m.mu.RUnlock()
	
	if createFunc != nil {
		return createFunc(ctx, spot)
	}
	
	if spot.ID == "" {
		return stdErrors.New("spot.ID cannot be empty")
	}
	
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Check if spot with the same ID already exists
	if _, exists := m.spots[spot.ID]; exists {
		return fmt.Errorf("spot with ID already exists: id=%s", spot.ID)
	}
	
	m.spots[spot.ID] = spot
	return nil
}

func (m *MockSpotRepository) GetByID(ctx context.Context, id string) (*entities.Spot, error) {
	m.mu.RLock()
	getFunc := m.getSpotFunc
	spot, exists := m.spots[id]
	m.mu.RUnlock()
	
	if getFunc != nil {
		return getFunc(ctx, id)
	}
	
	if !exists {
		return nil, fmt.Errorf("spot not found: id=%s", id)
	}
	
	return spot, nil
}

func (m *MockSpotRepository) GetByCoordinates(ctx context.Context, lat, lng, radiusKm float64) ([]*entities.Spot, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	if m.getByCoordinatesFunc != nil {
		return m.getByCoordinatesFunc(ctx, lat, lng, radiusKm)
	}
	
	result := make([]*entities.Spot, 0)
	for _, spot := range m.spots {
		// Simple distance check - in real implementation would use proper geospatial calculation
		if distanceInKm(lat, lng, spot.Latitude, spot.Longitude) <= radiusKm {
			result = append(result, spot)
		}
	}
	
	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})
	
	return result, nil
}

func (m *MockSpotRepository) List(ctx context.Context, offset, limit int) ([]*entities.Spot, int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	if m.listSpotsFunc != nil {
		return m.listSpotsFunc(ctx, offset, limit)
	}
	
	result := make([]*entities.Spot, 0, len(m.spots))
	for _, spot := range m.spots {
		result = append(result, spot)
	}
	
	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})
	
	total := len(result)
	
	// Apply pagination
	if offset >= total {
		return []*entities.Spot{}, total, nil
	}
	
	end := offset + limit
	if end > total {
		end = total
	}
	
	return result[offset:end], total, nil
}

func (m *MockSpotRepository) Search(ctx context.Context, query string, lang string, offset, limit int) ([]*entities.Spot, int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	if m.searchSpotsFunc != nil {
		return m.searchSpotsFunc(ctx, query, lang, offset, limit)
	}
	
	result := make([]*entities.Spot, 0)
	lowerQuery := strings.ToLower(query)
	
	for _, spot := range m.spots {
		matched := false
		
		// Simple text search in name and address
		if strings.Contains(strings.ToLower(spot.Name), lowerQuery) ||
		   strings.Contains(strings.ToLower(spot.Address), lowerQuery) ||
		   strings.Contains(strings.ToLower(spot.Category), lowerQuery) {
			matched = true
		}
		
		// Also search in localized names and addresses if available
		if !matched && spot.NameI18n != nil {
			if localizedName, ok := spot.NameI18n[lang]; ok {
				if strings.Contains(strings.ToLower(localizedName), lowerQuery) {
					matched = true
				}
			}
		}
		
		if !matched && spot.AddressI18n != nil {
			if localizedAddress, ok := spot.AddressI18n[lang]; ok {
				if strings.Contains(strings.ToLower(localizedAddress), lowerQuery) {
					matched = true
				}
			}
		}
		
		if matched {
			result = append(result, spot)
		}
	}
	
	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})
	
	total := len(result)
	
	// Apply pagination
	if offset >= total {
		return []*entities.Spot{}, total, nil
	}
	
	end := offset + limit
	if end > total {
		end = total
	}
	
	return result[offset:end], total, nil
}

func (m *MockSpotRepository) Update(ctx context.Context, spot *entities.Spot) error {
	m.mu.RLock()
	updateFunc := m.updateSpotFunc
	_, exists := m.spots[spot.ID]
	m.mu.RUnlock()
	
	if updateFunc != nil {
		return updateFunc(ctx, spot)
	}
	
	if !exists {
		return fmt.Errorf("spot not found: id=%s", spot.ID)
	}
	
	m.mu.Lock()
	m.spots[spot.ID] = spot
	m.mu.Unlock()
	return nil
}

func (m *MockSpotRepository) Delete(ctx context.Context, id string) error {
	m.mu.RLock()
	deleteFunc := m.deleteSpotFunc
	_, exists := m.spots[id]
	m.mu.RUnlock()
	
	if deleteFunc != nil {
		return deleteFunc(ctx, id)
	}
	
	if !exists {
		return fmt.Errorf("spot not found: id=%s", id)
	}
	
	m.mu.Lock()
	delete(m.spots, id)
	m.mu.Unlock()
	return nil
}

func (m *MockSpotRepository) UpdateRating(ctx context.Context, spotID string, averageRating float64, reviewCount int) error {
	m.mu.RLock()
	updateRatingFunc := m.updateRatingFunc
	spot, exists := m.spots[spotID]
	m.mu.RUnlock()
	
	if updateRatingFunc != nil {
		return updateRatingFunc(ctx, spotID, averageRating, reviewCount)
	}
	
	if !exists {
		return fmt.Errorf("spot not found: id=%s", spotID)
	}
	
	m.mu.Lock()
	spot.AverageRating = averageRating
	spot.ReviewCount = reviewCount
	m.mu.Unlock()
	
	return nil
}

// Mock configuration methods
func (m *MockSpotRepository) SetCreateSpotFunc(fn func(ctx context.Context, spot *entities.Spot) error) {
	m.mu.Lock()
	m.createSpotFunc = fn
	m.mu.Unlock()
}

func (m *MockSpotRepository) SetGetSpotFunc(fn func(ctx context.Context, id string) (*entities.Spot, error)) {
	m.mu.Lock()
	m.getSpotFunc = fn
	m.mu.Unlock()
}

func (m *MockSpotRepository) SetListSpotsFunc(fn func(ctx context.Context, offset, limit int) ([]*entities.Spot, int, error)) {
	m.mu.Lock()
	m.listSpotsFunc = fn
	m.mu.Unlock()
}

func (m *MockSpotRepository) SetUpdateSpotFunc(fn func(ctx context.Context, spot *entities.Spot) error) {
	m.mu.Lock()
	m.updateSpotFunc = fn
	m.mu.Unlock()
}

func (m *MockSpotRepository) SetDeleteSpotFunc(fn func(ctx context.Context, id string) error) {
	m.mu.Lock()
	m.deleteSpotFunc = fn
	m.mu.Unlock()
}

func (m *MockSpotRepository) SetGetByCoordinatesFunc(fn func(ctx context.Context, lat, lng, radiusKm float64) ([]*entities.Spot, error)) {
	m.mu.Lock()
	m.getByCoordinatesFunc = fn
	m.mu.Unlock()
}

func (m *MockSpotRepository) SetSearchSpotsFunc(fn func(ctx context.Context, query string, lang string, offset, limit int) ([]*entities.Spot, int, error)) {
	m.mu.Lock()
	m.searchSpotsFunc = fn
	m.mu.Unlock()
}

func (m *MockSpotRepository) SetUpdateRatingFunc(fn func(ctx context.Context, spotID string, averageRating float64, reviewCount int) error) {
	m.mu.Lock()
	m.updateRatingFunc = fn
	m.mu.Unlock()
}

// Reset clears all mock data and configurations
func (m *MockSpotRepository) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Clear all stored spots
	m.spots = make(map[string]*entities.Spot)
	
	// Reset all mock functions to nil
	m.createSpotFunc = nil
	m.getSpotFunc = nil
	m.getByCoordinatesFunc = nil
	m.listSpotsFunc = nil
	m.searchSpotsFunc = nil
	m.updateSpotFunc = nil
	m.deleteSpotFunc = nil
	m.updateRatingFunc = nil
}

// MockUserRepository provides a mock implementation for user testing
type MockUserRepository struct {
	mu                  sync.RWMutex
	users               map[string]*entities.User
	usersByEmail        map[string]*entities.User
	usersByAuthProvider map[string]*entities.User
	createUserFunc      func(ctx context.Context, user *entities.User) error
	getUserFunc         func(ctx context.Context, id string) (*entities.User, error)
	getByEmailFunc      func(ctx context.Context, email string) (*entities.User, error)
	getByAuthProviderFunc func(ctx context.Context, provider, providerID string) (*entities.User, error)
	updateUserFunc      func(ctx context.Context, user *entities.User) error
	deleteUserFunc      func(ctx context.Context, id string) error
}

// NewMockUserRepository creates a new mock user repository
func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:               make(map[string]*entities.User),
		usersByEmail:        make(map[string]*entities.User),
		usersByAuthProvider: make(map[string]*entities.User),
	}
}

// Implement the UserRepository interface
func (m *MockUserRepository) Create(ctx context.Context, user *entities.User) error {
	m.mu.RLock()
	createFunc := m.createUserFunc
	m.mu.RUnlock()
	
	if createFunc != nil {
		return createFunc(ctx, user)
	}
	
	if user.ID == "" {
		return stdErrors.New("user ID is required")
	}
	
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Check for duplicate email
	if _, exists := m.usersByEmail[user.Email]; exists {
		return fmt.Errorf("user with this email already exists: email=%s", user.Email)
	}
	
	// Check for duplicate auth provider combination
	authKey := buildAuthKey(user.AuthProvider, user.AuthProviderID)
	if _, exists := m.usersByAuthProvider[authKey]; exists {
		return fmt.Errorf("user with this auth provider already exists: provider=%s, providerID=%s", user.AuthProvider, user.AuthProviderID)
	}
	
	// Insert the user
	m.users[user.ID] = user
	m.usersByEmail[user.Email] = user
	m.usersByAuthProvider[authKey] = user
	
	return nil
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*entities.User, error) {
	m.mu.RLock()
	getFunc := m.getUserFunc
	user, exists := m.users[id]
	m.mu.RUnlock()
	
	if getFunc != nil {
		return getFunc(ctx, id)
	}
	
	if !exists {
		return nil, fmt.Errorf("user not found: id=%s", id)
	}
	
	return user, nil
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	m.mu.RLock()
	getByEmailFunc := m.getByEmailFunc
	user, exists := m.usersByEmail[email]
	m.mu.RUnlock()
	
	if getByEmailFunc != nil {
		return getByEmailFunc(ctx, email)
	}
	
	if !exists {
		return nil, fmt.Errorf("user not found: email=%s", email)
	}
	return user, nil
}

func (m *MockUserRepository) GetByAuthProvider(ctx context.Context, provider, providerID string) (*entities.User, error) {
	authKey := buildAuthKey(entities.AuthProvider(provider), providerID)
	
	m.mu.RLock()
	getByAuthProviderFunc := m.getByAuthProviderFunc
	user, exists := m.usersByAuthProvider[authKey]
	m.mu.RUnlock()
	
	if getByAuthProviderFunc != nil {
		return getByAuthProviderFunc(ctx, provider, providerID)
	}
	
	if !exists {
		return nil, fmt.Errorf("user not found: provider=%s, providerID=%s", provider, providerID)
	}
	return user, nil
}

func (m *MockUserRepository) validateUniqueConstraints(newUser, oldUser *entities.User) error {
	// Check for duplicate email (if different from current user's email)
	if newUser.Email != oldUser.Email {
		if existingUserByEmail, emailExists := m.usersByEmail[newUser.Email]; emailExists && existingUserByEmail.ID != newUser.ID {
			return fmt.Errorf("user with this email already exists: email=%s", newUser.Email)
		}
	}
	
	// Check for duplicate AuthProvider combination (if different from current user's)
	newAuthKey := buildAuthKey(newUser.AuthProvider, newUser.AuthProviderID)
	oldAuthKey := buildAuthKey(oldUser.AuthProvider, oldUser.AuthProviderID)
	if newAuthKey != oldAuthKey {
		if existingUserByAuth, authExists := m.usersByAuthProvider[newAuthKey]; authExists && existingUserByAuth.ID != newUser.ID {
			return fmt.Errorf("user with this auth provider already exists: provider=%s, providerID=%s", newUser.AuthProvider, newUser.AuthProviderID)
		}
	}
	
	return nil
}

func (m *MockUserRepository) updateIndexes(newUser, oldUser *entities.User) {
	// Remove old mappings
	delete(m.usersByEmail, oldUser.Email)
	oldAuthKey := buildAuthKey(oldUser.AuthProvider, oldUser.AuthProviderID)
	delete(m.usersByAuthProvider, oldAuthKey)
	
	// Update with new data
	m.users[newUser.ID] = newUser
	m.usersByEmail[newUser.Email] = newUser
	newAuthKey := buildAuthKey(newUser.AuthProvider, newUser.AuthProviderID)
	m.usersByAuthProvider[newAuthKey] = newUser
}

func (m *MockUserRepository) Update(ctx context.Context, user *entities.User) error {
	m.mu.RLock()
	updateFunc := m.updateUserFunc
	oldUser, exists := m.users[user.ID]
	m.mu.RUnlock()
	
	if updateFunc != nil {
		return updateFunc(ctx, user)
	}
	
	if !exists {
		return fmt.Errorf("user not found: id=%s", user.ID)
	}
	
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if err := m.validateUniqueConstraints(user, oldUser); err != nil {
		return err
	}
	
	m.updateIndexes(user, oldUser)
	
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	m.mu.RLock()
	deleteFunc := m.deleteUserFunc
	user, exists := m.users[id]
	m.mu.RUnlock()
	
	if deleteFunc != nil {
		return deleteFunc(ctx, id)
	}
	
	if !exists {
		return fmt.Errorf("user not found: id=%s", id)
	}
	
	m.mu.Lock()
	delete(m.users, id)
	delete(m.usersByEmail, user.Email)
	
	authKey := buildAuthKey(user.AuthProvider, user.AuthProviderID)
	delete(m.usersByAuthProvider, authKey)
	m.mu.Unlock()
	
	return nil
}

// Mock configuration methods
func (m *MockUserRepository) SetCreateUserFunc(fn func(ctx context.Context, user *entities.User) error) {
	m.mu.Lock()
	m.createUserFunc = fn
	m.mu.Unlock()
}

func (m *MockUserRepository) SetGetUserFunc(fn func(ctx context.Context, id string) (*entities.User, error)) {
	m.mu.Lock()
	m.getUserFunc = fn
	m.mu.Unlock()
}

func (m *MockUserRepository) SetGetByEmailFunc(fn func(ctx context.Context, email string) (*entities.User, error)) {
	m.mu.Lock()
	m.getByEmailFunc = fn
	m.mu.Unlock()
}

func (m *MockUserRepository) SetGetByAuthProviderFunc(fn func(ctx context.Context, provider, providerID string) (*entities.User, error)) {
	m.mu.Lock()
	m.getByAuthProviderFunc = fn
	m.mu.Unlock()
}

// SetUsers sets multiple users for testing, maintaining consistency across all lookup maps
func (m *MockUserRepository) SetUsers(users []*entities.User) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Clear existing data
	m.users = make(map[string]*entities.User)
	m.usersByEmail = make(map[string]*entities.User)
	m.usersByAuthProvider = make(map[string]*entities.User)
	
	// Add all users ensuring consistency
	for _, user := range users {
		if user.ID != "" {
			m.users[user.ID] = user
		}
		if user.Email != "" {
			m.usersByEmail[user.Email] = user
		}
		if user.AuthProvider != "" && user.AuthProviderID != "" {
			authKey := buildAuthKey(user.AuthProvider, user.AuthProviderID)
			m.usersByAuthProvider[authKey] = user
		}
	}
}

// Reset clears all mock data and configurations
func (m *MockUserRepository) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Clear all stored users
	m.users = make(map[string]*entities.User)
	m.usersByEmail = make(map[string]*entities.User)
	m.usersByAuthProvider = make(map[string]*entities.User)
	
	// Reset all mock functions to nil
	m.createUserFunc = nil
	m.getUserFunc = nil
	m.getByEmailFunc = nil
	m.getByAuthProviderFunc = nil
	m.updateUserFunc = nil
	m.deleteUserFunc = nil
}

// BehaviorDrivenMocks provides scenario-based mock configurations for behavior-driven testing.
// It encapsulates mock repositories for different entities and offers pre-configured scenarios
// like happy path and failure path testing to streamline test setup and improve test readability.
//
// The BehaviorDrivenMocks supports:
// - Happy path scenarios with realistic default data
// - Failure scenarios (repository errors, validation failures, etc.)
// - Custom behavior injection through function setters
// - Concurrent-safe operations with proper mutex handling
//
// Example usage:
//
//	// Basic instantiation
//	mocks := NewBehaviorDrivenMocks()
//	
//	// Configure happy path scenario with default values
//	mocks.ConfigureHappyPath(DefaultHappyPathConfig())
//	
//	// Configure happy path with custom values
//	config := HappyPathConfig{
//		SpotName:     "Custom Cafe",
//		UserEmail:    "user@example.com",
//		DisplayName:  "John Doe",
//		AuthProvider: "google",
//	}
//	mocks.ConfigureHappyPath(config)
//	
//	// Configure failure scenarios
//	mocks.ConfigureRepositoryFailures()
//	
//	// Use in tests
//	spotService := application.NewSpotService(mocks.SpotRepo, mocks.UserRepo)
//	result, err := spotService.CreateSpot(ctx, request)
type BehaviorDrivenMocks struct {
	SpotRepo *MockSpotRepository
	UserRepo *MockUserRepository
}

// NewBehaviorDrivenMocks creates mocks for behavior-driven scenarios
func NewBehaviorDrivenMocks() *BehaviorDrivenMocks {
	return &BehaviorDrivenMocks{
		SpotRepo: NewMockSpotRepository(),
		UserRepo: NewMockUserRepository(),
	}
}

// HappyPathConfig provides configuration options for happy path mock behavior
type HappyPathConfig struct {
	SpotName      string
	SpotLatitude  float64
	SpotLongitude float64
	SpotCategory  string
	UserEmail     string
	DisplayName   string
	AuthProvider  string
	Language      string
	DarkMode      bool
	Timezone      string
	
	// Dynamic response configurators
	SpotGenerator func(id string) *entities.Spot
	UserGenerator func(id string) *entities.User
}

// DefaultHappyPathConfig returns default configuration for happy path scenarios
func DefaultHappyPathConfig() HappyPathConfig {
	return HappyPathConfig{
		SpotName:      "Mock Spot",
		SpotLatitude:  35.6762,
		SpotLongitude: 139.6503,
		SpotCategory:  "cafe",
		UserEmail:     "test@example.com",
		DisplayName:   "Test User",
		AuthProvider:  "google",
		Language:      "en",
		DarkMode:      false,
		Timezone:      "UTC",
	}
}

// configureSpotMocks sets up spot-related mock behavior for happy path scenarios
func (bdm *BehaviorDrivenMocks) configureSpotMocks(config HappyPathConfig) {
	bdm.SpotRepo.SetCreateSpotFunc(func(ctx context.Context, spot *entities.Spot) error {
		return nil
	})
	
	bdm.SpotRepo.SetGetSpotFunc(func(ctx context.Context, id string) (*entities.Spot, error) {
		if config.SpotGenerator != nil {
			return config.SpotGenerator(id), nil
		}
		return &entities.Spot{
			ID:        id,
			Name:      config.SpotName,
			Latitude:  config.SpotLatitude,
			Longitude: config.SpotLongitude,
			Category:  config.SpotCategory,
		}, nil
	})
}

// configureUserMocks sets up user-related mock behavior for happy path scenarios
func (bdm *BehaviorDrivenMocks) configureUserMocks(config HappyPathConfig) {
	bdm.UserRepo.SetCreateUserFunc(func(ctx context.Context, user *entities.User) error {
		return nil
	})
	
	bdm.UserRepo.SetGetUserFunc(func(ctx context.Context, id string) (*entities.User, error) {
		if config.UserGenerator != nil {
			return config.UserGenerator(id), nil
		}
		return &entities.User{
			ID:             id,
			Email:          config.UserEmail,
			DisplayName:    config.DisplayName,
			AuthProvider:   entities.AuthProvider(config.AuthProvider),
			AuthProviderID: DefaultMockAuthProviderID,
			Preferences: entities.UserPreferences{
				Language: config.Language,
				DarkMode: config.DarkMode,
				Timezone: config.Timezone,
			},
		}, nil
	})
}

// ConfigureHappyPath sets up mocks for successful scenarios with optional configuration
func (bdm *BehaviorDrivenMocks) ConfigureHappyPath(configs ...HappyPathConfig) {
	config := DefaultHappyPathConfig()
	if len(configs) > 0 {
		config = configs[0]
	}

	bdm.configureSpotMocks(config)
	bdm.configureUserMocks(config)
}

// ConfigureHappyPathWithCustomSpots sets up mocks with custom spot data generation
func (bdm *BehaviorDrivenMocks) ConfigureHappyPathWithCustomSpots(spotGen func(id string) *entities.Spot) {
	config := DefaultHappyPathConfig()
	config.SpotGenerator = spotGen
	bdm.ConfigureHappyPath(config)
}

// ConfigureHappyPathWithCustomUsers sets up mocks with custom user data generation
func (bdm *BehaviorDrivenMocks) ConfigureHappyPathWithCustomUsers(userGen func(id string) *entities.User) {
	config := DefaultHappyPathConfig()
	config.UserGenerator = userGen
	bdm.ConfigureHappyPath(config)
}

// ConfigureHappyPathWithMultipleScenarios allows different responses based on ID patterns
func (bdm *BehaviorDrivenMocks) ConfigureHappyPathWithMultipleScenarios(spotScenarios map[string]*entities.Spot, userScenarios map[string]*entities.User) {
	config := DefaultHappyPathConfig()
	
	if len(spotScenarios) > 0 {
		config.SpotGenerator = func(id string) *entities.Spot {
			if spot, exists := spotScenarios[id]; exists {
				return spot
			}
			// Fallback to default
			return &entities.Spot{
				ID:        id,
				Name:      config.SpotName,
				Latitude:  config.SpotLatitude,
				Longitude: config.SpotLongitude,
				Category:  config.SpotCategory,
			}
		}
	}
	
	if len(userScenarios) > 0 {
		config.UserGenerator = func(id string) *entities.User {
			if user, exists := userScenarios[id]; exists {
				return user
			}
			// Fallback to default
			return &entities.User{
				ID:             id,
				Email:          config.UserEmail,
				DisplayName:    config.DisplayName,
				AuthProvider:   entities.AuthProvider(config.AuthProvider),
				AuthProviderID: DefaultMockAuthProviderID,
				Preferences: entities.UserPreferences{
					Language: config.Language,
					DarkMode: config.DarkMode,
					Timezone: config.Timezone,
				},
			}
		}
	}
	
	bdm.ConfigureHappyPath(config)
}

// ConfigureFailurePath sets up mocks for error scenarios
func (bdm *BehaviorDrivenMocks) ConfigureFailurePath() {
	// Spots fail
	bdm.SpotRepo.SetCreateSpotFunc(func(ctx context.Context, spot *entities.Spot) error {
		return fmt.Errorf("database connection failed for spot creation: spot_id=%s", spot.ID)
	})
	
	bdm.SpotRepo.SetGetSpotFunc(func(ctx context.Context, id string) (*entities.Spot, error) {
		return nil, fmt.Errorf("spot not found: id=%s", id)
	})
	
	// Users fail
	bdm.UserRepo.SetCreateUserFunc(func(ctx context.Context, user *entities.User) error {
		return fmt.Errorf("user creation failed: user_id=%s, email=%s", user.ID, user.Email)
	})
	
	bdm.UserRepo.SetGetUserFunc(func(ctx context.Context, id string) (*entities.User, error) {
		return nil, fmt.Errorf("user not found: id=%s", id)
	})
}

// ConfigurePartialFailure sets up mocks for mixed scenarios
func (bdm *BehaviorDrivenMocks) ConfigurePartialFailure() {
	// Spots succeed
	bdm.ConfigureHappyPath()
	
	// But users fail
	bdm.UserRepo.SetCreateUserFunc(func(ctx context.Context, user *entities.User) error {
		return fmt.Errorf("user creation failed in partial failure scenario: user_id=%s, email=%s", user.ID, user.Email)
	})
}

// Reset clears all mock configurations
func (bdm *BehaviorDrivenMocks) Reset() {
	bdm.SpotRepo.Reset()
	bdm.UserRepo.Reset()
}

// AuthTestData represents authentication test data structure
type AuthTestData struct {
	ValidToken      string
	InvalidToken    string
	ExpiredToken    string
	ValidUserID     string
	InvalidUserID   string
	TestUser        *TestUserInfo
	AdminUser       *TestUserInfo
}

// TestUserInfo represents test user information
type TestUserInfo struct {
	Email          string
	DisplayName    string
	AuthProvider   entities.AuthProvider
	AuthProviderID string
	Preferences    entities.UserPreferences
	IsAdmin        bool
	Permissions    []string // JWT permissions for authentication testing
	Role           string   // User role for authorization testing
}

// AuthHelper provides authentication testing utilities
type AuthHelper struct {
	// Default test tokens and user data
	mockValidator  *MockJWTValidator
	mockMiddleware *MockAuthMiddleware
	mockService    *MockAuthService
}

// NewAuthHelper creates a new authentication helper for testing
func NewAuthHelper() *AuthHelper {
	return &AuthHelper{
		mockValidator:  NewMockJWTValidator(),
		mockMiddleware: NewMockAuthMiddleware(),
		mockService:    NewMockAuthService(),
	}
}

// NewAuthTestData creates new authentication test data
func (ah *AuthHelper) NewAuthTestData() *AuthTestData {
	return &AuthTestData{
		ValidToken:   "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0LXVzZXItMTIzIiwiZW1haWwiOiJ0ZXN0QGV4YW1wbGUuY29tIiwiaXNzIjoiaHR0cHM6Ly90ZXN0LmF1dGgwLmNvbS8iLCJhdWQiOlsiYm9jY2hpLXRoZS1tYXAtYXBpIl0sImlhdCI6MTcwNDEwMDgwMCwiZXhwIjozNzEwMjIwODAwLCJzY29wZSI6InJlYWQ6c3BvdHMgd3JpdGU6c3BvdHMgcmVhZDpyZXZpZXdzIHdyaXRlOnJldmlld3MifQ.test-signature",
		InvalidToken: "invalid.token.value",
		ExpiredToken: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0LXVzZXItMTIzIiwiZW1haWwiOiJ0ZXN0QGV4YW1wbGUuY29tIiwiaXNzIjoiaHR0cHM6Ly90ZXN0LmF1dGgwLmNvbS8iLCJhdWQiOlsiYm9jY2hpLXRoZS1tYXAtYXBpIl0sImlhdCI6MTcwNDEwMDgwMCwiZXhwIjoxNzA0MTAwODAwLCJzY29wZSI6InJlYWQ6c3BvdHMgd3JpdGU6c3BvdHMgcmVhZDpyZXZpZXdzIHdyaXRlOnJldmlld3MifQ.expired-signature",
		ValidUserID:  "test-user-123",
		InvalidUserID: "non-existent-user",
		TestUser: &TestUserInfo{
			Email:          "test@example.com",
			DisplayName:    "Test User",
			AuthProvider:   entities.AuthProviderGoogle,
			AuthProviderID: "google_test_123",
			Preferences: entities.UserPreferences{
				Language: "en",
				DarkMode: false,
				Timezone: "UTC",
			},
			IsAdmin:     false,
			Permissions: []string{"read:spots", "write:reviews", "edit:profile"},
			Role:        "user",
		},
		AdminUser: &TestUserInfo{
			Email:          "admin@example.com",
			DisplayName:    "Admin User",
			AuthProvider:   entities.AuthProviderGoogle,
			AuthProviderID: "google_admin_456",
			Preferences: entities.UserPreferences{
				Language: "en",
				DarkMode: true,
				Timezone: "UTC",
			},
			IsAdmin:     true,
			Permissions: []string{"read:spots", "write:spots", "delete:spots", "admin:users", "admin:system"},
			Role:        "admin",
		},
	}
}

// CreateAuthenticatedContext creates a context with authentication data for testing
func (ah *AuthHelper) CreateAuthenticatedContext(baseCtx context.Context, userID string, userEmail string) context.Context {
	// In a real implementation, this would set up the context with authentication middleware data
	// For testing purposes, we simulate the context that would be created by the auth middleware
	ctx := context.WithValue(baseCtx, "user_id", userID)
	ctx = context.WithValue(ctx, "user_email", userEmail)
	ctx = context.WithValue(ctx, "user_info", map[string]interface{}{
		"sub":            userID,
		"email":          userEmail,
		"email_verified": true,
		"name":           "Test User",
		"picture":        "https://example.com/avatar.jpg",
		"scope":          "read:spots write:spots read:reviews write:reviews",
	})
	return ctx
}

// CreateUnauthenticatedContext creates a context without authentication data
func (ah *AuthHelper) CreateUnauthenticatedContext(baseCtx context.Context) context.Context {
	// Return the base context without any authentication data
	return baseCtx
}

// GetValidAuthHeaders returns valid authentication headers for HTTP requests
func (ah *AuthHelper) GetValidAuthHeaders(authData *AuthTestData) map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + authData.ValidToken,
		"Content-Type":  "application/json",
	}
}

// GetInvalidAuthHeaders returns invalid authentication headers for testing error scenarios
func (ah *AuthHelper) GetInvalidAuthHeaders(authData *AuthTestData) map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + authData.InvalidToken,
		"Content-Type":  "application/json",
	}
}

// GetExpiredAuthHeaders returns expired authentication headers for testing token expiration
func (ah *AuthHelper) GetExpiredAuthHeaders(authData *AuthTestData) map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + authData.ExpiredToken,
		"Content-Type":  "application/json",
	}
}

// GetMockValidator returns the mock JWT validator
func (ah *AuthHelper) GetMockValidator() *MockJWTValidator {
	return ah.mockValidator
}

// GetMockMiddleware returns the mock auth middleware
func (ah *AuthHelper) GetMockMiddleware() *MockAuthMiddleware {
	return ah.mockMiddleware
}

// GetMockService returns the mock auth service
func (ah *AuthHelper) GetMockService() *MockAuthService {
	return ah.mockService
}

// Reset resets all mocks to their default state
func (ah *AuthHelper) Reset() {
	ah.mockValidator.Reset()
	ah.mockMiddleware.Reset()
	ah.mockService.Reset()
}

// CreateTestUserContext creates a context with test user information
func (ah *AuthHelper) CreateTestUserContext(ctx context.Context, userID, email string, permissions []string) context.Context {
	userInfo := map[string]interface{}{
		"user_id":        userID,
		"email":          email,
		"email_verified": true,
		"name":           fmt.Sprintf("Test User %s", userID),
		"nickname":       userID,
		"picture":        "https://example.com/avatar.jpg",
		"permissions":    permissions,
	}

	ctx = context.WithValue(ctx, "user", userInfo)
	ctx = context.WithValue(ctx, "user_id", userID)
	ctx = context.WithValue(ctx, "email", email)

	return ctx
}

// CreateAdminContext creates a context with admin user permissions
func (ah *AuthHelper) CreateAdminContext(ctx context.Context) context.Context {
	return ah.CreateTestUserContext(ctx, "admin_user_123", "admin@example.com",
		[]string{"read:spots", "write:spots", "delete:spots", "admin:users", "admin:system"})
}

// CreateUserContext creates a context with regular user permissions
func (ah *AuthHelper) CreateUserContext(ctx context.Context) context.Context {
	return ah.CreateTestUserContext(ctx, "regular_user_123", "user@example.com",
		[]string{"read:spots", "write:reviews", "edit:profile"})
}

// CreateViewerContext creates a context with viewer-only permissions
func (ah *AuthHelper) CreateViewerContext(ctx context.Context) context.Context {
	return ah.CreateTestUserContext(ctx, "viewer_user_123", "viewer@example.com",
		[]string{"read:spots"})
}

// MockJWTValidator provides a mock implementation for JWT validation testing
type MockJWTValidator struct {
	mu                          sync.RWMutex
	validateTokenFunc           func(tokenString string) (*auth.Claims, error)
	validateTokenFromRequestFunc func(r *http.Request) (*auth.Claims, error)
	extractTokenFromRequestFunc func(r *http.Request) (string, error)
	getUserContextFunc          func(ctx context.Context, claims *auth.Claims) context.Context
}

// NewMockJWTValidator creates a new mock JWT validator
func NewMockJWTValidator() *MockJWTValidator {
	return &MockJWTValidator{}
}

// ValidateToken validates a JWT token (mock implementation)
func (m *MockJWTValidator) ValidateToken(tokenString string) (*auth.Claims, error) {
	m.mu.RLock()
	validateFunc := m.validateTokenFunc
	m.mu.RUnlock()

	if validateFunc != nil {
		return validateFunc(tokenString)
	}

	// Default behavior: parse test token or create mock claims
	if tokenString == "" {
		return nil, errors.Unauthorized("token is required")
	}

	// Handle special test tokens
	return m.parseTestToken(tokenString)
}

// ValidateTokenFromRequest validates a JWT token from HTTP request (mock implementation)
func (m *MockJWTValidator) ValidateTokenFromRequest(r *http.Request) (*auth.Claims, error) {
	m.mu.RLock()
	validateFromRequestFunc := m.validateTokenFromRequestFunc
	m.mu.RUnlock()

	if validateFromRequestFunc != nil {
		return validateFromRequestFunc(r)
	}

	// Extract token using mock extractor
	token, err := m.ExtractTokenFromRequest(r)
	if err != nil {
		return nil, err
	}

	return m.ValidateToken(token)
}

// ExtractTokenFromRequest extracts JWT token from HTTP request (mock implementation)
func (m *MockJWTValidator) ExtractTokenFromRequest(r *http.Request) (string, error) {
	m.mu.RLock()
	extractFunc := m.extractTokenFromRequestFunc
	m.mu.RUnlock()

	if extractFunc != nil {
		return extractFunc(r)
	}

	// Default extraction logic
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			return parts[1], nil
		}
		return "", errors.Unauthorized("authorization header format must be 'Bearer {token}'")
	}

	// Check query parameter as fallback
	token := r.URL.Query().Get("token")
	if token != "" {
		return token, nil
	}

	return "", errors.Unauthorized("no authorization token found")
}

// GetUserContext creates a context with user information from claims (mock implementation)
func (m *MockJWTValidator) GetUserContext(ctx context.Context, claims *auth.Claims) context.Context {
	m.mu.RLock()
	getUserContextFunc := m.getUserContextFunc
	m.mu.RUnlock()

	if getUserContextFunc != nil {
		return getUserContextFunc(ctx, claims)
	}

	// Default behavior: create context with user info
	userInfo := map[string]interface{}{
		"user_id":        claims.Subject,
		"email":          claims.Email,
		"email_verified": claims.EmailVerified,
		"name":           claims.Name,
		"nickname":       claims.Nickname,
		"picture":        claims.Picture,
		"permissions":    claims.Permissions,
	}

	// Set both "user" object and individual context keys for compatibility
	ctx = context.WithValue(ctx, "user", userInfo)
	ctx = context.WithValue(ctx, "user_id", claims.Subject)
	ctx = context.WithValue(ctx, "email", claims.Email)

	return ctx
}

// parseTestToken parses special test tokens and returns mock claims
func (m *MockJWTValidator) parseTestToken(tokenString string) (*auth.Claims, error) {
	switch {
	case strings.HasPrefix(tokenString, "test_valid_"):
		return m.createValidTestClaims(tokenString), nil
	case strings.HasPrefix(tokenString, "test_expired_"):
		return nil, errors.Unauthorized("token has expired")
	case strings.HasPrefix(tokenString, "test_invalid_"):
		return nil, errors.Unauthorized("invalid token")
	case tokenString == "test_admin_token":
		return m.createAdminTestClaims(), nil
	case tokenString == "test_user_token":
		return m.createUserTestClaims(), nil
	case tokenString == "test_viewer_token":
		return m.createViewerTestClaims(), nil
	default:
		// Try to parse as real JWT (for integration tests)
		return m.parseRealJWT(tokenString)
	}
}

// createValidTestClaims creates valid test claims for a given token
func (m *MockJWTValidator) createValidTestClaims(tokenString string) *auth.Claims {
	userID := strings.TrimPrefix(tokenString, "test_valid_")
	if userID == "" {
		userID = "test_user_123"
	}

	return &auth.Claims{
		Subject:       userID,
		Email:         fmt.Sprintf("%s@example.com", userID),
		EmailVerified: true,
		Name:          fmt.Sprintf("Test User %s", userID),
		Nickname:      userID,
		Picture:       "https://example.com/avatar.jpg",
		Permissions:   []string{"read:spots", "write:reviews"},
		Audience:      []string{"bocchi-the-map-api"},
		Issuer:        "https://test.auth0.com/",
		ExpiresAt:     time.Now().Add(time.Hour).Unix(),
		IssuedAt:      time.Now().Unix(),
		NotBefore:     time.Now().Unix(),
	}
}

// createAdminTestClaims creates test claims with admin permissions
func (m *MockJWTValidator) createAdminTestClaims() *auth.Claims {
	return &auth.Claims{
		Subject:       "admin_user_123",
		Email:         "admin@example.com",
		EmailVerified: true,
		Name:          "Admin User",
		Nickname:      "admin",
		Picture:       "https://example.com/admin_avatar.jpg",
		Permissions:   []string{"read:spots", "write:spots", "delete:spots", "admin:users", "admin:system"},
		Audience:      []string{"bocchi-the-map-api"},
		Issuer:        "https://test.auth0.com/",
		ExpiresAt:     time.Now().Add(time.Hour).Unix(),
		IssuedAt:      time.Now().Unix(),
		NotBefore:     time.Now().Unix(),
	}
}

// createUserTestClaims creates test claims with regular user permissions
func (m *MockJWTValidator) createUserTestClaims() *auth.Claims {
	return &auth.Claims{
		Subject:       "regular_user_123",
		Email:         "user@example.com",
		EmailVerified: true,
		Name:          "Regular User",
		Nickname:      "user",
		Picture:       "https://example.com/user_avatar.jpg",
		Permissions:   []string{"read:spots", "write:reviews", "edit:profile"},
		Audience:      []string{"bocchi-the-map-api"},
		Issuer:        "https://test.auth0.com/",
		ExpiresAt:     time.Now().Add(time.Hour).Unix(),
		IssuedAt:      time.Now().Unix(),
		NotBefore:     time.Now().Unix(),
	}
}

// createViewerTestClaims creates test claims with viewer-only permissions
func (m *MockJWTValidator) createViewerTestClaims() *auth.Claims {
	return &auth.Claims{
		Subject:       "viewer_user_123",
		Email:         "viewer@example.com",
		EmailVerified: true,
		Name:          "Viewer User",
		Nickname:      "viewer",
		Picture:       "https://example.com/viewer_avatar.jpg",
		Permissions:   []string{"read:spots"},
		Audience:      []string{"bocchi-the-map-api"},
		Issuer:        "https://test.auth0.com/",
		ExpiresAt:     time.Now().Add(time.Hour).Unix(),
		IssuedAt:      time.Now().Unix(),
		NotBefore:     time.Now().Unix(),
	}
}

// parseRealJWT attempts to parse a real JWT token (for integration testing)
func (m *MockJWTValidator) parseRealJWT(tokenString string) (*auth.Claims, error) {
	// This is a simplified parser for testing purposes
	// In real tests, you might want to use the actual Auth0 validator
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &auth.Claims{})
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrTypeUnauthorized, "failed to parse JWT token")
	}

	claims, ok := token.Claims.(*auth.Claims)
	if !ok {
		return nil, errors.Unauthorized("invalid token claims")
	}

	return claims, nil
}

// Mock configuration methods for MockJWTValidator
func (m *MockJWTValidator) SetValidateTokenFunc(fn func(tokenString string) (*auth.Claims, error)) {
	m.mu.Lock()
	m.validateTokenFunc = fn
	m.mu.Unlock()
}

func (m *MockJWTValidator) SetValidateTokenFromRequestFunc(fn func(r *http.Request) (*auth.Claims, error)) {
	m.mu.Lock()
	m.validateTokenFromRequestFunc = fn
	m.mu.Unlock()
}

func (m *MockJWTValidator) SetExtractTokenFromRequestFunc(fn func(r *http.Request) (string, error)) {
	m.mu.Lock()
	m.extractTokenFromRequestFunc = fn
	m.mu.Unlock()
}

func (m *MockJWTValidator) SetGetUserContextFunc(fn func(ctx context.Context, claims *auth.Claims) context.Context) {
	m.mu.Lock()
	m.getUserContextFunc = fn
	m.mu.Unlock()
}

// Reset clears all mock configurations for MockJWTValidator
func (m *MockJWTValidator) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.validateTokenFunc = nil
	m.validateTokenFromRequestFunc = nil
	m.extractTokenFromRequestFunc = nil
	m.getUserContextFunc = nil
}

// MockAuthMiddleware provides a mock implementation for authentication middleware testing
type MockAuthMiddleware struct {
	mu                     sync.RWMutex
	requireAuthFunc        func() func(http.Handler) http.Handler
	optionalAuthFunc       func() func(http.Handler) http.Handler
	validateRequestFunc    func(r *http.Request) (*auth.Claims, error)
	getValidatorFunc       func() *MockJWTValidator
	shouldSkipPath         func(path string) bool
	development            bool
	skipPaths              map[string]bool
}

// NewMockAuthMiddleware creates a new mock auth middleware
func NewMockAuthMiddleware() *MockAuthMiddleware {
	return &MockAuthMiddleware{
		skipPaths: make(map[string]bool),
	}
}

// RequireAuth returns a Chi middleware that requires authentication (mock implementation)
func (m *MockAuthMiddleware) RequireAuth() func(http.Handler) http.Handler {
	m.mu.RLock()
	requireFunc := m.requireAuthFunc
	m.mu.RUnlock()

	if requireFunc != nil {
		return requireFunc()
	}

	// Default mock behavior: create middleware that validates test tokens
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if path should be skipped
			if m.shouldSkipTestPath(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			// Mock token validation
			claims, err := m.validateTestRequest(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			// Add user context to request
			userCtx := m.createUserContext(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(userCtx))
		})
	}
}

// OptionalAuth returns a Chi middleware that optionally validates authentication (mock implementation)
func (m *MockAuthMiddleware) OptionalAuth() func(http.Handler) http.Handler {
	m.mu.RLock()
	optionalFunc := m.optionalAuthFunc
	m.mu.RUnlock()

	if optionalFunc != nil {
		return optionalFunc()
	}

	// Default mock behavior: validate if token present, but don't block if missing
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := m.validateTestRequest(r)
			if err == nil && claims != nil {
				// Add user context if validation succeeded
				userCtx := m.createUserContext(r.Context(), claims)
				next.ServeHTTP(w, r.WithContext(userCtx))
				return
			}

			// Continue without authentication
			next.ServeHTTP(w, r)
		})
	}
}

// GetValidator returns the mock JWT validator
func (m *MockAuthMiddleware) GetValidator() *MockJWTValidator {
	m.mu.RLock()
	getValidatorFunc := m.getValidatorFunc
	m.mu.RUnlock()

	if getValidatorFunc != nil {
		return getValidatorFunc()
	}

	// Return a default mock validator
	return NewMockJWTValidator()
}

// shouldSkipTestPath checks if the given path should skip authentication in tests
func (m *MockAuthMiddleware) shouldSkipTestPath(path string) bool {
	m.mu.RLock()
	shouldSkipFunc := m.shouldSkipPath
	skipPaths := m.skipPaths
	m.mu.RUnlock()

	if shouldSkipFunc != nil {
		return shouldSkipFunc(path)
	}

	// Check exact match
	if skipPaths[path] {
		return true
	}

	// Check for path prefixes that should be skipped
	skipPrefixes := []string{
		"/health",
		"/metrics",
		"/debug",
		"/swagger",
		"/docs",
		"/test",
	}

	for _, prefix := range skipPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}

	return false
}

// validateTestRequest validates a test request and returns claims
func (m *MockAuthMiddleware) validateTestRequest(r *http.Request) (*auth.Claims, error) {
	m.mu.RLock()
	validateFunc := m.validateRequestFunc
	m.mu.RUnlock()

	if validateFunc != nil {
		return validateFunc(r)
	}

	// Default behavior: extract and validate test token
	validator := NewMockJWTValidator()
	return validator.ValidateTokenFromRequest(r)
}

// createUserContext creates a context with user information from claims
func (m *MockAuthMiddleware) createUserContext(ctx context.Context, claims *auth.Claims) context.Context {
	if claims == nil {
		return ctx
	}

	userInfo := map[string]interface{}{
		"user_id":        claims.Subject,
		"email":          claims.Email,
		"email_verified": claims.EmailVerified,
		"name":           claims.Name,
		"nickname":       claims.Nickname,
		"picture":        claims.Picture,
		"permissions":    claims.Permissions,
	}

	// Set both "user" object and individual context keys for compatibility
	ctx = context.WithValue(ctx, "user", userInfo)
	ctx = context.WithValue(ctx, "user_id", claims.Subject)
	ctx = context.WithValue(ctx, "email", claims.Email)

	return ctx
}

// Mock configuration methods for MockAuthMiddleware
func (m *MockAuthMiddleware) SetRequireAuthFunc(fn func() func(http.Handler) http.Handler) {
	m.mu.Lock()
	m.requireAuthFunc = fn
	m.mu.Unlock()
}

func (m *MockAuthMiddleware) SetOptionalAuthFunc(fn func() func(http.Handler) http.Handler) {
	m.mu.Lock()
	m.optionalAuthFunc = fn
	m.mu.Unlock()
}

func (m *MockAuthMiddleware) SetValidateRequestFunc(fn func(r *http.Request) (*auth.Claims, error)) {
	m.mu.Lock()
	m.validateRequestFunc = fn
	m.mu.Unlock()
}

func (m *MockAuthMiddleware) SetGetValidatorFunc(fn func() *MockJWTValidator) {
	m.mu.Lock()
	m.getValidatorFunc = fn
	m.mu.Unlock()
}

func (m *MockAuthMiddleware) SetShouldSkipPathFunc(fn func(path string) bool) {
	m.mu.Lock()
	m.shouldSkipPath = fn
	m.mu.Unlock()
}

func (m *MockAuthMiddleware) SetSkipPaths(paths []string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.skipPaths = make(map[string]bool)
	for _, path := range paths {
		m.skipPaths[path] = true
	}
}

func (m *MockAuthMiddleware) SetDevelopment(dev bool) {
	m.mu.Lock()
	m.development = dev
	m.mu.Unlock()
}

// Reset clears all mock configurations for MockAuthMiddleware
func (m *MockAuthMiddleware) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.requireAuthFunc = nil
	m.optionalAuthFunc = nil
	m.validateRequestFunc = nil
	m.getValidatorFunc = nil
	m.shouldSkipPath = nil
	m.development = false
	m.skipPaths = make(map[string]bool)
}

// MockAuthService provides a mock implementation for the authentication service
type MockAuthService struct {
	mu                    sync.RWMutex
	validateTokenFunc     func(ctx context.Context, token string) (*auth.Claims, error)
	checkPermissionFunc   func(ctx context.Context, permission string) error
	getUserInfoFunc       func(ctx context.Context) (map[string]interface{}, error)
	requireUserFunc       func(ctx context.Context) (map[string]interface{}, error)
	requireUserIDFunc     func(ctx context.Context) (string, error)
	requireUserEmailFunc  func(ctx context.Context) (string, error)
	healthFunc            func(ctx context.Context) error
	getStatsFunc          func() map[string]interface{}
}

// NewMockAuthService creates a new mock auth service
func NewMockAuthService() *MockAuthService {
	return &MockAuthService{}
}

// ValidateToken validates a JWT token and returns claims (mock implementation)
func (m *MockAuthService) ValidateToken(ctx context.Context, token string) (*auth.Claims, error) {
	m.mu.RLock()
	validateFunc := m.validateTokenFunc
	m.mu.RUnlock()

	if validateFunc != nil {
		return validateFunc(ctx, token)
	}

	// Default behavior: use mock validator
	validator := NewMockJWTValidator()
	return validator.ValidateToken(token)
}

// CheckPermission validates that a user has a specific permission (mock implementation)
func (m *MockAuthService) CheckPermission(ctx context.Context, permission string) error {
	m.mu.RLock()
	checkFunc := m.checkPermissionFunc
	m.mu.RUnlock()

	if checkFunc != nil {
		return checkFunc(ctx, permission)
	}

	// Default behavior: check if user has permission in context
	if !auth.HasPermission(ctx, permission) {
		return errors.Forbidden("permission", permission)
	}
	return nil
}

// GetUserInfo extracts user information from the request context (mock implementation)
func (m *MockAuthService) GetUserInfo(ctx context.Context) (map[string]interface{}, error) {
	m.mu.RLock()
	getUserFunc := m.getUserInfoFunc
	m.mu.RUnlock()

	if getUserFunc != nil {
		return getUserFunc(ctx)
	}

	// Default behavior: get user from context
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, errors.Unauthorized("user context not found")
	}
	return user, nil
}

// RequireUser ensures a user is authenticated and returns user info (mock implementation)
func (m *MockAuthService) RequireUser(ctx context.Context) (map[string]interface{}, error) {
	m.mu.RLock()
	requireFunc := m.requireUserFunc
	m.mu.RUnlock()

	if requireFunc != nil {
		return requireFunc(ctx)
	}

	return m.GetUserInfo(ctx)
}

// RequireUserID ensures a user is authenticated and returns the user ID (mock implementation)
func (m *MockAuthService) RequireUserID(ctx context.Context) (string, error) {
	m.mu.RLock()
	requireIDFunc := m.requireUserIDFunc
	m.mu.RUnlock()

	if requireIDFunc != nil {
		return requireIDFunc(ctx)
	}

	// Default behavior: get user ID from context
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return "", errors.Unauthorized("user ID not found in context")
	}
	return userID, nil
}

// RequireUserEmail ensures a user is authenticated and returns the user email (mock implementation)
func (m *MockAuthService) RequireUserEmail(ctx context.Context) (string, error) {
	m.mu.RLock()
	requireEmailFunc := m.requireUserEmailFunc
	m.mu.RUnlock()

	if requireEmailFunc != nil {
		return requireEmailFunc(ctx)
	}

	// Default behavior: get user email from context
	email, ok := auth.GetUserEmailFromContext(ctx)
	if !ok {
		return "", errors.Unauthorized("user email not found in context")
	}
	return email, nil
}

// Health checks the health of the authentication service (mock implementation)
func (m *MockAuthService) Health(ctx context.Context) error {
	m.mu.RLock()
	healthFunc := m.healthFunc
	m.mu.RUnlock()

	if healthFunc != nil {
		return healthFunc(ctx)
	}

	// Default behavior: always healthy
	return nil
}

// GetStats returns statistics about the authentication service (mock implementation)
func (m *MockAuthService) GetStats() map[string]interface{} {
	m.mu.RLock()
	getStatsFunc := m.getStatsFunc
	m.mu.RUnlock()

	if getStatsFunc != nil {
		return getStatsFunc()
	}

	// Default behavior: return mock stats
	return map[string]interface{}{
		"validator":     "mock",
		"middleware":    "mock",
		"rate_limiter":  "mock",
		"healthy":       true,
	}
}

// Mock configuration methods for MockAuthService
func (m *MockAuthService) SetValidateTokenFunc(fn func(ctx context.Context, token string) (*auth.Claims, error)) {
	m.mu.Lock()
	m.validateTokenFunc = fn
	m.mu.Unlock()
}

func (m *MockAuthService) SetCheckPermissionFunc(fn func(ctx context.Context, permission string) error) {
	m.mu.Lock()
	m.checkPermissionFunc = fn
	m.mu.Unlock()
}

func (m *MockAuthService) SetGetUserInfoFunc(fn func(ctx context.Context) (map[string]interface{}, error)) {
	m.mu.Lock()
	m.getUserInfoFunc = fn
	m.mu.Unlock()
}

func (m *MockAuthService) SetRequireUserFunc(fn func(ctx context.Context) (map[string]interface{}, error)) {
	m.mu.Lock()
	m.requireUserFunc = fn
	m.mu.Unlock()
}

func (m *MockAuthService) SetRequireUserIDFunc(fn func(ctx context.Context) (string, error)) {
	m.mu.Lock()
	m.requireUserIDFunc = fn
	m.mu.Unlock()
}

func (m *MockAuthService) SetRequireUserEmailFunc(fn func(ctx context.Context) (string, error)) {
	m.mu.Lock()
	m.requireUserEmailFunc = fn
	m.mu.Unlock()
}

func (m *MockAuthService) SetHealthFunc(fn func(ctx context.Context) error) {
	m.mu.Lock()
	m.healthFunc = fn
	m.mu.Unlock()
}

func (m *MockAuthService) SetGetStatsFunc(fn func() map[string]interface{}) {
	m.mu.Lock()
	m.getStatsFunc = fn
	m.mu.Unlock()
}

// Reset clears all mock configurations for MockAuthService
func (m *MockAuthService) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.validateTokenFunc = nil
	m.checkPermissionFunc = nil
	m.getUserInfoFunc = nil
	m.requireUserFunc = nil
	m.requireUserIDFunc = nil
	m.requireUserEmailFunc = nil
	m.healthFunc = nil
	m.getStatsFunc = nil
}