package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewHealthHandler(db *pgxpool.Pool, redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
		defer cancel()

		dbOK := true
		if err := db.Ping(ctx); err != nil {
			dbOK = false
		}

		redisOK := true
		if err := redisClient.Ping(ctx).Err(); err != nil {
			redisOK = false
		}

		status := http.StatusOK
		if !dbOK || !redisOK {
			status = http.StatusServiceUnavailable
		}

		c.JSON(status, gin.H{
			"status":    "ok",
			"postgres":  dbOK,
			"redis":     redisOK,
			"timestamp": time.Now().UTC(),
		})
	}
}


