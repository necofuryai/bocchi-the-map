//go:build e2e

package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/necofuryai/bocchi-the-map/api/application/clients"
	"github.com/necofuryai/bocchi-the-map/api/interfaces/http/handlers"
	"github.com/necofuryai/bocchi-the-map/api/pkg/auth"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("End-to-End Integration Scenarios", func() {
	var (
		api           huma.API
		testServer    *httptest.Server
		userClient    *clients.UserClient
		spotClient    *clients.SpotClient
		reviewClient  *clients.ReviewClient
		authMiddleware *auth.AuthMiddleware
		rateLimiter   *auth.RateLimiter
	)

	BeforeEach(func() {
		By("Setting up complete API environment for integration testing")
		
		// Create all clients
		var err error
		userClient, err = clients.NewUserClient("internal", testSuite.TestDB.DB)
		Expect(err).NotTo(HaveOccurred())
		
		spotClient, err = clients.NewSpotClient("internal", testSuite.TestDB.DB)
		Expect(err).NotTo(HaveOccurred())
		
		reviewClient, err = clients.NewReviewClient("internal", testSuite.TestDB.DB)
		Expect(err).NotTo(HaveOccurred())
		
		// Create middleware and rate limiter
		authMiddleware = auth.NewAuthMiddleware("test-secret-key", testSuite.TestDB.Queries)
		rateLimiter = auth.NewRateLimiter(10, 300) // More lenient for integration tests
		
		// Setup complete API with all handlers
		router := chi.NewRouter()
		api = humachi.New(router, huma.DefaultConfig("Integration Test API", "1.0.0"))
		testServer = httptest.NewServer(router)
		
		// Register all handlers (simulating production setup)
		authHandler := handlers.NewAuthHandler(authMiddleware, userClient)
		userHandler := handlers.NewUserHandler(userClient)
		spotHandler := handlers.NewSpotHandler(spotClient)
		reviewHandler := handlers.NewReviewHandler(reviewClient)
		
		// Register all routes
		authHandler.RegisterRoutesWithRateLimit(api, rateLimiter)
		userHandler.RegisterRoutesWithAuth(api, authMiddleware)
		spotHandler.RegisterRoutes(api)
		reviewHandler.RegisterRoutes(api)
	})

	AfterEach(func() {
		By("Cleaning up test resources")
		
		// Close test server to prevent resource leaks
		if testServer != nil {
			testServer.Close()
			testServer = nil
		}
		
		// Clean up clients if they have cleanup methods
		// Note: These clients don't currently have explicit cleanup methods,
		// but they should close their database connections properly
		userClient = nil
		spotClient = nil
		reviewClient = nil
		authMiddleware = nil
		rateLimiter = nil
	})

	Describe("Complete Solo Traveler Journey", func() {
		var (
			userAccessToken string
			createdUserID   string
			createdSpotID   string
			createdReviewID string
		)

		Context("Given a new solo traveler starting their journey", func() {
			It("Should complete the full user lifecycle successfully", func() {
				By("Step 1: User registers via OAuth (Google authentication)")
				registerRequestBody := map[string]interface{}{
					"email":            "solo.traveler@example.com",
					"display_name":     "Solo Traveler",
					"auth_provider":    "google",
					"auth_provider_id": "google_solo_traveler_123",
					"avatar_url":       "https://example.com/avatar.jpg",
				}
				
				bodyBytes, err := json.Marshal(registerRequestBody)
				Expect(err).NotTo(HaveOccurred())
				
				req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(bodyBytes))
				req.Header.Set("Content-Type", "application/json")
				
				resp := httptest.NewRecorder()
				testServer.Config.Handler.ServeHTTP(resp, req)
				
				Expect(resp.Code).To(Equal(http.StatusCreated), "User registration should succeed")
				
				var userResponse map[string]interface{}
				err = json.Unmarshal(resp.Body.Bytes(), &userResponse)
				Expect(err).NotTo(HaveOccurred())
				
				createdUserID = userResponse["id"].(string)
				Expect(createdUserID).To(Not(BeEmpty()), "User ID should be generated")
				
				By("Step 2: User generates API token for app access")
				tokenRequestBody := map[string]interface{}{
					"email":            "solo.traveler@example.com",
					"auth_provider":    "google",
					"auth_provider_id": "google_solo_traveler_123",
				}
				
				bodyBytes, err = json.Marshal(tokenRequestBody)
				Expect(err).NotTo(HaveOccurred())
				
				req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/token", bytes.NewReader(bodyBytes))
				req.Header.Set("Content-Type", "application/json")
				
				resp = httptest.NewRecorder()
				testServer.Config.Handler.ServeHTTP(resp, req)
				
				Expect(resp.Code).To(Equal(http.StatusOK), "Token generation should succeed")
				
				// Extract access token from cookies
				cookies := resp.Result().Cookies()
				for _, cookie := range cookies {
					if cookie.Name == "bocchi_access_token" {
						userAccessToken = cookie.Value
						break
					}
				}
				Expect(userAccessToken).To(Not(BeEmpty()), "Access token should be generated")
				
				By("Step 3: User retrieves their profile to verify login")
				req = httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
				req.Header.Set("Authorization", "Bearer "+userAccessToken)
				
				resp = httptest.NewRecorder()
				testServer.Config.Handler.ServeHTTP(resp, req)
				
				Expect(resp.Code).To(Equal(http.StatusOK), "Profile retrieval should succeed")
				
				var profileResponse map[string]interface{}
				err = json.Unmarshal(resp.Body.Bytes(), &profileResponse)
				Expect(err).NotTo(HaveOccurred())
				
				Expect(profileResponse["id"]).To(Equal(createdUserID), "Profile should match created user")
				Expect(profileResponse["email"]).To(Equal("solo.traveler@example.com"))
				
				By("Step 4: User updates their preferences for better experience")
				preferencesRequestBody := map[string]interface{}{
					"language":  "en",
					"dark_mode": true,
					"timezone":  "UTC",
				}
				
				bodyBytes, err = json.Marshal(preferencesRequestBody)
				Expect(err).NotTo(HaveOccurred())
				
				req = httptest.NewRequest(http.MethodPatch, "/api/v1/users/me/preferences", bytes.NewReader(bodyBytes))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+userAccessToken)
				
				resp = httptest.NewRecorder()
				testServer.Config.Handler.ServeHTTP(resp, req)
				
				Expect(resp.Code).To(Equal(http.StatusOK), "Preferences update should succeed")
				
				var prefsResponse map[string]interface{}
				err = json.Unmarshal(resp.Body.Bytes(), &prefsResponse)
				Expect(err).NotTo(HaveOccurred())
				
				preferences := prefsResponse["preferences"].(map[string]interface{})
				Expect(preferences["language"]).To(Equal("en"))
				Expect(preferences["dark_mode"]).To(Equal(true))
				
				By("Step 5: User creates their first solo-friendly spot")
				spotRequestBody := map[string]interface{}{
					"name":         "Solo Traveler's Cafe",
					"name_i18n":    map[string]string{"ja": "ソロ旅行者のカフェ"},
					"latitude":     35.6762,
					"longitude":    139.6503,
					"category":     "cafe",
					"address":      "1-2-3 Solo District, Tokyo",
					"address_i18n": map[string]string{"ja": "東京都ソロ区1-2-3"},
					"country_code": "JP",
				}
				
				bodyBytes, err = json.Marshal(spotRequestBody)
				Expect(err).NotTo(HaveOccurred())
				
				req = httptest.NewRequest(http.MethodPost, "/api/v1/spots", bytes.NewReader(bodyBytes))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+userAccessToken)
				
				resp = httptest.NewRecorder()
				testServer.Config.Handler.ServeHTTP(resp, req)
				
				Expect(resp.Code).To(Equal(http.StatusCreated), "Spot creation should succeed")
				
				var spotResponse map[string]interface{}
				err = json.Unmarshal(resp.Body.Bytes(), &spotResponse)
				Expect(err).NotTo(HaveOccurred())
				
				createdSpotID = spotResponse["id"].(string)
				Expect(createdSpotID).To(Not(BeEmpty()), "Spot ID should be generated")
				Expect(spotResponse["name"]).To(Equal("Solo Traveler's Cafe"))
				
				By("Step 6: User writes a detailed review for their spot")
				reviewRequestBody := map[string]interface{}{
					"spot_id": createdSpotID,
					"rating":  5,
					"comment": "Perfect spot for solo travelers! Quiet atmosphere, great wifi, and solo-friendly seating. The staff is understanding of people working alone.",
					"rating_aspects": map[string]int{
						"quietness":      5,
						"wifi_quality":   5,
						"solo_friendly":  5,
						"accessibility":  4,
					},
				}
				
				bodyBytes, err = json.Marshal(reviewRequestBody)
				Expect(err).NotTo(HaveOccurred())
				
				req = httptest.NewRequest(http.MethodPost, "/api/v1/reviews", bytes.NewReader(bodyBytes))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+userAccessToken)
				
				resp = httptest.NewRecorder()
				testServer.Config.Handler.ServeHTTP(resp, req)
				
				Expect(resp.Code).To(Equal(http.StatusCreated), "Review creation should succeed")
				
				var reviewResponse map[string]interface{}
				err = json.Unmarshal(resp.Body.Bytes(), &reviewResponse)
				Expect(err).NotTo(HaveOccurred())
				
				createdReviewID = reviewResponse["id"].(string)
				Expect(createdReviewID).To(Not(BeEmpty()), "Review ID should be generated")
				Expect(reviewResponse["rating"]).To(Equal(float64(5)))
				Expect(reviewResponse["user_id"]).To(Equal(createdUserID))
				Expect(reviewResponse["spot_id"]).To(Equal(createdSpotID))
				
				By("Step 7: User retrieves the spot to see their review reflected")
				req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/spots/%s", createdSpotID), nil)
				
				resp = httptest.NewRecorder()
				testServer.Config.Handler.ServeHTTP(resp, req)
				
				Expect(resp.Code).To(Equal(http.StatusOK), "Spot retrieval should succeed")
				
				var updatedSpotResponse map[string]interface{}
				err = json.Unmarshal(resp.Body.Bytes(), &updatedSpotResponse)
				Expect(err).NotTo(HaveOccurred())
				
				Expect(updatedSpotResponse["average_rating"]).To(Equal(float64(5)), "Spot should reflect the review rating")
				Expect(updatedSpotResponse["review_count"]).To(Equal(float64(1)), "Spot should show 1 review")
				
				By("Step 8: User checks their review history")
				req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/users/%s/reviews", createdUserID), nil)
				
				resp = httptest.NewRecorder()
				testServer.Config.Handler.ServeHTTP(resp, req)
				
				Expect(resp.Code).To(Equal(http.StatusOK), "User reviews retrieval should succeed")
				
				var userReviewsResponse map[string]interface{}
				err = json.Unmarshal(resp.Body.Bytes(), &userReviewsResponse)
				Expect(err).NotTo(HaveOccurred())
				
				reviews := userReviewsResponse["reviews"].([]interface{})
				Expect(len(reviews)).To(Equal(1), "User should have 1 review")
				
				userReview := reviews[0].(map[string]interface{})
				Expect(userReview["id"]).To(Equal(createdReviewID))
				Expect(userReview["spot"]).To(Not(BeNil()), "Review should include spot information")
				
				By("Step 9: User logs out to end session")
				req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
				req.AddCookie(&http.Cookie{
					Name:  "bocchi_access_token",
					Value: userAccessToken,
				})
				
				resp = httptest.NewRecorder()
				testServer.Config.Handler.ServeHTTP(resp, req)
				
				Expect(resp.Code).To(Equal(http.StatusOK), "Logout should succeed")
				
				// Verify cookies are cleared
				cookies = resp.Result().Cookies()
				for _, cookie := range cookies {
					if cookie.Name == "bocchi_access_token" {
						Expect(cookie.Value).To(BeEmpty(), "Access token cookie should be cleared")
						Expect(cookie.MaxAge).To(Equal(-1), "Cookie should be expired")
					}
				}
				
				By("Journey completed successfully! User has registered, authenticated, created content, and logged out.")
			})
		})
	})

	Describe("Multi-User Content Discovery Scenario", func() {
		var (
			user1Token, user2Token string
			user1ID, user2ID       string
			communitySpotID        string
		)

		Context("Given multiple users creating and discovering content", func() {
			It("Should demonstrate community-driven content discovery", func() {
				By("Setting up User 1: Content Creator")
				// Create User 1
				user1RequestBody := map[string]interface{}{
					"email":            "creator@example.com",
					"display_name":     "Content Creator",
					"auth_provider":    "google",
					"auth_provider_id": "google_creator_123",
				}
				
				bodyBytes, err := json.Marshal(user1RequestBody)
				Expect(err).NotTo(HaveOccurred())
				
				req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(bodyBytes))
				req.Header.Set("Content-Type", "application/json")
				
				resp := httptest.NewRecorder()
				testServer.Config.Handler.ServeHTTP(resp, req)
				Expect(resp.Code).To(Equal(http.StatusCreated))
				
				var user1Response map[string]interface{}
				err = json.Unmarshal(resp.Body.Bytes(), &user1Response)
				Expect(err).NotTo(HaveOccurred())
				user1ID = user1Response["id"].(string)
				
				// Generate token for User 1
				tokenRequestBody := map[string]interface{}{
					"email":            "creator@example.com",
					"auth_provider":    "google",
					"auth_provider_id": "google_creator_123",
				}
				
				bodyBytes, _ = json.Marshal(tokenRequestBody)
				req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/token", bytes.NewReader(bodyBytes))
				req.Header.Set("Content-Type", "application/json")
				resp = httptest.NewRecorder()
				testServer.Config.Handler.ServeHTTP(resp, req)
				
				user1Token = extractBocchiAccessToken(resp)
				
				By("Setting up User 2: Content Discoverer")
				// Create User 2
				user2RequestBody := map[string]interface{}{
					"email":            "discoverer@example.com",
					"display_name":     "Content Discoverer",
					"auth_provider":    "google",
					"auth_provider_id": "google_discoverer_456",
				}
				
				bodyBytes, _ = json.Marshal(user2RequestBody)
				req = httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(bodyBytes))
				req.Header.Set("Content-Type", "application/json")
				resp = httptest.NewRecorder()
				testServer.Config.Handler.ServeHTTP(resp, req)
				
				var user2Response map[string]interface{}
				err = json.Unmarshal(resp.Body.Bytes(), &user2Response)
				Expect(err).NotTo(HaveOccurred())
				user2ID = user2Response["id"].(string)
				Expect(user2ID).To(Not(BeEmpty()), "User 2 ID should be generated")
				
				// Generate token for User 2
				tokenRequestBody = map[string]interface{}{
					"email":            "discoverer@example.com",
					"auth_provider":    "google",
					"auth_provider_id": "google_discoverer_456",
				}
				
				bodyBytes, _ = json.Marshal(tokenRequestBody)
				req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/token", bytes.NewReader(bodyBytes))
				req.Header.Set("Content-Type", "application/json")
				resp = httptest.NewRecorder()
				testServer.Config.Handler.ServeHTTP(resp, req)
				
				user2Token = extractBocchiAccessToken(resp)
				
				By("User 1 creates a community spot")
				spotRequestBody := map[string]interface{}{
					"name":         "Community Solo Library",
					"latitude":     35.6895,
					"longitude":    139.6917,
					"category":     "library",
					"address":      "Community District, Tokyo",
					"country_code": "JP",
				}
				
				bodyBytes, _ = json.Marshal(spotRequestBody)
				req = httptest.NewRequest(http.MethodPost, "/api/v1/spots", bytes.NewReader(bodyBytes))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+user1Token)
				resp = httptest.NewRecorder()
				testServer.Config.Handler.ServeHTTP(resp, req)
				
				var spotResponse map[string]interface{}
				err = json.Unmarshal(resp.Body.Bytes(), &spotResponse)
				Expect(err).NotTo(HaveOccurred())
				communitySpotID = spotResponse["id"].(string)
				
				By("User 1 writes the first review")
				reviewRequestBody := map[string]interface{}{
					"spot_id": communitySpotID,
					"rating":  4,
					"comment": "Excellent for solo studying. Quiet and spacious.",
				}
				
				bodyBytes, _ = json.Marshal(reviewRequestBody)
				req = httptest.NewRequest(http.MethodPost, "/api/v1/reviews", bytes.NewReader(bodyBytes))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+user1Token)
				resp = httptest.NewRecorder()
				testServer.Config.Handler.ServeHTTP(resp, req)
				Expect(resp.Code).To(Equal(http.StatusCreated))
				
				By("User 2 discovers the spot through search")
				req = httptest.NewRequest(http.MethodGet, "/api/v1/spots?page=1&page_size=10", nil)
				resp = httptest.NewRecorder()
				testServer.Config.Handler.ServeHTTP(resp, req)
				
				var spotsResponse map[string]interface{}
				err = json.Unmarshal(resp.Body.Bytes(), &spotsResponse)
				Expect(err).NotTo(HaveOccurred())
				
				spots := spotsResponse["spots"].([]interface{})
				Expect(len(spots)).To(BeNumerically(">=", 1), "Should find at least one spot")
				
				// Verify community spot is in results
				foundCommunitySpot := false
				for _, spotInterface := range spots {
					spot := spotInterface.(map[string]interface{})
					if spot["id"] == communitySpotID {
						foundCommunitySpot = true
						Expect(spot["average_rating"]).To(Equal(float64(4)))
						Expect(spot["review_count"]).To(Equal(float64(1)))
						break
					}
				}
				Expect(foundCommunitySpot).To(BeTrue(), "Community spot should be discoverable")
				
				By("User 2 reads reviews for the discovered spot")
				req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/spots/%s/reviews", communitySpotID), nil)
				resp = httptest.NewRecorder()
				testServer.Config.Handler.ServeHTTP(resp, req)
				
				var reviewsResponse map[string]interface{}
				err = json.Unmarshal(resp.Body.Bytes(), &reviewsResponse)
				Expect(err).NotTo(HaveOccurred())
				
				reviews := reviewsResponse["reviews"].([]interface{})
				Expect(len(reviews)).To(Equal(1), "Should find User 1's review")
				
				review := reviews[0].(map[string]interface{})
				Expect(review["user_id"]).To(Equal(user1ID))
				Expect(review["comment"]).To(ContainSubstring("solo studying"))
				
				By("User 2 adds their own review based on their experience")
				user2ReviewRequestBody := map[string]interface{}{
					"spot_id": communitySpotID,
					"rating":  5,
					"comment": "Amazing place! User 1 was right - perfect for solo work. Great internet too!",
				}
				
				bodyBytes, _ = json.Marshal(user2ReviewRequestBody)
				req = httptest.NewRequest(http.MethodPost, "/api/v1/reviews", bytes.NewReader(bodyBytes))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+user2Token)
				resp = httptest.NewRecorder()
				testServer.Config.Handler.ServeHTTP(resp, req)
				Expect(resp.Code).To(Equal(http.StatusCreated))
				
				By("Verifying community statistics are updated")
				req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/spots/%s", communitySpotID), nil)
				resp = httptest.NewRecorder()
				testServer.Config.Handler.ServeHTTP(resp, req)
				
				var updatedSpotResponse map[string]interface{}
				err = json.Unmarshal(resp.Body.Bytes(), &updatedSpotResponse)
				Expect(err).NotTo(HaveOccurred())
				
				Expect(updatedSpotResponse["review_count"]).To(Equal(float64(2)), "Should have 2 reviews")
				Expect(updatedSpotResponse["average_rating"]).To(Equal(float64(4.5)), "Average should be (4+5)/2 = 4.5")
				
				By("Community content discovery scenario completed successfully!")
			})
		})
	})

	Describe("Data Consistency and Validation", func() {
		Context("Given cross-handler data relationships", func() {
			It("Should maintain referential integrity and proper validation", func() {
				By("Creating base user and spot for relationship testing")
				// Create user
				userRequestBody := map[string]interface{}{
					"email":            "consistency@example.com",
					"display_name":     "Consistency Tester",
					"auth_provider":    "google",
					"auth_provider_id": "google_consistency_789",
				}
				
				bodyBytes, err := json.Marshal(userRequestBody)
				Expect(err).NotTo(HaveOccurred())
				
				req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(bodyBytes))
				req.Header.Set("Content-Type", "application/json")
				resp := httptest.NewRecorder()
				testServer.Config.Handler.ServeHTTP(resp, req)
				
				var userResponse map[string]interface{}
				err = json.Unmarshal(resp.Body.Bytes(), &userResponse)
				Expect(err).NotTo(HaveOccurred())
				userID := userResponse["id"].(string)
				
				// Generate token
				tokenRequestBody := map[string]interface{}{
					"email":            "consistency@example.com",
					"auth_provider":    "google",
					"auth_provider_id": "google_consistency_789",
				}
				
				bodyBytes, err = json.Marshal(tokenRequestBody)
				Expect(err).NotTo(HaveOccurred())
				req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/token", bytes.NewReader(bodyBytes))
				req.Header.Set("Content-Type", "application/json")
				resp = httptest.NewRecorder()
				testServer.Config.Handler.ServeHTTP(resp, req)
				
				var userToken string
				cookies := resp.Result().Cookies()
				for _, cookie := range cookies {
					if cookie.Name == "bocchi_access_token" {
						userToken = cookie.Value
						break
					}
				}
				Expect(userToken).NotTo(BeEmpty(), "User token should be retrieved from cookies")
				
				// Create spot
				spotRequestBody := map[string]interface{}{
					"name":         "Consistency Test Spot",
					"latitude":     35.7000,
					"longitude":    139.7000,
					"category":     "cafe",
					"address":      "Consistency Address",
					"country_code": "JP",
				}
				
				bodyBytes, err = json.Marshal(spotRequestBody)
				Expect(err).NotTo(HaveOccurred())
				req = httptest.NewRequest(http.MethodPost, "/api/v1/spots", bytes.NewReader(bodyBytes))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+userToken)
				resp = httptest.NewRecorder()
				testServer.Config.Handler.ServeHTTP(resp, req)
				
				var spotResponse map[string]interface{}
				err = json.Unmarshal(resp.Body.Bytes(), &spotResponse)
				Expect(err).NotTo(HaveOccurred())
				spotID := spotResponse["id"].(string)
				
				By("Testing that reviews maintain referential integrity")
				// Attempt to create review for non-existent spot
				invalidReviewRequestBody := map[string]interface{}{
					"spot_id": "non-existent-spot-id",
					"rating":  4,
					"comment": "This should fail",
				}
				
				bodyBytes, err = json.Marshal(invalidReviewRequestBody)
				Expect(err).NotTo(HaveOccurred())
				req = httptest.NewRequest(http.MethodPost, "/api/v1/reviews", bytes.NewReader(bodyBytes))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+userToken)
				resp = httptest.NewRecorder()
				testServer.Config.Handler.ServeHTTP(resp, req)
				
				Expect(resp.Code).To(Equal(http.StatusNotFound), "Should reject review for non-existent spot")
				
				By("Testing successful review creation with valid relationships")
				validReviewRequestBody := map[string]interface{}{
					"spot_id": spotID,
					"rating":  3,
					"comment": "Valid review with proper relationships",
				}
				
				bodyBytes, err = json.Marshal(validReviewRequestBody)
				Expect(err).NotTo(HaveOccurred())
				req = httptest.NewRequest(http.MethodPost, "/api/v1/reviews", bytes.NewReader(bodyBytes))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+userToken)
				resp = httptest.NewRecorder()
				testServer.Config.Handler.ServeHTTP(resp, req)
				
				Expect(resp.Code).To(Equal(http.StatusCreated), "Should create review with valid relationships")
				
				By("Testing duplicate review prevention")
				// Attempt to create another review from same user for same spot
				duplicateReviewRequestBody := map[string]interface{}{
					"spot_id": spotID,
					"rating":  5,
					"comment": "Duplicate review attempt",
				}
				
				bodyBytes, err = json.Marshal(duplicateReviewRequestBody)
				Expect(err).NotTo(HaveOccurred())
				req = httptest.NewRequest(http.MethodPost, "/api/v1/reviews", bytes.NewReader(bodyBytes))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+userToken)
				resp = httptest.NewRecorder()
				testServer.Config.Handler.ServeHTTP(resp, req)
				
				Expect(resp.Code).To(Equal(http.StatusConflict), "Should prevent duplicate reviews")
				
				By("Verifying data consistency across all endpoints")
				// Check spot statistics are updated
				req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/spots/%s", spotID), nil)
				resp = httptest.NewRecorder()
				testServer.Config.Handler.ServeHTTP(resp, req)
				
				var finalSpotResponse map[string]interface{}
				err = json.Unmarshal(resp.Body.Bytes(), &finalSpotResponse)
				Expect(err).NotTo(HaveOccurred())
				
				Expect(finalSpotResponse["review_count"]).To(Equal(float64(1)), "Spot should show correct review count")
				Expect(finalSpotResponse["average_rating"]).To(Equal(float64(3)), "Spot should show correct average rating")
				
				// Check user review history
				req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/users/%s/reviews", userID), nil)
				resp = httptest.NewRecorder()
				testServer.Config.Handler.ServeHTTP(resp, req)
				
				var userReviewsResponse map[string]interface{}
				err = json.Unmarshal(resp.Body.Bytes(), &userReviewsResponse)
				Expect(err).NotTo(HaveOccurred())
				
				reviews := userReviewsResponse["reviews"].([]interface{})
				Expect(len(reviews)).To(Equal(1), "User should have exactly 1 review")
				
				review := reviews[0].(map[string]interface{})
				Expect(review["spot_id"]).To(Equal(spotID), "Review should reference correct spot")
				Expect(review["user_id"]).To(Equal(userID), "Review should reference correct user")
				
				By("Data consistency validation completed successfully!")
			})
		})
	})
})

// extractBocchiAccessToken extracts the bocchi_access_token from response cookies
func extractBocchiAccessToken(resp *httptest.ResponseRecorder) string {
	// Iterate through all cookies in the HTTP response to find the authentication token
	cookies := resp.Result().Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "bocchi_access_token" {
			return cookie.Value
		}
	}
	// Return empty string if token cookie is not found
	return ""
}