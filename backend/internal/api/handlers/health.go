package handlers

import (
	"context"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type HealthStatus struct {
	Status        string                   `json:"status"`
	Timestamp     time.Time                `json:"timestamp"`
	UptimeSeconds float64                  `json:"uptime_seconds"`
	Version       string                   `json:"version"`
	Services      map[string]ServiceHealth `json:"services"`
	System        SystemInfo               `json:"system"`
}

type ServiceHealth struct {
	Healthy      bool    `json:"healthy"`
	ResponseTime float64 `json:"response_time_ms,omitempty"`
	Error        string  `json:"error,omitempty"`
	Message      string  `json:"message,omitempty"`
}

type SystemInfo struct {
	GoVersion     string  `json:"go_version"`
	NumGoroutines int     `json:"num_goroutines"`
	MemoryUsageMB float64 `json:"memory_usage_mb"`
	NumCPU        int     `json:"num_cpu"`
}

func NewHealthHandler(db *pgxpool.Pool, redisClient *redis.Client, appStart time.Time, version string, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		start := time.Now()

		health := HealthStatus{
			Timestamp:     time.Now().UTC(),
			UptimeSeconds: time.Since(appStart).Seconds(),
			Version:       version,
			Services:      make(map[string]ServiceHealth),
			System:        getSystemInfo(),
		}

		// Check PostgreSQL
		health.Services["postgres"] = checkPostgresHealth(ctx, db)

		// Check Redis
		health.Services["redis"] = checkRedisHealth(ctx, redisClient)

		// Check overall status
		allHealthy := true
		for _, service := range health.Services {
			if !service.Healthy {
				allHealthy = false
				break
			}
		}

		if allHealthy {
			health.Status = "healthy"
			c.JSON(http.StatusOK, health)
		} else {
			health.Status = "unhealthy"
			logger.Warn("health check failed",
				zap.Any("services", health.Services),
				zap.Float64("response_time_ms", float64(time.Since(start).Nanoseconds())/1e6))
			c.JSON(http.StatusServiceUnavailable, health)
		}
	}
}

func checkPostgresHealth(ctx context.Context, db *pgxpool.Pool) ServiceHealth {
	start := time.Now()

	// Test connection
	if err := db.Ping(ctx); err != nil {
		return ServiceHealth{
			Healthy:      false,
			ResponseTime: float64(time.Since(start).Nanoseconds()) / 1e6,
			Error:        "connection_failed",
			Message:      err.Error(),
		}
	}

	// Test simple query
	var result int
	err := db.QueryRow(ctx, "SELECT 1").Scan(&result)
	if err != nil {
		return ServiceHealth{
			Healthy:      false,
			ResponseTime: float64(time.Since(start).Nanoseconds()) / 1e6,
			Error:        "query_failed",
			Message:      err.Error(),
		}
	}

	return ServiceHealth{
		Healthy:      true,
		ResponseTime: float64(time.Since(start).Nanoseconds()) / 1e6,
		Message:      "OK",
	}
}

func checkRedisHealth(ctx context.Context, client *redis.Client) ServiceHealth {
	start := time.Now()

	result := client.Ping(ctx)
	if err := result.Err(); err != nil {
		return ServiceHealth{
			Healthy:      false,
			ResponseTime: float64(time.Since(start).Nanoseconds()) / 1e6,
			Error:        "connection_failed",
			Message:      err.Error(),
		}
	}

	return ServiceHealth{
		Healthy:      true,
		ResponseTime: float64(time.Since(start).Nanoseconds()) / 1e6,
		Message:      result.Val(),
	}
}

func getSystemInfo() SystemInfo {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return SystemInfo{
		GoVersion:     runtime.Version(),
		NumGoroutines: runtime.NumGoroutine(),
		MemoryUsageMB: float64(m.Alloc) / 1024 / 1024,
		NumCPU:        runtime.NumCPU(),
	}
}

// DetailedMetricsHandler provides comprehensive application metrics
func DetailedMetricsHandler(db *pgxpool.Pool, redisClient *redis.Client, appStart time.Time, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		metrics := gin.H{
			"timestamp":      time.Now().UTC(),
			"uptime_seconds": time.Since(appStart).Seconds(),
			"system":         getSystemInfo(),
		}

		// Database metrics
		if db != nil {
			stats := db.Stat()
			metrics["database"] = gin.H{
				"total_connections":        stats.TotalConns(),
				"idle_connections":         stats.IdleConns(),
				"acquired_connections":     stats.AcquiredConns(),
				"constructing_connections": stats.ConstructingConns(),
				"new_connections_count":    stats.NewConnsCount(),
			}
		}

		// Redis metrics (basic)
		if redisClient != nil {
			info := redisClient.Info(ctx, "memory")
			if info.Err() == nil {
				metrics["redis"] = gin.H{
					"connected": true,
					"info":      info.Val(),
				}
			} else {
				metrics["redis"] = gin.H{
					"connected": false,
					"error":     info.Err().Error(),
				}
			}
		}

		c.JSON(http.StatusOK, metrics)
	}
}
