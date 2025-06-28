// +build integration

package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/necofuryai/bocchi-the-map/api/application/clients"
	"github.com/necofuryai/bocchi-the-map/api/domain/entities"
	"github.com/necofuryai/bocchi-the-map/api/pkg/auth"
	"github.com/necofuryai/bocchi-the-map/api/tests/helpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UserHandler BDD Tests", func() {
	var (
		api           huma.API
		testServer    *httptest.Server
		userHandler   *UserHandler
		userClient    *clients.UserClient
		authMiddleware *auth.AuthMiddleware
		authData      *helpers.AuthTestData
	)

	BeforeEach(func() {
		By("Setting up UserHandler test environment")
		
		// Create user client with test database
		var err error
		userClient, err = clients.NewUserClient("internal", testDB.DB)
		Expect(err).NotTo(HaveOccurred())
		
		// Create auth middleware
		authMiddleware = auth.NewAuthMiddleware("test-secret-key", testDB.Queries)
		
		// Create user handler
		userHandler = NewUserHandler(userClient)
		
		// Setup test API with chi router
		router := chi.NewRouter()
		api = humachi.New(router, huma.DefaultConfig("Test API", "1.0.0"))
		testServer = httptest.NewServer(router)
		
		// Register user endpoints with authentication middleware
		userHandler.RegisterRoutesWithAuth(api, authMiddleware)
		
		// Setup authentication test data
		authData = authHelper.NewAuthTestData()
	})

	Describe("User Creation and Update via OAuth", func() {
		Context("Given no existing user with OAuth provider ID", func() {
			Context("When creating a new user via OAuth", func() {
				It("Then a new user should be created successfully", func() {
					By("Preparing a valid OAuth user creation request")
					requestBody := map[string]interface{}{
						"email":            "newuser@example.com",
						"display_name":     "New Test User",
						"auth_provider":    "google",
						"auth_provider_id": "google_new_user_123",
						"avatar_url":       "https://example.com/avatar.jpg",
					}
					
					bodyBytes, err := json.Marshal(requestBody)
					Expect(err).NotTo(HaveOccurred())
					
					By("Sending the user creation request")
					req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(bodyBytes))
					req.Header.Set("Content-Type", "application/json")
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the successful creation response")
					Expect(resp.Code).To(Equal(http.StatusCreated), "Expected status 201 Created")
					
					var responseBody map[string]interface{}
					err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
					Expect(err).NotTo(HaveOccurred())
					
					By("Verifying the created user data")
					Expect(responseBody["id"]).To(Not(BeEmpty()), "User ID should be generated")
					Expect(responseBody["email"]).To(Equal("newuser@example.com"))
					Expect(responseBody["display_name"]).To(Equal("New Test User"))
					Expect(responseBody["auth_provider"]).To(Equal("google"))
					Expect(responseBody["auth_provider_id"]).To(Equal("google_new_user_123"))
					Expect(responseBody["avatar_url"]).To(Equal("https://example.com/avatar.jpg"))
					Expect(responseBody["created_at"]).To(Not(BeEmpty()), "Creation timestamp should be set")
					
					By("Verifying default Japanese preferences are set")
					preferences, exists := responseBody["preferences"]
					Expect(exists).To(BeTrue(), "Default preferences should be set")
					
					prefMap, ok := preferences.(map[string]interface{})
					Expect(ok).To(BeTrue(), "Preferences should be an object")
					Expect(prefMap["language"]).To(Equal("ja"), "Default language should be Japanese")
					Expect(prefMap["timezone"]).To(Equal("Asia/Tokyo"), "Default timezone should be Asia/Tokyo")
				})
			})
		})

		Context("Given an existing user with OAuth provider ID", func() {
			BeforeEach(func() {
				By("Creating an existing user in the database")
				fixtureManager.CreateUserFixture(helpers.UserFixture{
					ID:             "existing-user-123",
					Email:          "existing@example.com",
					DisplayName:    "Existing User",
					AuthProvider:   "google",
					AuthProviderID: "google_existing_123",
					Preferences: entities.UserPreferences{
						Language: "en",
						Timezone: "UTC",
					},
				})
			})
			
			Context("When updating user data via OAuth upsert", func() {
				It("Then the existing user should be updated", func() {
					By("Preparing an OAuth update request with same provider ID")
					requestBody := map[string]interface{}{
						"email":            "existing_updated@example.com",
						"display_name":     "Updated Existing User",
						"auth_provider":    "google",
						"auth_provider_id": "google_existing_123", // Same provider ID
						"avatar_url":       "https://example.com/new_avatar.jpg",
					}
					
					bodyBytes, err := json.Marshal(requestBody)
					Expect(err).NotTo(HaveOccurred())
					
					By("Sending the user update request")
					req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(bodyBytes))
					req.Header.Set("Content-Type", "application/json")
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the successful update response")
					Expect(resp.Code).To(Equal(http.StatusOK), "Expected status 200 OK for update")
					
					var responseBody map[string]interface{}
					err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
					Expect(err).NotTo(HaveOccurred())
					
					By("Verifying the updated user data")
					Expect(responseBody["id"]).To(Equal("existing-user-123"), "User ID should remain the same")
					Expect(responseBody["email"]).To(Equal("existing_updated@example.com"), "Email should be updated")
					Expect(responseBody["display_name"]).To(Equal("Updated Existing User"), "Display name should be updated")
					Expect(responseBody["avatar_url"]).To(Equal("https://example.com/new_avatar.jpg"), "Avatar URL should be updated")
					Expect(responseBody["updated_at"]).To(Not(BeEmpty()), "Update timestamp should be set")
					
					By("Verifying existing preferences are preserved")
					preferences, exists := responseBody["preferences"]
					Expect(exists).To(BeTrue(), "Existing preferences should be preserved")
					
					prefMap, ok := preferences.(map[string]interface{})
					Expect(ok).To(BeTrue(), "Preferences should be an object")
					Expect(prefMap["language"]).To(Equal("en"), "Existing language preference should be preserved")
					Expect(prefMap["timezone"]).To(Equal("UTC"), "Existing timezone preference should be preserved")
				})
			})
		})

		Context("Given an invalid OAuth provider", func() {
			Context("When attempting to create user with unsupported provider", func() {
				It("Then it should return a validation error", func() {
					By("Preparing a request with invalid OAuth provider")
					requestBody := map[string]interface{}{
						"email":            "test@example.com",
						"display_name":     "Test User",
						"auth_provider":    "invalid_provider",
						"auth_provider_id": "invalid_123",
					}
					
					bodyBytes, err := json.Marshal(requestBody)
					Expect(err).NotTo(HaveOccurred())
					
					By("Sending the invalid provider request")
					req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(bodyBytes))
					req.Header.Set("Content-Type", "application/json")
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the validation error response")
					Expect(resp.Code).To(Equal(http.StatusBadRequest), "Expected status 400 Bad Request")
					
					var errorResponse map[string]interface{}
					err = json.Unmarshal(resp.Body.Bytes(), &errorResponse)
					Expect(err).NotTo(HaveOccurred())
					
					Expect(errorResponse["title"]).To(ContainSubstring("validation"), "Should contain validation error message")
				})
			})
		})

		Context("Given missing required fields", func() {
			Context("When attempting to create user without email", func() {
				It("Then it should return a validation error", func() {
					By("Preparing a request missing email field")
					requestBody := map[string]interface{}{
						"display_name":     "Test User",
						"auth_provider":    "google",
						"auth_provider_id": "google_123",
					}
					
					bodyBytes, err := json.Marshal(requestBody)
					Expect(err).NotTo(HaveOccurred())
					
					By("Sending the incomplete request")
					req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(bodyBytes))
					req.Header.Set("Content-Type", "application/json")
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the validation error")
					Expect(resp.Code).To(Equal(http.StatusBadRequest), "Expected status 400 Bad Request")
					
					var errorResponse map[string]interface{}
					err = json.Unmarshal(resp.Body.Bytes(), &errorResponse)
					Expect(err).NotTo(HaveOccurred())
					
					Expect(errorResponse["title"]).To(ContainSubstring("validation"), "Should indicate validation failure")
				})
			})
		})
	})

	Describe("Current User Retrieval", func() {
		BeforeEach(func() {
			By("Creating test user for authentication testing")
			fixtureManager.CreateUserFixture(helpers.UserFixture{
				ID:             authData.ValidUserID,
				Email:          authData.TestUser.Email,
				DisplayName:    authData.TestUser.DisplayName,
				AuthProvider:   string(authData.TestUser.AuthProvider),
				AuthProviderID: authData.TestUser.AuthProviderID,
				Preferences:    authData.TestUser.Preferences,
			})
		})

		Context("Given an authenticated user with valid JWT token", func() {
			Context("When requesting current user information", func() {
				It("Then the user's profile information should be returned", func() {
					By("Sending authenticated request to get current user")
					req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
					req.Header.Set("Authorization", "Bearer "+authData.ValidToken)
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the successful response")
					Expect(resp.Code).To(Equal(http.StatusOK), "Expected status 200 OK")
					
					var responseBody map[string]interface{}
					err := json.Unmarshal(resp.Body.Bytes(), &responseBody)
					Expect(err).NotTo(HaveOccurred())
					
					By("Verifying the user profile information")
					Expect(responseBody["id"]).To(Equal(authData.ValidUserID))
					Expect(responseBody["email"]).To(Equal(authData.TestUser.Email))
					Expect(responseBody["display_name"]).To(Equal(authData.TestUser.DisplayName))
					Expect(responseBody["auth_provider"]).To(Equal(string(authData.TestUser.AuthProvider)))
					Expect(responseBody["auth_provider_id"]).To(Equal(authData.TestUser.AuthProviderID))
					
					By("Verifying preferences are included in response")
					preferences, exists := responseBody["preferences"]
					Expect(exists).To(BeTrue(), "Preferences should be included")
					
					prefMap, ok := preferences.(map[string]interface{})
					Expect(ok).To(BeTrue(), "Preferences should be an object")
					Expect(prefMap["language"]).To(Not(BeEmpty()), "Language preference should be set")
					Expect(prefMap["timezone"]).To(Not(BeEmpty()), "Timezone preference should be set")
				})
			})
		})

		Context("Given no authentication token", func() {
			Context("When attempting to get current user", func() {
				It("Then it should return an authentication error", func() {
					By("Sending unauthenticated request")
					req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
					// No Authorization header
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the authentication error")
					Expect(resp.Code).To(Equal(http.StatusUnauthorized), "Expected status 401 Unauthorized")
					
					var errorResponse map[string]interface{}
					err := json.Unmarshal(resp.Body.Bytes(), &errorResponse)
					Expect(err).NotTo(HaveOccurred())
					
					Expect(errorResponse["title"]).To(ContainSubstring("authentication"), "Should indicate authentication required")
				})
			})
		})

		Context("Given an expired JWT token", func() {
			Context("When attempting to get current user", func() {
				It("Then it should return an authentication error", func() {
					By("Sending request with expired token")
					req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
					req.Header.Set("Authorization", "Bearer "+authData.ExpiredToken)
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the authentication error")
					Expect(resp.Code).To(Equal(http.StatusUnauthorized), "Expected status 401 Unauthorized")
				})
			})
		})

		Context("Given an invalid JWT token", func() {
			Context("When attempting to get current user", func() {
				It("Then it should return an authentication error", func() {
					By("Sending request with invalid token")
					req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
					req.Header.Set("Authorization", "Bearer invalid.jwt.token")
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the authentication error")
					Expect(resp.Code).To(Equal(http.StatusUnauthorized), "Expected status 401 Unauthorized")
				})
			})
		})
	})

	Describe("User Preference Updates", func() {
		BeforeEach(func() {
			By("Creating test user for preference testing")
			fixtureManager.CreateUserFixture(helpers.UserFixture{
				ID:             authData.ValidUserID,
				Email:          authData.TestUser.Email,
				DisplayName:    authData.TestUser.DisplayName,
				AuthProvider:   string(authData.TestUser.AuthProvider),
				AuthProviderID: authData.TestUser.AuthProviderID,
				Preferences: entities.UserPreferences{
					Language: "ja",
					DarkMode: false,
					Timezone: "Asia/Tokyo",
				},
			})
		})

		Context("Given an authenticated user with existing preferences", func() {
			Context("When updating preferences with valid data", func() {
				It("Then the preferences should be updated successfully", func() {
					By("Preparing valid preferences update request")
					requestBody := map[string]interface{}{
						"language":  "en",
						"dark_mode": true,
						"timezone":  "UTC",
					}
					
					bodyBytes, err := json.Marshal(requestBody)
					Expect(err).NotTo(HaveOccurred())
					
					By("Sending authenticated preferences update request")
					req := httptest.NewRequest(http.MethodPatch, "/api/v1/users/me/preferences", bytes.NewReader(bodyBytes))
					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("Authorization", "Bearer "+authData.ValidToken)
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the successful update response")
					Expect(resp.Code).To(Equal(http.StatusOK), "Expected status 200 OK")
					
					var responseBody map[string]interface{}
					err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
					Expect(err).NotTo(HaveOccurred())
					
					By("Verifying the updated preferences")
					preferences, exists := responseBody["preferences"]
					Expect(exists).To(BeTrue(), "Updated preferences should be returned")
					
					prefMap, ok := preferences.(map[string]interface{})
					Expect(ok).To(BeTrue(), "Preferences should be an object")
					Expect(prefMap["language"]).To(Equal("en"), "Language should be updated to English")
					Expect(prefMap["dark_mode"]).To(Equal(true), "Dark mode should be enabled")
					Expect(prefMap["timezone"]).To(Equal("UTC"), "Timezone should be updated to UTC")
					
					By("Verifying the updated_at timestamp is refreshed")
					Expect(responseBody["updated_at"]).To(Not(BeEmpty()), "Update timestamp should be set")
				})
			})
		})

		Context("Given invalid preference values", func() {
			Context("When updating preferences with invalid language code", func() {
				It("Then it should return a validation error", func() {
					By("Preparing request with invalid language")
					requestBody := map[string]interface{}{
						"language":  "invalid_lang_code",
						"dark_mode": false,
						"timezone":  "Asia/Tokyo",
					}
					
					bodyBytes, err := json.Marshal(requestBody)
					Expect(err).NotTo(HaveOccurred())
					
					By("Sending request with invalid language")
					req := httptest.NewRequest(http.MethodPatch, "/api/v1/users/me/preferences", bytes.NewReader(bodyBytes))
					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("Authorization", "Bearer "+authData.ValidToken)
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the validation error")
					Expect(resp.Code).To(Equal(http.StatusBadRequest), "Expected status 400 Bad Request")
					
					var errorResponse map[string]interface{}
					err = json.Unmarshal(resp.Body.Bytes(), &errorResponse)
					Expect(err).NotTo(HaveOccurred())
					
					Expect(errorResponse["title"]).To(ContainSubstring("validation"), "Should specify language validation failure")
				})
			})

			Context("When updating preferences with invalid timezone", func() {
				It("Then it should return a validation error", func() {
					By("Preparing request with invalid timezone")
					requestBody := map[string]interface{}{
						"language":  "en",
						"dark_mode": false,
						"timezone":  "Invalid/Timezone",
					}
					
					bodyBytes, err := json.Marshal(requestBody)
					Expect(err).NotTo(HaveOccurred())
					
					By("Sending request with invalid timezone")
					req := httptest.NewRequest(http.MethodPatch, "/api/v1/users/me/preferences", bytes.NewReader(bodyBytes))
					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("Authorization", "Bearer "+authData.ValidToken)
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the validation error")
					Expect(resp.Code).To(Equal(http.StatusBadRequest), "Expected status 400 Bad Request")
					
					var errorResponse map[string]interface{}
					err = json.Unmarshal(resp.Body.Bytes(), &errorResponse)
					Expect(err).NotTo(HaveOccurred())
					
					Expect(errorResponse["title"]).To(ContainSubstring("validation"), "Should specify timezone validation failure")
				})
			})
		})

		Context("Given an unauthenticated request", func() {
			Context("When attempting to update preferences", func() {
				It("Then it should return an authentication error", func() {
					By("Preparing preferences update request")
					requestBody := map[string]interface{}{
						"language":  "en",
						"dark_mode": true,
						"timezone":  "UTC",
					}
					
					bodyBytes, err := json.Marshal(requestBody)
					Expect(err).NotTo(HaveOccurred())
					
					By("Sending unauthenticated request")
					req := httptest.NewRequest(http.MethodPatch, "/api/v1/users/me/preferences", bytes.NewReader(bodyBytes))
					req.Header.Set("Content-Type", "application/json")
					// No Authorization header
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the authentication error")
					Expect(resp.Code).To(Equal(http.StatusUnauthorized), "Expected status 401 Unauthorized")
				})
			})
		})

		Context("Given partial preference updates", func() {
			Context("When updating only some preference fields", func() {
				It("Then only specified fields should be updated", func() {
					By("Preparing partial preferences update (only language)")
					requestBody := map[string]interface{}{
						"language": "en",
						// Omitting dark_mode and timezone
					}
					
					bodyBytes, err := json.Marshal(requestBody)
					Expect(err).NotTo(HaveOccurred())
					
					By("Sending partial preferences update")
					req := httptest.NewRequest(http.MethodPatch, "/api/v1/users/me/preferences", bytes.NewReader(bodyBytes))
					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("Authorization", "Bearer "+authData.ValidToken)
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the successful partial update")
					Expect(resp.Code).To(Equal(http.StatusOK), "Expected status 200 OK")
					
					var responseBody map[string]interface{}
					err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
					Expect(err).NotTo(HaveOccurred())
					
					By("Verifying only language was updated, others preserved")
					preferences, exists := responseBody["preferences"]
					Expect(exists).To(BeTrue(), "Preferences should be returned")
					
					prefMap, ok := preferences.(map[string]interface{})
					Expect(ok).To(BeTrue(), "Preferences should be an object")
					Expect(prefMap["language"]).To(Equal("en"), "Language should be updated")
					Expect(prefMap["dark_mode"]).To(Equal(false), "Dark mode should be preserved")
					Expect(prefMap["timezone"]).To(Equal("Asia/Tokyo"), "Timezone should be preserved")
				})
			})
		})
	})
})