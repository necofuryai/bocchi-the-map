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
	"bocchi/api/domain/entities"
	"bocchi/api/pkg/auth"
	"bocchi/api/tests/helpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Helper functions for user handler tests
func verifyPublicUserResponse(responseBody map[string]interface{}, userID, displayName string) {
	Expect(responseBody["id"]).To(Equal(userID), "User ID should match")
	Expect(responseBody["display_name"]).To(Equal(displayName), "Display name should match")
	Expect(responseBody["created_at"]).To(Not(BeEmpty()), "Created at timestamp should be present")

	// Verify that sensitive information is NOT exposed in public endpoint
	Expect(responseBody["email"]).To(BeEmpty(), "Email should not be exposed in public endpoint")
	Expect(responseBody["auth_provider"]).To(BeEmpty(), "Auth provider should not be exposed in public endpoint")
	Expect(responseBody["auth_provider_id"]).To(BeEmpty(), "Auth provider ID should not be exposed in public endpoint")
	Expect(responseBody["preferences"]).To(BeEmpty(), "Preferences should not be exposed in public endpoint")
}

func verifyFullUserResponse(responseBody map[string]interface{}, userID, email, displayName string) {
	Expect(responseBody["id"]).To(Equal(userID), "User ID should match")
	Expect(responseBody["email"]).To(Equal(email), "Email should be present for authenticated user")
	Expect(responseBody["display_name"]).To(Equal(displayName), "Display name should match")
	Expect(responseBody["auth_provider"]).To(Not(BeEmpty()), "Auth provider should be present for authenticated user")
	Expect(responseBody["auth_provider_id"]).To(Not(BeEmpty()), "Auth provider ID should be present for authenticated user")
	Expect(responseBody["created_at"]).To(Not(BeEmpty()), "Created at timestamp should be present")
	Expect(responseBody["updated_at"]).To(Not(BeEmpty()), "Updated at timestamp should be present")
}

func verifyErrorResponse(resp *httptest.ResponseRecorder, expectedStatus int, expectedTitleContent string) {
	Expect(resp.Code).To(Equal(expectedStatus), fmt.Sprintf("Expected status %d", expectedStatus))

	var errorResponse map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &errorResponse)
	Expect(err).NotTo(HaveOccurred())

	title, exists := errorResponse["title"]
	Expect(exists).To(BeTrue(), "Error response should have a title")
	Expect(title).To(ContainSubstring(expectedTitleContent), fmt.Sprintf("Error title should contain '%s'", expectedTitleContent))
}

var _ = Describe("UserHandler BDD Tests", func() {
	var (
		api         huma.API
		testServer  *httptest.Server
		userHandler *UserHandler
		userClient  *clients.UserClient
		authData    *helpers.AuthTestData
		otherUser   *helpers.UserFixture
	)

	BeforeEach(func() {
		By("Setting up UserHandler test environment")

		// Create user client with test database
		var err error
		userClient, err = clients.NewUserClient("internal", testSuite.TestDB.DB)
		Expect(err).NotTo(HaveOccurred())

		// Create user handler
		userHandler = NewUserHandler(userClient)

		// Setup test API with chi router
		router := chi.NewRouter()
		api = humachi.New(router, huma.DefaultConfig("Test API", "1.0.0"))
		testServer = httptest.NewServer(router)

		// Create auth middleware for protected endpoints
		authMiddleware := auth.NewAuthMiddleware("test-jwt-secret", testSuite.TestDB.Queries)

		// Register user endpoints (both public and authenticated)
		userHandler.RegisterRoutesWithAuth(api, authMiddleware)

		// Setup authentication test data
		authData = testSuite.AuthHelper.NewAuthTestData()

		// Create primary test user in database
		testSuite.FixtureManager.CreateUserFixture(context.Background(), helpers.UserFixture{
			ID:             authData.ValidUserID,
			Email:          authData.TestUser.Email,
			DisplayName:    authData.TestUser.DisplayName,
			AuthProvider:   string(authData.TestUser.AuthProvider),
			AuthProviderID: authData.TestUser.AuthProviderID,
			Preferences:    authData.TestUser.Preferences,
		})

		// Create another user for permission testing
		otherUser = &helpers.UserFixture{
			ID:             "other-user-123",
			Email:          "otheruser@example.com",
			DisplayName:    "Other User",
			AuthProvider:   "google",
			AuthProviderID: "google_other_user_123",
			Preferences: entities.UserPreferences{
				Language: "en",
				DarkMode: true,
				Timezone: "UTC",
			},
		}
		testSuite.FixtureManager.CreateUserFixture(context.Background(), *otherUser)
	})

	Describe("Public User Information Retrieval", func() {
		Context("Given a user exists in the system", func() {
			Context("When requesting public user information by ID", func() {
				It("Then limited user data should be returned without sensitive information", func() {
					By("Sending request for public user information")
					req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/users/%s", authData.ValidUserID), nil)

					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)

					By("Verifying successful response")
					Expect(resp.Code).To(Equal(http.StatusOK), "Expected status 200 OK")

					var responseBody map[string]interface{}
					err := json.Unmarshal(resp.Body.Bytes(), &responseBody)
					Expect(err).NotTo(HaveOccurred())

					By("Verifying public user data contains only non-sensitive information")
					verifyPublicUserResponse(responseBody, authData.ValidUserID, authData.TestUser.DisplayName)

					By("Verifying avatar URL structure")
					// Avatar URL should be included in public response if present
					if avatarURL, exists := responseBody["avatar_url"]; exists {
						Expect(avatarURL).To(BeAssignableToTypeOf(""), "Avatar URL should be a string if present")
					}
				})
			})
		})

		Context("Given a non-existent user ID", func() {
			Context("When requesting public user information", func() {
				It("Then a not found error should be returned", func() {
					By("Sending request for non-existent user")
					req := httptest.NewRequest(http.MethodGet, "/api/v1/users/non-existent-user-id", nil)

					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)

					By("Verifying not found error")
					verifyErrorResponse(resp, http.StatusNotFound, "not found")
				})
			})
		})

		Context("Given an invalid user ID format", func() {
			Context("When requesting public user information", func() {
				It("Then appropriate error handling should occur", func() {
					By("Sending request with malformed user ID")
					req := httptest.NewRequest(http.MethodGet, "/api/v1/users/", nil)

					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)

					By("Verifying error response for malformed ID")
					Expect(resp.Code).To(BeNumerically(">=", 400), "Should return client error for malformed ID")
				})
			})
		})
	})

	Describe("Current User Information Retrieval", func() {
		Context("Given an authenticated user", func() {
			Context("When requesting current user information", func() {
				It("Then full user data should be returned including sensitive information", func() {
					By("Sending authenticated request for current user")
					req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
					req.Header.Set("Authorization", "Bearer "+authData.ValidToken)

					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)

					By("Verifying successful response")
					Expect(resp.Code).To(Equal(http.StatusOK), "Expected status 200 OK")

					var responseBody map[string]interface{}
					err := json.Unmarshal(resp.Body.Bytes(), &responseBody)
					Expect(err).NotTo(HaveOccurred())

					By("Verifying full user data includes all information")
					verifyFullUserResponse(responseBody, authData.ValidUserID, authData.TestUser.Email, authData.TestUser.DisplayName)

					By("Verifying preferences are included for authenticated user")
					preferences, exists := responseBody["preferences"]
					Expect(exists).To(BeTrue(), "Preferences should be present for authenticated user")
					Expect(preferences).To(Not(BeNil()), "Preferences should not be nil")
				})
			})
		})

		Context("Given an unauthenticated request", func() {
			Context("When requesting current user information", func() {
				It("Then an authentication error should be returned", func() {
					By("Sending unauthenticated request")
					req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
					// No Authorization header

					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)

					By("Verifying authentication error")
					verifyErrorResponse(resp, http.StatusUnauthorized, "authentication")
				})
			})
		})

		Context("Given an invalid authentication token", func() {
			Context("When requesting current user information", func() {
				It("Then an authentication error should be returned", func() {
					By("Sending request with invalid token")
					req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
					req.Header.Set("Authorization", "Bearer invalid-token-123")

					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)

					By("Verifying authentication error for invalid token")
					verifyErrorResponse(resp, http.StatusUnauthorized, "authentication")
				})
			})
		})
	})

	Describe("User Profile Updates - Self vs Other Users", func() {
		Context("Given an authenticated user", func() {
			Context("When updating their own profile via user ID endpoint", func() {
				It("Then the update should be successful", func() {
					By("Preparing user profile update data")
					updateData := map[string]interface{}{
						"display_name": "Updated Display Name",
						"avatar_url":   "https://example.com/new-avatar.jpg",
						"preferences": map[string]interface{}{
							"theme":         "dark",
							"notifications": true,
						},
					}

					bodyBytes, err := json.Marshal(updateData)
					Expect(err).NotTo(HaveOccurred())

					By("Sending authenticated update request for own profile")
					req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/users/%s", authData.ValidUserID), bytes.NewReader(bodyBytes))
					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("Authorization", "Bearer "+authData.ValidToken)

					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)

					By("Verifying successful update")
					Expect(resp.Code).To(Equal(http.StatusOK), "Expected status 200 OK")

					var responseBody map[string]interface{}
					err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
					Expect(err).NotTo(HaveOccurred())

					By("Verifying updated user data")
					Expect(responseBody["id"]).To(Equal(authData.ValidUserID))
					Expect(responseBody["display_name"]).To(Equal("Updated Display Name"))
					Expect(responseBody["avatar_url"]).To(Equal("https://example.com/new-avatar.jpg"))

					preferences, exists := responseBody["preferences"]
					Expect(exists).To(BeTrue(), "Preferences should be present")
					prefMap, ok := preferences.(map[string]interface{})
					Expect(ok).To(BeTrue(), "Preferences should be a map")
					Expect(prefMap["theme"]).To(Equal("dark"))
					Expect(prefMap["notifications"]).To(Equal(true))
				})
			})

			Context("When updating their own profile via /me endpoint", func() {
				It("Then the update should be successful", func() {
					By("Preparing user profile update data")
					updateData := map[string]interface{}{
						"display_name": "Me Updated Name",
						"preferences": map[string]interface{}{
							"language": "ja",
						},
					}

					bodyBytes, err := json.Marshal(updateData)
					Expect(err).NotTo(HaveOccurred())

					By("Sending authenticated update request to /me endpoint")
					req := httptest.NewRequest(http.MethodPut, "/api/v1/users/me", bytes.NewReader(bodyBytes))
					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("Authorization", "Bearer "+authData.ValidToken)

					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)

					By("Verifying successful update")
					Expect(resp.Code).To(Equal(http.StatusOK), "Expected status 200 OK")

					var responseBody map[string]interface{}
					err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
					Expect(err).NotTo(HaveOccurred())

					By("Verifying updated user data")
					Expect(responseBody["id"]).To(Equal(authData.ValidUserID))
					Expect(responseBody["display_name"]).To(Equal("Me Updated Name"))
				})
			})

			Context("When attempting to update another user's profile", func() {
				It("Then access should be denied with forbidden error", func() {
					By("Preparing update data for another user")
					updateData := map[string]interface{}{
						"display_name": "Malicious Update",
					}

					bodyBytes, err := json.Marshal(updateData)
					Expect(err).NotTo(HaveOccurred())

					By("Sending authenticated update request for another user")
					req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/users/%s", otherUser.ID), bytes.NewReader(bodyBytes))
					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("Authorization", "Bearer "+authData.ValidToken)

					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)

					By("Verifying access denied")
					verifyErrorResponse(resp, http.StatusForbidden, "insufficient permissions")
				})
			})
		})

		Context("Given invalid update data", func() {
			Context("When updating profile with invalid display name", func() {
				It("Then validation error should be returned", func() {
					By("Preparing invalid update data")
					updateData := map[string]interface{}{
						"display_name": "", // Empty string should be invalid
					}

					bodyBytes, err := json.Marshal(updateData)
					Expect(err).NotTo(HaveOccurred())

					By("Sending update request with invalid data")
					req := httptest.NewRequest(http.MethodPut, "/api/v1/users/me", bytes.NewReader(bodyBytes))
					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("Authorization", "Bearer "+authData.ValidToken)

					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)

					By("Verifying validation error")
					verifyErrorResponse(resp, http.StatusBadRequest, "validation")
				})
			})

			Context("When updating profile with overly long avatar URL", func() {
				It("Then validation error should be returned", func() {
					By("Preparing update data with overly long URL")
					longURL := "https://example.com/" + string(make([]byte, 600)) // Exceeds 500 character limit
					updateData := map[string]interface{}{
						"avatar_url": longURL,
					}

					bodyBytes, err := json.Marshal(updateData)
					Expect(err).NotTo(HaveOccurred())

					By("Sending update request with overly long URL")
					req := httptest.NewRequest(http.MethodPut, "/api/v1/users/me", bytes.NewReader(bodyBytes))
					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("Authorization", "Bearer "+authData.ValidToken)

					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)

					By("Verifying validation error")
					verifyErrorResponse(resp, http.StatusBadRequest, "validation")
				})
			})
		})
	})

	Describe("Permission Validation for User Operations", func() {
		Context("Given various authentication scenarios", func() {
			Context("When accessing protected endpoints without authentication", func() {
				It("Then all protected endpoints should require authentication", func() {
					protectedEndpoints := []struct {
						method  string
						path    string
						hasBody bool
					}{
						{http.MethodGet, "/api/v1/users/me", false},
						{http.MethodPut, "/api/v1/users/me", true},
						{http.MethodPut, fmt.Sprintf("/api/v1/users/%s", authData.ValidUserID), true},
					}

					for _, endpoint := range protectedEndpoints {
						By(fmt.Sprintf("Testing %s %s without authentication", endpoint.method, endpoint.path))

						var req *http.Request
						if endpoint.hasBody {
							bodyData := map[string]interface{}{"display_name": "Test"}
							bodyBytes, err := json.Marshal(bodyData)
							Expect(err).NotTo(HaveOccurred())
							req = httptest.NewRequest(endpoint.method, endpoint.path, bytes.NewReader(bodyBytes))
							req.Header.Set("Content-Type", "application/json")
						} else {
							req = httptest.NewRequest(endpoint.method, endpoint.path, nil)
						}

						resp := httptest.NewRecorder()
						testServer.Config.Handler.ServeHTTP(resp, req)

						Expect(resp.Code).To(Equal(http.StatusUnauthorized), fmt.Sprintf("%s %s should require authentication", endpoint.method, endpoint.path))
					}
				})
			})
		})

		Context("Given different user roles and permissions", func() {
			Context("When checking self vs other user access patterns", func() {
				It("Then access should be properly restricted", func() {
					By("Verifying that user can only update their own profile")
					// This is already tested above, but verifying the pattern

					updateData := map[string]interface{}{
						"display_name": "Should Not Work",
					}
					bodyBytes, err := json.Marshal(updateData)
					Expect(err).NotTo(HaveOccurred())

					// Try to update other user
					req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/users/%s", otherUser.ID), bytes.NewReader(bodyBytes))
					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("Authorization", "Bearer "+authData.ValidToken)

					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)

					By("Verifying proper permission denial")
					Expect(resp.Code).To(Equal(http.StatusForbidden), "Should deny access to other user's profile")

					By("Verifying that the user can still update their own profile")
					req2 := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/users/%s", authData.ValidUserID), bytes.NewReader(bodyBytes))
					req2.Header.Set("Content-Type", "application/json")
					req2.Header.Set("Authorization", "Bearer "+authData.ValidToken)

					resp2 := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp2, req2)

					Expect(resp2.Code).To(Equal(http.StatusOK), "Should allow access to own profile")
				})
			})
		})
	})

	Describe("Privacy Controls and Data Protection", func() {
		Context("Given different access levels", func() {
			Context("When comparing public vs authenticated user data", func() {
				It("Then sensitive data should only be available to authenticated users", func() {
					By("Retrieving public user data")
					publicReq := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/users/%s", authData.ValidUserID), nil)
					publicResp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(publicResp, publicReq)

					By("Retrieving authenticated user data")
					authReq := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
					authReq.Header.Set("Authorization", "Bearer "+authData.ValidToken)
					authResp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(authResp, authReq)

					By("Verifying both requests succeed")
					Expect(publicResp.Code).To(Equal(http.StatusOK))
					Expect(authResp.Code).To(Equal(http.StatusOK))

					var publicData, authData map[string]interface{}
					err := json.Unmarshal(publicResp.Body.Bytes(), &publicData)
					Expect(err).NotTo(HaveOccurred())
					err = json.Unmarshal(authResp.Body.Bytes(), &authData)
					Expect(err).NotTo(HaveOccurred())

					By("Verifying public data excludes sensitive information")
					sensitiveFields := []string{"email", "auth_provider", "auth_provider_id", "preferences"}
					for _, field := range sensitiveFields {
						Expect(publicData[field]).To(BeEmpty(), fmt.Sprintf("Public endpoint should not expose %s", field))
						Expect(authData[field]).To(Not(BeEmpty()), fmt.Sprintf("Authenticated endpoint should include %s", field))
					}

					By("Verifying public data includes appropriate non-sensitive fields")
					publicFields := []string{"id", "display_name", "created_at"}
					for _, field := range publicFields {
						Expect(publicData[field]).To(Not(BeEmpty()), fmt.Sprintf("Public endpoint should include %s", field))
						Expect(authData[field]).To(Equal(publicData[field]), fmt.Sprintf("Field %s should match between endpoints", field))
					}
				})
			})
		})

		Context("Given user preference data", func() {
			Context("When accessing user preferences", func() {
				It("Then preferences should only be accessible to the user themselves", func() {
					By("Setting up user with specific preferences")
					updateData := map[string]interface{}{
						"preferences": map[string]interface{}{
							"private_setting": "secret_value",
							"theme":           "dark",
						},
					}
					bodyBytes, err := json.Marshal(updateData)
					Expect(err).NotTo(HaveOccurred())

					// Update user preferences
					updateReq := httptest.NewRequest(http.MethodPut, "/api/v1/users/me", bytes.NewReader(bodyBytes))
					updateReq.Header.Set("Content-Type", "application/json")
					updateReq.Header.Set("Authorization", "Bearer "+authData.ValidToken)
					updateResp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(updateResp, updateReq)
					Expect(updateResp.Code).To(Equal(http.StatusOK))

					By("Verifying preferences are not exposed in public endpoint")
					publicReq := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/users/%s", authData.ValidUserID), nil)
					publicResp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(publicResp, publicReq)

					var publicData map[string]interface{}
					err = json.Unmarshal(publicResp.Body.Bytes(), &publicData)
					Expect(err).NotTo(HaveOccurred())

					Expect(publicData["preferences"]).To(BeEmpty(), "Preferences should not be exposed in public endpoint")

					By("Verifying preferences are accessible to authenticated user")
					authReq := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
					authReq.Header.Set("Authorization", "Bearer "+authData.ValidToken)
					authResp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(authResp, authReq)

					var authData map[string]interface{}
					err = json.Unmarshal(authResp.Body.Bytes(), &authData)
					Expect(err).NotTo(HaveOccurred())

					preferences, exists := authData["preferences"]
					Expect(exists).To(BeTrue(), "Preferences should be present for authenticated user")
					prefMap, ok := preferences.(map[string]interface{})
					Expect(ok).To(BeTrue(), "Preferences should be a map")
					Expect(prefMap["private_setting"]).To(Equal("secret_value"))
					Expect(prefMap["theme"]).To(Equal("dark"))
				})
			})
		})
	})

	Describe("Error Handling for User Operations", func() {
		Context("Given various error conditions", func() {
			Context("When handling network or database errors", func() {
				It("Then appropriate error responses should be returned", func() {
					By("Testing with non-existent user ID in update operation")
					updateData := map[string]interface{}{
						"display_name": "Won't Work",
					}
					bodyBytes, err := json.Marshal(updateData)
					Expect(err).NotTo(HaveOccurred())

					req := httptest.NewRequest(http.MethodPut, "/api/v1/users/nonexistent-user-id", bytes.NewReader(bodyBytes))
					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("Authorization", "Bearer "+authData.ValidToken)

					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)

					// This should fail with forbidden (since user tries to update another user)
					// or not found if the permission check passes
					Expect(resp.Code).To(BeNumerically(">=", 400), "Should return error for nonexistent user")
				})
			})

			Context("When handling malformed request data", func() {
				It("Then validation errors should be returned", func() {
					By("Sending request with malformed JSON")
					malformedJSON := `{"display_name": "test"`

					req := httptest.NewRequest(http.MethodPut, "/api/v1/users/me", bytes.NewReader([]byte(malformedJSON)))
					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("Authorization", "Bearer "+authData.ValidToken)

					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)

					By("Verifying appropriate error handling")
					Expect(resp.Code).To(Equal(http.StatusBadRequest), "Should return bad request for malformed JSON")
				})
			})

			Context("When handling missing required fields", func() {
				It("Then appropriate validation errors should be returned", func() {
					By("Testing with excessively long display name")
					updateData := map[string]interface{}{
						"display_name": string(make([]byte, 300)), // Exceeds 255 character limit
					}
					bodyBytes, err := json.Marshal(updateData)
					Expect(err).NotTo(HaveOccurred())

					req := httptest.NewRequest(http.MethodPut, "/api/v1/users/me", bytes.NewReader(bodyBytes))
					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("Authorization", "Bearer "+authData.ValidToken)

					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)

					By("Verifying validation error")
					verifyErrorResponse(resp, http.StatusBadRequest, "validation")
				})
			})
		})

		Context("Given edge cases in user data", func() {
			Context("When handling special characters in user data", func() {
				It("Then data should be properly handled and validated", func() {
					By("Testing with special characters in display name")
					updateData := map[string]interface{}{
						"display_name": "Test User 🚀 Special & Chars",
					}
					bodyBytes, err := json.Marshal(updateData)
					Expect(err).NotTo(HaveOccurred())

					req := httptest.NewRequest(http.MethodPut, "/api/v1/users/me", bytes.NewReader(bodyBytes))
					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("Authorization", "Bearer "+authData.ValidToken)

					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)

					By("Verifying special characters are handled properly")
					if resp.Code == http.StatusOK {
						var responseBody map[string]interface{}
						err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
						Expect(err).NotTo(HaveOccurred())
						Expect(responseBody["display_name"]).To(Equal("Test User 🚀 Special & Chars"))
					} else {
						// If validation fails, it should be a clear validation error
						verifyErrorResponse(resp, http.StatusBadRequest, "validation")
					}
				})
			})
		})
	})
})
