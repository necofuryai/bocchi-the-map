package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/necofuryai/bocchi-the-map/api/interfaces/http/handlers"
	"github.com/necofuryai/bocchi-the-map/api/pkg/config"
	"github.com/necofuryai/bocchi-the-map/api/pkg/logger"
)

// Options for the CLI
type Options struct {
	Port string `help:"Port to listen on" default:"8080"`
}

// HealthCheckOutput represents the health check response
type HealthCheckOutput struct {
	Body struct {
		Status  string `json:"status" example:"ok" doc:"Health status"`
		Version string `json:"version" example:"1.0.0" doc:"API version"`
	}
}

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load configuration", err)
	}

	// Initialize logger
	logger.Init(logger.Level(cfg.App.LogLevel))
	logger.Info("Starting Bocchi The Map API")

	// Create CLI
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		// Create chi router
		router := chi.NewRouter()

		// Add middleware
		router.Use(middleware.RequestID)
		router.Use(middleware.RealIP)
		router.Use(middleware.Logger)
		router.Use(middleware.Recoverer)
		router.Use(middleware.Compress(5))

		// Create Huma API
		api := humachi.New(router, huma.DefaultConfig("Bocchi The Map API", "1.0.0"))

		// Register routes
		registerRoutes(api)

		// Start server
		hooks.OnStart(func() {
			logger.Info(fmt.Sprintf("Server starting on port %s", options.Port))
			if err := http.ListenAndServe(":"+options.Port, router); err != nil {
				logger.Fatal("Server failed to start", err)
			}
		})
	})

	// Run CLI
	cli.Run()
}

// registerRoutes registers all API routes
func registerRoutes(api huma.API) {
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
		resp.Body.Version = "1.0.0"
		return resp, nil
	})

	// API v1 routes
	v1 := api.Group("/api/v1")

	// Spot routes
	registerSpotRoutes(v1)

	// Review routes
	registerReviewRoutes(v1)

	// User routes
	registerUserRoutes(v1)
}

// registerSpotRoutes registers spot-related routes
func registerSpotRoutes(api huma.API) {
	spotHandler := handlers.NewSpotHandler()
	spotHandler.RegisterRoutes(api)
	logger.Info("Spot routes registered")
}

func registerReviewRoutes(api huma.API) {
	// TODO: Implement review routes
	logger.Info("Review routes registered")
}

func registerUserRoutes(api huma.API) {
	// TODO: Implement user routes
	logger.Info("User routes registered")
}