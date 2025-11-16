package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewHealthHandler(db *pgxpool.Pool, redisClient *redis.Client, appStart time.Time, version string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
		defer cancel()

		dbOK := true
		t0 := time.Now()
		dbErr := db.Ping(ctx)
		dbLatency := time.Since(t0).Milliseconds()
		if dbErr != nil {
			dbOK = false
		}

		redisOK := true
		t1 := time.Now()
		rErr := redisClient.Ping(ctx).Err()
		redisLatency := time.Since(t1).Milliseconds()
		if rErr != nil {
			redisOK = false
		}

		status := http.StatusOK
		if !dbOK || !redisOK {
			status = http.StatusServiceUnavailable
		}

		c.JSON(status, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().UTC(),
			"version":   version,
			"uptime_seconds": int(time.Since(appStart).Seconds()),
			"services": gin.H{
				"database": gin.H{
					"status":            ternaryStatus(dbOK),
					"response_time_ms":  dbLatency,
				},
				"redis": gin.H{
					"status":            ternaryStatus(redisOK),
					"response_time_ms":  redisLatency,
				},
			},
		})
	}
}

func ternaryStatus(ok bool) string {
	if ok {
		return "up"
	}
	return "down"
}


