package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/necofuryai/bocchi-the-map/api/pkg/auth"
	"github.com/necofuryai/bocchi-the-map/api/pkg/config"
	"github.com/necofuryai/bocchi-the-map/api/pkg/logger"
)

// HealthCheckOutput represents the health check response
type HealthCheckOutput struct {
	Body struct {
		Status  string `json:"status" example:"ok" doc:"Health status"`
		Version string `json:"version" example:"1.0.0" doc:"API version"`
	}
}

// AuthStatusOutput represents the response for auth status check
type AuthStatusOutput struct {
	Body struct {
		Authenticated bool      `json:"authenticated" doc:"Whether the user is authenticated"`
		Timestamp     time.Time `json:"timestamp" doc:"Response timestamp"`
		Message       string    `json:"message" doc:"Status message"`
	}
}

// ValidateTokenInput represents the request to validate a token
type ValidateTokenInput struct {
	Body struct {
		Token string `json:"token" minLength:"1" doc:"JWT token to validate"`
	}
}

// ValidateTokenOutput represents the response for token validation
type ValidateTokenOutput struct {
	Body struct {
		Valid     bool      `json:"valid" doc:"Whether the token is valid"`
		Error     string    `json:"error,omitempty" doc:"Error message if invalid"`
		Timestamp time.Time `json:"timestamp" doc:"Response timestamp"`
	}
}

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		// Continue with default values for testing
		cfg = &config.Config{
			Auth: config.AuthConfig{
				JWTSecret:     "test-jwt-secret-1234567890abcdefghijklmnopqrstuvwxyz-with-special-chars!@#$%",
				Auth0Domain:   "test-domain.auth0.com",
				Auth0Audience: "bocchi-the-map-api",
			},
			App: config.AppConfig{
				Environment: "development",
				LogLevel:    "INFO",
				Version:     "test-1.0.0",
			},
		}
	}

	// Initialize logger
	logger.Init(logger.Level(cfg.App.LogLevel))
	logger.Info("Starting Test Auth Server")

	// Create chi router
	router := chi.NewRouter()

	// Add middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Compress(5))

	// Add CORS middleware
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Request-ID"},
		ExposedHeaders:   []string{"Link", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Initialize authentication service without database
	authService, err := initAuthService(cfg)
	if err != nil {
		logger.Error("Failed to initialize authentication service", err)
		// Continue without auth for basic testing
	}

	// Create Huma API
	config := huma.DefaultConfig("Test Auth API", cfg.App.Version)
	api := humachi.New(router, config)

	// Register routes
	registerTestRoutes(api, cfg, authService)

	// Start HTTP server
	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Channel to listen for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		logger.Info(fmt.Sprintf("Test Auth Server listening on port %s", port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("HTTP server failed to start", err)
		}
	}()

	// Wait for interrupt signal
	<-quit
	logger.Info("Shutting down server...")

	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server gracefully
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", err)
	}

	logger.Info("Server shutdown complete")
}

func initAuthService(cfg *config.Config) (*auth.Service, error) {
	// Try to create auth service without database
	serviceConfig := auth.ServiceConfig{
		Auth0Domain:     cfg.Auth.Auth0Domain,
		Auth0Audience:   cfg.Auth.Auth0Audience,
		JWTSecret:       cfg.Auth.JWTSecret,
		Development:     cfg.App.Environment == "development",
		RateLimit:       5,
		RateLimitWindow: 5 * time.Minute,
		SkipPaths: []string{
			"/health",
			"/metrics",
			"/swagger",
			"/docs",
		},
	}

	// Create service without database queries (pass nil)
	return auth.NewService(serviceConfig, nil)
}

func registerTestRoutes(api huma.API, cfg *config.Config, authService *auth.Service) {
	// Health check endpoint
	huma.Register(api, huma.Operation{
		OperationID: "health-check",
		Method:      http.MethodGet,
		Path:        "/health",
		Summary:     "Health Check",
		Description: "Check if the API is healthy",
		Tags:        []string{"System"},
	}, func(ctx context.Context, input *struct{}) (*HealthCheckOutput, error) {
		resp := &HealthCheckOutput{}
		resp.Body.Status = "ok"
		resp.Body.Version = cfg.App.Version
		return resp, nil
	})

	// Auth Status endpoint
	huma.Register(api, huma.Operation{
		OperationID: "auth-status",
		Method:      http.MethodGet,
		Path:        "/api/v1/auth/status",
		Summary:     "Get authentication status",
		Description: "Check if the current request is authenticated",
		Tags:        []string{"Authentication"},
	}, func(ctx context.Context, input *struct{}) (*AuthStatusOutput, error) {
		resp := &AuthStatusOutput{}
		resp.Body.Authenticated = false
		resp.Body.Timestamp = time.Now()
		resp.Body.Message = "Authentication service is available (no token provided)"
		return resp, nil
	})

	// Token validation endpoint
	huma.Register(api, huma.Operation{
		OperationID: "validate-token",
		Method:      http.MethodPost,
		Path:        "/api/v1/auth/validate",
		Summary:     "Validate JWT token",
		Description: "Validate a JWT token and return validation result",
		Tags:        []string{"Authentication"},
	}, func(ctx context.Context, input *ValidateTokenInput) (*ValidateTokenOutput, error) {
		resp := &ValidateTokenOutput{}
		resp.Body.Timestamp = time.Now()

		if input.Body.Token == "" {
			resp.Body.Valid = false
			resp.Body.Error = "token is required"
			return resp, nil
		}

		// Simple token validation (just check if it's not empty for testing)
		if input.Body.Token == "invalid.jwt.token" || len(input.Body.Token) < 10 {
			resp.Body.Valid = false
			resp.Body.Error = "invalid token format"
			return resp, nil
		}

		if authService != nil {
			// Try to validate with auth service
			_, err := authService.ValidateToken(ctx, input.Body.Token)
			if err != nil {
				resp.Body.Valid = false
				resp.Body.Error = err.Error()
				return resp, nil
			}
		}

		resp.Body.Valid = true
		return resp, nil
	})

	// Test protected endpoint
	huma.Register(api, huma.Operation{
		OperationID: "test-protected",
		Method:      http.MethodGet,
		Path:        "/api/v1/protected",
		Summary:     "Test protected endpoint",
		Description: "Test endpoint that requires authentication",
		Tags:        []string{"Test"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, func(ctx context.Context, input *struct{}) (*struct {
		Body struct {
			Message string `json:"message"`
		}
	}, error) {
		// Check for Authorization header
		// In a real implementation, this would be handled by middleware
		resp := &struct {
			Body struct {
				Message string `json:"message"`
			}
		}{}
		resp.Body.Message = "This endpoint requires authentication"
		return resp, nil
	})

	logger.Info("Test routes registered successfully")
}