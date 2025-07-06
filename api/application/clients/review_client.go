package clients

import (
	"context"
	"database/sql"
	"fmt"

	"google.golang.org/grpc"

	grpcSvc "bocchi/api/infrastructure/grpc"
)

// ReviewClient wraps gRPC client calls for review operations
type ReviewClient struct {
	service *grpcSvc.ReviewService
	conn    *grpc.ClientConn
}

// NewReviewClient creates a new review client
func NewReviewClient(serviceAddr string, db *sql.DB) (*ReviewClient, error) {
	// For internal communication in monolith, we can use direct service calls
	// In a true microservice setup, this would connect to remote gRPC service
	if serviceAddr == "internal" {
		return &ReviewClient{
			service: grpcSvc.NewReviewService(db),
		}, nil
	}

	// TODO: Implement external gRPC service connection when protobuf client is ready
	// For now, return error for external services to avoid silent failures
	return nil, fmt.Errorf("external gRPC service not implemented yet: %s", serviceAddr)
}

// Close closes the gRPC connection
func (c *ReviewClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// CreateReview creates a new review via gRPC
func (c *ReviewClient) CreateReview(ctx context.Context, req *grpcSvc.CreateReviewRequest) (*grpcSvc.CreateReviewResponse, error) {
	return c.service.CreateReview(ctx, req)
}

// GetSpotReviews retrieves reviews for a spot via gRPC
func (c *ReviewClient) GetSpotReviews(ctx context.Context, req *grpcSvc.GetSpotReviewsRequest) (*grpcSvc.GetSpotReviewsResponse, error) {
	return c.service.GetSpotReviews(ctx, req)
}

// GetUserReviews retrieves reviews by user via gRPC
func (c *ReviewClient) GetUserReviews(ctx context.Context, req *grpcSvc.GetUserReviewsRequest) (*grpcSvc.GetUserReviewsResponse, error) {
	return c.service.GetUserReviews(ctx, req)
}