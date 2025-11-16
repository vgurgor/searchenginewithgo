package config

import (
	"os"
)

type Config struct {
	DatabaseURL string
	RedisURL    string
	APIPort     string
	LogLevel    string
}

func Load() Config {
	cfg := Config{
		DatabaseURL: getenv("DATABASE_URL", "postgres://postgres:postgres@db:5432/searchdb?sslmode=disable"),
		RedisURL:    getenv("REDIS_URL", "redis://redis:6379"),
		APIPort:     getenv("API_PORT", "8080"),
		LogLevel:    getenv("LOG_LEVEL", "info"),
	}
	return cfg
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}


