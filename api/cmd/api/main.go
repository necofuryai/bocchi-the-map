package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"

	"github.com/necofuryai/bocchi-the-map/api/application/clients"
	"github.com/necofuryai/bocchi-the-map/api/infrastructure/database"
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
		// Initialize database connection
		db, err := sql.Open("mysql", cfg.Database.GetDSN())
		if err != nil {
			logger.Fatal("Failed to connect to database", err)
		}
		// Note: Database connection will be closed when the application shuts down

		// Configure connection pool
		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(25)
		db.SetConnMaxLifetime(5 * time.Minute)

		// Test database connection
		if err := db.Ping(); err != nil {
			logger.Fatal("Failed to ping database", err)
		}
		logger.Info("Database connection established")

		// Ensure database connection is closed on shutdown
		hooks.OnStop(func() {
			logger.Info("Closing database connection")
			if err := db.Close(); err != nil {
				logger.Error("Failed to close database connection", err)
			}
		})

		// Create database queries instance
		queries := database.New(db)

		// Initialize gRPC clients (using internal communication for monolith)
		spotClient, err := clients.NewSpotClient("internal")
		if err != nil {
			logger.Fatal("Failed to create spot client", err)
		}

		userClient, err := clients.NewUserClient("internal", db)
		if err != nil {
			spotClient.Close()
			logger.Fatal("Failed to create user client", err)
		}

		reviewClient, err := clients.NewReviewClient("internal")
		if err != nil {
			spotClient.Close()
			userClient.Close()
			logger.Fatal("Failed to create review client", err)
		}

		defer spotClient.Close()
		defer userClient.Close()
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

		// Register routes with gRPC clients and database queries
		registerRoutes(api, spotClient, userClient, reviewClient, queries)

		// Start gRPC server in a goroutine
		errChan := make(chan error, 1)
		go func() {
			if err := startGRPCServer(options.GRPCPort); err != nil {
				errChan <- fmt.Errorf("gRPC server failed: %w", err)
			}
		}()

		// Check for immediate startup errors
		select {
		case err := <-errChan:
			logger.Fatal("Server startup failed", err)
		case <-time.After(100 * time.Millisecond):
			// Continue if no immediate errors
		}

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
func registerRoutes(api huma.API, spotClient *clients.SpotClient, userClient *clients.UserClient, reviewClient *clients.ReviewClient, queries *database.Queries) {
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

	// Spot routes
	registerSpotRoutes(api, spotClient)

	// Review routes
	registerReviewRoutes(api, reviewClient)

	// User routes
	registerUserRoutes(api, userClient, queries)
	
	// Auth.js compatible routes (outside of v1 prefix)
	registerAuthRoutes(api, userClient, queries)
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

func registerUserRoutes(api huma.API, userClient *clients.UserClient, queries *database.Queries) {
	userHandler := handlers.NewUserHandler(queries)
	
	// Register standard API routes (under /api/v1/users)
	userHandler.RegisterRoutes(api)
	logger.Info("User routes registered")
}

func registerAuthRoutes(api huma.API, userClient *clients.UserClient, queries *database.Queries) {
	userHandler := handlers.NewUserHandler(queries)
	
	// Register Auth.js compatible routes (POST /api/users)
	huma.Register(api, huma.Operation{
		OperationID: "auth-create-user",
		Method:      http.MethodPost,
		Path:        "/api/users", // Auth.js expects this exact path
		Summary:     "Create or update user (Auth.js)",
		Description: "OAuth authentication endpoint for Auth.js user creation/update",
		Tags:        []string{"Authentication"},
	}, userHandler.CreateUser)
	
	logger.Info("Auth.js compatible routes registered")
}