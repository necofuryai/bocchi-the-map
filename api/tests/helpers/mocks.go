package helpers

import (
	"context"
	"errors"

	"github.com/necofuryai/bocchi-the-map/api/domain/entities"
)

// MockSpotRepository provides a mock implementation for testing
type MockSpotRepository struct {
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
	if m.createSpotFunc != nil {
		return m.createSpotFunc(ctx, spot)
	}
	
	if spot.ID == "" {
		return errors.New("spot ID is required")
	}
	
	m.spots[spot.ID] = spot
	return nil
}

func (m *MockSpotRepository) GetByID(ctx context.Context, id string) (*entities.Spot, error) {
	if m.getSpotFunc != nil {
		return m.getSpotFunc(ctx, id)
	}
	
	spot, exists := m.spots[id]
	if !exists {
		return nil, errors.New("spot not found")
	}
	
	return spot, nil
}

func (m *MockSpotRepository) List(ctx context.Context, filters map[string]interface{}) ([]*entities.Spot, error) {
	if m.listSpotsFunc != nil {
		return m.listSpotsFunc(ctx, filters)
	}
	
	var result []*entities.Spot
	for _, spot := range m.spots {
		result = append(result, spot)
	}
	
	return result, nil
}

func (m *MockSpotRepository) Update(ctx context.Context, spot *entities.Spot) error {
	if m.updateSpotFunc != nil {
		return m.updateSpotFunc(ctx, spot)
	}
	
	if _, exists := m.spots[spot.ID]; !exists {
		return errors.New("spot not found")
	}
	
	m.spots[spot.ID] = spot
	return nil
}

func (m *MockSpotRepository) Delete(ctx context.Context, id string) error {
	if m.deleteSpotFunc != nil {
		return m.deleteSpotFunc(ctx, id)
	}
	
	if _, exists := m.spots[id]; !exists {
		return errors.New("spot not found")
	}
	
	delete(m.spots, id)
	return nil
}

// Mock configuration methods
func (m *MockSpotRepository) SetCreateSpotFunc(fn func(ctx context.Context, spot *entities.Spot) error) {
	m.createSpotFunc = fn
}

func (m *MockSpotRepository) SetGetSpotFunc(fn func(ctx context.Context, id string) (*entities.Spot, error)) {
	m.getSpotFunc = fn
}

func (m *MockSpotRepository) SetListSpotsFunc(fn func(ctx context.Context, filters map[string]interface{}) ([]*entities.Spot, error)) {
	m.listSpotsFunc = fn
}

func (m *MockSpotRepository) SetUpdateSpotFunc(fn func(ctx context.Context, spot *entities.Spot) error) {
	m.updateSpotFunc = fn
}

func (m *MockSpotRepository) SetDeleteSpotFunc(fn func(ctx context.Context, id string) error) {
	m.deleteSpotFunc = fn
}

// MockUserRepository provides a mock implementation for user testing
type MockUserRepository struct {
	users           map[string]*entities.User
	createUserFunc  func(ctx context.Context, user *entities.User) error
	getUserFunc     func(ctx context.Context, id string) (*entities.User, error)
	updateUserFunc  func(ctx context.Context, user *entities.User) error
	deleteUserFunc  func(ctx context.Context, id string) error
}

// NewMockUserRepository creates a new mock user repository
func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]*entities.User),
	}
}

// Implement the UserRepository interface
func (m *MockUserRepository) Create(ctx context.Context, user *entities.User) error {
	if m.createUserFunc != nil {
		return m.createUserFunc(ctx, user)
	}
	
	if user.ID == "" {
		return errors.New("user ID is required")
	}
	
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*entities.User, error) {
	if m.getUserFunc != nil {
		return m.getUserFunc(ctx, id)
	}
	
	user, exists := m.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	
	return user, nil
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *MockUserRepository) GetByAuthProvider(ctx context.Context, provider, providerID string) (*entities.User, error) {
	for _, user := range m.users {
		if string(user.AuthProvider) == provider && user.AuthProviderID == providerID {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *MockUserRepository) Update(ctx context.Context, user *entities.User) error {
	if m.updateUserFunc != nil {
		return m.updateUserFunc(ctx, user)
	}
	
	if _, exists := m.users[user.ID]; !exists {
		return errors.New("user not found")
	}
	
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	if m.deleteUserFunc != nil {
		return m.deleteUserFunc(ctx, id)
	}
	
	if _, exists := m.users[id]; !exists {
		return errors.New("user not found")
	}
	
	delete(m.users, id)
	return nil
}

// Mock configuration methods
func (m *MockUserRepository) SetCreateUserFunc(fn func(ctx context.Context, user *entities.User) error) {
	m.createUserFunc = fn
}

func (m *MockUserRepository) SetGetUserFunc(fn func(ctx context.Context, id string) (*entities.User, error)) {
	m.getUserFunc = fn
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

// ConfigureHappyPath sets up mocks for successful scenarios
func (bdm *BehaviorDrivenMocks) ConfigureHappyPath() {
	// Spots always succeed
	bdm.SpotRepo.SetCreateSpotFunc(func(ctx context.Context, spot *entities.Spot) error {
		return nil
	})
	
	bdm.SpotRepo.SetGetSpotFunc(func(ctx context.Context, id string) (*entities.Spot, error) {
		return &entities.Spot{
			ID:        id,
			Name:      "Mock Spot",
			Latitude:  35.6762,
			Longitude: 139.6503,
			Category:  "cafe",
		}, nil
	})
	
	// Users always succeed
	bdm.UserRepo.SetCreateUserFunc(func(ctx context.Context, user *entities.User) error {
		return nil
	})
	
	bdm.UserRepo.SetGetUserFunc(func(ctx context.Context, id string) (*entities.User, error) {
		return &entities.User{
			ID:             id,
			Email:          "test@example.com",
			DisplayName:    "Test User",
			AuthProvider:   "google",
			AuthProviderID: "mock_123",
			Preferences: entities.UserPreferences{
				Language: "en",
				DarkMode: false,
				Timezone: "UTC",
			},
		}, nil
	})
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
	bdm.SpotRepo = NewMockSpotRepository()
	bdm.UserRepo = NewMockUserRepository()
}