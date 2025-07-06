package application_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/masyusakai/bocchi-the-map/api/internal/application"
	"github.com/masyusakai/bocchi-the-map/api/internal/domain/entities"
	"github.com/masyusakai/bocchi-the-map/api/internal/domain/rating"
	"github.com/masyusakai/bocchi-the-map/api/pb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock repositories for TDD
type MockRatingRepository struct {
	mock.Mock
}

func (m *MockRatingRepository) Create(ctx context.Context, rating *rating.Rating) error {
	args := m.Called(ctx, rating)
	return args.Error(0)
}

func (m *MockRatingRepository) Update(ctx context.Context, rating *rating.Rating) error {
	args := m.Called(ctx, rating)
	return args.Error(0)
}

func (m *MockRatingRepository) GetBySpotAndUser(ctx context.Context, spotID, userID string) (*rating.Rating, error) {
	args := m.Called(ctx, spotID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*rating.Rating), args.Error(1)
}

func (m *MockRatingRepository) GetBySpot(ctx context.Context, spotID string) ([]*rating.Rating, error) {
	args := m.Called(ctx, spotID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*rating.Rating), args.Error(1)
}

type MockSpotRepository struct {
	mock.Mock
}

func (m *MockSpotRepository) GetByID(ctx context.Context, id string) (*entities.Spot, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Spot), args.Error(1)
}

func (m *MockSpotRepository) UpdateSoloFriendlyStats(ctx context.Context, spotID string, avgRating float64, totalRatings int) error {
	args := m.Called(ctx, spotID, avgRating, totalRatings)
	return args.Error(0)
}

// TDD Service Layer Tests
func TestRatingService_CreateRating(t *testing.T) {
	// Arrange
	mockRatingRepo := &MockRatingRepository{}
	mockSpotRepo := &MockSpotRepository{}
	service := application.NewRatingService(mockRatingRepo, mockSpotRepo)

	ctx := context.Background()
	testSpot := &entities.Spot{
		ID:   "spot-123",
		Name: "Test Cafe",
	}

	request := &pb.CreateRatingRequest{
		SpotId:             "spot-123",
		UserId:             "user-456",
		SoloFriendlyRating: 5,
		Categories:         []string{"quiet_atmosphere", "wifi_available"},
		Comment:            "Great spot for solo work",
	}

	// Setup mocks
	mockSpotRepo.On("GetByID", ctx, "spot-123").Return(testSpot, nil)
	mockRatingRepo.On("GetBySpotAndUser", ctx, "spot-123", "user-456").Return(nil, errors.New("not found"))
	mockRatingRepo.On("Create", ctx, mock.AnythingOfType("*rating.Rating")).Return(nil)
	mockRatingRepo.On("GetBySpot", ctx, "spot-123").Return([]*rating.Rating{}, nil)
	mockSpotRepo.On("UpdateSoloFriendlyStats", ctx, "spot-123", 5.0, 1).Return(nil)

	// Act
	result, err := service.CreateRating(ctx, request)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "spot-123", result.SpotId)
	assert.Equal(t, "user-456", result.UserId)
	assert.Equal(t, int32(5), result.SoloFriendlyRating)
	assert.Equal(t, []string{"quiet_atmosphere", "wifi_available"}, result.Categories)
	assert.Equal(t, "Great spot for solo work", result.Comment)
	assert.NotEmpty(t, result.Id)

	// Verify mocks
	mockSpotRepo.AssertExpectations(t)
	mockRatingRepo.AssertExpectations(t)
}

func TestRatingService_CreateRating_UpdateExisting(t *testing.T) {
	// Arrange
	mockRatingRepo := &MockRatingRepository{}
	mockSpotRepo := &MockSpotRepository{}
	service := application.NewRatingService(mockRatingRepo, mockSpotRepo)

	ctx := context.Background()
	testSpot := &entities.Spot{
		ID:   "spot-123",
		Name: "Test Cafe",
	}

	existingRating, err := rating.NewRating("spot-123", "user-456", 3, []string{"quiet_atmosphere"}, "Initial rating")
	require.NoError(t, err)

	request := &pb.CreateRatingRequest{
		SpotId:             "spot-123",
		UserId:             "user-456",
		SoloFriendlyRating: 5,
		Categories:         []string{"quiet_atmosphere", "wifi_available"},
		Comment:            "Updated rating",
	}

	// Setup mocks
	mockSpotRepo.On("GetByID", ctx, "spot-123").Return(testSpot, nil)
	mockRatingRepo.On("GetBySpotAndUser", ctx, "spot-123", "user-456").Return(existingRating, nil)
	mockRatingRepo.On("Update", ctx, existingRating).Return(nil)
	mockRatingRepo.On("GetBySpot", ctx, "spot-123").Return([]*rating.Rating{existingRating}, nil)
	mockSpotRepo.On("UpdateSoloFriendlyStats", ctx, "spot-123", 5.0, 1).Return(nil)

	// Act
	result, err := service.CreateRating(ctx, request)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, existingRating.ID, result.Id)
	assert.Equal(t, int32(5), result.SoloFriendlyRating)
	assert.Equal(t, "Updated rating", result.Comment)

	// Verify mocks
	mockSpotRepo.AssertExpectations(t)
	mockRatingRepo.AssertExpectations(t)
}

func TestRatingService_CreateRating_SpotNotFound(t *testing.T) {
	// Arrange
	mockRatingRepo := &MockRatingRepository{}
	mockSpotRepo := &MockSpotRepository{}
	service := application.NewRatingService(mockRatingRepo, mockSpotRepo)

	ctx := context.Background()
	request := &pb.CreateRatingRequest{
		SpotId:             "nonexistent-spot",
		UserId:             "user-456",
		SoloFriendlyRating: 5,
		Categories:         []string{"quiet_atmosphere"},
		Comment:            "Great spot",
	}

	// Setup mocks
	mockSpotRepo.On("GetByID", ctx, "nonexistent-spot").Return(nil, errors.New("spot not found"))

	// Act
	result, err := service.CreateRating(ctx, request)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "spot not found")

	// Verify mocks
	mockSpotRepo.AssertExpectations(t)
}

func TestRatingService_CreateRating_InvalidInput(t *testing.T) {
	// Arrange
	mockRatingRepo := &MockRatingRepository{}
	mockSpotRepo := &MockSpotRepository{}
	service := application.NewRatingService(mockRatingRepo, mockSpotRepo)

	ctx := context.Background()
	
	tests := []struct {
		name           string
		request        *pb.CreateRatingRequest
		wantErrMessage string
	}{
		{
			name: "empty spot ID",
			request: &pb.CreateRatingRequest{
				SpotId:             "",
				UserId:             "user-456",
				SoloFriendlyRating: 5,
				Categories:         []string{"quiet_atmosphere"},
				Comment:            "Great spot",
			},
			wantErrMessage: "spot ID cannot be empty",
		},
		{
			name: "empty user ID",
			request: &pb.CreateRatingRequest{
				SpotId:             "spot-123",
				UserId:             "",
				SoloFriendlyRating: 5,
				Categories:         []string{"quiet_atmosphere"},
				Comment:            "Great spot",
			},
			wantErrMessage: "user ID cannot be empty",
		},
		{
			name: "invalid rating",
			request: &pb.CreateRatingRequest{
				SpotId:             "spot-123",
				UserId:             "user-456",
				SoloFriendlyRating: 10,
				Categories:         []string{"quiet_atmosphere"},
				Comment:            "Great spot",
			},
			wantErrMessage: "solo-friendly rating must be between 1 and 5",
		},
		{
			name: "invalid category",
			request: &pb.CreateRatingRequest{
				SpotId:             "spot-123",
				UserId:             "user-456",
				SoloFriendlyRating: 5,
				Categories:         []string{"invalid_category"},
				Comment:            "Great spot",
			},
			wantErrMessage: "invalid category: invalid_category",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result, err := service.CreateRating(ctx, tt.request)

			// Assert
			assert.Error(t, err)
			assert.Nil(t, result)
			assert.Contains(t, err.Error(), tt.wantErrMessage)
		})
	}
}

func TestRatingService_GetSpotRatings(t *testing.T) {
	// Arrange
	mockRatingRepo := &MockRatingRepository{}
	mockSpotRepo := &MockSpotRepository{}
	service := application.NewRatingService(mockRatingRepo, mockSpotRepo)

	ctx := context.Background()
	spotID := "spot-123"

	rating1, err := rating.NewRating(spotID, "user-1", 5, []string{"quiet_atmosphere"}, "Great spot")
	require.NoError(t, err)
	rating2, err := rating.NewRating(spotID, "user-2", 4, []string{"wifi_available"}, "Good wifi")
	require.NoError(t, err)

	ratings := []*rating.Rating{rating1, rating2}

	// Setup mocks
	mockRatingRepo.On("GetBySpot", ctx, spotID).Return(ratings, nil)

	// Act
	result, err := service.GetSpotRatings(ctx, &pb.GetSpotRatingsRequest{SpotId: spotID})

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Ratings, 2)
	
	// Verify first rating
	assert.Equal(t, rating1.ID, result.Ratings[0].Id)
	assert.Equal(t, rating1.UserID, result.Ratings[0].UserId)
	assert.Equal(t, int32(rating1.SoloFriendlyRating), result.Ratings[0].SoloFriendlyRating)
	
	// Verify second rating
	assert.Equal(t, rating2.ID, result.Ratings[1].Id)
	assert.Equal(t, rating2.UserID, result.Ratings[1].UserId)
	assert.Equal(t, int32(rating2.SoloFriendlyRating), result.Ratings[1].SoloFriendlyRating)

	// Verify mocks
	mockRatingRepo.AssertExpectations(t)
}

func TestRatingService_CalculateSpotStatistics(t *testing.T) {
	// Arrange
	mockRatingRepo := &MockRatingRepository{}
	mockSpotRepo := &MockSpotRepository{}
	service := application.NewRatingService(mockRatingRepo, mockSpotRepo)

	rating1, err := rating.NewRating("spot-123", "user-1", 5, []string{"quiet_atmosphere"}, "Great")
	require.NoError(t, err)
	rating2, err := rating.NewRating("spot-123", "user-2", 3, []string{"wifi_available"}, "OK")
	require.NoError(t, err)
	rating3, err := rating.NewRating("spot-123", "user-3", 4, []string{"good_lighting"}, "Good")
	require.NoError(t, err)

	ratings := []*rating.Rating{rating1, rating2, rating3}

	// Act
	avgRating, totalRatings := service.CalculateSpotStatistics(ratings)

	// Assert
	expectedAvg := float64(5+3+4) / 3 // = 4.0
	assert.Equal(t, expectedAvg, avgRating)
	assert.Equal(t, 3, totalRatings)
}

func TestRatingService_CalculateSpotStatistics_EmptyRatings(t *testing.T) {
	// Arrange
	mockRatingRepo := &MockRatingRepository{}
	mockSpotRepo := &MockSpotRepository{}
	service := application.NewRatingService(mockRatingRepo, mockSpotRepo)

	ratings := []*rating.Rating{}

	// Act
	avgRating, totalRatings := service.CalculateSpotStatistics(ratings)

	// Assert
	assert.Equal(t, 0.0, avgRating)
	assert.Equal(t, 0, totalRatings)
}