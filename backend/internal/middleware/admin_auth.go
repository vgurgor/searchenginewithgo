package middleware

import (
	"github.com/gin-gonic/gin"
	"search_engine/internal/api"
)

func AdminAuthMiddleware(enabled bool, apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !enabled {
			api.SendError(c, api.ErrForbidden().WithDetails("reason", "admin API is disabled"))
			return
		}
		key := c.GetHeader("X-API-Key")
		if key == "" {
			api.SendError(c, api.ErrUnauthorized().WithDetails("reason", "X-API-Key header is required"))
			return
		}
		if key != apiKey || apiKey == "" {
			api.SendError(c, api.ErrInvalidAPIKey())
			return
		}
		c.Next()
	}
}


