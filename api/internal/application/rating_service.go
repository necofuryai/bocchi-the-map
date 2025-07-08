package application

import (
	"context"
	"fmt"

	"bocchi/api/domain/entities"
	"bocchi/api/internal/domain/rating"
)

// RatingRepository defines the interface for rating data access
type RatingRepository interface {
	Create(ctx context.Context, rating *rating.Rating) error
	Update(ctx context.Context, rating *rating.Rating) error
	GetBySpotAndUser(ctx context.Context, spotID, userID string) (*rating.Rating, error)
	GetBySpot(ctx context.Context, spotID string) ([]*rating.Rating, error)
}

// SpotRepository defines the interface for spot data access
type SpotRepository interface {
	GetByID(ctx context.Context, id string) (*entities.Spot, error)
	UpdateSoloFriendlyStats(ctx context.Context, spotID string, avgRating float64, totalRatings int) error
}

// RatingService handles business logic for solo-friendly ratings
type RatingService struct {
	ratingRepo RatingRepository
	spotRepo   SpotRepository
}

// NewRatingService creates a new RatingService instance
func NewRatingService(ratingRepo RatingRepository, spotRepo SpotRepository) *RatingService {
	return &RatingService{
		ratingRepo: ratingRepo,
		spotRepo:   spotRepo,
	}
}

// CalculateSpotStatistics calculates average rating and total count
func (s *RatingService) CalculateSpotStatistics(ratings []*rating.Rating) (float64, int) {
	if len(ratings) == 0 {
		return 0.0, 0
	}

	totalRating := 0
	for _, rating := range ratings {
		totalRating += rating.SoloFriendlyRating
	}

	avgRating := float64(totalRating) / float64(len(ratings))
	return avgRating, len(ratings)
}

// updateSpotStatistics recalculates and updates spot statistics
func (s *RatingService) updateSpotStatistics(ctx context.Context, spotID string) error {
	ratings, err := s.ratingRepo.GetBySpot(ctx, spotID)
	if err != nil {
		return fmt.Errorf("failed to get spot ratings: %w", err)
	}

	avgRating, totalRatings := s.CalculateSpotStatistics(ratings)

	err = s.spotRepo.UpdateSoloFriendlyStats(ctx, spotID, avgRating, totalRatings)
	if err != nil {
		return fmt.Errorf("failed to update spot statistics: %w", err)
	}

	return nil
}

// isNotFoundError checks if an error represents a "not found" condition
// This is a helper function that should be implemented based on your error handling strategy
func isNotFoundError(err error) bool {
	// This is a simplified implementation
	// In a real application, you might check for specific error types
	return err != nil && (err.Error() == "not found" || err.Error() == "record not found")
}