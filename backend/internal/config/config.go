package config

import (
	"os"
)

type Config struct {
	DatabaseURL string
	RedisURL    string
	APIPort     string
	LogLevel    string
	// Public API rate limiting
	PublicRateLimit       string // requests per window (per IP)
	PublicRateLimitWindow string // duration, e.g., "1m"
	// Search cache
	SearchCacheEnabled string // "true" | "false"
	SearchCacheTTL     string // duration, e.g., "60s"
	Provider1BaseURL   string
	Provider2BaseURL   string
	ProviderTimeout    string
	RateLimitEnabled   string
	// Scoring
	ScoreRecalcEnabled  string
	ScoreRecalcInterval string
	ScoreBatchSize      string
	VideoTypeMultiplier string
	TextTypeMultiplier  string
	Freshness1Week      string
	Freshness1Month     string
	Freshness3Months    string
	// Sync
	ContentSyncEnabled                 string
	ContentSyncInterval                string
	ContentSyncRetryCount              string
	ContentSyncRetryDelay              string
	MetricsChangeThresholdPercent      string
	MetricsChangeThresholdAbsViews     string
	MetricsChangeThresholdAbsLikes     string
	MetricsChangeThresholdAbsReactions string
	AdminAPIKey                        string
	// API pagination
	DefaultPageSize string
	MaxPageSize     string
	// Admin API
	AdminAPIEnabled         string
	AdminAPIKeyRotationDays string
	AdminRateLimit          string
	AdminAuditEnabled       string
	// Async Jobs
	AsyncJobsEnabled  string
	MaxConcurrentJobs string
	JobTimeout        string
}

func Load() Config {
	cfg := Config{
		DatabaseURL:                        getenv("DATABASE_URL", "postgres://postgres:postgres@db:5432/searchdb?sslmode=disable"),
		RedisURL:                           getenv("REDIS_URL", "redis://redis:6379"),
		APIPort:                            getenv("API_PORT", "8080"),
		LogLevel:                           getenv("LOG_LEVEL", "info"),
		PublicRateLimit:                    getenv("PUBLIC_RATE_LIMIT", "300"),
		PublicRateLimitWindow:              getenv("PUBLIC_RATE_LIMIT_WINDOW", "1m"),
		SearchCacheEnabled:                 getenv("SEARCH_CACHE_ENABLED", "true"),
		SearchCacheTTL:                     getenv("SEARCH_CACHE_TTL", "60s"),
		Provider1BaseURL:                   getenv("PROVIDER1_BASE_URL", "http://localhost:8080/mock/provider1"),
		Provider2BaseURL:                   getenv("PROVIDER2_BASE_URL", "http://localhost:8080/mock/provider2"),
		ProviderTimeout:                    getenv("PROVIDER_TIMEOUT", "10s"),
		RateLimitEnabled:                   getenv("RATE_LIMIT_ENABLED", "true"),
		ScoreRecalcEnabled:                 getenv("SCORE_RECALCULATION_ENABLED", "true"),
		ScoreRecalcInterval:                getenv("SCORE_RECALCULATION_INTERVAL", "24h"),
		ScoreBatchSize:                     getenv("SCORE_BATCH_SIZE", "100"),
		VideoTypeMultiplier:                getenv("VIDEO_TYPE_MULTIPLIER", "1.5"),
		TextTypeMultiplier:                 getenv("TEXT_TYPE_MULTIPLIER", "1.0"),
		Freshness1Week:                     getenv("FRESHNESS_1_WEEK", "5"),
		Freshness1Month:                    getenv("FRESHNESS_1_MONTH", "3"),
		Freshness3Months:                   getenv("FRESHNESS_3_MONTHS", "1"),
		ContentSyncEnabled:                 getenv("CONTENT_SYNC_ENABLED", "true"),
		ContentSyncInterval:                getenv("CONTENT_SYNC_INTERVAL", "6h"),
		ContentSyncRetryCount:              getenv("CONTENT_SYNC_RETRY_COUNT", "3"),
		ContentSyncRetryDelay:              getenv("CONTENT_SYNC_RETRY_DELAY", "30s"),
		MetricsChangeThresholdPercent:      getenv("METRICS_CHANGE_THRESHOLD_PERCENT", "5"),
		MetricsChangeThresholdAbsViews:     getenv("METRICS_CHANGE_THRESHOLD_ABS_VIEWS", "100"),
		MetricsChangeThresholdAbsLikes:     getenv("METRICS_CHANGE_THRESHOLD_ABS_LIKES", "10"),
		MetricsChangeThresholdAbsReactions: getenv("METRICS_CHANGE_THRESHOLD_ABS_REACTIONS", "5"),
		AdminAPIKey:                        getenv("ADMIN_API_KEY", ""),
		DefaultPageSize:                    getenv("DEFAULT_PAGE_SIZE", "20"),
		MaxPageSize:                        getenv("MAX_PAGE_SIZE", "100"),
		AdminAPIEnabled:                    getenv("ADMIN_API_ENABLED", "true"),
		AdminAPIKeyRotationDays:            getenv("ADMIN_API_KEY_ROTATION_DAYS", "90"),
		AdminRateLimit:                     getenv("ADMIN_RATE_LIMIT", "200"),
		AdminAuditEnabled:                  getenv("ADMIN_AUDIT_ENABLED", "false"),
		AsyncJobsEnabled:                   getenv("ASYNC_JOBS_ENABLED", "true"),
		MaxConcurrentJobs:                  getenv("MAX_CONCURRENT_JOBS", "3"),
		JobTimeout:                         getenv("JOB_TIMEOUT", "30m"),
	}
	return cfg
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
