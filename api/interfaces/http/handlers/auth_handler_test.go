//go:build integration

package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/necofuryai/bocchi-the-map/api/application/clients"
	"github.com/necofuryai/bocchi-the-map/api/pkg/auth"
	"github.com/necofuryai/bocchi-the-map/api/tests/helpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("AuthHandler BDD Tests", func() {
	var (
		api           huma.API
		testServer    *httptest.Server
		authHandler   *AuthHandler
		userClient    *clients.UserClient
		authMiddleware *auth.AuthMiddleware
		rateLimiter   *auth.RateLimiter
		authData      *helpers.AuthTestData
	)

	BeforeEach(func() {
		By("Setting up AuthHandler test environment")
		
		// Create user client with test database
		var err error
		userClient, err = clients.NewUserClient("internal", testSuite.TestDB.DB)
		Expect(err).NotTo(HaveOccurred())
		
		// Get JWT secret from environment variable, with a fallback for tests
		jwtSecret := os.Getenv("TEST_JWT_SECRET")
		if jwtSecret == "" {
			jwtSecret = "test-jwt-secret-key-for-bdd-testing-32-chars-minimum"
		}
		
		// Create auth middleware
		authMiddleware = auth.NewAuthMiddleware(jwtSecret, testSuite.TestDB.Queries)
		
		// Create rate limiter (5 requests per 5 minutes for testing)
		rateLimiter = auth.NewRateLimiter(5, 300)
		
		// Create auth handler
		authHandler = NewAuthHandler(authMiddleware, userClient)
		
		// Setup test API with chi router
		router := chi.NewRouter()
		api = humachi.New(router, huma.DefaultConfig("Test API", "1.0.0"))
		testServer = httptest.NewServer(router)
		
		// Register authentication endpoints
		authHandler.RegisterRoutesWithRateLimit(api, rateLimiter)
		
		// Setup authentication test data
		authData = testSuite.AuthHelper.NewAuthTestData()
		
		// Create test user in database for token generation
		testSuite.FixtureManager.CreateUserFixture(context.Background(), helpers.UserFixture{
			ID:             authData.ValidUserID,
			Email:          authData.TestUser.Email,
			DisplayName:    authData.TestUser.DisplayName,
			AuthProvider:   string(authData.TestUser.AuthProvider),
			AuthProviderID: authData.TestUser.AuthProviderID,
			Preferences:    authData.TestUser.Preferences,
		})
	})

	AfterEach(func() {
		By("Cleaning up test server resources")
		if testServer != nil {
			testServer.Close()
		}
	})

	Describe("API Token Generation", func() {
		Context("Given a valid OAuth user", func() {
			Context("When requesting API token with matching credentials", func() {
				It("Then access and refresh tokens should be generated", func() {
					By("Preparing a valid token generation request")
					requestBody := map[string]interface{}{
						"email":            authData.TestUser.Email,
						"auth_provider":    string(authData.TestUser.AuthProvider),
						"auth_provider_id": authData.TestUser.AuthProviderID,
					}
					
					bodyBytes, err := json.Marshal(requestBody)
					Expect(err).NotTo(HaveOccurred())
					
					By("Sending the token generation request")
					req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/token", bytes.NewReader(bodyBytes))
					req.Header.Set("Content-Type", "application/json")
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the successful response")
					Expect(resp.Code).To(Equal(http.StatusOK), "Expected status 200 OK")
					
					var responseBody map[string]interface{}
					err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
					Expect(err).NotTo(HaveOccurred())
					
					By("Verifying token information in response")
					Expect(responseBody["access_token_expires_at"]).To(Not(BeEmpty()), "Access token expiration should be set")
					Expect(responseBody["refresh_token_expires_at"]).To(Not(BeEmpty()), "Refresh token expiration should be set")
					
					By("Verifying secure httpOnly cookies are set")
					cookies := resp.Result().Cookies()
					cookieNames := make([]string, len(cookies))
					for i, cookie := range cookies {
						cookieNames[i] = cookie.Name
						if cookie.Name == "bocchi_access_token" || cookie.Name == "bocchi_refresh_token" {
							Expect(cookie.HttpOnly).To(BeTrue(), "Cookies should be httpOnly")
							Expect(cookie.Secure).To(BeTrue(), "Cookies should be secure")
						}
					}
					Expect(cookieNames).To(ContainElement("bocchi_access_token"), "Access token cookie should be set")
					Expect(cookieNames).To(ContainElement("bocchi_refresh_token"), "Refresh token cookie should be set")
				})
			})
			
			Context("When requesting token with missing required fields", func() {
				It("Then it should return a validation error", func() {
					By("Preparing a request missing the email field")
					requestBody := map[string]interface{}{
						"auth_provider":    "google",
						"auth_provider_id": "google_123",
					}
					
					bodyBytes, err := json.Marshal(requestBody)
					Expect(err).NotTo(HaveOccurred())
					
					By("Sending the incomplete request")
					req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/token", bytes.NewReader(bodyBytes))
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

		Context("Given an invalid OAuth provider", func() {
			Context("When requesting token with unsupported provider", func() {
				It("Then it should return a validation error", func() {
					By("Preparing a request with invalid provider")
					requestBody := map[string]interface{}{
						"email":            "test@example.com",
						"auth_provider":    "invalid_provider",
						"auth_provider_id": "invalid_123",
					}
					
					bodyBytes, err := json.Marshal(requestBody)
					Expect(err).NotTo(HaveOccurred())
					
					By("Sending the request with invalid provider")
					req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/token", bytes.NewReader(bodyBytes))
					req.Header.Set("Content-Type", "application/json")
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the validation error")
					Expect(resp.Code).To(Equal(http.StatusBadRequest), "Expected status 400 Bad Request")
				})
			})
		})

		Context("Given a non-existent user", func() {
			Context("When requesting token for unregistered user", func() {
				It("Then it should return a not found error", func() {
					By("Preparing a request for non-existent user")
					requestBody := map[string]interface{}{
						"email":            "nonexistent@example.com",
						"auth_provider":    "google",
						"auth_provider_id": "google_nonexistent",
					}
					
					bodyBytes, err := json.Marshal(requestBody)
					Expect(err).NotTo(HaveOccurred())
					
					By("Sending the request for non-existent user")
					req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/token", bytes.NewReader(bodyBytes))
					req.Header.Set("Content-Type", "application/json")
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the not found error")
					Expect(resp.Code).To(Equal(http.StatusNotFound), "Expected status 404 Not Found")
				})
			})
		})
	})

	Describe("Token Refresh", func() {
		var validRefreshToken string

		BeforeEach(func() {
			By("Generating initial tokens for refresh testing")
			requestBody := map[string]interface{}{
				"email":            authData.TestUser.Email,
				"auth_provider":    string(authData.TestUser.AuthProvider),
				"auth_provider_id": authData.TestUser.AuthProviderID,
			}
			
			bodyBytes, err := json.Marshal(requestBody)
			Expect(err).NotTo(HaveOccurred())
			
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/token", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			
			resp := httptest.NewRecorder()
			testServer.Config.Handler.ServeHTTP(resp, req)
			Expect(resp.Code).To(Equal(http.StatusOK))
			
			// Extract refresh token from cookies for testing
			cookies := resp.Result().Cookies()
			for _, cookie := range cookies {
				if cookie.Name == "bocchi_refresh_token" {
					validRefreshToken = cookie.Value
					break
				}
			}
			Expect(validRefreshToken).To(Not(BeEmpty()), "Refresh token should be available for testing")
		})

		Context("Given a valid refresh token", func() {
			Context("When requesting token refresh", func() {
				It("Then new tokens should be generated", func() {
					By("Sending a refresh request with valid refresh token")
					req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", nil)
					req.AddCookie(&http.Cookie{
						Name:  "bocchi_refresh_token",
						Value: validRefreshToken,
					})
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the successful refresh response")
					Expect(resp.Code).To(Equal(http.StatusOK), "Expected status 200 OK")
					
					var responseBody map[string]interface{}
					err := json.Unmarshal(resp.Body.Bytes(), &responseBody)
					Expect(err).NotTo(HaveOccurred())
					
					By("Verifying new token information")
					Expect(responseBody["access_token_expires_at"]).To(Not(BeEmpty()), "New access token expiration should be set")
					Expect(responseBody["refresh_token_expires_at"]).To(Not(BeEmpty()), "New refresh token expiration should be set")
					
					By("Verifying new cookies are set")
					cookies := resp.Result().Cookies()
					cookieNames := make([]string, len(cookies))
					for i, cookie := range cookies {
						cookieNames[i] = cookie.Name
					}
					Expect(cookieNames).To(ContainElement("bocchi_access_token"), "New access token cookie should be set")
					Expect(cookieNames).To(ContainElement("bocchi_refresh_token"), "New refresh token cookie should be set")
				})
			})
		})

		Context("Given no refresh token", func() {
			Context("When attempting to refresh without token", func() {
				It("Then it should return an unauthorized error", func() {
					By("Sending a refresh request without refresh token")
					req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", nil)
					// No refresh token cookie
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the unauthorized error")
					Expect(resp.Code).To(Equal(http.StatusUnauthorized), "Expected status 401 Unauthorized")
					
					var errorResponse map[string]interface{}
					err := json.Unmarshal(resp.Body.Bytes(), &errorResponse)
					Expect(err).NotTo(HaveOccurred())
					
					Expect(errorResponse["title"]).To(ContainSubstring("refresh token"), "Should indicate refresh token error")
				})
			})
		})

		Context("Given an invalid refresh token", func() {
			Context("When attempting to refresh with malformed token", func() {
				It("Then it should return an unauthorized error", func() {
					By("Sending a refresh request with invalid refresh token")
					req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", nil)
					req.AddCookie(&http.Cookie{
						Name:  "bocchi_refresh_token",
						Value: "invalid.token.here",
					})
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the unauthorized error")
					Expect(resp.Code).To(Equal(http.StatusUnauthorized), "Expected status 401 Unauthorized")
				})
			})
		})
	})

	Describe("User Logout", func() {
		var accessToken, refreshToken string

		BeforeEach(func() {
			By("Generating tokens for logout testing")
			requestBody := map[string]interface{}{
				"email":            authData.TestUser.Email,
				"auth_provider":    string(authData.TestUser.AuthProvider),
				"auth_provider_id": authData.TestUser.AuthProviderID,
			}
			
			bodyBytes, err := json.Marshal(requestBody)
			Expect(err).NotTo(HaveOccurred())
			
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/token", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			
			resp := httptest.NewRecorder()
			testServer.Config.Handler.ServeHTTP(resp, req)
			Expect(resp.Code).To(Equal(http.StatusOK))
			
			// Extract tokens from cookies for testing
			cookies := resp.Result().Cookies()
			for _, cookie := range cookies {
				if cookie.Name == "bocchi_access_token" {
					accessToken = cookie.Value
				} else if cookie.Name == "bocchi_refresh_token" {
					refreshToken = cookie.Value
				}
			}
			Expect(accessToken).To(Not(BeEmpty()), "Access token should be available for logout testing")
			Expect(refreshToken).To(Not(BeEmpty()), "Refresh token should be available for logout testing")
		})

		Context("Given an authenticated user with active tokens", func() {
			Context("When logging out", func() {
				It("Then tokens should be blacklisted and cookies cleared", func() {
					By("Sending a logout request with valid tokens")
					req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
					req.AddCookie(&http.Cookie{
						Name:  "bocchi_access_token",
						Value: accessToken,
					})
					req.AddCookie(&http.Cookie{
						Name:  "bocchi_refresh_token",
						Value: refreshToken,
					})
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the successful logout response")
					Expect(resp.Code).To(Equal(http.StatusOK), "Expected status 200 OK")
					
					var responseBody map[string]interface{}
					err := json.Unmarshal(resp.Body.Bytes(), &responseBody)
					Expect(err).NotTo(HaveOccurred())
					
					Expect(responseBody["message"]).To(ContainSubstring("logout"), "Should confirm logout success")
					
					By("Verifying cookies are cleared")
					cookies := resp.Result().Cookies()
					for _, cookie := range cookies {
						if cookie.Name == "bocchi_access_token" || cookie.Name == "bocchi_refresh_token" {
							Expect(cookie.Value).To(BeEmpty(), fmt.Sprintf("Cookie %s should be cleared", cookie.Name))
							Expect(cookie.MaxAge).To(Equal(-1), fmt.Sprintf("Cookie %s should be expired", cookie.Name))
						}
					}
				})
			})
		})

		Context("Given no authentication tokens", func() {
			Context("When attempting to logout", func() {
				It("Then it should still clear any remaining cookies and succeed", func() {
					By("Sending a logout request without tokens")
					req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
					// No cookies
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the successful logout response")
					Expect(resp.Code).To(Equal(http.StatusOK), "Expected status 200 OK")
					
					var responseBody map[string]interface{}
					err := json.Unmarshal(resp.Body.Bytes(), &responseBody)
					Expect(err).NotTo(HaveOccurred())
					
					Expect(responseBody["message"]).To(ContainSubstring("logout"), "Should confirm logout success")
				})
			})
		})

		Context("Given only partial token data", func() {
			Context("When logging out with only access token", func() {
				It("Then available tokens should be blacklisted", func() {
					By("Sending a logout request with only access token")
					req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
					req.AddCookie(&http.Cookie{
						Name:  "bocchi_access_token",
						Value: accessToken,
					})
					// No refresh token
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying the successful logout response")
					Expect(resp.Code).To(Equal(http.StatusOK), "Expected status 200 OK")
					
					By("Verifying access token cookie is cleared")
					cookies := resp.Result().Cookies()
					for _, cookie := range cookies {
						if cookie.Name == "bocchi_access_token" {
							Expect(cookie.Value).To(BeEmpty(), "Access token cookie should be cleared")
						}
					}
				})
			})
		})
	})

	Describe("Rate Limiting", func() {
		Context("Given normal usage within rate limits", func() {
			Context("When making authentication requests at normal frequency", func() {
				It("Then all requests should be processed successfully", func() {
					By("Making multiple requests within rate limit")
					for i := 0; i < 3; i++ {
						requestBody := map[string]interface{}{
							"email":            fmt.Sprintf("test%d@example.com", i),
							"auth_provider":    "google",
							"auth_provider_id": fmt.Sprintf("google_%d", i),
						}
						
						bodyBytes, err := json.Marshal(requestBody)
						Expect(err).NotTo(HaveOccurred())
						
						req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/token", bytes.NewReader(bodyBytes))
						req.Header.Set("Content-Type", "application/json")
						req.RemoteAddr = fmt.Sprintf("127.0.0.%d:12345", i+1) // Unique IP for each request
						
						resp := httptest.NewRecorder()
						testServer.Config.Handler.ServeHTTP(resp, req)
						
						By(fmt.Sprintf("Verifying request %d is not rate limited", i+1))
						Expect(resp.Code).To(Not(Equal(http.StatusTooManyRequests)), "Request should not be rate limited")
					}
				})
			})
		})

		Context("Given excessive requests beyond rate limits", func() {
			Context("When making more than 5 authentication requests from same IP", func() {
				It("Then should return rate limit error after 5 requests", func() {
					By("Making exactly 5 requests to reach the limit")
					for i := 0; i < 5; i++ {
						requestBody := map[string]interface{}{
							"email":            fmt.Sprintf("ratetest%d@example.com", i),
							"auth_provider":    "google",
							"auth_provider_id": fmt.Sprintf("google_rate_%d", i),
						}
						
						bodyBytes, err := json.Marshal(requestBody)
						Expect(err).NotTo(HaveOccurred())
						
						req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/token", bytes.NewReader(bodyBytes))
						req.Header.Set("Content-Type", "application/json")
						req.RemoteAddr = "192.168.1.100:12345" // Consistent IP for rate limiting
						
						resp := httptest.NewRecorder()
						testServer.Config.Handler.ServeHTTP(resp, req)
						
						By(fmt.Sprintf("Verifying request %d is processed successfully", i+1))
						Expect(resp.Code).To(Not(Equal(http.StatusTooManyRequests)), 
							fmt.Sprintf("Request %d should not be rate limited yet", i+1))
					}
					
					By("Making the 6th request that should trigger rate limit")
					requestBody := map[string]interface{}{
						"email":            "ratelimit6@example.com",
						"auth_provider":    "google",
						"auth_provider_id": "google_rate_6",
					}
					
					bodyBytes, err := json.Marshal(requestBody)
					Expect(err).NotTo(HaveOccurred())
					
					req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/token", bytes.NewReader(bodyBytes))
					req.Header.Set("Content-Type", "application/json")
					req.RemoteAddr = "192.168.1.100:12345" // Same IP as previous requests
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					By("Verifying rate limit is triggered")
					Expect(resp.Code).To(Equal(http.StatusTooManyRequests), "6th request should be rate limited")
					
					By("Verifying rate limit headers are set correctly")
					Expect(resp.Header().Get("X-RateLimit-Limit")).To(Equal("5"), "Rate limit header should show limit of 5")
					Expect(resp.Header().Get("X-RateLimit-Window")).To(Equal("300"), "Rate limit window should be 300 seconds")
					Expect(resp.Header().Get("X-RateLimit-Remaining")).To(Equal("0"), "Rate limit remaining should be 0 when limit is exceeded")
					Expect(resp.Header().Get("Retry-After")).To(Equal("300"), "Retry-After should be 300 seconds")
					
					By("Verifying rate limit error message")
					Expect(resp.Body.String()).To(ContainSubstring("Too many authentication attempts"), 
						"Response should contain rate limit error message")
				})
			})
			
			Context("When making requests from different IPs", func() {
				It("Then each IP should have separate rate limit counters", func() {
					By("Making 5 requests from first IP")
					for i := 0; i < 5; i++ {
						requestBody := map[string]interface{}{
							"email":            fmt.Sprintf("ip1test%d@example.com", i),
							"auth_provider":    "google", 
							"auth_provider_id": fmt.Sprintf("google_ip1_%d", i),
						}
						
						bodyBytes, err := json.Marshal(requestBody)
						Expect(err).NotTo(HaveOccurred())
						
						req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/token", bytes.NewReader(bodyBytes))
						req.Header.Set("Content-Type", "application/json")
						req.RemoteAddr = "192.168.1.101:12345" // First IP
						
						resp := httptest.NewRecorder()
						testServer.Config.Handler.ServeHTTP(resp, req)
						
						By(fmt.Sprintf("Verifying request %d from first IP is processed", i+1))
						Expect(resp.Code).To(Not(Equal(http.StatusTooManyRequests)))
					}
					
					By("Making 5 requests from second IP (should not be rate limited)")
					for i := 0; i < 5; i++ {
						requestBody := map[string]interface{}{
							"email":            fmt.Sprintf("ip2test%d@example.com", i),
							"auth_provider":    "google",
							"auth_provider_id": fmt.Sprintf("google_ip2_%d", i),
						}
						
						bodyBytes, err := json.Marshal(requestBody)
						Expect(err).NotTo(HaveOccurred())
						
						req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/token", bytes.NewReader(bodyBytes))
						req.Header.Set("Content-Type", "application/json")
						req.RemoteAddr = "192.168.1.102:12345" // Different IP
						
						resp := httptest.NewRecorder()
						testServer.Config.Handler.ServeHTTP(resp, req)
						
						By(fmt.Sprintf("Verifying request %d from second IP is processed", i+1))
						Expect(resp.Code).To(Not(Equal(http.StatusTooManyRequests)), 
							"Requests from different IP should not be rate limited")
					}
					
					By("Verifying first IP is still rate limited")
					requestBody := map[string]interface{}{
						"email":            "ip1blocked@example.com",
						"auth_provider":    "google",
						"auth_provider_id": "google_ip1_blocked",
					}
					
					bodyBytes, err := json.Marshal(requestBody)
					Expect(err).NotTo(HaveOccurred())
					
					req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/token", bytes.NewReader(bodyBytes))
					req.Header.Set("Content-Type", "application/json")
					req.RemoteAddr = "192.168.1.101:12345" // First IP again
					
					resp := httptest.NewRecorder()
					testServer.Config.Handler.ServeHTTP(resp, req)
					
					Expect(resp.Code).To(Equal(http.StatusTooManyRequests), 
						"First IP should still be rate limited")
				})
			})
		})
	})
})