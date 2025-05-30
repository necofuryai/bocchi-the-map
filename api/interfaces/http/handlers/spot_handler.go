package handlers

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// SpotHandler handles spot-related HTTP requests
type SpotHandler struct {
	// TODO: Add spot service dependency
}

// NewSpotHandler creates a new spot handler
func NewSpotHandler() *SpotHandler {
	return &SpotHandler{}
}

// CreateSpotInput represents the request to create a spot
type CreateSpotInput struct {
	Body struct {
		Name        string            `json:"name" minLength:"1" maxLength:"255" doc:"Spot name"`
		NameI18n    map[string]string `json:"name_i18n,omitempty" doc:"Localized names"`
		Latitude    float64           `json:"latitude" minimum:"-90" maximum:"90" doc:"Latitude"`
		Longitude   float64           `json:"longitude" minimum:"-180" maximum:"180" doc:"Longitude"`
		Category    string            `json:"category" minLength:"1" maxLength:"50" doc:"Category"`
		Address     string            `json:"address" minLength:"1" maxLength:"500" doc:"Address"`
		AddressI18n map[string]string `json:"address_i18n,omitempty" doc:"Localized addresses"`
		CountryCode string            `json:"country_code" minLength:"2" maxLength:"2" pattern:"^[A-Z]{2}$" doc:"ISO 3166-1 alpha-2 country code"`
	}
}

// CreateSpotOutput represents the response for spot creation
type CreateSpotOutput struct {
	Body struct {
		ID          string            `json:"id" doc:"Spot ID"`
		Name        string            `json:"name" doc:"Spot name"`
		NameI18n    map[string]string `json:"name_i18n,omitempty" doc:"Localized names"`
		Latitude    float64           `json:"latitude" doc:"Latitude"`
		Longitude   float64           `json:"longitude" doc:"Longitude"`
		Category    string            `json:"category" doc:"Category"`
		Address     string            `json:"address" doc:"Address"`
		AddressI18n map[string]string `json:"address_i18n,omitempty" doc:"Localized addresses"`
		CountryCode string            `json:"country_code" doc:"ISO country code"`
		CreatedAt   string            `json:"created_at" doc:"Creation timestamp"`
	}
}

// GetSpotInput represents the request to get a spot
type GetSpotInput struct {
	ID string `path:"id" doc:"Spot ID"`
}

// GetSpotOutput represents the response for getting a spot
type GetSpotOutput struct {
	Body struct {
		ID            string            `json:"id" doc:"Spot ID"`
		Name          string            `json:"name" doc:"Spot name"`
		NameI18n      map[string]string `json:"name_i18n,omitempty" doc:"Localized names"`
		Latitude      float64           `json:"latitude" doc:"Latitude"`
		Longitude     float64           `json:"longitude" doc:"Longitude"`
		Category      string            `json:"category" doc:"Category"`
		Address       string            `json:"address" doc:"Address"`
		AddressI18n   map[string]string `json:"address_i18n,omitempty" doc:"Localized addresses"`
		CountryCode   string            `json:"country_code" doc:"ISO country code"`
		AverageRating float64           `json:"average_rating" doc:"Average rating"`
		ReviewCount   int               `json:"review_count" doc:"Number of reviews"`
		CreatedAt     string            `json:"created_at" doc:"Creation timestamp"`
		UpdatedAt     string            `json:"updated_at" doc:"Last update timestamp"`
	}
}

// ListSpotsInput represents the request to list spots
type ListSpotsInput struct {
	Page        int     `query:"page" default:"1" minimum:"1" doc:"Page number"`
	PageSize    int     `query:"page_size" default:"20" minimum:"1" maximum:"100" doc:"Items per page"`
	Latitude    float64 `query:"lat,omitempty" minimum:"-90" maximum:"90" doc:"Center latitude"`
	Longitude   float64 `query:"lng,omitempty" minimum:"-180" maximum:"180" doc:"Center longitude"`
	RadiusKm    float64 `query:"radius_km,omitempty" minimum:"0.1" maximum:"50" doc:"Search radius in km"`
	Category    string  `query:"category,omitempty" doc:"Filter by category"`
	CountryCode string  `query:"country_code,omitempty" doc:"Filter by country code"`
}

// ListSpotsOutput represents the response for listing spots
type ListSpotsOutput struct {
	Body struct {
		Spots []struct {
			ID            string  `json:"id" doc:"Spot ID"`
			Name          string  `json:"name" doc:"Spot name"`
			Latitude      float64 `json:"latitude" doc:"Latitude"`
			Longitude     float64 `json:"longitude" doc:"Longitude"`
			Category      string  `json:"category" doc:"Category"`
			AverageRating float64 `json:"average_rating" doc:"Average rating"`
			ReviewCount   int     `json:"review_count" doc:"Number of reviews"`
		} `json:"spots" doc:"List of spots"`
		Pagination struct {
			TotalCount int `json:"total_count" doc:"Total number of items"`
			Page       int `json:"page" doc:"Current page"`
			PageSize   int `json:"page_size" doc:"Items per page"`
			TotalPages int `json:"total_pages" doc:"Total number of pages"`
		} `json:"pagination" doc:"Pagination metadata"`
	}
}

// RegisterRoutes registers spot routes
func (h *SpotHandler) RegisterRoutes(api huma.API) {
	// Create spot
	huma.Register(api, huma.Operation{
		OperationID: "create-spot",
		Method:      http.MethodPost,
		Path:        "/spots",
		Summary:     "Create a new spot",
		Description: "Create a new reviewable spot on the map",
		Tags:        []string{"Spots"},
	}, h.CreateSpot)

	// Get spot
	huma.Register(api, huma.Operation{
		OperationID: "get-spot",
		Method:      http.MethodGet,
		Path:        "/spots/{id}",
		Summary:     "Get a spot",
		Description: "Get details of a specific spot",
		Tags:        []string{"Spots"},
	}, h.GetSpot)

	// List spots
	huma.Register(api, huma.Operation{
		OperationID: "list-spots",
		Method:      http.MethodGet,
		Path:        "/spots",
		Summary:     "List spots",
		Description: "List spots with optional filters",
		Tags:        []string{"Spots"},
	}, h.ListSpots)
}

// CreateSpot creates a new spot
func (h *SpotHandler) CreateSpot(ctx context.Context, input *CreateSpotInput) (*CreateSpotOutput, error) {
	// TODO: Implement spot creation logic
	resp := &CreateSpotOutput{}
	resp.Body.ID = "spot_dummy_id"
	resp.Body.Name = input.Body.Name
	resp.Body.NameI18n = input.Body.NameI18n
	resp.Body.Latitude = input.Body.Latitude
	resp.Body.Longitude = input.Body.Longitude
	resp.Body.Category = input.Body.Category
	resp.Body.Address = input.Body.Address
	resp.Body.AddressI18n = input.Body.AddressI18n
	resp.Body.CountryCode = input.Body.CountryCode
	resp.Body.CreatedAt = "2024-01-01T00:00:00Z"
	
	return resp, nil
}

// GetSpot gets a specific spot
func (h *SpotHandler) GetSpot(ctx context.Context, input *GetSpotInput) (*GetSpotOutput, error) {
	// TODO: Implement get spot logic
	resp := &GetSpotOutput{}
	resp.Body.ID = input.ID
	resp.Body.Name = "Sample Spot"
	resp.Body.Latitude = 35.6762
	resp.Body.Longitude = 139.6503
	resp.Body.Category = "cafe"
	resp.Body.Address = "Tokyo, Japan"
	resp.Body.CountryCode = "JP"
	resp.Body.AverageRating = 4.5
	resp.Body.ReviewCount = 10
	resp.Body.CreatedAt = "2024-01-01T00:00:00Z"
	resp.Body.UpdatedAt = "2024-01-01T00:00:00Z"
	
	return resp, nil
}

// ListSpots lists spots
func (h *SpotHandler) ListSpots(ctx context.Context, input *ListSpotsInput) (*ListSpotsOutput, error) {
	// TODO: Implement list spots logic
	resp := &ListSpotsOutput{}
	
	// Dummy data
	spot := struct {
		ID            string  `json:"id" doc:"Spot ID"`
		Name          string  `json:"name" doc:"Spot name"`
		Latitude      float64 `json:"latitude" doc:"Latitude"`
		Longitude     float64 `json:"longitude" doc:"Longitude"`
		Category      string  `json:"category" doc:"Category"`
		AverageRating float64 `json:"average_rating" doc:"Average rating"`
		ReviewCount   int     `json:"review_count" doc:"Number of reviews"`
	}{
		ID:            "spot_1",
		Name:          "Sample Cafe",
		Latitude:      35.6762,
		Longitude:     139.6503,
		Category:      "cafe",
		AverageRating: 4.5,
		ReviewCount:   10,
	}
	
	resp.Body.Spots = append(resp.Body.Spots, spot)
	resp.Body.Pagination.TotalCount = 1
	resp.Body.Pagination.Page = input.Page
	resp.Body.Pagination.PageSize = input.PageSize
	resp.Body.Pagination.TotalPages = 1
	
	return resp, nil
}