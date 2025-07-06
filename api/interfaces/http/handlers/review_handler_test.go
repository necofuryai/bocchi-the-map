//go:build integration

package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"bocchi/api/application/clients"
	"bocchi/api/tests/helpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Test helper functions
func verifyPaginationInfo(paginationMap map[string]interface{}, expectedPage, expectedTotal, expectedTotalPages int) {
	Expect(paginationMap["page"]).To(Equal(float64(expectedPage)), fmt.Sprintf("Current page should be %d", expectedPage))
	Expect(paginationMap["total_count"]).To(Equal(float64(expectedTotal)), fmt.Sprintf("Total should be %d", expectedTotal))
	if expectedTotalPages > 0 {
		Expect(paginationMap["total_pages"]).To(Equal(float64(expectedTotalPages)), fmt.Sprintf("Should have %d pages total", expectedTotalPages))
	}
}

func verifyReviewStatistics(statsMap map[string]interface{}, expectedTotal int, expectedAvgMin, expectedAvgMax float64) {
	Expect(statsMap["total_reviews"]).To(Equal(float64(expectedTotal)), fmt.Sprintf("Total reviews count should be %d", expectedTotal))
	if expectedTotal > 0 {
		Expect(statsMap["average_rating"]).To(BeNumerically(">=", expectedAvgMin), "Average rating should be above minimum")
		Expect(statsMap["average_rating"]).To(BeNumerically("<=", expectedAvgMax), "Average rating should not exceed maximum")
	} else {
		Expect(statsMap["average_rating"]).To(Equal(float64(0)), "Average rating should be 0 for no reviews")
	}
}

func verifyRatingDistribution(statsMap map[string]interface{}, expectedRatingsPerLevel int) {
	distribution, exists := statsMap["rating_distribution"]
	Expect(exists).To(BeTrue(), "Rating distribution should be present")
	
	distMap, ok := distribution.(map[string]interface{})
	Expect(ok).To(BeTrue(), "Rating distribution should be an object")
	
	for rating := 1; rating <= 5; rating++ {
		key := fmt.Sprintf("%d", rating)
		Expect(distMap[key]).To(Equal(float64(expectedRatingsPerLevel)), fmt.Sprintf("Rating %d should have %d reviews", rating, expectedRatingsPerLevel))
	}
}

func verifyResponseBody(resp *httptest.ResponseRecorder) map[string]interface{} {
	var responseBody map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &responseBody)
	Expect(err).NotTo(HaveOccurred())
	return responseBody
}

func verifyReviewsArray(responseBody map[string]interface{}, expectedLength int) []interface{} {
	reviews, exists := responseBody["reviews"]
	Expect(exists).To(BeTrue(), "Reviews array should be present")
	
	reviewsArray, ok := reviews.([]interface{})
	Expect(ok).To(BeTrue(), "Reviews should be an array")
	
	if expectedLength >= 0 {
		Expect(len(reviewsArray)).To(Equal(expectedLength), fmt.Sprintf("Should have %d reviews", expectedLength))
	}
	
	return reviewsArray
}

func verifyPaginationExists(responseBody map[string]interface{}) map[string]interface{} {
	pagination, exists := responseBody["pagination"]
	Expect(exists).To(BeTrue(), "Pagination info should be present")
	
	paginationMap, ok := pagination.(map[string]interface{})
	Expect(ok).To(BeTrue(), "Pagination should be an object")
	
	return paginationMap
}

func verifyStatisticsExists(responseBody map[string]interface{}) map[string]interface{} {
	statistics, exists := responseBody["statistics"]
	Expect(exists).To(BeTrue(), "Statistics should be present")
	
	statsMap, ok := statistics.(map[string]interface{})
	Expect(ok).To(BeTrue(), "Statistics should be an object")
	
	return statsMap
}

var _ = Describe("ReviewHandler BDD Tests", func() {
	var (
		api           huma.API
		testServer    *httptest.Server
		reviewHandler *ReviewHandler
		reviewClient  *clients.ReviewClient
		authData      *helpers.AuthTestData
		spotFixture   *helpers.SpotFixture
		userFixture   *helpers.UserFixture
	)

	BeforeEach(func() {
		By("Setting up ReviewHandler test environment")
		
		// Create review client with test database
		var err error
		reviewClient, err = clients.NewReviewClient("internal", testSuite.TestDB.DB)
		Expect(err).NotTo(HaveOccurred())
		
		// Create review handler
		reviewHandler = NewReviewHandler(reviewClient)
		
		// Setup test API with chi router
		router := chi.NewRouter()
		api = humachi.New(router, huma.DefaultConfig("Test API", "1.0.0"))
		testServer = httptest.NewServer(router)
		
		// Register review endpoints
		reviewHandler.RegisterRoutes(api)
		
		// Setup authentication test data
		authData = testSuite.AuthHelper.NewAuthTestData()
		
		// Create test user in database
		userFixture = &helpers.UserFixture{
			ID:             authData.ValidUserID,
			Email:          authData.TestUser.Email,
			DisplayName:    authData.TestUser.DisplayName,
			AuthProvider:   string(authData.TestUser.AuthProvider),
			AuthProviderID: authData.TestUser.AuthProviderID,
			Preferences:    authData.TestUser.Preferences,
		}
		testSuite.FixtureManager.CreateUserFixture(context.Background(), *userFixture)
		
		// Create test spot for reviews
		spotFixture = &helpers.SpotFixture{
			ID:          "test-spot-review",
			Name:        "Test Cafe for Reviews",
			Latitude:    35.6762,
			Longitude:   139.6503,
			Category:    "cafe",
			Address:     "Test Address",
			CountryCode: "JP",
		}
		testSuite.FixtureManager.CreateSpotFixture(context.Background(), *spotFixture)
	})

	Describe("Review Creation", func() {
		Context("Given an authenticated user and existing spot", func() {
			Context("When creating a review with valid data", func() {
				It("Then the review should be created successfully", func() {
					By("Preparing a valid review creation request")
					requestBody := map[string]interface{}{
						"spot_id": spotFixture.ID,
						"rating":  4,
						"comment": "Great quiet spot for solo work. Perfect for studying and has excellent wifi.",
					}
					
					bodyBytes, err := json.Marshal(requestBody)
					Expect(err).NotTo(HaveOccurred())
					
					By("Sending the authenticated review creation request")
					req := httptest.NewRequest(http.MethodPost, "/api/v1/reviews", bytes.NewReader(bodyBytes))
					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("Authorization", "Bearer "+authData.ValidToken)
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the successful creation response")
					Expect(resp.Code).To(Equal(http.StatusCreated), "Expected status 201 Created")
					
					var responseBody map[string]interface{}
					err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
					Expect(err).NotTo(HaveOccurred())
					
					By("Verifying the created review data")
					Expect(responseBody["id"]).To(Not(BeEmpty()), "Review ID should be generated")
					Expect(responseBody["spot_id"]).To(Equal(spotFixture.ID))
					Expect(responseBody["user_id"]).To(Equal(authData.ValidUserID))
					Expect(responseBody["rating"]).To(Equal(float64(4)))
					Expect(responseBody["comment"]).To(Equal("Great quiet spot for solo work. Perfect for studying and has excellent wifi."))
					Expect(responseBody["created_at"]).To(Not(BeEmpty()), "Creation timestamp should be set")
				})
			})

			Context("When creating a review with aspect ratings", func() {
				It("Then the review should be created with aspect ratings stored", func() {
					By("Preparing a review with aspect ratings")
					requestBody := map[string]interface{}{
						"spot_id": spotFixture.ID,
						"rating":  5,
						"comment": "Perfect for solo travelers!",
						"rating_aspects": map[string]int{
							"quietness":      5,
							"wifi_quality":   4,
							"solo_friendly":  5,
							"accessibility":  3,
						},
					}
					
					bodyBytes, err := json.Marshal(requestBody)
					Expect(err).NotTo(HaveOccurred())
					
					By("Sending the review with aspect ratings")
					req := httptest.NewRequest(http.MethodPost, "/api/v1/reviews", bytes.NewReader(bodyBytes))
					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("Authorization", "Bearer "+authData.ValidToken)
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the successful creation with aspects")
					Expect(resp.Code).To(Equal(http.StatusCreated), "Expected status 201 Created")
					
					var responseBody map[string]interface{}
					err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
					Expect(err).NotTo(HaveOccurred())
					
					By("Verifying the aspect ratings are included")
					aspects, exists := responseBody["rating_aspects"]
					Expect(exists).To(BeTrue(), "Aspect ratings should be included")
					
					aspectMap, ok := aspects.(map[string]interface{})
					Expect(ok).To(BeTrue(), "Aspect ratings should be a map")
					Expect(aspectMap["quietness"]).To(Equal(float64(5)))
					Expect(aspectMap["wifi_quality"]).To(Equal(float64(4)))
					Expect(aspectMap["solo_friendly"]).To(Equal(float64(5)))
					Expect(aspectMap["accessibility"]).To(Equal(float64(3)))
				})
			})
		})

		Context("Given a non-existent spot", func() {
			Context("When creating a review for non-existent spot", func() {
				It("Then it should return a not found error", func() {
					By("Preparing a review for non-existent spot")
					requestBody := map[string]interface{}{
						"spot_id": "non-existent-spot-id",
						"rating":  4,
						"comment": "This spot doesn't exist",
					}
					
					bodyBytes, err := json.Marshal(requestBody)
					Expect(err).NotTo(HaveOccurred())
					
					By("Sending the review for non-existent spot")
					req := httptest.NewRequest(http.MethodPost, "/api/v1/reviews", bytes.NewReader(bodyBytes))
					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("Authorization", "Bearer "+authData.ValidToken)
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the not found error")
					Expect(resp.Code).To(Equal(http.StatusNotFound), "Expected status 404 Not Found")
					
					var errorResponse map[string]interface{}
					err = json.Unmarshal(resp.Body.Bytes(), &errorResponse)
					Expect(err).NotTo(HaveOccurred())
					
					Expect(errorResponse["title"]).To(ContainSubstring("not found"), "Should indicate spot not found")
				})
			})
		})

		Context("Given invalid rating value", func() {
			Context("When creating review with rating outside 1-5 range", func() {
				It("Then it should return a validation error", func() {
					By("Preparing a review with invalid rating")
					requestBody := map[string]interface{}{
						"spot_id": spotFixture.ID,
						"rating":  6, // Invalid: outside 1-5 range
						"comment": "Rating too high",
					}
					
					bodyBytes, err := json.Marshal(requestBody)
					Expect(err).NotTo(HaveOccurred())
					
					By("Sending the review with invalid rating")
					req := httptest.NewRequest(http.MethodPost, "/api/v1/reviews", bytes.NewReader(bodyBytes))
					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("Authorization", "Bearer "+authData.ValidToken)
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the validation error")
					Expect(resp.Code).To(Equal(http.StatusBadRequest), "Expected status 400 Bad Request")
					
					var errorResponse map[string]interface{}
					err = json.Unmarshal(resp.Body.Bytes(), &errorResponse)
					Expect(err).NotTo(HaveOccurred())
					
					Expect(errorResponse["title"]).To(ContainSubstring("validation"), "Should specify rating validation failure")
				})
			})
		})

		Context("Given an unauthenticated request", func() {
			Context("When attempting to create a review", func() {
				It("Then it should return an authentication error", func() {
					By("Preparing a review creation request")
					requestBody := map[string]interface{}{
						"spot_id": spotFixture.ID,
						"rating":  4,
						"comment": "Should fail without auth",
					}
					
					bodyBytes, err := json.Marshal(requestBody)
					Expect(err).NotTo(HaveOccurred())
					
					By("Sending unauthenticated request")
					req := httptest.NewRequest(http.MethodPost, "/api/v1/reviews", bytes.NewReader(bodyBytes))
					req.Header.Set("Content-Type", "application/json")
					// No Authorization header
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the authentication error")
					Expect(resp.Code).To(Equal(http.StatusUnauthorized), "Expected status 401 Unauthorized")
				})
			})
		})

		Context("Given a duplicate review attempt", func() {
			BeforeEach(func() {
				By("Creating an existing review first")
				testSuite.FixtureManager.CreateReviewFixture(context.Background(), helpers.ReviewFixture{
					ID:     "existing-review-123",
					SpotID: spotFixture.ID,
					UserID: authData.ValidUserID,
					Rating: 3,
					Comment: "Already reviewed this spot",
				})
			})
			
			Context("When user tries to review same spot twice", func() {
				It("Then it should return a conflict error", func() {
					By("Preparing a duplicate review attempt")
					requestBody := map[string]interface{}{
						"spot_id": spotFixture.ID,
						"rating":  5,
						"comment": "Trying to review again",
					}
					
					bodyBytes, err := json.Marshal(requestBody)
					Expect(err).NotTo(HaveOccurred())
					
					By("Sending the duplicate review request")
					req := httptest.NewRequest(http.MethodPost, "/api/v1/reviews", bytes.NewReader(bodyBytes))
					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("Authorization", "Bearer "+authData.ValidToken)
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the conflict error")
					Expect(resp.Code).To(Equal(http.StatusConflict), "Expected status 409 Conflict")
					
					var errorResponse map[string]interface{}
					err = json.Unmarshal(resp.Body.Bytes(), &errorResponse)
					Expect(err).NotTo(HaveOccurred())
					
					Expect(errorResponse["title"]).To(ContainSubstring("already reviewed"), "Should indicate user already reviewed this spot")
				})
			})
		})
	})

	Describe("Spot Review Retrieval", func() {
		const (
			totalTestReviews = 15
			pageLimit        = 10
			remainingReviews = 5
			minRating        = 1
			maxRating        = 5
			expectedRatingsPerLevel = 3
		)

		BeforeEach(func() {
			By("Creating multiple reviews for the test spot")
			// Create additional users for diverse reviews
			for i := 1; i <= totalTestReviews; i++ {
				userID := fmt.Sprintf("review-user-%d", i)
				testSuite.FixtureManager.CreateUserFixture(context.Background(), helpers.UserFixture{
					ID:             userID,
					Email:          fmt.Sprintf("reviewer%d@example.com", i),
					DisplayName:    fmt.Sprintf("Reviewer %d", i),
					AuthProvider:   "google",
					AuthProviderID: fmt.Sprintf("google_%d", i),
				})
				
				// Create reviews with varying ratings
				rating := (i % maxRating) + minRating // Ratings from 1-5
				testSuite.FixtureManager.CreateReviewFixture(context.Background(), helpers.ReviewFixture{
					ID:      fmt.Sprintf("review-%d", i),
					SpotID:  spotFixture.ID,
					UserID:  userID,
					Rating:  rating,
					Comment: fmt.Sprintf("Review number %d with rating %d", i, rating),
				})
			}
		})

		Context("Given existing reviews for a spot", func() {
			Context("When requesting spot reviews with pagination", func() {
				It("Then paginated reviews and statistics should be returned", func() {
					By("Sending request for first page of reviews")
					req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/spots/%s/reviews?page=1&limit=%d", spotFixture.ID, pageLimit), nil)
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the successful response")
					Expect(resp.Code).To(Equal(http.StatusOK), "Expected status 200 OK")
					
					responseBody := verifyResponseBody(resp)
					
					By("Verifying the paginated reviews")
					verifyReviewsArray(responseBody, pageLimit)
					
					By("Verifying pagination information")
					paginationMap := verifyPaginationExists(responseBody)
					verifyPaginationInfo(paginationMap, 1, totalTestReviews, 2)
					
					By("Verifying review statistics")
					statsMap := verifyStatisticsExists(responseBody)
					verifyReviewStatistics(statsMap, totalTestReviews, 0.1, maxRating)
					
					By("Verifying rating distribution")
					verifyRatingDistribution(statsMap, expectedRatingsPerLevel)
				})
			})

			Context("When requesting second page of reviews", func() {
				It("Then correct page results should be returned", func() {
					By("Sending request for second page")
					req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/spots/%s/reviews?page=2&limit=%d", spotFixture.ID, pageLimit), nil)
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the successful response")
					Expect(resp.Code).To(Equal(http.StatusOK), "Expected status 200 OK")
					
					responseBody := verifyResponseBody(resp)
					
					By("Verifying second page results")
					verifyReviewsArray(responseBody, remainingReviews)
					
					By("Verifying pagination for second page")
					paginationMap := verifyPaginationExists(responseBody)
					verifyPaginationInfo(paginationMap, 2, totalTestReviews, -1) // -1 to skip total_pages check
				})
			})
		})

		Context("Given invalid pagination parameters", func() {
			Context("When requesting reviews with invalid page number", func() {
				It("Then it should return a validation error", func() {
					By("Sending request with invalid page")
					req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/spots/%s/reviews?page=0&limit=10", spotFixture.ID), nil)
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the validation error")
					Expect(resp.Code).To(Equal(http.StatusBadRequest), "Expected status 400 Bad Request")
				})
			})

			Context("When requesting reviews with excessive limit", func() {
				It("Then it should return a validation error", func() {
					By("Sending request with excessive limit")
					req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/spots/%s/reviews?page=1&limit=100", spotFixture.ID), nil)
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the validation error")
					Expect(resp.Code).To(Equal(http.StatusBadRequest), "Expected status 400 Bad Request")
				})
			})
		})

		Context("Given a non-existent spot", func() {
			Context("When requesting reviews for non-existent spot", func() {
				It("Then it should return a not found error", func() {
					By("Sending request for non-existent spot")
					req := httptest.NewRequest(http.MethodGet, "/api/v1/spots/non-existent-spot/reviews", nil)
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the not found error")
					Expect(resp.Code).To(Equal(http.StatusNotFound), "Expected status 404 Not Found")
				})
			})
		})

		Context("Given a spot with no reviews", func() {
			BeforeEach(func() {
				By("Creating a spot without reviews")
				testSuite.FixtureManager.CreateSpotFixture(context.Background(), helpers.SpotFixture{
					ID:          "empty-spot",
					Name:        "Spot Without Reviews",
					Latitude:    35.6762,
					Longitude:   139.6503,
					Category:    "cafe",
					Address:     "Empty Address",
					CountryCode: "JP",
				})
			})
			
			Context("When requesting reviews for empty spot", func() {
				It("Then empty list with statistics should be returned", func() {
					By("Sending request for spot without reviews")
					req := httptest.NewRequest(http.MethodGet, "/api/v1/spots/empty-spot/reviews", nil)
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the successful empty response")
					Expect(resp.Code).To(Equal(http.StatusOK), "Expected status 200 OK")
					
					responseBody := verifyResponseBody(resp)
					
					By("Verifying empty reviews list")
					verifyReviewsArray(responseBody, 0)
					
					By("Verifying zero statistics")
					statsMap := verifyStatisticsExists(responseBody)
					verifyReviewStatistics(statsMap, 0, 0, 0)
				})
			})
		})
	})

	Describe("User Review History", func() {
		var reviewerUserID string

		BeforeEach(func() {
			By("Creating a user with multiple reviews")
			reviewerUserID = "multi-reviewer-user"
			testSuite.FixtureManager.CreateUserFixture(context.Background(), helpers.UserFixture{
				ID:             reviewerUserID,
				Email:          "multireviewer@example.com",
				DisplayName:    "Multi Reviewer",
				AuthProvider:   "google",
				AuthProviderID: "google_multi_reviewer",
			})
			
			// Create multiple spots and reviews for this user
			for i := 1; i <= 8; i++ {
				spotID := fmt.Sprintf("user-spot-%d", i)
				testSuite.FixtureManager.CreateSpotFixture(context.Background(), helpers.SpotFixture{
					ID:          spotID,
					Name:        fmt.Sprintf("User Spot %d", i),
					Latitude:    35.6762 + float64(i)*0.01,
					Longitude:   139.6503 + float64(i)*0.01,
					Category:    "cafe",
					Address:     fmt.Sprintf("Address %d", i),
					CountryCode: "JP",
				})
				
				testSuite.FixtureManager.CreateReviewFixture(context.Background(), helpers.ReviewFixture{
					ID:      fmt.Sprintf("user-review-%d", i),
					SpotID:  spotID,
					UserID:  reviewerUserID,
					Rating:  (i % 5) + 1,
					Comment: fmt.Sprintf("User review %d", i),
				})
			}
		})

		Context("Given a user with existing reviews", func() {
			Context("When requesting user's reviews", func() {
				It("Then paginated reviews should be returned", func() {
					By("Sending request for user's reviews")
					req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/users/%s/reviews", reviewerUserID), nil)
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the successful response")
					Expect(resp.Code).To(Equal(http.StatusOK), "Expected status 200 OK")
					
					responseBody := verifyResponseBody(resp)
					
					By("Verifying user's reviews")
					reviewsArray := verifyReviewsArray(responseBody, -1) // -1 to skip exact length check
					Expect(len(reviewsArray)).To(BeNumerically(">=", 1), "Should contain user's reviews")
					Expect(len(reviewsArray)).To(BeNumerically("<=", 8), "Should not exceed total reviews")
					
					By("Verifying each review includes spot information")
					if len(reviewsArray) > 0 {
						firstReview, ok := reviewsArray[0].(map[string]interface{})
						Expect(ok).To(BeTrue(), "Review should be an object")
						Expect(firstReview["user_id"]).To(Equal(reviewerUserID), "Review should belong to the user")
						Expect(firstReview["spot"]).To(Not(BeNil()), "Spot information should be included")
						
						spot, ok := firstReview["spot"].(map[string]interface{})
						Expect(ok).To(BeTrue(), "Spot should be an object")
						Expect(spot["name"]).To(Not(BeEmpty()), "Spot name should be included")
					}
					
					By("Verifying pagination information")
					paginationMap := verifyPaginationExists(responseBody)
					verifyPaginationInfo(paginationMap, 1, 8, -1) // -1 to skip total_pages check
				})
			})
		})

		Context("Given a non-existent user", func() {
			Context("When requesting reviews for non-existent user", func() {
				It("Then it should return a not found error", func() {
					By("Sending request for non-existent user")
					req := httptest.NewRequest(http.MethodGet, "/api/v1/users/non-existent-user/reviews", nil)
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the not found error")
					Expect(resp.Code).To(Equal(http.StatusNotFound), "Expected status 404 Not Found")
				})
			})
		})

		Context("Given a user with no reviews", func() {
			BeforeEach(func() {
				By("Creating a user without reviews")
				testSuite.FixtureManager.CreateUserFixture(context.Background(), helpers.UserFixture{
					ID:             "empty-reviewer",
					Email:          "empty@example.com",
					DisplayName:    "Empty Reviewer",
					AuthProvider:   "google",
					AuthProviderID: "google_empty",
				})
			})
			
			Context("When requesting reviews for user without reviews", func() {
				It("Then empty list should be returned", func() {
					By("Sending request for user without reviews")
					req := httptest.NewRequest(http.MethodGet, "/api/v1/users/empty-reviewer/reviews", nil)
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the successful empty response")
					Expect(resp.Code).To(Equal(http.StatusOK), "Expected status 200 OK")
					
					responseBody := verifyResponseBody(resp)
					
					By("Verifying empty reviews list")
					verifyReviewsArray(responseBody, 0)
					
					By("Verifying zero pagination count")
					paginationMap := verifyPaginationExists(responseBody)
					verifyPaginationInfo(paginationMap, 1, 0, -1) // -1 to skip total_pages check
				})
			})
		})
	})
})