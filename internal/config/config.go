package config

import (
	"os"
)

type Config struct {
	DatabaseURL string
	RedisURL    string
	APIPort     string
	LogLevel    string
	Provider1BaseURL string
	Provider2BaseURL string
	ProviderTimeout  string
	RateLimitEnabled string
}

func Load() Config {
	cfg := Config{
		DatabaseURL: getenv("DATABASE_URL", "postgres://postgres:postgres@db:5432/searchdb?sslmode=disable"),
		RedisURL:    getenv("REDIS_URL", "redis://redis:6379"),
		APIPort:     getenv("API_PORT", "8080"),
		LogLevel:    getenv("LOG_LEVEL", "info"),
		Provider1BaseURL: getenv("PROVIDER1_BASE_URL", "http://localhost:8080/mock/provider1"),
		Provider2BaseURL: getenv("PROVIDER2_BASE_URL", "http://localhost:8080/mock/provider2"),
		ProviderTimeout:  getenv("PROVIDER_TIMEOUT", "10s"),
		RateLimitEnabled: getenv("RATE_LIMIT_ENABLED", "true"),
	}
	return cfg
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}


