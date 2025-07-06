package application_test

import (
	"context"
	"testing"
	"time"

	"bocchi/api/internal/application"
	"bocchi/api/domain/entities"
	"bocchi/api/internal/domain/rating"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

// Test statistics calculation functionality
func TestRatingService_CalculateSpotStatistics(t *testing.T) {
	// Arrange
	service := application.NewRatingService(nil, nil)
	
	ratings := []*rating.Rating{
		{SoloFriendlyRating: 5},
		{SoloFriendlyRating: 4},
		{SoloFriendlyRating: 3},
	}

	// Act
	avgRating, totalCount := service.CalculateSpotStatistics(ratings)

	// Assert
	assert.Equal(t, 4.0, avgRating)
	assert.Equal(t, 3, totalCount)
}

func TestRatingService_CalculateSpotStatistics_EmptyRatings(t *testing.T) {
	// Arrange
	service := application.NewRatingService(nil, nil)
	ratings := []*rating.Rating{}

	// Act
	avgRating, totalCount := service.CalculateSpotStatistics(ratings)

	// Assert
	assert.Equal(t, 0.0, avgRating)
	assert.Equal(t, 0, totalCount)
}