package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminAuthMiddleware(enabled bool, apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !enabled {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "error": gin.H{"code": "UNAUTHORIZED", "message": "admin api disabled"}})
			return
		}
		key := c.GetHeader("X-API-Key")
		if key == "" || key != apiKey || apiKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "error": gin.H{"code": "UNAUTHORIZED", "message": "invalid api key"}})
			return
		}
		c.Next()
	}
}


