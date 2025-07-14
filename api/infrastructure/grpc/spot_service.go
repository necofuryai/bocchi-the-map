package grpc

import (
	"context"
	"database/sql"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	commonv1 "bocchi/api/gen/common/v1"
	spotv1 "bocchi/api/gen/spot/v1"
	"bocchi/api/infrastructure/database"
	"bocchi/api/pkg/logger"
)

// SpotService implements the gRPC SpotService
type SpotService struct {
	queries *database.Queries
}

// NewSpotService creates a new SpotService instance
func NewSpotService(db *sql.DB) *SpotService {
	return &SpotService{
		queries: database.New(db),
	}
}

// Use Protocol Buffers generated types
type (
	Spot            = spotv1.Spot
	GetSpotRequest  = spotv1.GetSpotRequest
	GetSpotResponse = spotv1.GetSpotResponse
)


// GetSpot retrieves a spot by ID
func (s *SpotService) GetSpot(ctx context.Context, req *GetSpotRequest) (*GetSpotResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	// Get spot from database
	dbSpot, err := s.queries.GetSpotByID(ctx, req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "spot not found")
		}
		return nil, status.Error(codes.Internal, "failed to get spot")
	}

	// Convert database spot to gRPC response
	// Convert GetSpotByIDRow to Spot type
	spotData := database.Spot{
		ID:            dbSpot.ID,
		Name:          dbSpot.Name,
		Category:      dbSpot.Category,
		Address:       dbSpot.Address,
		CountryCode:   dbSpot.CountryCode,
		AverageRating: dbSpot.AverageRating,
		ReviewCount:   dbSpot.ReviewCount,
		CreatedAt:     dbSpot.CreatedAt,
		UpdatedAt:     dbSpot.UpdatedAt,
	}
	spot := s.convertDatabaseSpotToGRPC(spotData)
	return &GetSpotResponse{Spot: spot}, nil
}



// convertDatabaseSpotToGRPC converts database spot model to gRPC spot struct
func (s *SpotService) convertDatabaseSpotToGRPC(dbSpot database.Spot) *Spot {
	averageRating, err := strconv.ParseFloat(dbSpot.AverageRating, 64)
	if err != nil {
		logger.ErrorWithFields("Failed to parse average rating", err, map[string]interface{}{
			"spot_id": dbSpot.ID,
			"rating_value": dbSpot.AverageRating,
		})
		averageRating = 0.0 // Use default value for invalid rating
	}


	return &Spot{
		Id:   dbSpot.ID,
		Name: dbSpot.Name,
		Category:      dbSpot.Category,
		Address:       dbSpot.Address,
		CountryCode:   dbSpot.CountryCode,
		AverageRating: averageRating,
		ReviewCount:   dbSpot.ReviewCount,
		CreatedAt:     timestamppb.New(dbSpot.CreatedAt),
		UpdatedAt:     timestamppb.New(dbSpot.UpdatedAt),
	}
}