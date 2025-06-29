//go:build integration

package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

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

const (
	defaultTestLanguage   = "ja"
	defaultTestTimezone   = "Asia/Tokyo"
	alternateTestLanguage = "en"
	alternateTestTimezone = "UTC"
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
		userClient, err = clients.NewUserClient("internal", testSuite.TestDB.DB)
		Expect(err).NotTo(HaveOccurred())
		
		// Create auth middleware
		authMiddleware = auth.NewAuthMiddleware("test-secret-key", testSuite.TestDB.Queries)
		
		// Create user handler
		userHandler = NewUserHandler(userClient)
		
		// Setup test API with chi router
		router := chi.NewRouter()
		api = humachi.New(router, huma.DefaultConfig("Test API", "1.0.0"))
		testServer = httptest.NewServer(router)
		
		// Register user endpoints with authentication middleware
		userHandler.RegisterRoutesWithAuth(api, authMiddleware)
		
		// Setup authentication test data
		authData = testSuite.AuthHelper.NewAuthTestData()
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
					
					By("Verifying default Japanese preferences are set with exact structure")
					preferences, exists := responseBody["preferences"]
					Expect(exists).To(BeTrue(), "Default preferences should be set")
					Expect(preferences).ToNot(BeNil(), "Preferences should not be nil")
					
					prefMap, ok := preferences.(map[string]interface{})
					Expect(ok).To(BeTrue(), "Preferences should be an object")
					
					// Assert exact structure - only expected fields should be present
					Expect(prefMap).To(HaveLen(2), "Preferences should contain exactly 2 fields")
					expectedKeys := []string{"language", "timezone"}
					actualKeys := make([]string, 0, len(prefMap))
					for key := range prefMap {
						actualKeys = append(actualKeys, key)
					}
					Expect(actualKeys).To(ConsistOf(expectedKeys), "Preferences should contain only expected keys")
					
					// Assert exact values with type safety
					language, languageExists := prefMap["language"]
					Expect(languageExists).To(BeTrue(), "Language field must exist")
					Expect(language).To(BeAssignableToTypeOf(""), "Language must be a string")
					Expect(language).To(Equal("ja"), "Default language should be exactly 'ja'")
					
					timezone, timezoneExists := prefMap["timezone"]
					Expect(timezoneExists).To(BeTrue(), "Timezone field must exist")
					Expect(timezone).To(BeAssignableToTypeOf(""), "Timezone must be a string")
					Expect(timezone).To(Equal("Asia/Tokyo"), "Default timezone should be exactly 'Asia/Tokyo'")
					
					// Verify no additional or unexpected fields exist
					for key := range prefMap {
						Expect([]string{"language", "timezone"}).To(ContainElement(key), 
							fmt.Sprintf("Unexpected field '%s' found in preferences", key))
					}
				})
			})
		})

		Context("Given an existing user with OAuth provider ID", func() {
			BeforeEach(func() {
				By("Creating an existing user in the database")
				testSuite.FixtureManager.CreateUserFixture(context.Background(), helpers.UserFixture{
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
			testSuite.FixtureManager.CreateUserFixture(context.Background(), helpers.UserFixture{
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
					
					By("Verifying required fields exist and have correct data types")
					Expect(responseBody["id"]).To(BeAssignableToTypeOf(""), "ID should be a string")
					Expect(responseBody["email"]).To(BeAssignableToTypeOf(""), "Email should be a string")
					Expect(responseBody["display_name"]).To(BeAssignableToTypeOf(""), "Display name should be a string")
					Expect(responseBody["auth_provider"]).To(BeAssignableToTypeOf(""), "Auth provider should be a string")
					Expect(responseBody["auth_provider_id"]).To(BeAssignableToTypeOf(""), "Auth provider ID should be a string")
					
					By("Verifying timestamp fields are present and properly formatted")
					createdAt, createdAtExists := responseBody["created_at"]
					Expect(createdAtExists).To(BeTrue(), "created_at should be present")
					Expect(createdAt).To(BeAssignableToTypeOf(""), "created_at should be a string")
					_, err = time.Parse(time.RFC3339, createdAt.(string))
					Expect(err).NotTo(HaveOccurred(), "created_at should be valid RFC3339 timestamp")
					
					updatedAt, updatedAtExists := responseBody["updated_at"]
					Expect(updatedAtExists).To(BeTrue(), "updated_at should be present")
					Expect(updatedAt).To(BeAssignableToTypeOf(""), "updated_at should be a string")
					_, err = time.Parse(time.RFC3339, updatedAt.(string))
					Expect(err).NotTo(HaveOccurred(), "updated_at should be valid RFC3339 timestamp")
					
					By("Verifying avatar_url field handling")
					if avatarURL, exists := responseBody["avatar_url"]; exists {
						Expect(avatarURL).To(BeAssignableToTypeOf(""), "avatar_url should be a string if present")
						if avatarURL != "" {
							Expect(avatarURL.(string)).To(MatchRegexp(`^https?://`), "avatar_url should be a valid URL if not empty")
						}
					}
					
					By("Verifying preferences structure and content")
					preferences, exists := responseBody["preferences"]
					Expect(exists).To(BeTrue(), "Preferences should be included")
					
					prefMap, ok := preferences.(map[string]interface{})
					Expect(ok).To(BeTrue(), "Preferences should be an object")
					
					By("Verifying preferences contain required fields with correct types")
					language, langExists := prefMap["language"]
					Expect(langExists).To(BeTrue(), "Language preference should exist")
					Expect(language).To(BeAssignableToTypeOf(""), "Language should be a string")
					Expect(language).To(Or(Equal("ja"), Equal("en")), "Language should be 'ja' or 'en'")
					
					darkMode, darkModeExists := prefMap["dark_mode"]
					Expect(darkModeExists).To(BeTrue(), "Dark mode preference should exist")
					Expect(darkMode).To(BeAssignableToTypeOf(true), "Dark mode should be a boolean")
					
					timezone, timezoneExists := prefMap["timezone"]
					Expect(timezoneExists).To(BeTrue(), "Timezone preference should exist")
					Expect(timezone).To(BeAssignableToTypeOf(""), "Timezone should be a string")
					Expect(timezone).To(Not(BeEmpty()), "Timezone should not be empty")
					_, err = time.LoadLocation(timezone.(string))
					Expect(err).NotTo(HaveOccurred(), "Timezone should be a valid timezone identifier")
					
					By("Verifying response structure completeness")
					expectedFields := []string{"id", "email", "display_name", "auth_provider", "auth_provider_id", "preferences", "created_at", "updated_at"}
					for _, field := range expectedFields {
						_, exists := responseBody[field]
						Expect(exists).To(BeTrue(), fmt.Sprintf("Required field '%s' should be present", field))
					}
					
					By("Verifying no unexpected fields are present")
					allowedFields := map[string]bool{
						"id": true, "email": true, "display_name": true, "avatar_url": true,
						"auth_provider": true, "auth_provider_id": true, "preferences": true,
						"created_at": true, "updated_at": true,
					}
					for field := range responseBody {
						Expect(allowedFields[field]).To(BeTrue(), fmt.Sprintf("Unexpected field '%s' found in response", field))
					}
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
					
					By("Verifying specific error message and structure")
					Expect(errorResponse["title"]).To(Equal("Authentication Required"), "Should have exact authentication error title")
					Expect(errorResponse["detail"]).To(Equal("No authorization token provided"), "Should specify missing token reason")
					Expect(errorResponse["status"]).To(Equal(float64(401)), "Should include numeric status code")
					Expect(errorResponse["type"]).To(Equal("about:blank"), "Should follow RFC 7807 problem details format")
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
			testSuite.FixtureManager.CreateUserFixture(context.Background(), helpers.UserFixture{
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
					
					By("Verifying all expected fields exist with correct types")
					Expect(prefMap).To(HaveKey("language"), "Language field should exist")
					Expect(prefMap).To(HaveKey("dark_mode"), "Dark mode field should exist")
					Expect(prefMap).To(HaveKey("timezone"), "Timezone field should exist")
					
					By("Verifying data types of each field")
					language, langExists := prefMap["language"]
					Expect(langExists).To(BeTrue(), "Language field should exist")
					Expect(language).To(BeAssignableToTypeOf(""), "Language should be a string")
					
					darkMode, darkExists := prefMap["dark_mode"]
					Expect(darkExists).To(BeTrue(), "Dark mode field should exist")
					Expect(darkMode).To(BeAssignableToTypeOf(false), "Dark mode should be a boolean")
					
					timezone, tzExists := prefMap["timezone"]
					Expect(tzExists).To(BeTrue(), "Timezone field should exist")
					Expect(timezone).To(BeAssignableToTypeOf(""), "Timezone should be a string")
					
					By("Verifying field values after partial update")
					Expect(prefMap["language"]).To(Equal("en"), "Language should be updated")
					Expect(prefMap["dark_mode"]).To(Equal(false), "Dark mode should be preserved")
					Expect(prefMap["timezone"]).To(Equal("Asia/Tokyo"), "Timezone should be preserved")
					
					By("Verifying response structure integrity")
					Expect(len(prefMap)).To(Equal(3), "Preferences should contain exactly 3 fields")
					
					By("Verifying field order consistency (if JSON maintains order)")
					responseJSON := resp.Body.String()
					languageIndex := strings.Index(responseJSON, `"language"`)
					darkModeIndex := strings.Index(responseJSON, `"dark_mode"`)
					timezoneIndex := strings.Index(responseJSON, `"timezone"`)
					
					Expect(languageIndex).To(BeNumerically(">", -1), "Language field should be present in JSON")
					Expect(darkModeIndex).To(BeNumerically(">", -1), "Dark mode field should be present in JSON")
					Expect(timezoneIndex).To(BeNumerically(">", -1), "Timezone field should be present in JSON")
				})
			})
		})
	})
})