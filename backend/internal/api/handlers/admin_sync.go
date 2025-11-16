package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"search_engine/internal/infrastructure/services"
)

type AdminSyncRequest struct {
	ProviderID *string `json:"provider_id"`
	APIKey     string  `json:"api_key"`
}

func AdminSyncHandler(logger *zap.Logger, apiKey string, svc *services.ContentSyncService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req AdminSyncRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}
		if req.APIKey != apiKey || apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		if req.ProviderID == nil || *req.ProviderID == "" {
			results, err := svc.SyncAllProviders(c.Request.Context())
			if err != nil {
				logger.Error("admin sync all failed", zap.Error(err))
			}
			c.JSON(http.StatusOK, gin.H{"success": err == nil, "results": results})
			return
		}
		res, err := svc.SyncProvider(c.Request.Context(), *req.ProviderID)
		if err != nil {
			logger.Error("admin sync provider failed", zap.String("provider", *req.ProviderID), zap.Error(err))
		}
		c.JSON(http.StatusOK, gin.H{"success": err == nil, "results": []services.SyncResult{res}})
	}
}
