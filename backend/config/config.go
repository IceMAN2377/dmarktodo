package config

import (
	"os"
	"strconv"
)

type Config struct {
	PostgresMigrate  bool
	PostgresHost     string
	PostgresPort     int
	PostgresUser     string
	PostgresPassword string
	PostgresDb       string
	PostgresSslMode  string
}

// NewConfig creates a new Config with default values for Docker Postgres
func NewConfig() *Config {
	cfg := &Config{
		PostgresMigrate:  true,
		PostgresHost:     getEnvOrDefault("POSTGRES_HOST", "localhost"),
		PostgresPort:     getEnvIntOrDefault("POSTGRES_PORT", 5432),
		PostgresUser:     getEnvOrDefault("POSTGRES_USER", "myuser"),
		PostgresPassword: getEnvOrDefault("POSTGRES_PASSWORD", "mypassword"),
		PostgresDb:       getEnvOrDefault("POSTGRES_DB", "todolist"),
		PostgresSslMode:  getEnvOrDefault("POSTGRES_SSL_MODE", "disable"),
	}
	return cfg
}

// getEnvOrDefault gets environment variable or returns default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvIntOrDefault gets environment variable as int or returns default value
func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
