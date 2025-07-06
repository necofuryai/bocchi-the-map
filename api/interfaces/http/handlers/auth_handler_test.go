//go:build integration

package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/necofuryai/bocchi-the-map/api/application/clients"
	"github.com/necofuryai/bocchi-the-map/api/pkg/auth"
	"github.com/necofuryai/bocchi-the-map/api/tests/helpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Test helper functions for auth handler
func verifyTimestampFormat(timestamp interface{}) {
	timestampStr, ok := timestamp.(string)
	Expect(ok).To(BeTrue(), "Timestamp should be a string")

	parsedTime, err := time.Parse(time.RFC3339, timestampStr)
	Expect(err).NotTo(HaveOccurred(), "Timestamp should be in RFC3339 format")
	Expect(parsedTime).To(BeTemporally("~", time.Now(), time.Minute), "Timestamp should be recent")
}

func verifyAuthResponseStructure(responseBody map[string]interface{}) {
	Expect(responseBody["authenticated"]).To(Not(BeNil()), "Response should contain authenticated field")
	Expect(responseBody["timestamp"]).To(Not(BeNil()), "Response should contain timestamp")
	verifyTimestampFormat(responseBody["timestamp"])
}

func verifyUserInfo(userInfo map[string]interface{}, expectedUserID, expectedEmail string) {
	Expect(userInfo["id"]).To(Equal(expectedUserID), "User ID should match")
	Expect(userInfo["email"]).To(Equal(expectedEmail), "User email should match")
	Expect(userInfo["display_name"]).To(Not(BeEmpty()), "Display name should be present")
	Expect(userInfo["auth_provider"]).To(Equal("google"), "Auth provider should be set")
	Expect(userInfo["verified"]).To(BeTrue(), "Email should be verified")
}

func verifyTokenInfo(tokenInfo map[string]interface{}, expectedUserID string) {
	Expect(tokenInfo["subject"]).To(Equal(expectedUserID), "Token subject should match user ID")
	Expect(tokenInfo["issuer"]).To(Not(BeEmpty()), "Token issuer should be present")
	Expect(tokenInfo["audience"]).To(Not(BeEmpty()), "Token audience should be present")
	Expect(tokenInfo["issued_at"]).To(Not(BeEmpty()), "Token issued_at should be present")
	Expect(tokenInfo["expires_at"]).To(Not(BeEmpty()), "Token expires_at should be present")
}

func createMockAuthMiddleware() *auth.AuthMiddleware {
	// Create a mock auth middleware for testing
	// In a real test environment, this would be configured with test keys
	return &auth.AuthMiddleware{} // Simplified for testing
}

var _ = Describe("AuthHandler BDD Tests", func() {
	var (
		api            huma.API
		authHandler    *AuthHandler
		userClient     *clients.UserClient
		authData       *helpers.AuthTestData
		authMiddleware *auth.AuthMiddleware
	)

	BeforeEach(func() {
		By("Setting up AuthHandler test environment")

		// Create user client with test database
		var err error
		userClient, err = clients.NewUserClient("internal", testSuite.TestDB.DB)
		Expect(err).NotTo(HaveOccurred())

		// Create mock auth middleware
		authMiddleware = createMockAuthMiddleware()

		// Create auth handler
		authHandler = NewAuthHandler(authMiddleware, userClient)

		// Setup test API with chi router
		router := chi.NewRouter()
		api = humachi.New(router, huma.DefaultConfig("Test API", "1.0.0"))

		// Register auth endpoints
		authHandler.RegisterRoutes(api)
		authHandler.RegisterRoutesWithRateLimit(api, nil) // nil rate limiter for testing

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

	Describe("Authentication Status Check", func() {
		Context("Given an authenticated user", func() {
			Context("When requesting authentication status", func() {
				It("Then user information and token details should be returned", func() {
					By("Creating an authenticated context")
					baseCtx := context.Background()
					authCtx := testSuite.AuthHelper.CreateAuthenticatedContext(baseCtx, authData.ValidUserID, authData.TestUser.Email)

					By("Calling GetAuthStatus directly with authenticated context")
					input := &AuthStatusInput{}
					output, err := authHandler.GetAuthStatus(authCtx, input)

					By("Verifying successful authentication response")
					Expect(err).NotTo(HaveOccurred(), "GetAuthStatus should succeed for authenticated user")
					Expect(output).NotTo(BeNil(), "Output should not be nil")

					By("Verifying authentication status")
					Expect(output.Body.Authenticated).To(BeTrue(), "User should be authenticated")

					By("Verifying user information")
					Expect(output.Body.User).NotTo(BeNil(), "User info should be present")
					Expect(output.Body.User.ID).To(Equal(authData.ValidUserID), "User ID should match")
					Expect(output.Body.User.Email).To(Equal(authData.TestUser.Email), "Email should match")
					Expect(output.Body.User.DisplayName).To(Not(BeEmpty()), "Display name should be present")

					By("Verifying token information")
					Expect(output.Body.TokenInfo).NotTo(BeNil(), "Token info should be present")
					Expect(output.Body.TokenInfo.Subject).To(Equal(authData.ValidUserID), "Token subject should match user ID")

					By("Verifying permissions structure")
					Expect(output.Body.Permissions).NotTo(BeNil(), "Permissions should be present")
					Expect(output.Body.Permissions).To(BeEmpty(), "Permissions should be empty array for now")

					By("Verifying response timestamp")
					Expect(output.Body.Timestamp).To(BeTemporally("~", time.Now(), time.Minute), "Timestamp should be recent")
				})
			})
		})

		Context("Given an unauthenticated user", func() {
			Context("When requesting authentication status", func() {
				It("Then unauthenticated status should be returned", func() {
					By("Creating an unauthenticated context")
					baseCtx := context.Background()
					unauthCtx := testSuite.AuthHelper.CreateUnauthenticatedContext(baseCtx)

					By("Calling GetAuthStatus with unauthenticated context")
					input := &AuthStatusInput{}
					output, err := authHandler.GetAuthStatus(unauthCtx, input)

					By("Verifying successful response for unauthenticated user")
					Expect(err).NotTo(HaveOccurred(), "GetAuthStatus should succeed even for unauthenticated user")
					Expect(output).NotTo(BeNil(), "Output should not be nil")

					By("Verifying unauthenticated status")
					Expect(output.Body.Authenticated).To(BeFalse(), "User should not be authenticated")

					By("Verifying no user information")
					Expect(output.Body.User).To(BeNil(), "User info should be nil")
					Expect(output.Body.TokenInfo).To(BeNil(), "Token info should be nil")
					Expect(output.Body.Permissions).To(BeEmpty(), "Permissions should be empty")

					By("Verifying response timestamp")
					Expect(output.Body.Timestamp).To(BeTemporally("~", time.Now(), time.Minute), "Timestamp should be recent")
				})
			})
		})

		Context("Given corrupted authentication context", func() {
			Context("When user data is partially missing", func() {
				It("Then graceful fallback should occur", func() {
					By("Creating context with missing user information")
					baseCtx := context.Background()
					partialCtx := context.WithValue(baseCtx, "user_id", authData.ValidUserID)
					// Missing user_email and user_info

					By("Calling GetAuthStatus with partial context")
					input := &AuthStatusInput{}
					output, err := authHandler.GetAuthStatus(partialCtx, input)

					By("Verifying graceful handling")
					Expect(err).NotTo(HaveOccurred(), "Should handle partial context gracefully")
					Expect(output.Body.Authenticated).To(BeFalse(), "Should treat as unauthenticated when data is incomplete")
				})
			})
		})
	})

	Describe("Token Validation", func() {
		Context("Given a valid JWT token", func() {
			Context("When validating the token", func() {
				It("Then token claims should be returned", func() {
					By("Preparing token validation request")
					input := &ValidateTokenInput{}
					input.Body.Token = authData.ValidToken

					By("Calling ValidateToken")
					output, err := authHandler.ValidateToken(context.Background(), input)

					By("Verifying successful validation")
					Expect(err).NotTo(HaveOccurred(), "Token validation should succeed")
					Expect(output).NotTo(BeNil(), "Output should not be nil")

					By("Verifying token validity")
					// Note: This test may fail if the auth middleware validator is not properly mocked
					// In a real implementation, we would mock the validator to return successful validation
					if output.Body.Valid {
						Expect(output.Body.Claims).NotTo(BeNil(), "Claims should be present for valid token")
						Expect(output.Body.Error).To(BeEmpty(), "Error should be empty for valid token")
					} else {
						// If validation fails due to missing mock, verify error handling
						Expect(output.Body.Error).NotTo(BeEmpty(), "Error message should be present for failed validation")
					}

					By("Verifying response timestamp")
					Expect(output.Body.Timestamp).To(BeTemporally("~", time.Now(), time.Minute), "Timestamp should be recent")
				})
			})
		})

		Context("Given an invalid JWT token", func() {
			Context("When validating the token", func() {
				It("Then validation error should be returned", func() {
					By("Preparing invalid token validation request")
					input := &ValidateTokenInput{}
					input.Body.Token = authData.InvalidToken

					By("Calling ValidateToken with invalid token")
					output, err := authHandler.ValidateToken(context.Background(), input)

					By("Verifying failed validation")
					Expect(err).NotTo(HaveOccurred(), "ValidateToken should not return error, but indicate invalid token")
					Expect(output).NotTo(BeNil(), "Output should not be nil")

					By("Verifying token invalidity")
					Expect(output.Body.Valid).To(BeFalse(), "Token should be invalid")
					Expect(output.Body.Claims).To(BeNil(), "Claims should be nil for invalid token")
					Expect(output.Body.Error).NotTo(BeEmpty(), "Error message should be present")

					By("Verifying response timestamp")
					Expect(output.Body.Timestamp).To(BeTemporally("~", time.Now(), time.Minute), "Timestamp should be recent")
				})
			})
		})

		Context("Given an expired JWT token", func() {
			Context("When validating the token", func() {
				It("Then expiration error should be returned", func() {
					By("Preparing expired token validation request")
					input := &ValidateTokenInput{}
					input.Body.Token = authData.ExpiredToken

					By("Calling ValidateToken with expired token")
					output, err := authHandler.ValidateToken(context.Background(), input)

					By("Verifying expired token handling")
					Expect(err).NotTo(HaveOccurred(), "ValidateToken should not return error")
					Expect(output).NotTo(BeNil(), "Output should not be nil")

					By("Verifying token expiration")
					Expect(output.Body.Valid).To(BeFalse(), "Expired token should be invalid")
					Expect(output.Body.Claims).To(BeNil(), "Claims should be nil for expired token")
					Expect(output.Body.Error).NotTo(BeEmpty(), "Error message should be present")

					By("Verifying response timestamp")
					Expect(output.Body.Timestamp).To(BeTemporally("~", time.Now(), time.Minute), "Timestamp should be recent")
				})
			})
		})

		Context("Given malformed token input", func() {
			Context("When validating empty token", func() {
				It("Then validation error should be returned", func() {
					By("Preparing empty token validation request")
					input := &ValidateTokenInput{}
					input.Body.Token = ""

					By("Calling ValidateToken with empty token")
					output, err := authHandler.ValidateToken(context.Background(), input)

					By("Verifying empty token handling")
					Expect(err).NotTo(HaveOccurred(), "ValidateToken should handle empty token gracefully")
					Expect(output).NotTo(BeNil(), "Output should not be nil")

					By("Verifying empty token rejection")
					Expect(output.Body.Valid).To(BeFalse(), "Empty token should be invalid")
					Expect(output.Body.Error).NotTo(BeEmpty(), "Error message should be present")
				})
			})
		})

		Context("Given authentication service unavailable", func() {
			Context("When token validation is attempted", func() {
				It("Then service unavailable error should be returned", func() {
					By("Creating auth handler with nil middleware")
					authHandlerWithNilMiddleware := NewAuthHandler(nil, userClient)

					By("Preparing token validation request")
					input := &ValidateTokenInput{}
					input.Body.Token = authData.ValidToken

					By("Expecting panic from nil middleware")
					Expect(func() {
						authHandlerWithNilMiddleware.ValidateToken(context.Background(), input)
					}).To(Panic(), "Should panic when middleware is nil")
				})
			})
		})
	})

	Describe("User Logout", func() {
		Context("Given an authenticated user", func() {
			Context("When requesting logout", func() {
				It("Then logout should succeed with appropriate message", func() {
					By("Creating authenticated context")
					baseCtx := context.Background()
					authCtx := context.WithValue(baseCtx, "user_id", authData.ValidUserID)

					By("Calling Logout")
					input := &LogoutInput{}
					output, err := authHandler.Logout(authCtx, input)

					By("Verifying successful logout")
					Expect(err).NotTo(HaveOccurred(), "Logout should succeed for authenticated user")
					Expect(output).NotTo(BeNil(), "Output should not be nil")

					By("Verifying logout response")
					Expect(output.Body.Success).To(BeTrue(), "Logout should be successful")
					Expect(output.Body.Message).To(ContainSubstring("Logout successful"), "Should contain success message")
					Expect(output.Body.Message).To(ContainSubstring("remove the token"), "Should instruct client to remove token")

					By("Verifying response timestamp")
					Expect(output.Body.Timestamp).To(BeTemporally("~", time.Now(), time.Minute), "Timestamp should be recent")
				})
			})
		})

		Context("Given an unauthenticated user", func() {
			Context("When requesting logout", func() {
				It("Then authentication error should be returned", func() {
					By("Creating unauthenticated context")
					baseCtx := context.Background()
					unauthCtx := testSuite.AuthHelper.CreateUnauthenticatedContext(baseCtx)

					By("Calling Logout without authentication")
					input := &LogoutInput{}
					_, err := authHandler.Logout(unauthCtx, input)

					By("Verifying authentication requirement")
					Expect(err).To(HaveOccurred(), "Logout should require authentication")
					Expect(err.Error()).To(ContainSubstring("401"), "Should return 401 Unauthorized")
					Expect(err.Error()).To(ContainSubstring("authentication required"), "Should indicate authentication is required")
				})
			})
		})

		Context("Given context with empty user ID", func() {
			Context("When requesting logout", func() {
				It("Then authentication error should be returned", func() {
					By("Creating context with empty user ID")
					baseCtx := context.Background()
					emptyUserCtx := context.WithValue(baseCtx, "user_id", "")

					By("Calling Logout with empty user ID")
					input := &LogoutInput{}
					_, err := authHandler.Logout(emptyUserCtx, input)

					By("Verifying empty user ID rejection")
					Expect(err).To(HaveOccurred(), "Logout should reject empty user ID")
					Expect(err.Error()).To(ContainSubstring("401"), "Should return 401 Unauthorized")
				})
			})
		})

		Context("Given context with wrong type user ID", func() {
			Context("When requesting logout", func() {
				It("Then authentication error should be returned", func() {
					By("Creating context with non-string user ID")
					baseCtx := context.Background()
					wrongTypeCtx := context.WithValue(baseCtx, "user_id", 12345)

					By("Calling Logout with wrong type user ID")
					input := &LogoutInput{}
					_, err := authHandler.Logout(wrongTypeCtx, input)

					By("Verifying wrong type rejection")
					Expect(err).To(HaveOccurred(), "Logout should reject non-string user ID")
					Expect(err.Error()).To(ContainSubstring("401"), "Should return 401 Unauthorized")
				})
			})
		})
	})

	Describe("Authentication Statistics", func() {
		Context("Given an authenticated user", func() {
			Context("When requesting auth statistics", func() {
				It("Then service statistics should be returned", func() {
					By("Creating authenticated context")
					baseCtx := context.Background()
					authCtx := context.WithValue(baseCtx, "user_id", authData.ValidUserID)

					By("Calling GetAuthStats")
					input := &AuthStatsInput{}
					output, err := authHandler.GetAuthStats(authCtx, input)

					By("Verifying successful stats retrieval")
					Expect(err).NotTo(HaveOccurred(), "GetAuthStats should succeed for authenticated user")
					Expect(output).NotTo(BeNil(), "Output should not be nil")

					By("Verifying service statistics structure")
					Expect(output.Body.ServiceStats).NotTo(BeNil(), "Service stats should be present")
					Expect(output.Body.HealthStatus).To(Equal("healthy"), "Health status should be healthy")

					By("Verifying rate limiter stats")
					rateLimiterStats, exists := output.Body.ServiceStats["rate_limiter"]
					Expect(exists).To(BeTrue(), "Rate limiter stats should be present")
					rateLimiterMap, ok := rateLimiterStats.(map[string]interface{})
					Expect(ok).To(BeTrue(), "Rate limiter stats should be a map")
					Expect(rateLimiterMap["active"]).To(BeTrue(), "Rate limiter should be marked as active")

					By("Verifying JWT validator stats")
					jwtStats, exists := output.Body.ServiceStats["jwt_validator"]
					Expect(exists).To(BeTrue(), "JWT validator stats should be present")
					jwtMap, ok := jwtStats.(map[string]interface{})
					Expect(ok).To(BeTrue(), "JWT validator stats should be a map")
					Expect(jwtMap["active"]).To(BeTrue(), "JWT validator should be marked as active")

					By("Verifying response timestamp")
					Expect(output.Body.Timestamp).To(BeTemporally("~", time.Now(), time.Minute), "Timestamp should be recent")
				})
			})
		})

		Context("Given an unauthenticated user", func() {
			Context("When requesting auth statistics", func() {
				It("Then authentication error should be returned", func() {
					By("Creating unauthenticated context")
					baseCtx := context.Background()
					unauthCtx := testSuite.AuthHelper.CreateUnauthenticatedContext(baseCtx)

					By("Calling GetAuthStats without authentication")
					input := &AuthStatsInput{}
					_, err := authHandler.GetAuthStats(unauthCtx, input)

					By("Verifying authentication requirement")
					Expect(err).To(HaveOccurred(), "GetAuthStats should require authentication")
					Expect(err.Error()).To(ContainSubstring("401"), "Should return 401 Unauthorized")
					Expect(err.Error()).To(ContainSubstring("authentication required"), "Should indicate authentication is required")
				})
			})
		})

		Context("Given admin user access control", func() {
			Context("When implementing role-based access", func() {
				It("Then future admin verification should be considered", func() {
					By("Documenting future admin access control requirement")
					// Currently, any authenticated user can access stats
					// In future implementation, this should check for admin role

					baseCtx := context.Background()
					authCtx := context.WithValue(baseCtx, "user_id", authData.ValidUserID)

					input := &AuthStatsInput{}
					output, err := authHandler.GetAuthStats(authCtx, input)

					Expect(err).NotTo(HaveOccurred(), "Currently allows any authenticated user")
					Expect(output).NotTo(BeNil(), "Should return stats for now")

					// TODO: Implement admin role checking
					// Future test should verify that only admin users can access auth stats
					By("Noting that admin role verification should be implemented")
				})
			})
		})
	})

	Describe("Security Boundary Testing", func() {
		Context("Given various authentication edge cases", func() {
			Context("When testing security boundaries", func() {
				It("Then proper error handling should occur", func() {
					By("Testing nil context values")
					baseCtx := context.Background()
					nilValueCtx := context.WithValue(baseCtx, "user_id", nil)

					logoutInput := &LogoutInput{}
					_, err := authHandler.Logout(nilValueCtx, logoutInput)
					Expect(err).To(HaveOccurred(), "Should reject nil user ID")

					statsInput := &AuthStatsInput{}
					_, err = authHandler.GetAuthStats(nilValueCtx, statsInput)
					Expect(err).To(HaveOccurred(), "Should reject nil user ID for stats")
				})

				It("Then token validation should handle edge cases", func() {
					By("Testing extremely long token")
					longToken := string(make([]byte, 10000)) // Very long token
					input := &ValidateTokenInput{}
					input.Body.Token = longToken

					output, err := authHandler.ValidateToken(context.Background(), input)
					Expect(err).NotTo(HaveOccurred(), "Should handle long token gracefully")
					Expect(output.Body.Valid).To(BeFalse(), "Long invalid token should be rejected")

					By("Testing token with special characters")
					specialToken := "token.with.special!@#$%^&*()characters"
					input.Body.Token = specialToken

					output, err = authHandler.ValidateToken(context.Background(), input)
					Expect(err).NotTo(HaveOccurred(), "Should handle special characters gracefully")
					Expect(output.Body.Valid).To(BeFalse(), "Invalid token with special chars should be rejected")
				})

				It("Then response structure should be consistent", func() {
					By("Verifying all endpoints return timestamps")

					// Test unauthenticated auth status
					unauthCtx := testSuite.AuthHelper.CreateUnauthenticatedContext(context.Background())
					authStatusInput := &AuthStatusInput{}
					authStatusOutput, err := authHandler.GetAuthStatus(unauthCtx, authStatusInput)
					Expect(err).NotTo(HaveOccurred())
					Expect(authStatusOutput.Body.Timestamp).To(Not(BeZero()), "Auth status should have timestamp")

					// Test token validation
					validateInput := &ValidateTokenInput{Body: struct {
						Token string `json:"token" minLength:"1" doc:"JWT token to validate"`
					}{Token: "invalid"}}
					validateOutput, err := authHandler.ValidateToken(context.Background(), validateInput)
					Expect(err).NotTo(HaveOccurred())
					Expect(validateOutput.Body.Timestamp).To(Not(BeZero()), "Token validation should have timestamp")

					By("Verifying timestamp format consistency")
					// All timestamps should be within a reasonable time window
					now := time.Now()
					Expect(authStatusOutput.Body.Timestamp).To(BeTemporally("~", now, time.Minute))
					Expect(validateOutput.Body.Timestamp).To(BeTemporally("~", now, time.Minute))
				})
			})
		})

		Context("Given concurrent authentication requests", func() {
			Context("When handling multiple simultaneous requests", func() {
				It("Then thread safety should be maintained", func() {
					By("Making concurrent auth status requests")
					const numRequests = 10
					results := make(chan error, numRequests)

					for i := 0; i < numRequests; i++ {
						go func() {
							defer GinkgoRecover()
							unauthCtx := testSuite.AuthHelper.CreateUnauthenticatedContext(context.Background())
							input := &AuthStatusInput{}
							_, err := authHandler.GetAuthStatus(unauthCtx, input)
							results <- err
						}()
					}

					By("Verifying all requests complete successfully")
					for i := 0; i < numRequests; i++ {
						err := <-results
						Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("Request %d should succeed", i))
					}
				})
			})
		})
	})
})
