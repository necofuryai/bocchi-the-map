package clients

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	grpcSvc "github.com/necofuryai/bocchi-the-map/api/infrastructure/grpc"
)

// SpotClient wraps gRPC client calls for spot operations
type SpotClient struct {
	service *grpcSvc.SpotService
	conn    *grpc.ClientConn
}

// NewSpotClient creates a new spot client
func NewSpotClient(serviceAddr string) (*SpotClient, error) {
	// For internal communication in monolith, we can use direct service calls
	// In a true microservice setup, this would connect to remote gRPC service
	if serviceAddr == "internal" {
		return &SpotClient{
			service: grpcSvc.NewSpotService(),
		}, nil
	}

	// For external gRPC service connection
	conn, err := grpc.NewClient(serviceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to spot service: %w", err)
	}

	return &SpotClient{
		conn: conn,
		// TODO: Use generated gRPC client when protobuf is available
		service: grpcSvc.NewSpotService(),
	}, nil
}

// Close closes the gRPC connection
func (c *SpotClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// CreateSpot creates a new spot via gRPC
func (c *SpotClient) CreateSpot(ctx context.Context, req *grpcSvc.CreateSpotRequest) (*grpcSvc.CreateSpotResponse, error) {
	return c.service.CreateSpot(ctx, req)
}

// GetSpot retrieves a spot by ID via gRPC
func (c *SpotClient) GetSpot(ctx context.Context, req *grpcSvc.GetSpotRequest) (*grpcSvc.GetSpotResponse, error) {
	return c.service.GetSpot(ctx, req)
}

// ListSpots lists spots with filters via gRPC
func (c *SpotClient) ListSpots(ctx context.Context, req *grpcSvc.ListSpotsRequest) (*grpcSvc.ListSpotsResponse, error) {
	return c.service.ListSpots(ctx, req)
}

// SearchSpots searches spots via gRPC
func (c *SpotClient) SearchSpots(ctx context.Context, req *grpcSvc.SearchSpotsRequest) (*grpcSvc.SearchSpotsResponse, error) {
	return c.service.SearchSpots(ctx, req)
}