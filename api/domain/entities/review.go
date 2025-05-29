package entities

import (
	"errors"
	"time"
)

// Review represents a user's rating of a spot
type Review struct {
	ID            string         `json:"id"`
	SpotID        string         `json:"spot_id"`
	UserID        string         `json:"user_id"` // Anonymous ID
	Rating        int            `json:"rating"`  // 1-5 stars
	Comment       string         `json:"comment,omitempty"`
	RatingAspects map[string]int `json:"rating_aspects,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

// NewReview creates a new Review instance
func NewReview(spotID, userID string, rating int) (*Review, error) {
	if err := validateRating(rating); err != nil {
		return nil, err
	}

	now := time.Now()
	return &Review{
		SpotID:        spotID,
		UserID:        userID,
		Rating:        rating,
		RatingAspects: make(map[string]int),
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

// UpdateRating updates the review rating
func (r *Review) UpdateRating(rating int) error {
	if err := validateRating(rating); err != nil {
		return err
	}
	r.Rating = rating
	r.UpdatedAt = time.Now()
	return nil
}

// SetComment sets the review comment
func (r *Review) SetComment(comment string) {
	r.Comment = comment
	r.UpdatedAt = time.Now()
}

// validateRating validates that rating is between 1 and 5
func validateRating(rating int) error {
	if rating < 1 || rating > 5 {
		return errors.New("rating must be between 1 and 5")
	}
	return nil
}