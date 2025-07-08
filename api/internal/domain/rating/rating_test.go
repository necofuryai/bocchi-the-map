package rating_test

import (
	"testing"
	"time"

	"bocchi/api/internal/domain/rating"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TDD Unit Tests - Inner Loop

func TestRating_NewRating(t *testing.T) {
	tests := []struct {
		name           string
		spotID         string
		userID         string
		soloRating     int
		categories     []string
		comment        string
		wantErr        bool
		wantErrMessage string
	}{
		{
			name:       "valid rating creation",
			spotID:     "spot-123",
			userID:     "user-456",
			soloRating: 5,
			categories: []string{"quiet_atmosphere", "wifi_available"},
			comment:    "Great spot for solo work",
			wantErr:    false,
		},
		{
			name:           "empty spot ID",
			spotID:         "",
			userID:         "user-456",
			soloRating:     5,
			categories:     []string{"quiet_atmosphere"},
			comment:        "Good spot",
			wantErr:        true,
			wantErrMessage: "spot ID cannot be empty",
		},
		{
			name:           "empty user ID",
			spotID:         "spot-123",
			userID:         "",
			soloRating:     5,
			categories:     []string{"quiet_atmosphere"},
			comment:        "Good spot",
			wantErr:        true,
			wantErrMessage: "user ID cannot be empty",
		},
		{
			name:           "rating too low",
			spotID:         "spot-123",
			userID:         "user-456",
			soloRating:     0,
			categories:     []string{"quiet_atmosphere"},
			comment:        "Bad spot",
			wantErr:        true,
			wantErrMessage: "solo-friendly rating must be between 1 and 5",
		},
		{
			name:           "rating too high",
			spotID:         "spot-123",
			userID:         "user-456",
			soloRating:     6,
			categories:     []string{"quiet_atmosphere"},
			comment:        "Amazing spot",
			wantErr:        true,
			wantErrMessage: "solo-friendly rating must be between 1 and 5",
		},
		{
			name:           "invalid category",
			spotID:         "spot-123",
			userID:         "user-456",
			soloRating:     5,
			categories:     []string{"invalid_category"},
			comment:        "Good spot",
			wantErr:        true,
			wantErrMessage: "invalid category: invalid_category",
		},
		{
			name:       "empty comment is allowed",
			spotID:     "spot-123",
			userID:     "user-456",
			soloRating: 4,
			categories: []string{"quiet_atmosphere"},
			comment:    "",
			wantErr:    false,
		},
		{
			name:       "no categories is allowed",
			spotID:     "spot-123",
			userID:     "user-456",
			soloRating: 3,
			categories: []string{},
			comment:    "Decent spot",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result, err := rating.NewRating(
				tt.spotID,
				tt.userID,
				tt.soloRating,
				tt.categories,
				tt.comment,
			)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErrMessage)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.spotID, result.SpotID)
				assert.Equal(t, tt.userID, result.UserID)
				assert.Equal(t, tt.soloRating, result.SoloFriendlyRating)
				assert.Equal(t, tt.categories, result.Categories)
				assert.Equal(t, tt.comment, result.Comment)
				assert.False(t, result.CreatedAt.IsZero())
				assert.False(t, result.UpdatedAt.IsZero())
				assert.NotEmpty(t, result.ID)
			}
		})
	}
}

func TestRating_UpdateRating(t *testing.T) {
	// Arrange
	originalRating, err := rating.NewRating(
		"spot-123",
		"user-456",
		3,
		[]string{"quiet_atmosphere"},
		"Original comment",
	)
	require.NoError(t, err)
	
	originalCreatedAt := originalRating.CreatedAt
	originalUpdatedAt := originalRating.UpdatedAt
	
	// Sleep to ensure timestamp difference
	time.Sleep(1 * time.Millisecond)

	// Act
	err = originalRating.UpdateRating(5, []string{"quiet_atmosphere", "wifi_available"}, "Updated comment")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 5, originalRating.SoloFriendlyRating)
	assert.Equal(t, []string{"quiet_atmosphere", "wifi_available"}, originalRating.Categories)
	assert.Equal(t, "Updated comment", originalRating.Comment)
	assert.Equal(t, originalCreatedAt, originalRating.CreatedAt) // Should not change
	assert.True(t, originalRating.UpdatedAt.After(originalUpdatedAt)) // Should be updated
}

func TestRating_UpdateRating_InvalidData(t *testing.T) {
	// Arrange
	originalRating, err := rating.NewRating(
		"spot-123",
		"user-456",
		3,
		[]string{"quiet_atmosphere"},
		"Original comment",
	)
	require.NoError(t, err)

	tests := []struct {
		name           string
		newRating      int
		newCategories  []string
		newComment     string
		wantErr        bool
		wantErrMessage string
	}{
		{
			name:           "invalid rating too low",
			newRating:      0,
			newCategories:  []string{"quiet_atmosphere"},
			newComment:     "Updated comment",
			wantErr:        true,
			wantErrMessage: "solo-friendly rating must be between 1 and 5",
		},
		{
			name:           "invalid rating too high",
			newRating:      6,
			newCategories:  []string{"quiet_atmosphere"},
			newComment:     "Updated comment",
			wantErr:        true,
			wantErrMessage: "solo-friendly rating must be between 1 and 5",
		},
		{
			name:           "invalid category",
			newRating:      4,
			newCategories:  []string{"invalid_category"},
			newComment:     "Updated comment",
			wantErr:        true,
			wantErrMessage: "invalid category: invalid_category",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := originalRating.UpdateRating(tt.newRating, tt.newCategories, tt.newComment)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErrMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRating_ValidCategories(t *testing.T) {
	validCategories := []string{
		"quiet_atmosphere",
		"wifi_available",
		"single_seating",
		"good_lighting",
		"power_outlets",
		"comfortable_seating",
		"minimal_noise",
		"study_friendly",
		"work_friendly",
		"reading_friendly",
	}

	for _, category := range validCategories {
		t.Run("valid_category_"+category, func(t *testing.T) {
			// Act
			result, err := rating.NewRating(
				"spot-123",
				"user-456",
				5,
				[]string{category},
				"Test comment",
			)

			// Assert
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Contains(t, result.Categories, category)
		})
	}
}

func TestRating_DuplicateCategories(t *testing.T) {
	// Act
	result, err := rating.NewRating(
		"spot-123",
		"user-456",
		5,
		[]string{"quiet_atmosphere", "wifi_available", "quiet_atmosphere"}, // Duplicate
		"Test comment",
	)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	// Should deduplicate categories
	assert.Equal(t, []string{"quiet_atmosphere", "wifi_available"}, result.Categories)
}