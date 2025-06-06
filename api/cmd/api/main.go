package main

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"

	"github.com/necofuryai/bocchi-the-map/api/application/clients"
	"github.com/necofuryai/bocchi-the-map/api/infrastructure/grpc"
	"github.com/necofuryai/bocchi-the-map/api/interfaces/http/handlers"
	"github.com/necofuryai/bocchi-the-map/api/pkg/config"
	"github.com/necofuryai/bocchi-the-map/api/pkg/logger"
)

// Options for the CLI
type Options struct {
	Port     string `help:"HTTP port to listen on" default:"8080"`
	GRPCPort string `help:"gRPC port to listen on" default:"9090"`
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
		// Initialize gRPC clients (using internal communication for monolith)
		spotClient, err := clients.NewSpotClient("internal")
		if err != nil {
			logger.Fatal("Failed to create spot client", err)
		}
		defer spotClient.Close()

		userClient, err := clients.NewUserClient("internal")
		if err != nil {
			logger.Fatal("Failed to create user client", err)
		}
		defer userClient.Close()

		reviewClient, err := clients.NewReviewClient("internal")
		if err != nil {
			logger.Fatal("Failed to create review client", err)
		}
		defer reviewClient.Close()

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

		// Register routes with gRPC clients
		registerRoutes(api, spotClient, userClient, reviewClient)

		// Start gRPC server in a goroutine
		go func() {
			if err := startGRPCServer(options.GRPCPort); err != nil {
				logger.Fatal("gRPC server failed to start", err)
			}
		}()

		// Start HTTP server
		hooks.OnStart(func() {
			logger.Info(fmt.Sprintf("HTTP server starting on port %s", options.Port))
			logger.Info(fmt.Sprintf("gRPC server starting on port %s", options.GRPCPort))
			if err := http.ListenAndServe(":"+options.Port, router); err != nil {
				logger.Fatal("HTTP server failed to start", err)
			}
		})
	})

	// Run CLI
	cli.Run()
}

// startGRPCServer starts the gRPC server
func startGRPCServer(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", port, err)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register gRPC services
	// TODO: Register actual gRPC service implementations when protobuf is generated
	// For now, the services are used internally via clients

	logger.Info(fmt.Sprintf("gRPC server listening on port %s", port))
	return grpcServer.Serve(lis)
}

// registerRoutes registers all API routes
func registerRoutes(api huma.API, spotClient *clients.SpotClient, userClient *clients.UserClient, reviewClient *clients.ReviewClient) {
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
	registerSpotRoutes(v1, spotClient)

	// Review routes
	registerReviewRoutes(v1, reviewClient)

	// User routes
	registerUserRoutes(v1, userClient)
}

// registerSpotRoutes registers spot-related routes
func registerSpotRoutes(api huma.API, spotClient *clients.SpotClient) {
	spotHandler := handlers.NewSpotHandler(spotClient)
	spotHandler.RegisterRoutes(api)
	logger.Info("Spot routes registered")
}

func registerReviewRoutes(api huma.API, reviewClient *clients.ReviewClient) {
	// TODO: Implement review routes with reviewClient
	logger.Info("Review routes registered")
}

func registerUserRoutes(api huma.API, userClient *clients.UserClient) {
	// TODO: Implement user routes with userClient
	logger.Info("User routes registered")
}