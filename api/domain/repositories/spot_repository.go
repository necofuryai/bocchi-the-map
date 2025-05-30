package repositories

import (
	"context"

	"github.com/necofuryai/bocchi-the-map/api/domain/entities"
)

// SpotRepository defines the interface for spot data operations
type SpotRepository interface {
	// Create creates a new spot
	Create(ctx context.Context, spot *entities.Spot) error

	// GetByID retrieves a spot by its ID
	GetByID(ctx context.Context, id string) (*entities.Spot, error)

	// GetByCoordinates retrieves spots within a radius from coordinates
	GetByCoordinates(ctx context.Context, lat, lng, radiusKm float64) ([]*entities.Spot, error)

	// List retrieves spots with pagination
	List(ctx context.Context, offset, limit int) ([]*entities.Spot, int, error)

	// Search searches spots by query
	Search(ctx context.Context, query string, lang string, offset, limit int) ([]*entities.Spot, int, error)

	// Update updates a spot
	Update(ctx context.Context, spot *entities.Spot) error

	// Delete deletes a spot
	Delete(ctx context.Context, id string) error

	// UpdateRating updates the average rating and review count
	UpdateRating(ctx context.Context, spotID string, averageRating float64, reviewCount int) error
}