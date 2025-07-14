//go:build integration

package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
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

const (
	InvalidLatitude  = 91.0
	InvalidLongitude = 181.0
)


var _ = Describe("SpotHandler BDD Tests", func() {
	var (
		api        huma.API
		testServer *httptest.Server
		spotHandler *SpotHandler
		spotClient *clients.SpotClient
		authData   *helpers.AuthTestData
	)

	BeforeEach(func() {
		By("Setting up SpotHandler test environment")
		
		// Create spot client with test database
		var err error
		spotClient, err = clients.NewSpotClient("internal", testSuite.TestDB.DB)
		Expect(err).NotTo(HaveOccurred())
		
		// Create spot handler
		spotHandler = NewSpotHandler(spotClient)
		
		// Setup test API with chi router (same as production)
		router := chi.NewRouter()
		api = humachi.New(router, huma.DefaultConfig("Test API", "1.0.0"))
		testServer = httptest.NewServer(router)
		
		// Register spot endpoints using the handler's RegisterRoutes method
		spotHandler.RegisterRoutes(api)
		
		// Setup authentication test data
		authData = testSuite.AuthHelper.NewAuthTestData()
		
		// Create test user in database
		testSuite.FixtureManager.CreateUserFixture(context.Background(), helpers.UserFixture{
			ID:             authData.ValidUserID,
			Email:          authData.TestUser.Email,
			DisplayName:    authData.TestUser.DisplayName,
			AuthProvider:   string(authData.TestUser.AuthProvider),
			AuthProviderID: authData.TestUser.AuthProviderID,
			Preferences:    authData.TestUser.Preferences,
		})
	})

	Describe("Creating a new spot", func() {
		Context("Given a valid authenticated user", func() {
			Context("When they submit a valid spot creation request", func() {
				It("Then the spot should be created successfully", func() {
					By("Preparing a valid spot creation request")
					requestBody := map[string]interface{}{
						"name":         "Test Solo Cafe",
						"latitude":     35.6762,
						"longitude":    139.6503,
						"category":     "cafe",
						"address":      "1-1-1 Test District, Tokyo",
						"country_code": "JP",
					}
					
					bodyBytes, err := json.Marshal(requestBody)
					Expect(err).NotTo(HaveOccurred())
					
					By("Sending the spot creation request")
					req := httptest.NewRequest(http.MethodPost, "/api/v1/spots", bytes.NewReader(bodyBytes))
					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("Authorization", "Bearer "+authData.ValidToken)
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the response")
					Expect(resp.Code).To(Equal(http.StatusCreated), "Expected status 201 Created")
					
					var responseBody map[string]interface{}
					err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
					Expect(err).NotTo(HaveOccurred())
					
					By("Verifying the created spot data")
					Expect(responseBody["id"]).To(Not(BeEmpty()), "Spot ID should be generated")
					Expect(responseBody["name"]).To(Equal("Test Solo Cafe"))
					Expect(responseBody["latitude"]).To(Equal(35.6762))
					Expect(responseBody["longitude"]).To(Equal(139.6503))
					Expect(responseBody["category"]).To(Equal("cafe"))
					Expect(responseBody["country_code"]).To(Equal("JP"))
					Expect(responseBody["created_at"]).To(Not(BeEmpty()), "Creation timestamp should be set")
				})
			})

			Context("When they submit a request with missing required fields", func() {
				It("Then it should return a validation error", func() {
					By("Preparing an invalid spot creation request (missing name)")
					requestBody := map[string]interface{}{
						"latitude":     35.6762,
						"longitude":    139.6503,
						"category":     "cafe",
						"country_code": "JP",
					}
					
					bodyBytes, err := json.Marshal(requestBody)
					Expect(err).NotTo(HaveOccurred())
					
					By("Sending the invalid request")
					req := httptest.NewRequest(http.MethodPost, "/api/v1/spots", bytes.NewReader(bodyBytes))
					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("Authorization", "Bearer "+authData.ValidToken)
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the validation error response")
					Expect(resp.Code).To(Equal(http.StatusUnprocessableEntity), "Expected status 422 Unprocessable Entity")
					
					var errorResponse map[string]interface{}
					err = json.Unmarshal(resp.Body.Bytes(), &errorResponse)
					Expect(err).NotTo(HaveOccurred())
					
					Expect(errorResponse["title"]).To(Equal("Unprocessable Entity"))
					Expect(errorResponse["detail"]).To(Equal("validation failed"))
					
					errors, exists := errorResponse["errors"].([]interface{})
					Expect(exists).To(BeTrue(), "Response should contain errors array")
					Expect(len(errors)).To(BeNumerically(">=", 1), "Should have at least one validation error")
					
					firstError := errors[0].(map[string]interface{})
					Expect(firstError["message"]).To(ContainSubstring("name"))
					Expect(firstError["message"]).To(ContainSubstring("required"))
					Expect(firstError["location"]).To(Equal("body"))
				})
			})

			Context("When they submit a request with invalid coordinates", func() {
				It("Then it should return a validation error for latitude", func() {
					By("Preparing a request with invalid latitude")
					requestBody := map[string]interface{}{
						"name":         "Test Cafe",
						"latitude":     InvalidLatitude, // Invalid: > 90
						"longitude":    139.6503,
						"category":     "cafe",
						"address":      "Test Address",
						"country_code": "JP",
					}
					
					bodyBytes, err := json.Marshal(requestBody)
					Expect(err).NotTo(HaveOccurred())
					
					By("Sending the request with invalid coordinates")
					req := httptest.NewRequest(http.MethodPost, "/api/v1/spots", bytes.NewReader(bodyBytes))
					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("Authorization", "Bearer "+authData.ValidToken)
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the coordinate validation error")
					Expect(resp.Code).To(Equal(http.StatusUnprocessableEntity), "Expected status 422 Unprocessable Entity")
					
					var errorResponse map[string]interface{}
					err = json.Unmarshal(resp.Body.Bytes(), &errorResponse)
					Expect(err).NotTo(HaveOccurred())
					
					Expect(errorResponse["title"]).To(Equal("Unprocessable Entity"))
					Expect(errorResponse["detail"]).To(Equal("validation failed"))
					
					errors, exists := errorResponse["errors"].([]interface{})
					Expect(exists).To(BeTrue(), "Response should contain errors array")
					Expect(len(errors)).To(BeNumerically(">=", 1), "Should have at least one validation error")
					
					firstError := errors[0].(map[string]interface{})
					Expect(firstError["message"]).To(Equal("expected number <= 90"))
					Expect(firstError["location"]).To(Equal("body.latitude"))
				})
			})

			Context("When they submit a request with invalid longitude", func() {
				It("Then it should return a validation error for longitude", func() {
					By("Preparing a request with invalid longitude")
					requestBody := map[string]interface{}{
						"name":         "Test Cafe",
						"latitude":     35.6762,
						"longitude":    InvalidLongitude, // Invalid: > 180
						"category":     "cafe",
						"address":      "Test Address",
						"country_code": "JP",
					}
					
					bodyBytes, err := json.Marshal(requestBody)
					Expect(err).NotTo(HaveOccurred())
					
					By("Sending the request with invalid longitude")
					req := httptest.NewRequest(http.MethodPost, "/api/v1/spots", bytes.NewReader(bodyBytes))
					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("Authorization", "Bearer "+authData.ValidToken)
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the longitude validation error")
					Expect(resp.Code).To(Equal(http.StatusUnprocessableEntity), "Expected status 422 Unprocessable Entity")
					
					var errorResponse map[string]interface{}
					err = json.Unmarshal(resp.Body.Bytes(), &errorResponse)
					Expect(err).NotTo(HaveOccurred())
					
					Expect(errorResponse["title"]).To(Equal("Unprocessable Entity"))
					Expect(errorResponse["detail"]).To(Equal("validation failed"))
					
					errors, exists := errorResponse["errors"].([]interface{})
					Expect(exists).To(BeTrue(), "Response should contain errors array")
					Expect(len(errors)).To(BeNumerically(">=", 1), "Should have at least one validation error")
					
					firstError := errors[0].(map[string]interface{})
					Expect(firstError["message"]).To(Equal("expected number <= 180"))
					Expect(firstError["location"]).To(Equal("body.longitude"))
				})
			})

			Context("When they submit a request with invalid country code format", func() {
				It("Then it should return a validation error for country code", func() {
					By("Preparing a request with invalid country code")
					requestBody := map[string]interface{}{
						"name":         "Test Cafe",
						"latitude":     35.6762,
						"longitude":    139.6503,
						"category":     "cafe",
						"address":      "Test Address",
						"country_code": "INVALID", // Invalid: not 2 chars
					}
					
					bodyBytes, err := json.Marshal(requestBody)
					Expect(err).NotTo(HaveOccurred())
					
					By("Sending the request with invalid country code")
					req := httptest.NewRequest(http.MethodPost, "/api/v1/spots", bytes.NewReader(bodyBytes))
					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("Authorization", "Bearer "+authData.ValidToken)
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the country code validation error")
					Expect(resp.Code).To(Equal(http.StatusUnprocessableEntity), "Expected status 422 Unprocessable Entity")
					
					var errorResponse map[string]interface{}
					err = json.Unmarshal(resp.Body.Bytes(), &errorResponse)
					Expect(err).NotTo(HaveOccurred())
					
					Expect(errorResponse["title"]).To(Equal("Unprocessable Entity"))
					Expect(errorResponse["detail"]).To(Equal("validation failed"))
					
					errors, exists := errorResponse["errors"].([]interface{})
					Expect(exists).To(BeTrue(), "Response should contain errors array")
					Expect(len(errors)).To(BeNumerically(">=", 1), "Should have at least one validation error")
					
					// Country code "INVALID" should trigger multiple validation errors:
					// 1. Length validation (expected length <= 2)
					// 2. Pattern validation (expected string to match pattern ^[A-Z]{2}$)
					errorMessages := make([]string, len(errors))
					for i, err := range errors {
						errMap := err.(map[string]interface{})
						errorMessages[i] = errMap["message"].(string)
						Expect(errMap["location"]).To(Equal("body.country_code"))
					}
					
					Expect(errorMessages).To(ContainElement("expected length <= 2"))
					Expect(errorMessages).To(ContainElement("expected string to match pattern ^[A-Z]{2}$"))
				})
			})
		})

		Context("Given an unauthenticated user", func() {
			Context("When they attempt to create a spot", func() {
				It("Then it should return an authentication error", func() {
					By("Preparing a spot creation request without authentication")
					requestBody := map[string]interface{}{
						"name":         "Test Cafe",
						"latitude":     35.6762,
						"longitude":    139.6503,
						"category":     "cafe",
						"address":      "Test Address",
						"country_code": "JP",
					}
					
					bodyBytes, err := json.Marshal(requestBody)
					Expect(err).NotTo(HaveOccurred())
					
					By("Sending the request without authentication")
					req := httptest.NewRequest(http.MethodPost, "/api/v1/spots", bytes.NewReader(bodyBytes))
					req.Header.Set("Content-Type", "application/json")
					// No Authorization header
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the authentication error")
					Expect(resp.Code).To(Equal(http.StatusUnauthorized), "Expected status 401 Unauthorized")
				})
			})
		})
	})

	Describe("Retrieving an existing spot", func() {
		var existingSpot *helpers.SpotFixture

		BeforeEach(func() {
			By("Creating a test spot in the database")
			existingSpot = &helpers.SpotFixture{
				ID:        "test-spot-retrieval",
				Name:      "Existing Test Cafe",
				Latitude:  35.6762,
				Longitude: 139.6503,
				Category:  "cafe",
				Address:   "Existing Address",
				CountryCode: "JP",
			}
			testSuite.FixtureManager.CreateSpotFixture(context.Background(), *existingSpot)
		})

		Context("Given a valid spot ID", func() {
			Context("When a user requests the spot", func() {
				It("Then the spot details should be returned", func() {
					By("Sending a request to retrieve the spot")
					req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/spots/%s", existingSpot.ID), nil)
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the successful response")
					Expect(resp.Code).To(Equal(http.StatusOK), "Expected status 200 OK")
					
					var responseBody map[string]interface{}
					err := json.Unmarshal(resp.Body.Bytes(), &responseBody)
					Expect(err).NotTo(HaveOccurred())
					
					By("Verifying the spot data")
					Expect(responseBody["id"]).To(Equal(existingSpot.ID))
					Expect(responseBody["name"]).To(Equal(existingSpot.Name))
					Expect(responseBody["latitude"]).To(Equal(existingSpot.Latitude))
					Expect(responseBody["longitude"]).To(Equal(existingSpot.Longitude))
					Expect(responseBody["category"]).To(Equal(existingSpot.Category))
					Expect(responseBody["country_code"]).To(Equal(existingSpot.CountryCode))
				})
			})
		})

		Context("Given a non-existent spot ID", func() {
			Context("When a user requests the spot", func() {
				It("Then it should return a not found error", func() {
					By("Sending a request for a non-existent spot")
					req := httptest.NewRequest(http.MethodGet, "/api/v1/spots/non-existent-id", nil)
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the not found response")
					Expect(resp.Code).To(Equal(http.StatusNotFound), "Expected status 404 Not Found")
				})
			})
		})
	})

	Describe("Listing spots", func() {
		BeforeEach(func() {
			By("Creating multiple test spots")
			testSuite.FixtureManager.SetupStandardFixtures(context.Background())
		})

		Context("Given existing spots in the database", func() {
			Context("When a user requests to list spots", func() {
				It("Then a paginated list of spots should be returned", func() {
					By("Sending a request to list spots")
					req := httptest.NewRequest(http.MethodGet, "/api/v1/spots?page=1&page_size=10", nil)
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the successful response")
					Expect(resp.Code).To(Equal(http.StatusOK), "Expected status 200 OK")
					
					var responseBody map[string]interface{}
					err := json.Unmarshal(resp.Body.Bytes(), &responseBody)
					Expect(err).NotTo(HaveOccurred())
					
					By("Verifying the spots list structure")
					spots, exists := responseBody["spots"]
					Expect(exists).To(BeTrue(), "Response should contain spots array")
					
					spotsArray, ok := spots.([]interface{})
					Expect(ok).To(BeTrue(), "Spots should be an array")
					Expect(len(spotsArray)).To(BeNumerically(">=", 1), "Should contain at least one spot")
					
					By("Verifying pagination information")
					pagination, exists := responseBody["pagination"]
					Expect(exists).To(BeTrue(), "Response should contain pagination info")
					
					paginationMap, ok := pagination.(map[string]interface{})
					Expect(ok).To(BeTrue(), "Pagination should be an object")
					Expect(paginationMap["page"]).To(Equal(float64(1)), "Page should be 1")
					Expect(paginationMap["total"]).To(BeNumerically(">=", 1), "Total should be at least 1")
				})
			})

		})

		Context("Given no spots in the database", func() {
			BeforeEach(func() {
				By("Ensuring database is clean")
				err := testSuite.TestDB.CleanDatabase()
				Expect(err).NotTo(HaveOccurred(), "Failed to clean database: %v", err)
			})

			Context("When a user requests to list spots", func() {
				It("Then an empty list should be returned", func() {
					By("Sending a request to list spots")
					req := httptest.NewRequest(http.MethodGet, "/api/v1/spots", nil)
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the successful empty response")
					Expect(resp.Code).To(Equal(http.StatusOK), "Expected status 200 OK")
					
					var responseBody map[string]interface{}
					err := json.Unmarshal(resp.Body.Bytes(), &responseBody)
					Expect(err).NotTo(HaveOccurred())
					
					By("Verifying empty spots list")
					spots, exists := responseBody["spots"]
					Expect(exists).To(BeTrue(), "Response should contain spots array")
					
					spotsArray, ok := spots.([]interface{})
					Expect(ok).To(BeTrue(), "Spots should be an array")
					Expect(len(spotsArray)).To(Equal(0), "Should be empty")
				})
			})
		})
	})
})