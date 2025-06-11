package entities

import (
	"time"
)

// Spot represents a reviewable location
type Spot struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	NameI18n     map[string]string `json:"name_i18n"`
	Latitude     float64           `json:"latitude"`
	Longitude    float64           `json:"longitude"`
	Category     string            `json:"category"`
	Address      string            `json:"address"`
	AddressI18n  map[string]string `json:"address_i18n"`
	CountryCode  string            `json:"country_code"`
	AverageRating float64          `json:"average_rating"`
	ReviewCount  int               `json:"review_count"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

// NewSpot creates a new Spot instance
func NewSpot(name string, lat, lng float64, category, address, countryCode string) *Spot {
	now := time.Now()
	return &Spot{
		Name:         name,
		NameI18n:     make(map[string]string),
		Latitude:     lat,
		Longitude:    lng,
		Category:     category,
		Address:      address,
		AddressI18n:  make(map[string]string),
		CountryCode:  countryCode,
		AverageRating: 0,
		ReviewCount:  0,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// UpdateRating updates the average rating and review count
func (s *Spot) UpdateRating(newRating float64, totalReviews int) {
	// Validate rating range
	if newRating < 0 || newRating > 5 {
		return // or return appropriate error
	}
	if totalReviews < 0 {
		return
	}
	s.AverageRating = newRating
	s.ReviewCount = totalReviews
	s.UpdatedAt = time.Now()
}