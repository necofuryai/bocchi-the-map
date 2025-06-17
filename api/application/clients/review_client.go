package clients

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	grpcSvc "github.com/necofuryai/bocchi-the-map/api/infrastructure/grpc"
)

// ReviewClient wraps gRPC client calls for review operations
type ReviewClient struct {
	service *grpcSvc.ReviewService
	conn    *grpc.ClientConn
}

// NewReviewClient creates a new review client
func NewReviewClient(serviceAddr string) (*ReviewClient, error) {
	// For internal communication in monolith, we can use direct service calls
	// In a true microservice setup, this would connect to remote gRPC service
	if serviceAddr == "internal" {
		return &ReviewClient{
			service: grpcSvc.NewReviewService(),
		}, nil
	}

	// For external gRPC service connection
	// TODO: Use TLS credentials in production
	var creds credentials.TransportCredentials
	if os.Getenv("GRPC_INSECURE") == "true" {
		creds = insecure.NewCredentials()
	} else {
		creds = credentials.NewTLS(&tls.Config{
			MinVersion: tls.VersionTLS13,
		})
	}
	conn, err := grpc.Dial(serviceAddr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to review service: %w", err)
	}

	return &ReviewClient{
		conn: conn,
		// TODO: Use generated gRPC client when protobuf is available
		service: grpcSvc.NewReviewService(),
	}, nil
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