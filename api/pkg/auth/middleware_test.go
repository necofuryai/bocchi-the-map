package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/necofuryai/bocchi-the-map/api/infrastructure/database"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Mock TokenBlacklistQuerier for testing
type mockTokenBlacklistQuerier struct {
	isBlacklistedFunc           func(ctx context.Context, jti string) (bool, error)
	blacklistAccessTokenFunc   func(ctx context.Context, arg database.BlacklistAccessTokenParams) error
	blacklistRefreshTokenFunc  func(ctx context.Context, arg database.BlacklistRefreshTokenParams) error
}

func (m *mockTokenBlacklistQuerier) IsTokenBlacklisted(ctx context.Context, jti string) (bool, error) {
	if m.isBlacklistedFunc != nil {
		return m.isBlacklistedFunc(ctx, jti)
	}
	return false, nil
}

func (m *mockTokenBlacklistQuerier) BlacklistAccessToken(ctx context.Context, arg database.BlacklistAccessTokenParams) error {
	if m.blacklistAccessTokenFunc != nil {
		return m.blacklistAccessTokenFunc(ctx, arg)
	}
	return nil
}

func (m *mockTokenBlacklistQuerier) BlacklistRefreshToken(ctx context.Context, arg database.BlacklistRefreshTokenParams) error {
	if m.blacklistRefreshTokenFunc != nil {
		return m.blacklistRefreshTokenFunc(ctx, arg)
	}
	return nil
}

// TestUserData holds common test data for JWT middleware tests
type TestUserData struct {
	UserID string
	Email  string
}

// NewTestUserData creates a new TestUserData instance with default test values
func NewTestUserData() TestUserData {
	return TestUserData{
		UserID: "test-user-id",
		Email:  "test@example.com",
	}
}

func TestJWTAuthMiddleware(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Auth Middleware")
}

var _ = Describe("JWT Token Generation and Claims", func() {
	var (
		middleware *AuthMiddleware
		mockQuerier *mockTokenBlacklistQuerier
		testData TestUserData
	)

	BeforeEach(func() {
		By("Setting up test environment")
		mockQuerier = &mockTokenBlacklistQuerier{}
		middleware = NewAuthMiddleware("test-secret", mockQuerier)
		testData = NewTestUserData()
	})

	Describe("Access token generation", func() {
		Context("When generating an access token with valid user data", func() {
			It("should create a token with correct claims and token type", func() {
				By("Generating the access token")
				tokenString, err := middleware.GenerateToken(testData.UserID, testData.Email)
				Expect(err).NotTo(HaveOccurred())
				Expect(tokenString).NotTo(BeEmpty())

				By("Parsing and verifying token claims")
				token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
					return []byte("test-secret"), nil
				})
				Expect(err).NotTo(HaveOccurred())

				By("Verifying the claims contain correct data")
				claims, ok := token.Claims.(*JWTClaims)
				Expect(ok).To(BeTrue())
				Expect(claims.TokenType).To(Equal("access"))
				Expect(claims.UserID).To(Equal(testData.UserID))
				Expect(claims.Email).To(Equal(testData.Email))
			})
		})
	})

	Describe("Refresh token generation", func() {
		Context("When generating a refresh token with valid user data", func() {
			It("should create a token with correct claims and token type", func() {
				By("Generating the refresh token")
				tokenString, err := middleware.GenerateRefreshToken(testData.UserID, testData.Email)
				Expect(err).NotTo(HaveOccurred())
				Expect(tokenString).NotTo(BeEmpty())

				By("Parsing and verifying token claims")
				token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
					return []byte("test-secret"), nil
				})
				Expect(err).NotTo(HaveOccurred())

				By("Verifying the claims contain correct data")
				claims, ok := token.Claims.(*JWTClaims)
				Expect(ok).To(BeTrue())
				Expect(claims.TokenType).To(Equal("refresh"))
				Expect(claims.UserID).To(Equal(testData.UserID))
				Expect(claims.Email).To(Equal(testData.Email))
			})
		})
	})
})

var _ = Describe("Token Blacklisting", func() {
	var (
		middleware *AuthMiddleware
		mockQuerier *mockTokenBlacklistQuerier
		testData TestUserData
		accessTokenBlacklisted bool
		refreshTokenBlacklisted bool
	)

	BeforeEach(func() {
		By("Setting up blacklist test environment")
		accessTokenBlacklisted = false
		refreshTokenBlacklisted = false
		testData = NewTestUserData()

		mockQuerier = &mockTokenBlacklistQuerier{
			blacklistAccessTokenFunc: func(ctx context.Context, arg database.BlacklistAccessTokenParams) error {
				accessTokenBlacklisted = true
				return nil
			},
			blacklistRefreshTokenFunc: func(ctx context.Context, arg database.BlacklistRefreshTokenParams) error {
				refreshTokenBlacklisted = true
				return nil
			},
		}
		middleware = NewAuthMiddleware("test-secret", mockQuerier)
	})

	Describe("Blacklisting access tokens", func() {
		Context("When blacklisting a valid access token", func() {
			It("should call the access token blacklist method", func() {
				By("Generating an access token")
				accessToken, err := middleware.GenerateToken(testData.UserID, testData.Email)
				Expect(err).NotTo(HaveOccurred())

				By("Blacklisting the access token")
				err = middleware.BlacklistToken(context.Background(), accessToken, "test")
				Expect(err).NotTo(HaveOccurred())

				By("Verifying that access token blacklist was called")
				Expect(accessTokenBlacklisted).To(BeTrue())
				Expect(refreshTokenBlacklisted).To(BeFalse())
			})
		})

		Context("When database error occurs during access token blacklisting", func() {
			It("should return the error from blacklistAccessTokenFunc", func() {
				By("Setting up mock to return database error")
				expectedError := errors.New("database connection failed")
				mockQuerier.blacklistAccessTokenFunc = func(ctx context.Context, arg database.BlacklistAccessTokenParams) error {
					return expectedError
				}

				By("Generating an access token")
				accessToken, err := middleware.GenerateToken(testData.UserID, testData.Email)
				Expect(err).NotTo(HaveOccurred())

				By("Attempting to blacklist the access token")
				err = middleware.BlacklistToken(context.Background(), accessToken, "test")
				
				By("Verifying the database error is returned")
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(expectedError))
			})
		})

		Context("When network error occurs during access token blacklisting", func() {
			It("should return the network error from blacklistAccessTokenFunc", func() {
				By("Setting up mock to return network error")
				expectedError := errors.New("network timeout")
				mockQuerier.blacklistAccessTokenFunc = func(ctx context.Context, arg database.BlacklistAccessTokenParams) error {
					return expectedError
				}

				By("Generating an access token")
				accessToken, err := middleware.GenerateToken(testData.UserID, testData.Email)
				Expect(err).NotTo(HaveOccurred())

				By("Attempting to blacklist the access token")
				err = middleware.BlacklistToken(context.Background(), accessToken, "test")
				
				By("Verifying the network error is returned")
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(expectedError))
			})
		})
	})

	Describe("Blacklisting refresh tokens", func() {
		Context("When blacklisting a valid refresh token", func() {
			It("should call the refresh token blacklist method", func() {
				By("Generating a refresh token")
				refreshToken, err := middleware.GenerateRefreshToken(testData.UserID, testData.Email)
				Expect(err).NotTo(HaveOccurred())

				By("Blacklisting the refresh token")
				err = middleware.BlacklistToken(context.Background(), refreshToken, "test")
				Expect(err).NotTo(HaveOccurred())

				By("Verifying that refresh token blacklist was called")
				Expect(accessTokenBlacklisted).To(BeFalse())
				Expect(refreshTokenBlacklisted).To(BeTrue())
			})
		})

		Context("When database error occurs during refresh token blacklisting", func() {
			It("should return the error from blacklistRefreshTokenFunc", func() {
				By("Setting up mock to return database error")
				expectedError := errors.New("database transaction failed")
				mockQuerier.blacklistRefreshTokenFunc = func(ctx context.Context, arg database.BlacklistRefreshTokenParams) error {
					return expectedError
				}

				By("Generating a refresh token")
				refreshToken, err := middleware.GenerateRefreshToken(testData.UserID, testData.Email)
				Expect(err).NotTo(HaveOccurred())

				By("Attempting to blacklist the refresh token")
				err = middleware.BlacklistToken(context.Background(), refreshToken, "test")
				
				By("Verifying the database error is returned")
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(expectedError))
			})
		})

		Context("When network error occurs during refresh token blacklisting", func() {
			It("should return the network error from blacklistRefreshTokenFunc", func() {
				By("Setting up mock to return network error")
				expectedError := errors.New("connection refused")
				mockQuerier.blacklistRefreshTokenFunc = func(ctx context.Context, arg database.BlacklistRefreshTokenParams) error {
					return expectedError
				}

				By("Generating a refresh token")
				refreshToken, err := middleware.GenerateRefreshToken(testData.UserID, testData.Email)
				Expect(err).NotTo(HaveOccurred())

				By("Attempting to blacklist the refresh token")
				err = middleware.BlacklistToken(context.Background(), refreshToken, "test")
				
				By("Verifying the network error is returned")
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(expectedError))
			})
		})
	})
})