package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/necofuryai/bocchi-the-map/api/internal/domain/entities"
	"github.com/necofuryai/bocchi-the-map/api/internal/domain/rating"
	"github.com/necofuryai/bocchi-the-map/api/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
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

// CreateRating creates or updates a rating for a spot
func (s *RatingService) CreateRating(ctx context.Context, req *pb.CreateRatingRequest) (*pb.Rating, error) {
	// Validate spot exists
	spot, err := s.spotRepo.GetByID(ctx, req.SpotId)
	if err != nil {
		return nil, fmt.Errorf("spot not found: %w", err)
	}
	if spot == nil {
		return nil, errors.New("spot not found")
	}

	// Check if user has already rated this spot
	existingRating, err := s.ratingRepo.GetBySpotAndUser(ctx, req.SpotId, req.UserId)
	if err != nil && !isNotFoundError(err) {
		return nil, fmt.Errorf("failed to check existing rating: %w", err)
	}

	var resultRating *rating.Rating

	if existingRating != nil {
		// Update existing rating
		err = existingRating.UpdateRating(
			int(req.SoloFriendlyRating),
			req.Categories,
			req.Comment,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to update rating: %w", err)
		}

		err = s.ratingRepo.Update(ctx, existingRating)
		if err != nil {
			return nil, fmt.Errorf("failed to save updated rating: %w", err)
		}
		resultRating = existingRating
	} else {
		// Create new rating
		newRating, err := rating.NewRating(
			req.SpotId,
			req.UserId,
			int(req.SoloFriendlyRating),
			req.Categories,
			req.Comment,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create rating: %w", err)
		}

		err = s.ratingRepo.Create(ctx, newRating)
		if err != nil {
			return nil, fmt.Errorf("failed to save new rating: %w", err)
		}
		resultRating = newRating
	}

	// Update spot statistics
	err = s.updateSpotStatistics(ctx, req.SpotId)
	if err != nil {
		// Log error but don't fail the operation
		// In a real application, you might want to use a proper logger
		fmt.Printf("Warning: failed to update spot statistics: %v\n", err)
	}

	return s.convertToProtobuf(resultRating), nil
}

// GetSpotRatings retrieves all ratings for a spot
func (s *RatingService) GetSpotRatings(ctx context.Context, req *pb.GetSpotRatingsRequest) (*pb.GetSpotRatingsResponse, error) {
	ratings, err := s.ratingRepo.GetBySpot(ctx, req.SpotId)
	if err != nil {
		return nil, fmt.Errorf("failed to get spot ratings: %w", err)
	}

	pbRatings := make([]*pb.Rating, len(ratings))
	for i, rating := range ratings {
		pbRatings[i] = s.convertToProtobuf(rating)
	}

	return &pb.GetSpotRatingsResponse{
		Ratings: pbRatings,
	}, nil
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

// convertToProtobuf converts domain rating to protobuf rating
func (s *RatingService) convertToProtobuf(r *rating.Rating) *pb.Rating {
	return &pb.Rating{
		Id:                 r.ID,
		SpotId:             r.SpotID,
		UserId:             r.UserID,
		SoloFriendlyRating: int32(r.SoloFriendlyRating),
		Categories:         r.Categories,
		Comment:            r.Comment,
		CreatedAt:          timestamppb.New(r.CreatedAt),
		UpdatedAt:          timestamppb.New(r.UpdatedAt),
	}
}

// isNotFoundError checks if an error represents a "not found" condition
// This is a helper function that should be implemented based on your error handling strategy
func isNotFoundError(err error) bool {
	// This is a simplified implementation
	// In a real application, you might check for specific error types
	return err != nil && (err.Error() == "not found" || err.Error() == "record not found")
}