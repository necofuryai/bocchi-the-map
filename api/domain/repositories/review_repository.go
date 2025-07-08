package repositories

import (
	"context"

	"bocchi/api/domain/entities"
)

// ReviewRepository defines the interface for review data operations
type ReviewRepository interface {
	// Create creates a new review
	Create(ctx context.Context, review *entities.Review) error

	// GetByID retrieves a review by its ID
	GetByID(ctx context.Context, id string) (*entities.Review, error)

	// GetBySpotID retrieves reviews for a specific spot
	GetBySpotID(ctx context.Context, spotID string, offset, limit int) ([]*entities.Review, int, error)

	// GetByUserID retrieves reviews by a specific user
	GetByUserID(ctx context.Context, userID string, offset, limit int) ([]*entities.Review, int, error)

	// GetUserReviewForSpot checks if a user has already reviewed a spot
	GetUserReviewForSpot(ctx context.Context, userID, spotID string) (*entities.Review, error)

	// Update updates a review
	Update(ctx context.Context, review *entities.Review) error

	// Delete deletes a review
	Delete(ctx context.Context, id string) error

	// GetStatisticsBySpotID gets review statistics for a spot
	GetStatisticsBySpotID(ctx context.Context, spotID string) (*ReviewStatistics, error)
}

// ReviewStatistics holds aggregated review data
type ReviewStatistics struct {
	AverageRating      float64
	TotalCount         int
	RatingDistribution map[int]int // key: rating (1-5), value: count
}