package grpc

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/necofuryai/bocchi-the-map/api/infrastructure/database"
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

	// TODO: CRITICAL - Extract user ID from context/authentication
	// This is currently hardcoded for development purposes
	userID := "user_123"
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	// Check if user already reviewed this spot
	existingReview, err := s.queries.GetReviewByUserAndSpot(ctx, database.GetReviewByUserAndSpotParams{
		UserID: userID,
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
		UserID:        userID,
		Rating:        req.Rating,
		Comment:       comment,
		RatingAspects: ratingAspectsJSON,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create review")
	}

	// Update spot rating statistics
	go s.updateSpotRating(context.Background(), req.SpotID)

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
	totalCount, err := s.queries.CountReviewsByUser(ctx, req.UserID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to count user reviews")
	}

	// Get reviews from database
	dbReviews, err := s.queries.ListReviewsByUser(ctx, database.ListReviewsByUserParams{
		UserID: req.UserID,
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

// convertDatabaseReviewToGRPC converts database review model to gRPC review struct
func (s *ReviewService) convertDatabaseReviewToGRPC(dbReview database.Review) *Review {
	// Parse rating aspects from JSON
	var ratingAspects map[string]int32
	if len(dbReview.RatingAspects) > 0 {
		var aspectsMap map[string]interface{}
		if err := json.Unmarshal(dbReview.RatingAspects, &aspectsMap); err == nil {
			ratingAspects = make(map[string]int32)
			for key, value := range aspectsMap {
				if intVal, ok := value.(float64); ok {
					ratingAspects[key] = int32(intVal)
				}
			}
		}
	}

	return &Review{
		ID:            dbReview.ID,
		SpotID:        dbReview.SpotID,
		UserID:        dbReview.UserID,
		Rating:        dbReview.Rating,
		Comment:       dbReview.Comment.String,
		RatingAspects: ratingAspects,
		CreatedAt:     dbReview.CreatedAt,
		UpdatedAt:     dbReview.UpdatedAt,
	}
}

// convertDatabaseReviewRowToGRPC converts database review row (with user info) to gRPC review struct
func (s *ReviewService) convertDatabaseReviewRowToGRPC(dbReview database.ListReviewsBySpotRow) *Review {
	// Parse rating aspects from JSON
	var ratingAspects map[string]int32
	if len(dbReview.RatingAspects) > 0 {
		var aspectsMap map[string]interface{}
		if err := json.Unmarshal(dbReview.RatingAspects, &aspectsMap); err == nil {
			ratingAspects = make(map[string]int32)
			for key, value := range aspectsMap {
				if intVal, ok := value.(float64); ok {
					ratingAspects[key] = int32(intVal)
				}
			}
		}
	}

	return &Review{
		ID:            dbReview.ID,
		SpotID:        dbReview.SpotID,
		UserID:        dbReview.UserID,
		Rating:        dbReview.Rating,
		Comment:       dbReview.Comment.String,
		RatingAspects: ratingAspects,
		CreatedAt:     dbReview.CreatedAt,
		UpdatedAt:     dbReview.UpdatedAt,
	}
}

// convertDatabaseUserReviewRowToGRPC converts database user review row (with spot info) to gRPC review struct
func (s *ReviewService) convertDatabaseUserReviewRowToGRPC(dbReview database.ListReviewsByUserRow) *Review {
	// Parse rating aspects from JSON
	var ratingAspects map[string]int32
	if len(dbReview.RatingAspects) > 0 {
		var aspectsMap map[string]interface{}
		if err := json.Unmarshal(dbReview.RatingAspects, &aspectsMap); err == nil {
			ratingAspects = make(map[string]int32)
			for key, value := range aspectsMap {
				if intVal, ok := value.(float64); ok {
					ratingAspects[key] = int32(intVal)
				}
			}
		}
	}

	return &Review{
		ID:            dbReview.ID,
		SpotID:        dbReview.SpotID,
		UserID:        dbReview.UserID,
		Rating:        dbReview.Rating,
		Comment:       dbReview.Comment.String,
		RatingAspects: ratingAspects,
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
	if stats.AverageRating.Valid {
		averageRating = stats.AverageRating.Float64
	}

	return &ReviewStatistics{
		AverageRating: averageRating,
		TotalCount:    int32(stats.ReviewCount),
		RatingDistribution: map[int32]int32{
			1: int32(stats.OneStarCount),
			2: int32(stats.TwoStarCount),
			3: int32(stats.ThreeStarCount),
			4: int32(stats.FourStarCount),
			5: int32(stats.FiveStarCount),
		},
	}, nil
}

// updateSpotRating updates the spot's average rating and review count
func (s *ReviewService) updateSpotRating(ctx context.Context, spotID string) {
	stats, err := s.queries.GetSpotRatingStats(ctx, spotID)
	if err != nil {
		return // Silent failure for background update
	}

	averageRating := "0.0"
	if stats.AverageRating.Valid {
		averageRating = strconv.FormatFloat(stats.AverageRating.Float64, 'f', 1, 64)
	}

	// Update spot table with new rating statistics
	s.queries.UpdateSpotRating(ctx, database.UpdateSpotRatingParams{
		ID:            spotID,
		AverageRating: averageRating,
		ReviewCount:   int32(stats.ReviewCount),
	})
}