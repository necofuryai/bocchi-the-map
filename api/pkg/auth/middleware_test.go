package auth

import (
	"context"
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

func TestJWTAuthMiddleware(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "JWT AuthMiddleware Test Suite")
}

var _ = Describe("JWT Token Generation and Claims", func() {
	var (
		middleware *AuthMiddleware
		mockQuerier *mockTokenBlacklistQuerier
		testUserID string
		testEmail string
	)

	BeforeEach(func() {
		By("Setting up test environment")
		mockQuerier = &mockTokenBlacklistQuerier{}
		middleware = NewAuthMiddleware("test-secret", mockQuerier)
		testUserID = "test-user-id"
		testEmail = "test@example.com"
	})

	Describe("Access token generation", func() {
		Context("When generating an access token with valid user data", func() {
			It("should create a token with correct claims and token type", func() {
				By("Generating the access token")
				tokenString, err := middleware.GenerateToken(testUserID, testEmail)
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
				Expect(claims.UserID).To(Equal(testUserID))
				Expect(claims.Email).To(Equal(testEmail))
			})
		})
	})

	Describe("Refresh token generation", func() {
		Context("When generating a refresh token with valid user data", func() {
			It("should create a token with correct claims and token type", func() {
				By("Generating the refresh token")
				tokenString, err := middleware.GenerateRefreshToken(testUserID, testEmail)
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
				Expect(claims.UserID).To(Equal(testUserID))
				Expect(claims.Email).To(Equal(testEmail))
			})
		})
	})
})

var _ = Describe("Token Blacklisting", func() {
	var (
		middleware *AuthMiddleware
		mockQuerier *mockTokenBlacklistQuerier
		testUserID string
		testEmail string
		accessTokenBlacklisted bool
		refreshTokenBlacklisted bool
	)

	BeforeEach(func() {
		By("Setting up blacklist test environment")
		accessTokenBlacklisted = false
		refreshTokenBlacklisted = false
		testUserID = "test-user-id"
		testEmail = "test@example.com"

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
				accessToken, err := middleware.GenerateToken(testUserID, testEmail)
				Expect(err).NotTo(HaveOccurred())

				By("Blacklisting the access token")
				err = middleware.BlacklistToken(context.Background(), accessToken, "test")
				Expect(err).NotTo(HaveOccurred())

				By("Verifying that access token blacklist was called")
				Expect(accessTokenBlacklisted).To(BeTrue())
				Expect(refreshTokenBlacklisted).To(BeFalse())
			})
		})
	})

	Describe("Blacklisting refresh tokens", func() {
		Context("When blacklisting a valid refresh token", func() {
			It("should call the refresh token blacklist method", func() {
				By("Generating a refresh token")
				refreshToken, err := middleware.GenerateRefreshToken(testUserID, testEmail)
				Expect(err).NotTo(HaveOccurred())

				By("Blacklisting the refresh token")
				err = middleware.BlacklistToken(context.Background(), refreshToken, "test")
				Expect(err).NotTo(HaveOccurred())

				By("Verifying that refresh token blacklist was called")
				Expect(accessTokenBlacklisted).To(BeFalse())
				Expect(refreshTokenBlacklisted).To(BeTrue())
			})
		})
	})
})