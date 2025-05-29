package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds the application configuration
type Config struct {
	Server     ServerConfig
	Database   DatabaseConfig
	Monitoring MonitoringConfig
	App        AppConfig
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
		},
	}

	return cfg, nil
}

// GetDSN returns the database connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.Database)
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
		}
	}
	return defaultValue
}