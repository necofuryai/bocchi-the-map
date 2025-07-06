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

	"github.com/necofuryai/bocchi-the-map/api/infrastructure/database"
	"github.com/necofuryai/bocchi-the-map/api/pkg/errors"
	"github.com/necofuryai/bocchi-the-map/api/pkg/logger"
	"github.com/necofuryai/bocchi-the-map/api/pkg/monitoring"
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

// Temporary structs until protobuf generates them
type Review struct {
	ID            string
	SpotID        string
	UserID        string
	Rating        int32
	Comment       string
	RatingAspects map[string]int32
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type CreateReviewRequest struct {
	SpotID        string
	Rating        int32
	Comment       string
	RatingAspects map[string]int32
}

type CreateReviewResponse struct {
	Review *Review
}

type GetSpotReviewsRequest struct {
	SpotID     string
	Pagination *PaginationRequest
}

type ReviewStatistics struct {
	AverageRating      float64
	TotalCount         int32
	RatingDistribution map[int32]int32
}

type GetSpotReviewsResponse struct {
	Reviews    []*Review
	Pagination *PaginationResponse
	Statistics *ReviewStatistics
}

type GetUserReviewsRequest struct {
	UserID     string
	Pagination *PaginationRequest
}

type GetUserReviewsResponse struct {
	Reviews    []*Review
	Pagination *PaginationResponse
}

// CreateReview creates a new review
func (s *ReviewService) CreateReview(ctx context.Context, req *CreateReviewRequest) (*CreateReviewResponse, error) {
	// Validate request
	if req.SpotID == "" {
		return nil, status.Error(codes.InvalidArgument, "spot_id is required")
	}
	if req.Rating < 1 || req.Rating > 5 {
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
		SpotID: req.SpotID,
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
	ratingAspectsJSON, err := json.Marshal(req.RatingAspects)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to marshal rating aspects")
	}
	if req.RatingAspects == nil {
		ratingAspectsJSON = []byte("{}")
	}

	// Convert comment to nullable string
	var comment sql.NullString
	if req.Comment != "" {
		comment = sql.NullString{String: req.Comment, Valid: true}
	}

	// Create review in database
	err = s.queries.CreateReview(ctx, database.CreateReviewParams{
		ID:            reviewID,
		SpotID:        req.SpotID,
		UserID:        sql.NullString{String: userID, Valid: true},
		Rating:        req.Rating,
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
			
			err := s.updateSpotRating(ctx, req.SpotID)
			cancel()
			
			if err == nil {
				// Success - exit retry loop
				break
			}
			
			// Check if parent context was cancelled
			if parentCtx.Err() != nil {
				log.Printf("spot rating update cancelled for spot %s: %v", req.SpotID, parentCtx.Err())
				return
			}
			
			// Log the error
			log.Printf("failed to update spot rating for spot %s (attempt %d/%d): %v", req.SpotID, attempt+1, maxRetries, err)
			
			// If this was the last attempt, we're done
			if attempt == maxRetries-1 {
				log.Printf("exhausted all retry attempts for spot rating update for spot %s", req.SpotID)
				return
			}
			
			// Exponential backoff with jitter
			delay := baseDelay * time.Duration(1<<attempt)
			jitter := time.Duration(rand.Int63n(int64(delay / 2)))
			select {
			case <-time.After(delay + jitter):
				// Continue to next retry
			case <-parentCtx.Done():
				log.Printf("spot rating update cancelled during backoff for spot %s: %v", req.SpotID, parentCtx.Err())
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
	return &CreateReviewResponse{Review: review}, nil
}

// GetSpotReviews retrieves reviews for a specific spot
func (s *ReviewService) GetSpotReviews(ctx context.Context, req *GetSpotReviewsRequest) (*GetSpotReviewsResponse, error) {
	if req.SpotID == "" {
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
	totalCount, err := s.queries.CountReviewsBySpot(ctx, req.SpotID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to count reviews")
	}

	// Get reviews from database
	dbReviews, err := s.queries.ListReviewsBySpot(ctx, database.ListReviewsBySpotParams{
		SpotID: req.SpotID,
		Limit:  pageSize,
		Offset: offset,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get spot reviews")
	}

	// Convert database reviews to gRPC format
	reviews := make([]*Review, len(dbReviews))
	for i, dbReview := range dbReviews {
		reviews[i] = s.convertDatabaseReviewRowToGRPC(dbReview)
	}

	// Calculate total pages
	totalPages := (int32(totalCount) + pageSize - 1) / pageSize

	// Create pagination response
	pagination := &PaginationResponse{
		TotalCount: int32(totalCount),
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	// Get rating statistics
	statistics, err := s.getSpotStatistics(ctx, req.SpotID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get spot statistics")
	}

	return &GetSpotReviewsResponse{
		Reviews:    reviews,
		Pagination: pagination,
		Statistics: statistics,
	}, nil
}

// GetUserReviews retrieves reviews by a specific user
func (s *ReviewService) GetUserReviews(ctx context.Context, req *GetUserReviewsRequest) (*GetUserReviewsResponse, error) {
	if req.UserID == "" {
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
	totalCount, err := s.queries.CountReviewsByUser(ctx, sql.NullString{String: req.UserID, Valid: true})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to count user reviews")
	}

	// Get reviews from database
	dbReviews, err := s.queries.ListReviewsByUser(ctx, database.ListReviewsByUserParams{
		UserID: sql.NullString{String: req.UserID, Valid: true},
		Limit:  pageSize,
		Offset: offset,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get user reviews")
	}

	// Convert database reviews to gRPC format
	reviews := make([]*Review, len(dbReviews))
	for i, dbReview := range dbReviews {
		reviews[i] = s.convertDatabaseUserReviewRowToGRPC(dbReview)
	}

	// Calculate total pages
	totalPages := (int32(totalCount) + pageSize - 1) / pageSize

	// Create pagination response
	pagination := &PaginationResponse{
		TotalCount: int32(totalCount),
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	return &GetUserReviewsResponse{
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
func (s *ReviewService) convertDatabaseReviewToGRPC(dbReview database.Review) *Review {
	return &Review{
		ID:            dbReview.ID,
		SpotID:        dbReview.SpotID,
		UserID:        dbReview.UserID.String,
		Rating:        dbReview.Rating,
		Comment:       dbReview.Comment.String,
		RatingAspects: parseRatingAspects(dbReview.RatingAspects),
		CreatedAt:     dbReview.CreatedAt,
		UpdatedAt:     dbReview.UpdatedAt,
	}
}

// convertDatabaseReviewRowToGRPC converts database review row (with user info) to gRPC review struct
func (s *ReviewService) convertDatabaseReviewRowToGRPC(dbReview database.ListReviewsBySpotRow) *Review {
	return &Review{
		ID:            dbReview.ID,
		SpotID:        dbReview.SpotID,
		UserID:        dbReview.UserID.String,
		Rating:        dbReview.Rating,
		Comment:       dbReview.Comment.String,
		RatingAspects: parseRatingAspects(dbReview.RatingAspects),
		CreatedAt:     dbReview.CreatedAt,
		UpdatedAt:     dbReview.UpdatedAt,
	}
}

// convertDatabaseUserReviewRowToGRPC converts database user review row (with spot info) to gRPC review struct
func (s *ReviewService) convertDatabaseUserReviewRowToGRPC(dbReview database.ListReviewsByUserRow) *Review {
	return &Review{
		ID:            dbReview.ID,
		SpotID:        dbReview.SpotID,
		UserID:        dbReview.UserID.String,
		Rating:        dbReview.Rating,
		Comment:       dbReview.Comment.String,
		RatingAspects: parseRatingAspects(dbReview.RatingAspects),
		CreatedAt:     dbReview.CreatedAt,
		UpdatedAt:     dbReview.UpdatedAt,
	}
}

// getSpotStatistics retrieves rating statistics for a spot
func (s *ReviewService) getSpotStatistics(ctx context.Context, spotID string) (*ReviewStatistics, error) {
	stats, err := s.queries.GetSpotRatingStats(ctx, spotID)
	if err != nil {
		return nil, err
	}

	averageRating := 0.0
	if avgRating, ok := stats.AverageRating.(float64); ok {
		averageRating = avgRating
	}

	return &ReviewStatistics{
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