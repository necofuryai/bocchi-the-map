package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strconv"
)

// Config holds the application configuration
type Config struct {
	Server     ServerConfig
	Database   DatabaseConfig
	Monitoring MonitoringConfig
	App        AppConfig
	Auth       AuthConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port string
	Host string
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

// MonitoringConfig holds monitoring-related configuration
type MonitoringConfig struct {
	NewRelicLicenseKey string
	SentryDSN          string
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	Environment string
	LogLevel    string
	Version     string
}

// AuthConfig holds authentication-related configuration
type AuthConfig struct {
	JWTSecret       string
	Auth0Domain     string
	Auth0Audience   string
	Auth0ClientID   string
	Auth0ClientSecret string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port: getEnvWithDefault("PORT", "8080"),
			Host: getEnvWithDefault("HOST", "0.0.0.0"),
		},
		Database: DatabaseConfig{
			Host:     getEnvWithDefault("TIDB_HOST", "localhost"),
			Port:     getIntEnvWithDefault("TIDB_PORT", 4000),
			User:     getEnvWithDefault("TIDB_USER", "root"),
			Password: os.Getenv("TIDB_PASSWORD"),
			Database: getEnvWithDefault("TIDB_DATABASE", "bocchi_the_map"),
		},
		Monitoring: MonitoringConfig{
			NewRelicLicenseKey: os.Getenv("NEW_RELIC_LICENSE_KEY"),
			SentryDSN:          os.Getenv("SENTRY_DSN"),
		},
		App: AppConfig{
			Environment: getEnvWithDefault("ENV", "development"),
			LogLevel:    getEnvWithDefault("LOG_LEVEL", "INFO"),
			Version:     getEnvWithDefault("APP_VERSION", "1.0.0"),
		},
		Auth: AuthConfig{
			JWTSecret:       os.Getenv("JWT_SECRET"),
			Auth0Domain:     os.Getenv("AUTH0_DOMAIN"),
			Auth0Audience:   os.Getenv("AUTH0_AUDIENCE"),
			Auth0ClientID:   os.Getenv("AUTH0_CLIENT_ID"),
			Auth0ClientSecret: os.Getenv("AUTH0_CLIENT_SECRET"),
		},
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Database.Password == "" {
		return errors.New("TIDB_PASSWORD is required")
	}
	if err := c.validateJWTSecret(); err != nil {
		return err
	}
	if err := c.validateAuth0Config(); err != nil {
		return err
	}
	return nil
}

// validateJWTSecret validates the JWT secret with environment-specific requirements
func (c *Config) validateJWTSecret() error {
	secret := c.Auth.JWTSecret
	
	if secret == "" {
		return errors.New("JWT_SECRET is required")
	}
	
	// Minimum length requirement (32 characters)
	if len(secret) < 32 {
		return errors.New("JWT_SECRET must be at least 32 characters long")
	}
	
	// Relaxed validation for development environment
	if c.App.Environment == "development" || c.App.Environment == "dev" {
		return nil
	}
	
	// Strict validation for production and other environments
	// Check for lowercase letters
	hasLower, _ := regexp.MatchString(`[a-z]`, secret)
	if !hasLower {
		return errors.New("JWT_SECRET must contain at least one lowercase letter")
	}
	
	// Check for uppercase letters
	hasUpper, _ := regexp.MatchString(`[A-Z]`, secret)
	if !hasUpper {
		return errors.New("JWT_SECRET must contain at least one uppercase letter")
	}
	
	// Check for numbers
	hasNumber, _ := regexp.MatchString(`[0-9]`, secret)
	if !hasNumber {
		return errors.New("JWT_SECRET must contain at least one number")
	}
	
	// Check for special characters
	hasSpecial, _ := regexp.MatchString(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?~` + "`" + `]`, secret)
	if !hasSpecial {
		return errors.New("JWT_SECRET must contain at least one special character")
	}
	
	return nil
}

// validateAuth0Config validates the Auth0 configuration
func (c *Config) validateAuth0Config() error {
	auth := c.Auth
	
	// Auth0 Domain validation
	if auth.Auth0Domain == "" {
		return errors.New("AUTH0_DOMAIN is required")
	}
	
	// Basic domain format validation (should not include protocol)
	if match, _ := regexp.MatchString(`^https?://`, auth.Auth0Domain); match {
		return errors.New("AUTH0_DOMAIN should not include protocol (http:// or https://)")
	}
	
	// Auth0 Audience validation
	if auth.Auth0Audience == "" {
		return errors.New("AUTH0_AUDIENCE is required")
	}
	
	// Auth0 Client ID validation
	if auth.Auth0ClientID == "" {
		return errors.New("AUTH0_CLIENT_ID is required")
	}
	
	// Auth0 Client Secret validation - only required in production
	if auth.Auth0ClientSecret == "" && (c.App.Environment == "production" || c.App.Environment == "prod") {
		return errors.New("AUTH0_CLIENT_SECRET is required in production environment")
	}
	
	return nil
}

// GetAuth0Issuer returns the Auth0 issuer URL
func (c *AuthConfig) GetAuth0Issuer() string {
	return fmt.Sprintf("https://%s/", c.Auth0Domain)
}

// GetAuth0JWKSURL returns the Auth0 JWKS URL for token validation
func (c *AuthConfig) GetAuth0JWKSURL() string {
	return fmt.Sprintf("https://%s/.well-known/jwks.json", c.Auth0Domain)
}

// IsAuth0Configured returns true if Auth0 configuration is present
func (c *AuthConfig) IsAuth0Configured() bool {
	return c.Auth0Domain != "" && c.Auth0ClientID != "" && c.Auth0Audience != ""
}

// GetDSN returns the database connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		url.QueryEscape(c.User), url.QueryEscape(c.Password), c.Host, c.Port, c.Database)
}

// getEnvWithDefault gets an environment variable with a default value
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getIntEnvWithDefault gets an integer environment variable with a default value
func getIntEnvWithDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		} else {
			// Log warning about invalid integer value
			fmt.Fprintf(os.Stderr, "Warning: Invalid integer value for %s: %s, using default %d\n", key, value, defaultValue)
		}
	}
	return defaultValue
}