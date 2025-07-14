package entities

import (
	"time"
)

// Spot represents a reviewable location (MVP version)
type Spot struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	Address     string    `json:"address"`
	CountryCode string    `json:"country_code"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewSpot creates a new Spot instance (MVP version)
func NewSpot(name, category, address, countryCode string) *Spot {
	now := time.Now()
	return &Spot{
		Name:        name,
		Category:    category,
		Address:     address,
		CountryCode: countryCode,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

