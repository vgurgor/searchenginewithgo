package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ErrorHandler(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors[0]
			log.Error("request error", zap.String("path", c.Request.URL.Path), zap.Error(err))
			status := c.Writer.Status()
			if status < 400 {
				status = http.StatusInternalServerError
			}
			c.JSON(status, gin.H{
				"error":   err.Error(),
				"status":  status,
				"path":    c.Request.URL.Path,
				"method":  c.Request.Method,
				"request": "failed",
			})
			return
		}
	}
}


