package handlers

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"bocchi/api/application/clients"
	"bocchi/api/pkg/auth"
	spotv1 "bocchi/api/gen/spot/v1"
	commonv1 "bocchi/api/gen/common/v1"
)


// SpotHandler handles spot-related HTTP requests
type SpotHandler struct {
	spotClient *clients.SpotClient
}

// NewSpotHandler creates a new spot handler
func NewSpotHandler(spotClient *clients.SpotClient) *SpotHandler {
	return &SpotHandler{
		spotClient: spotClient,
	}
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

// CreateSpotOutput represents the response for spot creation (using protobuf Spot type)
type CreateSpotOutput struct {
	Body *spotv1.Spot `json:"spot" doc:"Created spot data"`
}

// GetSpotInput represents the request to get a spot
type GetSpotInput struct {
	ID string `path:"id" doc:"Spot ID"`
}

// GetSpotOutput represents the response for getting a spot (using protobuf Spot type)
type GetSpotOutput struct {
	Body *spotv1.Spot `json:"spot" doc:"Spot data"`
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

// ListSpotsOutput represents the response for listing spots (using protobuf types)
type ListSpotsOutput struct {
	Body struct {
		Spots      []*spotv1.Spot             `json:"spots" doc:"List of spots"`
		Pagination *commonv1.PaginationResponse `json:"pagination" doc:"Pagination metadata"`
	}
}

// UpdateSpotInput represents the request to update a spot
type UpdateSpotInput struct {
	ID string `path:"id" doc:"Spot ID"`
	Body struct {
		Name        *string            `json:"name,omitempty" minLength:"1" maxLength:"255" doc:"Spot name"`
		NameI18n    *map[string]string `json:"name_i18n,omitempty" doc:"Localized names"`
		Category    *string            `json:"category,omitempty" minLength:"1" maxLength:"50" doc:"Category"`
		Address     *string            `json:"address,omitempty" minLength:"1" maxLength:"500" doc:"Address"`
		AddressI18n *map[string]string `json:"address_i18n,omitempty" doc:"Localized addresses"`
	}
}

// UpdateSpotOutput represents the response for spot update
type UpdateSpotOutput struct {
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
		UpdatedAt   string            `json:"updated_at" doc:"Last update timestamp"`
	}
}

// RegisterRoutes registers spot routes
func (h *SpotHandler) RegisterRoutes(api huma.API) {
	// Get spot (public)
	huma.Register(api, huma.Operation{
		OperationID: "get-spot",
		Method:      http.MethodGet,
		Path:        "/api/v1/spots/{id}",
		Summary:     "Get a spot",
		Description: "Get details of a specific spot",
		Tags:        []string{"Spots"},
	}, h.GetSpot)

	// List spots (public)
	huma.Register(api, huma.Operation{
		OperationID: "list-spots",
		Method:      http.MethodGet,
		Path:        "/api/v1/spots",
		Summary:     "List spots",
		Description: "List spots with optional filters",
		Tags:        []string{"Spots"},
	}, h.ListSpots)
}

// RegisterRoutesWithAuth registers spot routes with authentication middleware
func (h *SpotHandler) RegisterRoutesWithAuth(api huma.API, authMiddleware *auth.AuthMiddleware) {
	// Register public routes first
	h.RegisterRoutes(api)

	// Create spot (protected - requires authentication)
	huma.Register(api, authMiddleware.CreateProtectedOperation(huma.Operation{
		OperationID: "create-spot",
		Method:      http.MethodPost,
		Path:        "/api/v1/spots",
		Summary:     "Create a new spot",
		Description: "Create a new reviewable spot on the map (requires authentication)",
		Tags:        []string{"Spots"},
	}), h.CreateSpot)

	// TODO: Update spot (protected - requires authentication and ownership or admin permission)
	// Will be implemented when UpdateSpotRequest is available in gRPC service
}

// CreateSpot creates a new spot
func (h *SpotHandler) CreateSpot(ctx context.Context, input *CreateSpotInput) (*CreateSpotOutput, error) {
	// Extract user ID from authentication context
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		return nil, huma.Error401Unauthorized("authentication required to create spot")
	}

	// Convert HTTP request to gRPC request
	grpcReq := &spotv1.CreateSpotRequest{
		Name:        input.Body.Name,
		NameI18N:    input.Body.NameI18n,
		Coordinates: &commonv1.Coordinates{
			Latitude:  input.Body.Latitude,
			Longitude: input.Body.Longitude,
		},
		Category:    input.Body.Category,
		Address:     input.Body.Address,
		AddressI18N: input.Body.AddressI18n,
		CountryCode: input.Body.CountryCode,
		// Note: Creator tracking will be handled by auth middleware
	}

	// Call gRPC service
	grpcResp, err := h.spotClient.CreateSpot(ctx, grpcReq)
	if err != nil {
		return nil, grpcToHTTPError(err, "failed to create spot")
	}

	// Convert gRPC response to HTTP response
	return &CreateSpotOutput{
		Body: grpcResp.Spot,
	}, nil
}

// GetSpot gets a specific spot
func (h *SpotHandler) GetSpot(ctx context.Context, input *GetSpotInput) (*GetSpotOutput, error) {
	// Convert HTTP request to gRPC request
	grpcReq := &spotv1.GetSpotRequest{
		Id: input.ID,
	}

	// Call gRPC service
	grpcResp, err := h.spotClient.GetSpot(ctx, grpcReq)
	if err != nil {
		return nil, err
	}

	// Convert gRPC response to HTTP response
	return &GetSpotOutput{
		Body: grpcResp.Spot,
	}, nil
}

// ListSpots lists spots
func (h *SpotHandler) ListSpots(ctx context.Context, input *ListSpotsInput) (*ListSpotsOutput, error) {
	// Convert HTTP request to gRPC request
	grpcReq := &spotv1.ListSpotsRequest{
		Pagination: &commonv1.PaginationRequest{
			Page:     int32(input.Page),
			PageSize: int32(input.PageSize),
		},
		Category:    input.Category,
		CountryCode: input.CountryCode,
	}

	// Add coordinates if provided (check for non-zero or explicit flag)
	if input.Latitude != 0 || input.Longitude != 0 {
		grpcReq.Center = &commonv1.Coordinates{
			Latitude:  input.Latitude,
			Longitude: input.Longitude,
		}
		grpcReq.RadiusKm = input.RadiusKm
	}

	// Call gRPC service
	grpcResp, err := h.spotClient.ListSpots(ctx, grpcReq)
	if err != nil {
		return nil, err
	}

	// Convert gRPC response to HTTP response
	return &ListSpotsOutput{
		Body: struct {
			Spots      []*spotv1.Spot              `json:"spots" doc:"List of spots"`
			Pagination *commonv1.PaginationResponse `json:"pagination" doc:"Pagination metadata"`
		}{
			Spots:      grpcResp.Spots,
			Pagination: grpcResp.Pagination,
		},
	}, nil
}

// TODO: UpdateSpot - will be implemented when UpdateSpotRequest is available in gRPC service