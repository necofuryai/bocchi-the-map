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
	JWTSecret string
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
			JWTSecret: os.Getenv("JWT_SECRET"),
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
	return nil
}

// validateJWTSecret validates the JWT secret with strict security requirements
func (c *Config) validateJWTSecret() error {
	secret := c.Auth.JWTSecret
	
	if secret == "" {
		return errors.New("JWT_SECRET is required")
	}
	
	// Minimum length requirement (32 characters)
	if len(secret) < 32 {
		return errors.New("JWT_SECRET must be at least 32 characters long")
	}
	
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