package handlers

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"bocchi/api/application/clients"
	"bocchi/api/pkg/auth"
	"bocchi/api/pkg/errors"
	reviewv1 "bocchi/api/gen/review/v1"
	commonv1 "bocchi/api/gen/common/v1"
)

// ReviewHandler handles review-related HTTP requests
type ReviewHandler struct {
	reviewClient *clients.ReviewClient
}

// NewReviewHandler creates a new review handler
func NewReviewHandler(reviewClient *clients.ReviewClient) *ReviewHandler {
	if reviewClient == nil {
		panic("reviewClient cannot be nil")
	}
	return &ReviewHandler{
		reviewClient: reviewClient,
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

// CreateReviewOutput represents the response for review creation (using protobuf Review type)
type CreateReviewOutput struct {
	Body *reviewv1.Review `json:"review" doc:"Created review data"`
}

// GetSpotReviewsInput represents the request to get reviews for a spot
type GetSpotReviewsInput struct {
	SpotID string `path:"spot_id" maxLength:"36" doc:"Spot ID"`
	Page   int32  `query:"page" minimum:"1" default:"1" doc:"Page number"`
	Limit  int32  `query:"limit" minimum:"1" maximum:"50" default:"20" doc:"Number of reviews per page"`
}

// GetSpotReviewsOutput represents the response for getting spot reviews (using protobuf types)
type GetSpotReviewsOutput struct {
	Body struct {
		Reviews    []*reviewv1.Review          `json:"reviews" doc:"List of reviews"`
		Pagination *commonv1.PaginationResponse `json:"pagination" doc:"Pagination information"`
		Statistics *reviewv1.ReviewStatistics   `json:"statistics" doc:"Review statistics"`
	}
}

// GetUserReviewsInput represents the request to get reviews by a user
type GetUserReviewsInput struct {
	UserID string `path:"user_id" maxLength:"36" doc:"User ID"`
	Page   int32  `query:"page" minimum:"1" default:"1" doc:"Page number"`
	Limit  int32  `query:"limit" minimum:"1" maximum:"50" default:"20" doc:"Number of reviews per page"`
}

// GetUserReviewsOutput represents the response for getting user reviews (using protobuf types)
type GetUserReviewsOutput struct {
	Body struct {
		Reviews    []*reviewv1.Review          `json:"reviews" doc:"List of reviews"`
		Pagination *commonv1.PaginationResponse `json:"pagination" doc:"Pagination information"`
	}
}


// RegisterRoutes registers review routes
func (h *ReviewHandler) RegisterRoutes(api huma.API) {
	// Get reviews for a spot (public)
	huma.Register(api, huma.Operation{
		OperationID: "get-spot-reviews",
		Method:      http.MethodGet,
		Path:        "/api/v1/spots/{spot_id}/reviews",
		Summary:     "Get reviews for a spot",
		Description: "Get paginated reviews for a specific spot with statistics",
		Tags:        []string{"Reviews"},
	}, h.GetSpotReviews)

	// Get reviews by a user (public)
	huma.Register(api, huma.Operation{
		OperationID: "get-user-reviews",
		Method:      http.MethodGet,
		Path:        "/api/v1/users/{user_id}/reviews",
		Summary:     "Get reviews by a user",
		Description: "Get paginated reviews created by a specific user",
		Tags:        []string{"Reviews"},
	}, h.GetUserReviews)
}

// RegisterRoutesWithAuth registers review routes with authentication middleware
func (h *ReviewHandler) RegisterRoutesWithAuth(api huma.API, authMiddleware *auth.AuthMiddleware) {
	// Register public routes first
	h.RegisterRoutes(api)

	// Create review (protected - requires authentication)
	huma.Register(api, authMiddleware.CreateProtectedOperation(huma.Operation{
		OperationID: "create-review",
		Method:      http.MethodPost,
		Path:        "/api/v1/reviews",
		Summary:     "Create a review",
		Description: "Create a new review for a spot (requires authentication)",
		Tags:        []string{"Reviews"},
	}), h.CreateReview)
}

// CreateReview creates a new review
func (h *ReviewHandler) CreateReview(ctx context.Context, input *CreateReviewInput) (*CreateReviewOutput, error) {
	// Extract user ID from Huma v2 authentication context
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		return nil, huma.Error401Unauthorized("authentication required to create review")
	}

	// Add user ID to context for gRPC service access
	ctx = errors.WithUserID(ctx, userID)

	// Call gRPC service via client with authenticated user context
	resp, err := h.reviewClient.CreateReview(ctx, &reviewv1.CreateReviewRequest{
		SpotId:        input.Body.SpotID,
		Rating:        input.Body.Rating,
		Comment:       input.Body.Comment,
		RatingAspects: input.Body.RatingAspects,
	})
	if err != nil {
		return nil, grpcToHTTPError(err, "failed to create review")
	}

	// Convert gRPC response to HTTP response
	return &CreateReviewOutput{
		Body: resp.Review,
	}, nil
}

// GetSpotReviews gets reviews for a specific spot
func (h *ReviewHandler) GetSpotReviews(ctx context.Context, input *GetSpotReviewsInput) (*GetSpotReviewsOutput, error) {
	// Call gRPC service via client
	resp, err := h.reviewClient.GetSpotReviews(ctx, &reviewv1.GetSpotReviewsRequest{
		SpotId: input.SpotID,
		Pagination: &commonv1.PaginationRequest{
			Page:     input.Page,
			PageSize: input.Limit,
		},
	})
	if err != nil {
		return nil, grpcToHTTPError(err, "failed to get spot reviews")
	}

	return &GetSpotReviewsOutput{
		Body: struct {
			Reviews    []*reviewv1.Review          `json:"reviews" doc:"List of reviews"`
			Pagination *commonv1.PaginationResponse `json:"pagination" doc:"Pagination information"`
			Statistics *reviewv1.ReviewStatistics   `json:"statistics" doc:"Review statistics"`
		}{
			Reviews:    resp.Reviews,
			Pagination: resp.Pagination,
			Statistics: resp.Statistics,
		},
	}, nil
}

// GetUserReviews gets reviews by a specific user
func (h *ReviewHandler) GetUserReviews(ctx context.Context, input *GetUserReviewsInput) (*GetUserReviewsOutput, error) {
	// Call gRPC service via client
	resp, err := h.reviewClient.GetUserReviews(ctx, &reviewv1.GetUserReviewsRequest{
		UserId: input.UserID,
		Pagination: &commonv1.PaginationRequest{
			Page:     input.Page,
			PageSize: input.Limit,
		},
	})
	if err != nil {
		return nil, grpcToHTTPError(err, "failed to get user reviews")
	}

	return &GetUserReviewsOutput{
		Body: struct {
			Reviews    []*reviewv1.Review          `json:"reviews" doc:"List of reviews"`
			Pagination *commonv1.PaginationResponse `json:"pagination" doc:"Pagination information"`
		}{
			Reviews:    resp.Reviews,
			Pagination: resp.Pagination,
		},
	}, nil
}