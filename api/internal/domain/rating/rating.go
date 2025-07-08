package rating

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Rating represents a solo-friendly rating for a spot
type Rating struct {
	ID                 string    `json:"id"`
	SpotID             string    `json:"spot_id"`
	UserID             string    `json:"user_id"`
	SoloFriendlyRating int       `json:"solo_friendly_rating"`
	Categories         []string  `json:"categories"`
	Comment            string    `json:"comment"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// Valid categories for solo-friendly ratings
var ValidCategories = map[string]bool{
	"quiet_atmosphere":     true,
	"wifi_available":       true,
	"single_seating":       true,
	"good_lighting":        true,
	"power_outlets":        true,
	"comfortable_seating":  true,
	"minimal_noise":        true,
	"study_friendly":       true,
	"work_friendly":        true,
	"reading_friendly":     true,
}

// NewRating creates a new rating with validation
func NewRating(spotID, userID string, soloRating int, categories []string, comment string) (*Rating, error) {
	// Validate inputs
	if err := validateRatingInputs(spotID, userID, soloRating, categories); err != nil {
		return nil, err
	}

	// Deduplicate categories
	uniqueCategories := deduplicateCategories(categories)

	now := time.Now()
	
	return &Rating{
		ID:                 uuid.New().String(),
		SpotID:             spotID,
		UserID:             userID,
		SoloFriendlyRating: soloRating,
		Categories:         uniqueCategories,
		Comment:            comment,
		CreatedAt:          now,
		UpdatedAt:          now,
	}, nil
}

// UpdateRating updates an existing rating with validation
func (r *Rating) UpdateRating(newRating int, newCategories []string, newComment string) error {
	// Validate inputs (reuse validation logic, but spotID and userID are already set)
	if err := validateRatingInputs(r.SpotID, r.UserID, newRating, newCategories); err != nil {
		return err
	}

	// Deduplicate categories
	uniqueCategories := deduplicateCategories(newCategories)

	// Update fields
	r.SoloFriendlyRating = newRating
	r.Categories = uniqueCategories
	r.Comment = newComment
	r.UpdatedAt = time.Now()

	return nil
}

// validateRatingInputs validates common input parameters
func validateRatingInputs(spotID, userID string, soloRating int, categories []string) error {
	if spotID == "" {
		return errors.New("spot ID cannot be empty")
	}

	if userID == "" {
		return errors.New("user ID cannot be empty")
	}

	if soloRating < 1 || soloRating > 5 {
		return errors.New("solo-friendly rating must be between 1 and 5")
	}

	// Validate categories
	for _, category := range categories {
		if !ValidCategories[category] {
			return fmt.Errorf("invalid category: %s", category)
		}
	}

	return nil
}

// deduplicateCategories removes duplicate categories while preserving order
func deduplicateCategories(categories []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(categories))

	for _, category := range categories {
		if !seen[category] {
			seen[category] = true
			result = append(result, category)
		}
	}

	return result
}

// GetCategoriesAsString returns categories as a comma-separated string
func (r *Rating) GetCategoriesAsString() string {
	return strings.Join(r.Categories, ", ")
}

// IsRecentlyUpdated checks if the rating was updated within the last 24 hours
func (r *Rating) IsRecentlyUpdated() bool {
	return time.Since(r.UpdatedAt) < 24*time.Hour
}

// GetValidCategoriesList returns all valid categories as a slice
func GetValidCategoriesList() []string {
	categories := make([]string, 0, len(ValidCategories))
	for category := range ValidCategories {
		categories = append(categories, category)
	}
	return categories
}