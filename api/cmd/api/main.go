package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"google.golang.org/grpc"

	"github.com/necofuryai/bocchi-the-map/api/application/clients"
	"github.com/necofuryai/bocchi-the-map/api/infrastructure/database"
	"github.com/necofuryai/bocchi-the-map/api/interfaces/http/handlers"
	"github.com/necofuryai/bocchi-the-map/api/pkg/auth"
	"github.com/necofuryai/bocchi-the-map/api/pkg/config"
	"github.com/necofuryai/bocchi-the-map/api/pkg/logger"
	"github.com/necofuryai/bocchi-the-map/api/pkg/monitoring"
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

	// Initialize monitoring services
	if err := monitoring.InitMonitoring(
		cfg.Monitoring.NewRelicLicenseKey,
		cfg.Monitoring.SentryDSN,
		"bocchi-the-map-api",
		cfg.App.Environment,
		cfg.App.Version,
	); err != nil {
		logger.Error("Failed to initialize monitoring", err)
		// Don't exit - monitoring is not critical for basic functionality
	}

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

		// Create database queries instance
		queries := database.New(db)

		// Initialize gRPC clients (using internal communication for monolith)
		spotClient, err := clients.NewSpotClient("internal", db)
		if err != nil {
			logger.Fatal("Failed to create spot client", err)
		}

		userClient, err := clients.NewUserClient("internal", db)
		if err != nil {
			spotClient.Close()
			logger.Fatal("Failed to create user client", err)
		}

		// TODO: Review client creation is temporarily disabled due to compilation errors.
		// Re-enable once the review service implementation is fixed.
		// reviewClient, err := clients.NewReviewClient("internal", db)
		// if err != nil {
		// 	spotClient.Close()
		// 	userClient.Close()
		// 	logger.Fatal("Failed to create review client", err)
		// }

		// Ensure proper cleanup on shutdown
		hooks.OnStop(func() {
			logger.Info("Shutting down application...")
			
			// Shutdown monitoring services
			monitoring.ShutdownMonitoring()
			
			// Close gRPC clients
			spotClient.Close()
			userClient.Close()
			// reviewClient.Close()
			
			// Close database connection
			logger.Info("Closing database connection")
			if err := db.Close(); err != nil {
				logger.Error("Failed to close database connection", err)
			}
			
			logger.Info("Application shutdown complete")
		})

		// Create chi router
		router := chi.NewRouter()

		// Add middleware
		router.Use(middleware.RequestID)
		router.Use(middleware.RealIP)
		router.Use(middleware.Logger)
		router.Use(middleware.Recoverer)
		router.Use(middleware.Compress(5))
		
		// Add CORS middleware for frontend integration
		router.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"http://localhost:3000", "https://bocchi-the-map.vercel.app"}, // Next.js dev and production
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Request-ID"},
			ExposedHeaders:   []string{"Link", "X-Request-ID"},
			AllowCredentials: true,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}))
		
		// Add monitoring middleware
		router.Use(monitoring.RequestIDMiddleware())
		router.Use(monitoring.MonitoringMiddleware())
		router.Use(monitoring.PerformanceMiddleware())

		// Initialize authentication middleware with token blacklist support
		authMiddleware := auth.NewAuthMiddleware(cfg.Auth.JWTSecret, queries)
		
		// Initialize rate limiter (5 requests per 5 minutes for auth endpoints)
		// Rate limiting is applied per IP address, not per user
		rateLimiter := auth.NewRateLimiter(5, 5*time.Minute)

		// Create Huma API
		api := humachi.New(router, huma.DefaultConfig("Bocchi The Map API", cfg.App.Version))

		// Register routes with gRPC clients and database queries
		registerRoutes(api, spotClient, userClient, nil, queries, cfg, authMiddleware, rateLimiter)

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

		// Start HTTP server with graceful shutdown
		hooks.OnStart(func() {
			logger.Info(fmt.Sprintf("HTTP server starting on port %s", options.Port))
			logger.Info(fmt.Sprintf("gRPC server starting on port %s", options.GRPCPort))
			
			// Create HTTP server
			server := &http.Server{
				Addr:    ":" + options.Port,
				Handler: router,
				ReadTimeout:  15 * time.Second,
				WriteTimeout: 15 * time.Second,
				IdleTimeout:  60 * time.Second,
			}
			
			// Channel to listen for interrupt signal
			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
			
			// Start server in a goroutine
			go func() {
				logger.Info("Server is ready to handle requests")
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Fatal("HTTP server failed to start", err)
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
func registerRoutes(api huma.API, spotClient *clients.SpotClient, userClient *clients.UserClient, reviewClient *clients.ReviewClient, queries *database.Queries, cfg *config.Config, authMiddleware *auth.AuthMiddleware, rateLimiter *auth.RateLimiter) {
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

	// Spot routes
	registerSpotRoutes(api, spotClient)

	// Review routes (temporarily disabled due to compile issues)
	// registerReviewRoutes(api, reviewClient)

	// User routes
	registerUserRoutes(api, userClient, queries, authMiddleware)
	
	// Authentication routes
	registerAuthRoutes(api, authMiddleware, userClient, rateLimiter)
}

// registerSpotRoutes registers spot-related routes
func registerSpotRoutes(api huma.API, spotClient *clients.SpotClient) {
	spotHandler := handlers.NewSpotHandler(spotClient)
	spotHandler.RegisterRoutes(api)
	logger.Info("Spot routes registered")
}

func registerReviewRoutes(api huma.API, reviewClient *clients.ReviewClient) {
	reviewHandler := handlers.NewReviewHandler(reviewClient)
	
	// Register standard API routes (under /api/v1/reviews)
	reviewHandler.RegisterRoutes(api)
	logger.Info("Review routes registered")
}

func registerUserRoutes(api huma.API, userClient *clients.UserClient, queries *database.Queries, authMiddleware *auth.AuthMiddleware) {
	userHandler := handlers.NewUserHandler(userClient)
	
	// Register standard API routes (under /api/v1/users) with authentication middleware
	userHandler.RegisterRoutesWithAuth(api, authMiddleware)
	logger.Info("User routes registered with authentication")
}

// registerAuthRoutes registers authentication-related routes
func registerAuthRoutes(api huma.API, authMiddleware *auth.AuthMiddleware, userClient *clients.UserClient, rateLimiter *auth.RateLimiter) {
	authHandler := handlers.NewAuthHandler(authMiddleware, userClient)
	
	// Register authentication routes with rate limiting
	authHandler.RegisterRoutesWithRateLimit(api, rateLimiter)
	logger.Info("Authentication routes registered with rate limiting")
}