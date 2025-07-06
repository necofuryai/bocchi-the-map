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

// Temporary structs until protobuf generates them
type Coordinates struct {
	Latitude  float64
	Longitude float64
}

type Spot struct {
	ID            string
	Name          string
	NameI18n      map[string]string
	Coordinates   *Coordinates
	Category      string
	Address       string
	AddressI18n   map[string]string
	CountryCode   string
	AverageRating float64
	ReviewCount   int32
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type CreateSpotRequest struct {
	Name        string
	NameI18n    map[string]string
	Coordinates *Coordinates
	Category    string
	Address     string
	AddressI18n map[string]string
	CountryCode string
}

type CreateSpotResponse struct {
	Spot *Spot
}

type GetSpotRequest struct {
	ID string
}

type GetSpotResponse struct {
	Spot *Spot
}

type PaginationRequest struct {
	Page     int32
	PageSize int32
}

type PaginationResponse struct {
	TotalCount int32
	Page       int32
	PageSize   int32
	TotalPages int32
}

type ListSpotsRequest struct {
	Pagination  *PaginationRequest
	Center      *Coordinates
	RadiusKm    float64
	Category    string
	CountryCode string
}

type ListSpotsResponse struct {
	Spots      []*Spot
	Pagination *PaginationResponse
}

type SearchSpotsRequest struct {
	Query      string
	Language   string // Simplified language handling
	Center     *Coordinates
	RadiusKm   float64
	Pagination *PaginationRequest
}

type SearchSpotsResponse struct {
	Spots      []*Spot
	Pagination *PaginationResponse
}

// CreateSpot creates a new spot
func (s *SpotService) CreateSpot(ctx context.Context, req *CreateSpotRequest) (*CreateSpotResponse, error) {
	// Validate request
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	if req.Coordinates == nil {
		return nil, status.Error(codes.InvalidArgument, "coordinates are required")
	}
	if req.Category == "" {
		return nil, status.Error(codes.InvalidArgument, "category is required")
	}
	if req.Address == "" {
		return nil, status.Error(codes.InvalidArgument, "address is required")
	}
	if req.CountryCode == "" {
		return nil, status.Error(codes.InvalidArgument, "country code is required")
	}

	// Generate UUID for new spot
	spotID := uuid.New().String()

	// Convert coordinates to strings (as expected by database)
	latitude := strconv.FormatFloat(req.Coordinates.Latitude, 'f', 8, 64)
	longitude := strconv.FormatFloat(req.Coordinates.Longitude, 'f', 8, 64)

	// Convert i18n maps to JSON
	var nameI18nJSON []byte
	if req.NameI18n == nil {
		nameI18nJSON = []byte("{}")
	} else {
		var err error
		nameI18nJSON, err = json.Marshal(req.NameI18n)
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to marshal name i18n")
		}
	}

	var addressI18nJSON []byte
	if req.AddressI18n == nil {
		addressI18nJSON = []byte("{}")
	} else {
		var err error
		addressI18nJSON, err = json.Marshal(req.AddressI18n)
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to marshal address i18n")
		}
	}

	// Create spot in database
	err := s.queries.CreateSpot(ctx, database.CreateSpotParams{
		ID:          spotID,
		Name:        req.Name,
		NameI18n:    nameI18nJSON,
		Latitude:    latitude,
		Longitude:   longitude,
		Category:    req.Category,
		Address:     req.Address,
		AddressI18n: addressI18nJSON,
		CountryCode: req.CountryCode,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create spot")
	}

	// Retrieve the created spot to get accurate timestamps and data
	dbSpot, err := s.queries.GetSpotByID(ctx, spotID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to retrieve created spot")
	}

	// Convert database spot to gRPC response
	spot := s.convertDatabaseSpotToGRPC(dbSpot)
	return &CreateSpotResponse{Spot: spot}, nil
}

// GetSpot retrieves a spot by ID
func (s *SpotService) GetSpot(ctx context.Context, req *GetSpotRequest) (*GetSpotResponse, error) {
	if req.ID == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	// Get spot from database
	dbSpot, err := s.queries.GetSpotByID(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "spot not found")
		}
		return nil, status.Error(codes.Internal, "failed to get spot")
	}

	// Convert database spot to gRPC response
	spot := s.convertDatabaseSpotToGRPC(dbSpot)
	return &GetSpotResponse{Spot: spot}, nil
}

// ListSpots lists spots with optional filters
func (s *SpotService) ListSpots(ctx context.Context, req *ListSpotsRequest) (*ListSpotsResponse, error) {
	// TODO: Implement actual listing logic
	// For now, return dummy data
	spot := &Spot{
		ID:            "spot_1",
		Name:          "Sample Cafe",
		Coordinates:   &Coordinates{Latitude: 35.6762, Longitude: 139.6503},
		Category:      "cafe",
		AverageRating: 4.5,
		ReviewCount:   10,
		CreatedAt:     time.Now().Add(-24 * time.Hour),
		UpdatedAt:     time.Now(),
	}

	pagination := &PaginationResponse{
		TotalCount: 1,
		Page:       1,
		PageSize:   20,
		TotalPages: 1,
	}
	if req.Pagination != nil {
		pagination.Page = req.Pagination.Page
		pagination.PageSize = req.Pagination.PageSize
	}

	return &ListSpotsResponse{
		Spots:      []*Spot{spot},
		Pagination: pagination,
	}, nil
}

// SearchSpots searches spots by query
func (s *SpotService) SearchSpots(ctx context.Context, req *SearchSpotsRequest) (*SearchSpotsResponse, error) {
	if req.Query == "" {
		return nil, status.Error(codes.InvalidArgument, "query is required")
	}

	// TODO: Implement actual search logic
	// For now, return dummy data similar to ListSpots
	spot := &Spot{
		ID:            "spot_search_1",
		Name:          fmt.Sprintf("Search Result for: %s", req.Query),
		Coordinates:   &Coordinates{Latitude: 35.6762, Longitude: 139.6503},
		Category:      "cafe",
		AverageRating: 4.5,
		ReviewCount:   10,
		CreatedAt:     time.Now().Add(-24 * time.Hour),
		UpdatedAt:     time.Now(),
	}

	pagination := &PaginationResponse{
		TotalCount: 1,
		Page:       1,
		PageSize:   20,
		TotalPages: 1,
	}
	if req.Pagination != nil {
		pagination.Page = req.Pagination.Page
		pagination.PageSize = req.Pagination.PageSize
	}

	return &SearchSpotsResponse{
		Spots:      []*Spot{spot},
		Pagination: pagination,
	}, nil
}

// convertDatabaseSpotToGRPC converts database spot model to gRPC spot struct
func (s *SpotService) convertDatabaseSpotToGRPC(dbSpot database.Spot) *Spot {
	// Parse coordinates from strings with error handling
	latitude, err := strconv.ParseFloat(dbSpot.Latitude, 64)
	if err != nil {
		logger.ErrorWithFields("Failed to parse latitude", err, map[string]interface{}{
			"spot_id": dbSpot.ID,
			"latitude_value": dbSpot.Latitude,
		})
		latitude = 0.0 // Use default value for invalid latitude
	}
	
	longitude, err := strconv.ParseFloat(dbSpot.Longitude, 64)
	if err != nil {
		logger.ErrorWithFields("Failed to parse longitude", err, map[string]interface{}{
			"spot_id": dbSpot.ID,
			"longitude_value": dbSpot.Longitude,
		})
		longitude = 0.0 // Use default value for invalid longitude
	}
	
	averageRating, err := strconv.ParseFloat(dbSpot.AverageRating, 64)
	if err != nil {
		logger.ErrorWithFields("Failed to parse average rating", err, map[string]interface{}{
			"spot_id": dbSpot.ID,
			"rating_value": dbSpot.AverageRating,
		})
		averageRating = 0.0 // Use default value for invalid rating
	}

	// Parse i18n JSON fields with error handling
	var nameI18n map[string]string
	if len(dbSpot.NameI18n) > 0 {
		if err := json.Unmarshal(dbSpot.NameI18n, &nameI18n); err != nil {
			logger.ErrorWithFields("Failed to parse name i18n JSON", err, map[string]interface{}{
				"spot_id": dbSpot.ID,
				"name_i18n_value": string(dbSpot.NameI18n),
			})
			nameI18n = nil // Use nil for invalid JSON data
		}
	}

	var addressI18n map[string]string
	if len(dbSpot.AddressI18n) > 0 {
		if err := json.Unmarshal(dbSpot.AddressI18n, &addressI18n); err != nil {
			logger.ErrorWithFields("Failed to parse address i18n JSON", err, map[string]interface{}{
				"spot_id": dbSpot.ID,
				"address_i18n_value": string(dbSpot.AddressI18n),
			})
			addressI18n = nil // Use nil for invalid JSON data
		}
	}

	return &Spot{
		ID:   dbSpot.ID,
		Name: dbSpot.Name,
		NameI18n: nameI18n,
		Coordinates: &Coordinates{
			Latitude:  latitude,
			Longitude: longitude,
		},
		Category:      dbSpot.Category,
		Address:       dbSpot.Address,
		AddressI18n:   addressI18n,
		CountryCode:   dbSpot.CountryCode,
		AverageRating: averageRating,
		ReviewCount:   dbSpot.ReviewCount,
		CreatedAt:     dbSpot.CreatedAt,
		UpdatedAt:     dbSpot.UpdatedAt,
	}
}