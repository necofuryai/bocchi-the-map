package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/necofuryai/bocchi-the-map/api/application/clients"
	"github.com/necofuryai/bocchi-the-map/api/pkg/errors"
	grpcSvc "github.com/necofuryai/bocchi-the-map/api/infrastructure/grpc"
)

// ReviewHandler handles review-related HTTP requests
type ReviewHandler struct {
	reviewClient *clients.ReviewClient
}

// NewReviewHandler creates a new review handler
func NewReviewHandler(reviewClient *clients.ReviewClient) *ReviewHandler {
	return &ReviewHandler{
		reviewClient: reviewClient,
	}
}

// grpcToHTTPError converts gRPC errors to appropriate HTTP error responses
func grpcToHTTPError(err error, defaultMessage string) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return huma.Error500InternalServerError(defaultMessage)
	}

	switch st.Code() {
	case codes.NotFound:
		return huma.Error404NotFound(st.Message())
	case codes.InvalidArgument:
		return huma.Error400BadRequest(st.Message())
	case codes.AlreadyExists:
		return huma.Error409Conflict(st.Message())
	case codes.PermissionDenied:
		return huma.Error403Forbidden(st.Message())
	case codes.Unauthenticated:
		return huma.Error401Unauthorized(st.Message())
	case codes.FailedPrecondition:
		return huma.Error412PreconditionFailed(st.Message())
	case codes.OutOfRange:
		return huma.Error400BadRequest(st.Message())
	case codes.Unimplemented:
		return huma.Error501NotImplemented(st.Message())
	case codes.Unavailable:
		return huma.Error503ServiceUnavailable(st.Message())
	case codes.DeadlineExceeded:
		return huma.Error503ServiceUnavailable(st.Message())
	default:
		return huma.Error500InternalServerError(defaultMessage)
	}
}

// CreateReviewInput represents the review creation request
type CreateReviewInput struct {
	Body struct {
		SpotID        string            `json:"spot_id" maxLength:"36" doc:"Spot ID to review"`
		Rating        int32             `json:"rating" minimum:"1" maximum:"5" doc:"Rating from 1 to 5"`
		Comment       string            `json:"comment,omitempty" maxLength:"1000" doc:"Optional review comment"`
		RatingAspects map[string]int32  `json:"rating_aspects,omitempty" doc:"Optional aspect ratings"`
	}
}

// CreateReviewOutput represents the response for review creation
type CreateReviewOutput struct {
	Body struct {
		ID            string            `json:"id" doc:"Review ID"`
		SpotID        string            `json:"spot_id" doc:"Spot ID"`
		UserID        string            `json:"user_id" doc:"User ID"`
		Rating        int32             `json:"rating" doc:"Rating from 1 to 5"`
		Comment       string            `json:"comment,omitempty" doc:"Review comment"`
		RatingAspects map[string]int32  `json:"rating_aspects,omitempty" doc:"Aspect ratings"`
		CreatedAt     time.Time         `json:"created_at" doc:"Creation timestamp"`
		UpdatedAt     time.Time         `json:"updated_at" doc:"Last update timestamp"`
	}
}

// GetSpotReviewsInput represents the request to get reviews for a spot
type GetSpotReviewsInput struct {
	SpotID string `path:"spot_id" maxLength:"36" doc:"Spot ID"`
	Page   int32  `query:"page" minimum:"1" default:"1" doc:"Page number"`
	Limit  int32  `query:"limit" minimum:"1" maximum:"50" default:"20" doc:"Number of reviews per page"`
}

// GetSpotReviewsOutput represents the response for getting spot reviews
type GetSpotReviewsOutput struct {
	Body struct {
		Reviews    []ReviewResponse     `json:"reviews" doc:"List of reviews"`
		Pagination PaginationResponse   `json:"pagination" doc:"Pagination information"`
		Statistics ReviewStatistics     `json:"statistics" doc:"Review statistics"`
	}
}

// GetUserReviewsInput represents the request to get reviews by a user
type GetUserReviewsInput struct {
	UserID string `path:"user_id" maxLength:"36" doc:"User ID"`
	Page   int32  `query:"page" minimum:"1" default:"1" doc:"Page number"`
	Limit  int32  `query:"limit" minimum:"1" maximum:"50" default:"20" doc:"Number of reviews per page"`
}

// GetUserReviewsOutput represents the response for getting user reviews
type GetUserReviewsOutput struct {
	Body struct {
		Reviews    []ReviewResponse     `json:"reviews" doc:"List of reviews"`
		Pagination PaginationResponse   `json:"pagination" doc:"Pagination information"`
	}
}

// ReviewResponse represents a review in HTTP responses
type ReviewResponse struct {
	ID            string            `json:"id" doc:"Review ID"`
	SpotID        string            `json:"spot_id" doc:"Spot ID"`
	UserID        string            `json:"user_id" doc:"User ID"`
	Rating        int32             `json:"rating" doc:"Rating from 1 to 5"`
	Comment       string            `json:"comment,omitempty" doc:"Review comment"`
	RatingAspects map[string]int32  `json:"rating_aspects,omitempty" doc:"Aspect ratings"`
	CreatedAt     time.Time         `json:"created_at" doc:"Creation timestamp"`
	UpdatedAt     time.Time         `json:"updated_at" doc:"Last update timestamp"`
}

// PaginationResponse represents pagination information
type PaginationResponse struct {
	TotalCount int32 `json:"total_count" doc:"Total number of items"`
	Page       int32 `json:"page" doc:"Current page number"`
	PageSize   int32 `json:"page_size" doc:"Number of items per page"`
	TotalPages int32 `json:"total_pages" doc:"Total number of pages"`
}

// ReviewStatistics represents review statistics for a spot
type ReviewStatistics struct {
	AverageRating      float64           `json:"average_rating" doc:"Average rating"`
	TotalCount         int32             `json:"total_count" doc:"Total number of reviews"`
	RatingDistribution map[int32]int32   `json:"rating_distribution" doc:"Distribution of ratings (1-5)"`
}

// RegisterRoutes registers review routes
func (h *ReviewHandler) RegisterRoutes(api huma.API) {
	// Create review
	huma.Register(api, huma.Operation{
		OperationID: "create-review",
		Method:      http.MethodPost,
		Path:        "/api/v1/reviews",
		Summary:     "Create a review",
		Description: "Create a new review for a spot",
		Tags:        []string{"Reviews"},
	}, h.CreateReview)

	// Get reviews for a spot
	huma.Register(api, huma.Operation{
		OperationID: "get-spot-reviews",
		Method:      http.MethodGet,
		Path:        "/api/v1/spots/{spot_id}/reviews",
		Summary:     "Get reviews for a spot",
		Description: "Get paginated reviews for a specific spot with statistics",
		Tags:        []string{"Reviews"},
	}, h.GetSpotReviews)

	// Get reviews by a user
	huma.Register(api, huma.Operation{
		OperationID: "get-user-reviews",
		Method:      http.MethodGet,
		Path:        "/api/v1/users/{user_id}/reviews",
		Summary:     "Get reviews by a user",
		Description: "Get paginated reviews created by a specific user",
		Tags:        []string{"Reviews"},
	}, h.GetUserReviews)
}

// CreateReview creates a new review
func (h *ReviewHandler) CreateReview(ctx context.Context, input *CreateReviewInput) (*CreateReviewOutput, error) {
	// Extract user ID from authentication context
	userID := errors.GetUserID(ctx)
	if userID == "" {
		return nil, huma.Error401Unauthorized("authentication required to create review")
	}

	// Call gRPC service via client with authenticated user context
	resp, err := h.reviewClient.CreateReview(ctx, &grpcSvc.CreateReviewRequest{
		SpotID:        input.Body.SpotID,
		Rating:        input.Body.Rating,
		Comment:       input.Body.Comment,
		RatingAspects: input.Body.RatingAspects,
	})
	if err != nil {
		return nil, grpcToHTTPError(err, "failed to create review")
	}

	// Convert gRPC response to HTTP response
	return &CreateReviewOutput{
		Body: ReviewResponse{
			ID:            resp.Review.ID,
			SpotID:        resp.Review.SpotID,
			UserID:        resp.Review.UserID,
			Rating:        resp.Review.Rating,
			Comment:       resp.Review.Comment,
			RatingAspects: resp.Review.RatingAspects,
			CreatedAt:     resp.Review.CreatedAt,
			UpdatedAt:     resp.Review.UpdatedAt,
		},
	}, nil
}

// GetSpotReviews gets reviews for a specific spot
func (h *ReviewHandler) GetSpotReviews(ctx context.Context, input *GetSpotReviewsInput) (*GetSpotReviewsOutput, error) {
	// Call gRPC service via client
	resp, err := h.reviewClient.GetSpotReviews(ctx, &grpcSvc.GetSpotReviewsRequest{
		SpotID: input.SpotID,
		Pagination: &grpcSvc.PaginationRequest{
			Page:     input.Page,
			PageSize: input.Limit,
		},
	})
	if err != nil {
		return nil, grpcToHTTPError(err, "failed to get spot reviews")
	}

	// Convert gRPC reviews to HTTP format
	reviews := make([]ReviewResponse, len(resp.Reviews))
	for i, review := range resp.Reviews {
		reviews[i] = ReviewResponse{
			ID:            review.ID,
			SpotID:        review.SpotID,
			UserID:        review.UserID,
			Rating:        review.Rating,
			Comment:       review.Comment,
			RatingAspects: review.RatingAspects,
			CreatedAt:     review.CreatedAt,
			UpdatedAt:     review.UpdatedAt,
		}
	}

	// Convert pagination and statistics
	pagination := PaginationResponse{
		TotalCount: resp.Pagination.TotalCount,
		Page:       resp.Pagination.Page,
		PageSize:   resp.Pagination.PageSize,
		TotalPages: resp.Pagination.TotalPages,
	}

	statistics := ReviewStatistics{
		AverageRating:      resp.Statistics.AverageRating,
		TotalCount:         resp.Statistics.TotalCount,
		RatingDistribution: resp.Statistics.RatingDistribution,
	}

	return &GetSpotReviewsOutput{
		Body: struct {
			Reviews    []ReviewResponse     `json:"reviews" doc:"List of reviews"`
			Pagination PaginationResponse   `json:"pagination" doc:"Pagination information"`
			Statistics ReviewStatistics     `json:"statistics" doc:"Review statistics"`
		}{
			Reviews:    reviews,
			Pagination: pagination,
			Statistics: statistics,
		},
	}, nil
}

// GetUserReviews gets reviews by a specific user
func (h *ReviewHandler) GetUserReviews(ctx context.Context, input *GetUserReviewsInput) (*GetUserReviewsOutput, error) {
	// Call gRPC service via client
	resp, err := h.reviewClient.GetUserReviews(ctx, &grpcSvc.GetUserReviewsRequest{
		UserID: input.UserID,
		Pagination: &grpcSvc.PaginationRequest{
			Page:     input.Page,
			PageSize: input.Limit,
		},
	})
	if err != nil {
		return nil, grpcToHTTPError(err, "failed to get user reviews")
	}

	// Convert gRPC reviews to HTTP format
	reviews := make([]ReviewResponse, len(resp.Reviews))
	for i, review := range resp.Reviews {
		reviews[i] = ReviewResponse{
			ID:            review.ID,
			SpotID:        review.SpotID,
			UserID:        review.UserID,
			Rating:        review.Rating,
			Comment:       review.Comment,
			RatingAspects: review.RatingAspects,
			CreatedAt:     review.CreatedAt,
			UpdatedAt:     review.UpdatedAt,
		}
	}

	// Convert pagination
	pagination := PaginationResponse{
		TotalCount: resp.Pagination.TotalCount,
		Page:       resp.Pagination.Page,
		PageSize:   resp.Pagination.PageSize,
		TotalPages: resp.Pagination.TotalPages,
	}

	return &GetUserReviewsOutput{
		Body: struct {
			Reviews    []ReviewResponse     `json:"reviews" doc:"List of reviews"`
			Pagination PaginationResponse   `json:"pagination" doc:"Pagination information"`
		}{
			Reviews:    reviews,
			Pagination: pagination,
		},
	}, nil
}