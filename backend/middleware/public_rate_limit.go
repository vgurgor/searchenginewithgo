package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// PublicRateLimit applies a simple IP-based rate limit per window across public endpoints.
// If enabled=false or rdb is nil or limit <=0, middleware becomes a no-op.
func PublicRateLimit(rdb *redis.Client, enabled bool, limit int, window time.Duration) gin.HandlerFunc {
	if !enabled || rdb == nil || limit <= 0 || window <= 0 {
		return func(c *gin.Context) { c.Next() }
	}
	return func(c *gin.Context) {
		ip := c.ClientIP()
		route := c.FullPath()
		if route == "" {
			route = c.Request.URL.Path
		}
		bucket := time.Now().UTC().Unix() / int64(window.Seconds())
		key := fmt.Sprintf("prl:%s:%s:%d", ip, route, bucket)
		// Use INCR with TTL to count requests in current window bucket.
		pipe := rdb.TxPipeline()
		inc := pipe.Incr(c, key)
		pipe.Expire(c, key, window+2*time.Second)
		_, _ = pipe.Exec(c)
		if inc.Val() > int64(limit) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "RATE_LIMIT_EXCEEDED",
					"message": "Too many requests. Please try again later.",
				},
			})
			return
		}
		c.Next()
	}
}


