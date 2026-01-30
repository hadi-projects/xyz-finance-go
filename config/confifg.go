package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort    string
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
}

// NewConfig menerapkan Constructor Pattern.
func NewConfig() (*Config, error) {
	// 1. Load .env file
	// Kita tidak panic di sini, karena di production mungkin environment variable sudah di-set di sistem (Docker/K8s)
	// Jadi jika file .env tidak ada, kita lanjut saja mengecek environment variable sistem.
	_ = godotenv.Load()

	cfg := &Config{
		AppPort:    getEnv("APP_PORT", "8080"),
		DBUser:     getEnv("DB_USER", ""),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBHost:     getEnv("DB_HOST", ""),
		DBPort:     getEnv("DB_PORT", ""),
		DBName:     getEnv("DB_NAME", ""),
	}

	// 2. Validasi (Opsional tapi Recommended)
	// Pastikan konfigurasi kritikal tidak kosong
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
