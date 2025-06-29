package helpers

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
	"sync"

	"github.com/necofuryai/bocchi-the-map/api/domain/entities"
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
		return errors.New("spot.ID cannot be empty")
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
		return errors.New("user ID is required")
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