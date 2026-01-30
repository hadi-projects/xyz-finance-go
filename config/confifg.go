package config

import (
	"errors"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort    string
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
	Security   SecurityConfig
}

type SecurityConfig struct {
	CORSAllowedOrigins   []string
	CORSAllowCredentials bool
}

// NewConfig menerapkan Constructor Pattern.
func NewConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		AppPort:    getEnv("APP_PORT", "8080"),
		DBUser:     getEnv("DB_USER", ""),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBHost:     getEnv("DB_HOST", ""),
		DBPort:     getEnv("DB_PORT", ""),
		DBName:     getEnv("DB_NAME", ""),
		Security: SecurityConfig{
			CORSAllowedOrigins:   getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{"*"}),
			CORSAllowCredentials: getEnvAsBool("CORS_ALLOW_CREDENTIALS", true),
		},
	}

	if cfg.DBHost == "" || cfg.DBPort == "" {
		return nil, errors.New("database configuration (HOST/PORT) is missing")
	}

	return cfg, nil
}

// Helper function untuk membaca env dengan default value
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsSlice(key string, fallback []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return strings.Split(value, ",")
}

func getEnvAsBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value == "true"
}
