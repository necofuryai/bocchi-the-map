package grpc

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"bocchi/api/infrastructure/database"
	"bocchi/api/pkg/errors"
	"bocchi/api/pkg/logger"
	"bocchi/api/pkg/monitoring"
	commonv1 "bocchi/api/gen/common/v1"
	reviewv1 "bocchi/api/gen/review/v1"
)

// ReviewService implements the gRPC ReviewService
type ReviewService struct {
	queries *database.Queries
}

// NewReviewService creates a new ReviewService instance
func NewReviewService(db *sql.DB) *ReviewService {
	return &ReviewService{
		queries: database.New(db),
	}
}


// CreateReview creates a new review
func (s *ReviewService) CreateReview(ctx context.Context, req *reviewv1.CreateReviewRequest) (*reviewv1.CreateReviewResponse, error) {
	// Validate request
	if req.GetSpotId() == "" {
		return nil, status.Error(codes.InvalidArgument, "spot_id is required")
	}
	if req.GetRating() < 1 || req.GetRating() > 5 {
		return nil, status.Error(codes.InvalidArgument, "rating must be between 1 and 5")
	}

	// Extract user ID from authentication context
	userID := errors.GetUserID(ctx)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	// Check if user already reviewed this spot
	_, err := s.queries.GetReviewByUserAndSpot(ctx, database.GetReviewByUserAndSpotParams{
		UserID: sql.NullString{String: userID, Valid: true},
		SpotID: req.GetSpotId(),
	})
	if err == nil {
		// User already reviewed this spot, return error
		return nil, status.Error(codes.AlreadyExists, "user has already reviewed this spot")
	} else if err != sql.ErrNoRows {
		return nil, status.Error(codes.Internal, "failed to check existing review")
	}

	// Generate UUID for new review
	reviewID := uuid.New().String()

	// Convert rating aspects to JSON
	ratingAspectsJSON, err := json.Marshal(req.GetRatingAspects())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to marshal rating aspects")
	}
	if req.GetRatingAspects() == nil {
		ratingAspectsJSON = []byte("{}")
	}

	// Convert comment to nullable string
	var comment sql.NullString
	if req.GetComment() != "" {
		comment = sql.NullString{String: req.GetComment(), Valid: true}
	}

	// Create review in database
	err = s.queries.CreateReview(ctx, database.CreateReviewParams{
		ID:            reviewID,
		SpotID:        req.GetSpotId(),
		UserID:        sql.NullString{String: userID, Valid: true},
		Rating:        req.GetRating(),
		Comment:       comment,
		RatingAspects: ratingAspectsJSON,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create review")
	}

	// Update spot rating statistics asynchronously with retry logic
	go func(parentCtx context.Context) {
		const maxRetries = 3
		const baseDelay = time.Second
		
		for attempt := 0; attempt < maxRetries; attempt++ {
			// Create timeout context for each attempt
			ctx, cancel := context.WithTimeout(parentCtx, 30*time.Second)
			
			err := s.updateSpotRating(ctx, req.GetSpotId())
			cancel()
			
			if err == nil {
				// Success - exit retry loop
				break
			}
			
			// Check if parent context was cancelled
			if parentCtx.Err() != nil {
				log.Printf("spot rating update cancelled for spot %s: %v", req.GetSpotId(), parentCtx.Err())
				return
			}
			
			// Log the error
			log.Printf("failed to update spot rating for spot %s (attempt %d/%d): %v", req.GetSpotId(), attempt+1, maxRetries, err)
			
			// If this was the last attempt, we're done
			if attempt == maxRetries-1 {
				log.Printf("exhausted all retry attempts for spot rating update for spot %s", req.GetSpotId())
				return
			}
			
			// Exponential backoff with jitter
			delay := baseDelay * time.Duration(1<<attempt)
			jitter := time.Duration(rand.Int63n(int64(delay / 2)))
			select {
			case <-time.After(delay + jitter):
				// Continue to next retry
			case <-parentCtx.Done():
				log.Printf("spot rating update cancelled during backoff for spot %s: %v", req.GetSpotId(), parentCtx.Err())
				return
			}
		}
	}(ctx)

	// Retrieve the created review to get accurate timestamps
	dbReview, err := s.queries.GetReviewByID(ctx, reviewID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to retrieve created review")
	}

	// Convert database review to gRPC response
	review := s.convertDatabaseReviewToGRPC(dbReview)
	return &reviewv1.CreateReviewResponse{Review: review}, nil
}

// GetSpotReviews retrieves reviews for a specific spot
func (s *ReviewService) GetSpotReviews(ctx context.Context, req *reviewv1.GetSpotReviewsRequest) (*reviewv1.GetSpotReviewsResponse, error) {
	if req.GetSpotId() == "" {
		return nil, status.Error(codes.InvalidArgument, "spot_id is required")
	}

	// Set pagination defaults
	pageSize := int32(20)
	page := int32(1)
	if req.Pagination != nil {
		if req.Pagination.PageSize > 0 {
			pageSize = req.Pagination.PageSize
		}
		if req.Pagination.Page > 0 {
			page = req.Pagination.Page
		}
	}
	offset := (page - 1) * pageSize

	// Get total count of reviews for this spot
	totalCount, err := s.queries.CountReviewsBySpot(ctx, req.GetSpotId())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to count reviews")
	}

	// Get reviews from database
	dbReviews, err := s.queries.ListReviewsBySpot(ctx, database.ListReviewsBySpotParams{
		SpotID: req.GetSpotId(),
		Limit:  pageSize,
		Offset: offset,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get spot reviews")
	}

	// Convert database reviews to gRPC format
	reviews := make([]*reviewv1.Review, len(dbReviews))
	for i, dbReview := range dbReviews {
		reviews[i] = s.convertDatabaseReviewRowToGRPC(dbReview)
	}

	// Calculate total pages
	totalPages := (int32(totalCount) + pageSize - 1) / pageSize

	// Create pagination response
	pagination := &commonv1.PaginationResponse{
		TotalCount: int32(totalCount),
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	// Get rating statistics
	statistics, err := s.getSpotStatistics(ctx, req.GetSpotId())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get spot statistics")
	}

	return &reviewv1.GetSpotReviewsResponse{
		Reviews:    reviews,
		Pagination: pagination,
		Statistics: statistics,
	}, nil
}

// GetUserReviews retrieves reviews by a specific user
func (s *ReviewService) GetUserReviews(ctx context.Context, req *reviewv1.GetUserReviewsRequest) (*reviewv1.GetUserReviewsResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	// Set pagination defaults
	pageSize := int32(20)
	page := int32(1)
	if req.Pagination != nil {
		if req.Pagination.PageSize > 0 {
			pageSize = req.Pagination.PageSize
		}
		if req.Pagination.Page > 0 {
			page = req.Pagination.Page
		}
	}
	offset := (page - 1) * pageSize

	// Get total count of reviews by this user
	totalCount, err := s.queries.CountReviewsByUser(ctx, sql.NullString{String: req.GetUserId(), Valid: true})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to count user reviews")
	}

	// Get reviews from database
	dbReviews, err := s.queries.ListReviewsByUser(ctx, database.ListReviewsByUserParams{
		UserID: sql.NullString{String: req.GetUserId(), Valid: true},
		Limit:  pageSize,
		Offset: offset,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get user reviews")
	}

	// Convert database reviews to gRPC format
	reviews := make([]*reviewv1.Review, len(dbReviews))
	for i, dbReview := range dbReviews {
		reviews[i] = s.convertDatabaseUserReviewRowToGRPC(dbReview)
	}

	// Calculate total pages
	totalPages := (int32(totalCount) + pageSize - 1) / pageSize

	// Create pagination response
	pagination := &commonv1.PaginationResponse{
		TotalCount: int32(totalCount),
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	return &reviewv1.GetUserReviewsResponse{
		Reviews:    reviews,
		Pagination: pagination,
	}, nil
}

// convertToInt32 safely converts interface{} to int32
func convertToInt32(val interface{}) int32 {
	switch v := val.(type) {
	case int64:
		return int32(v)
	case int32:
		return v
	case int:
		return int32(v)
	case float64:
		return int32(v)
	default:
		return 0
	}
}

// parseRatingAspects converts JSON rating aspects to map[string]int32
func parseRatingAspects(rawData json.RawMessage) map[string]int32 {
	var ratingAspects map[string]int32
	if len(rawData) > 0 {
		var aspectsMap map[string]interface{}
		if err := json.Unmarshal(rawData, &aspectsMap); err == nil {
			ratingAspects = make(map[string]int32)
			for key, value := range aspectsMap {
				if intVal, ok := value.(float64); ok {
					ratingAspects[key] = int32(intVal)
				}
			}
		}
	}
	return ratingAspects
}

// convertDatabaseReviewToGRPC converts database review model to gRPC review struct
func (s *ReviewService) convertDatabaseReviewToGRPC(dbReview database.Review) *reviewv1.Review {
	return &reviewv1.Review{
		Id:            dbReview.ID,
		SpotId:        dbReview.SpotID,
		UserId:        dbReview.UserID.String,
		Rating:        dbReview.Rating,
		Comment:       dbReview.Comment.String,
		RatingAspects: parseRatingAspects(dbReview.RatingAspects),
		CreatedAt:     timestamppb.New(dbReview.CreatedAt),
		UpdatedAt:     timestamppb.New(dbReview.UpdatedAt),
	}
}

// convertDatabaseReviewRowToGRPC converts database review row (with user info) to gRPC review struct
func (s *ReviewService) convertDatabaseReviewRowToGRPC(dbReview database.ListReviewsBySpotRow) *reviewv1.Review {
	return &reviewv1.Review{
		Id:            dbReview.ID,
		SpotId:        dbReview.SpotID,
		UserId:        dbReview.UserID.String,
		Rating:        dbReview.Rating,
		Comment:       dbReview.Comment.String,
		RatingAspects: parseRatingAspects(dbReview.RatingAspects),
		CreatedAt:     timestamppb.New(dbReview.CreatedAt),
		UpdatedAt:     timestamppb.New(dbReview.UpdatedAt),
	}
}

// convertDatabaseUserReviewRowToGRPC converts database user review row (with spot info) to gRPC review struct
func (s *ReviewService) convertDatabaseUserReviewRowToGRPC(dbReview database.ListReviewsByUserRow) *reviewv1.Review {
	return &reviewv1.Review{
		Id:            dbReview.ID,
		SpotId:        dbReview.SpotID,
		UserId:        dbReview.UserID.String,
		Rating:        dbReview.Rating,
		Comment:       dbReview.Comment.String,
		RatingAspects: parseRatingAspects(dbReview.RatingAspects),
		CreatedAt:     timestamppb.New(dbReview.CreatedAt),
		UpdatedAt:     timestamppb.New(dbReview.UpdatedAt),
	}
}

// getSpotStatistics retrieves rating statistics for a spot
func (s *ReviewService) getSpotStatistics(ctx context.Context, spotID string) (*reviewv1.ReviewStatistics, error) {
	stats, err := s.queries.GetSpotRatingStats(ctx, spotID)
	if err != nil {
		return nil, err
	}

	averageRating := 0.0
	if avgRating, ok := stats.AverageRating.(float64); ok {
		averageRating = avgRating
	}

	return &reviewv1.ReviewStatistics{
		AverageRating: averageRating,
		TotalCount:    int32(stats.ReviewCount),
		RatingDistribution: map[int32]int32{
			1: convertToInt32(stats.OneStarCount),
			2: convertToInt32(stats.TwoStarCount),
			3: convertToInt32(stats.ThreeStarCount),
			4: convertToInt32(stats.FourStarCount),
			5: convertToInt32(stats.FiveStarCount),
		},
	}, nil
}

// updateSpotRating updates the spot's average rating and review count
func (s *ReviewService) updateSpotRating(ctx context.Context, spotID string) error {
	stats, err := s.queries.GetSpotRatingStats(ctx, spotID)
	if err != nil {
		logger.ErrorWithContextAndFields(ctx, "Failed to get spot rating stats", err, map[string]interface{}{
			"spot_id": spotID,
		})
		monitoring.CaptureError(ctx, err)
		return err
	}

	averageRating := "0.0"
	if avgRating, ok := stats.AverageRating.(float64); ok {
		averageRating = strconv.FormatFloat(avgRating, 'f', 1, 64)
	}

	// Update spot table with new rating statistics
	err = s.queries.UpdateSpotRating(ctx, database.UpdateSpotRatingParams{
		ID:            spotID,
		AverageRating: averageRating,
		ReviewCount:   int32(stats.ReviewCount),
	})
	if err != nil {
		logger.ErrorWithContextAndFields(ctx, "Failed to update spot rating", err, map[string]interface{}{
			"spot_id": spotID,
		})
		monitoring.CaptureError(ctx, err)
		return err
	}
	return nil
}