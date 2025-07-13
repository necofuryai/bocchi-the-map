package handlers

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"bocchi/api/application/clients"
	"bocchi/api/pkg/auth"
	spotv1 "bocchi/api/gen/spot/v1"
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


// GetSpotInput represents the request to get a spot
type GetSpotInput struct {
	ID string `path:"id" doc:"Spot ID"`
}

// GetSpotOutput represents the response for getting a spot (using protobuf Spot type)
type GetSpotOutput struct {
	Body *spotv1.Spot `json:"spot" doc:"Spot data"`
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

}

// RegisterRoutesWithAuth registers spot routes with authentication middleware
func (h *SpotHandler) RegisterRoutesWithAuth(api huma.API, authMiddleware *auth.AuthMiddleware) {
	// Register public routes first
	h.RegisterRoutes(api)


	// TODO: Update spot (protected - requires authentication and ownership or admin permission)
	// Will be implemented when UpdateSpotRequest is available in gRPC service
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


// TODO: UpdateSpot - will be implemented when UpdateSpotRequest is available in gRPC service