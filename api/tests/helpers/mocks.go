package helpers

import (
	"context"
	"errors"
	"fmt"
	"sort"
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

// MockSpotRepository provides a mock implementation for testing
type MockSpotRepository struct {
	mu              sync.RWMutex
	spots           map[string]*entities.Spot
	createSpotFunc  func(ctx context.Context, spot *entities.Spot) error
	getSpotFunc     func(ctx context.Context, id string) (*entities.Spot, error)
	listSpotsFunc   func(ctx context.Context, filters map[string]interface{}) ([]*entities.Spot, error)
	updateSpotFunc  func(ctx context.Context, spot *entities.Spot) error
	deleteSpotFunc  func(ctx context.Context, id string) error
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
		return errors.New("spot ID is required")
	}
	
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Check if spot with the same ID already exists
	if _, exists := m.spots[spot.ID]; exists {
		return errors.New("spot with ID already exists")
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

func (m *MockSpotRepository) List(ctx context.Context, filters map[string]interface{}) ([]*entities.Spot, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	if m.listSpotsFunc != nil {
		return m.listSpotsFunc(ctx, filters)
	}
	
	result := make([]*entities.Spot, 0, len(m.spots))
	for _, spot := range m.spots {
		result = append(result, spot)
	}
	
	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})
	
	return result, nil
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

func (m *MockSpotRepository) SetListSpotsFunc(fn func(ctx context.Context, filters map[string]interface{}) ([]*entities.Spot, error)) {
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

// Reset clears all mock data and configurations
func (m *MockSpotRepository) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Clear all stored spots  
	m.spots = make(map[string]*proto.Spot)
	
	// Reset all mock functions to nil
	m.createSpotFunc = nil
	m.getSpotFunc = nil
	m.listSpotsFunc = nil
	m.updateSpotFunc = nil
	m.deleteSpotFunc = nil
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
		return errors.New("user with this email already exists")
	}
	
	// Check for duplicate auth provider combination
	authKey := buildAuthKey(user.AuthProvider, user.AuthProviderID)
	if _, exists := m.usersByAuthProvider[authKey]; exists {
		return errors.New("user with this auth provider already exists")
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
	
	// Check for duplicate email (if different from current user's email)
	if user.Email != oldUser.Email {
		if existingUserByEmail, emailExists := m.usersByEmail[user.Email]; emailExists && existingUserByEmail.ID != user.ID {
			return errors.New("user with this email already exists")
		}
	}
	
	// Check for duplicate AuthProvider combination (if different from current user's)
	newAuthKey := buildAuthKey(user.AuthProvider, user.AuthProviderID)
	oldAuthKey := buildAuthKey(oldUser.AuthProvider, oldUser.AuthProviderID)
	if newAuthKey != oldAuthKey {
		if existingUserByAuth, authExists := m.usersByAuthProvider[newAuthKey]; authExists && existingUserByAuth.ID != user.ID {
			return errors.New("user with this auth provider already exists")
		}
	}
	
	// Remove old mappings
	delete(m.usersByEmail, oldUser.Email)
	delete(m.usersByAuthProvider, oldAuthKey)
	
	// Update with new data
	m.users[user.ID] = user
	m.usersByEmail[user.Email] = user
	m.usersByAuthProvider[newAuthKey] = user
	
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
		return errors.New("user not found")
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
	for k := range m.users {
		delete(m.users, k)
	}
	for k := range m.usersByEmail {
		delete(m.usersByEmail, k)
	}
	for k := range m.usersByAuthProvider {
		delete(m.usersByAuthProvider, k)
	}
	
	// Reset all mock functions to nil
	m.createUserFunc = nil
	m.getUserFunc = nil
	m.getByEmailFunc = nil
	m.getByAuthProviderFunc = nil
	m.updateUserFunc = nil
	m.deleteUserFunc = nil
}

// BehaviorDrivenMocks provides scenario-based mock configurations
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
		return errors.New("database connection failed")
	})
	
	bdm.SpotRepo.SetGetSpotFunc(func(ctx context.Context, id string) (*entities.Spot, error) {
		return nil, errors.New("spot not found")
	})
	
	// Users fail
	bdm.UserRepo.SetCreateUserFunc(func(ctx context.Context, user *entities.User) error {
		return errors.New("user creation failed")
	})
	
	bdm.UserRepo.SetGetUserFunc(func(ctx context.Context, id string) (*entities.User, error) {
		return nil, errors.New("user not found")
	})
}

// ConfigurePartialFailure sets up mocks for mixed scenarios
func (bdm *BehaviorDrivenMocks) ConfigurePartialFailure() {
	// Spots succeed
	bdm.ConfigureHappyPath()
	
	// But users fail
	bdm.UserRepo.SetCreateUserFunc(func(ctx context.Context, user *entities.User) error {
		return errors.New("user creation failed")
	})
}

// Reset clears all mock configurations
func (bdm *BehaviorDrivenMocks) Reset() {
	bdm.SpotRepo.Reset()
	bdm.UserRepo.Reset()
}