package config

import (
	"errors"
	"os"
	"strconv"
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
	JWT        JWTConfig
}

type SecurityConfig struct {
	RateLimitRPS         float64
	RateLimitBurst       int
	CORSAllowedOrigins   []string
	CORSAllowCredentials bool
	RequestTimeout       int
}

type JWTConfig struct {
	Secret      string
	ExpiryHours int
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
			RateLimitRPS:         getEnvAsFloat("RATE_LIMIT_RPS", 10),
			RateLimitBurst:       getEnvAsInt("RATE_LIMIT_BURST", 10),
			CORSAllowedOrigins:   getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{"*"}),
			CORSAllowCredentials: getEnvAsBool("CORS_ALLOW_CREDENTIALS", true),
			RequestTimeout:       getEnvAsInt("REQUEST_TIMEOUT", 10),
		},
		JWT: JWTConfig{
			Secret:      getEnv("JWT_SECRET", ""),
			ExpiryHours: getEnvAsInt("JWT_EXPIRY_HOURS", 1),
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

func getEnvAsFloat(key string, defaultValue float64) float64 {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseFloat(valueStr, 64); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
