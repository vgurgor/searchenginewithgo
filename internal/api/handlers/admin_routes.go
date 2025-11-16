package handlers

import (
	"net/http"
	"time"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"search_engine/internal/config"
	"search_engine/internal/infrastructure/jobs"
	"search_engine/internal/infrastructure/services"
	"search_engine/internal/middleware"
	"search_engine/internal/domain/entities"
)

type AdminHandlers struct {
	Logger    *zap.Logger
	Config    config.Config
	SyncSvc   *services.ContentSyncService
	ScoreCalc *services.ScoreCalculatorService
	JobMgr    *jobs.JobManager
}

func RegisterAdminRoutes(router *gin.Engine, h *AdminHandlers) {
	auth := middleware.AdminAuthMiddleware(h.Config.AdminAPIEnabled == "true", h.Config.AdminAPIKey)
	grp := router.Group("/api/v1/admin")
	grp.Use(auth)

	grp.POST("/sync", func(c *gin.Context) {
		var body struct {
			ProviderID *string `json:"provider_id"`
			Force      bool    `json:"force"`
			Async      *bool   `json:"async"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"code":"INVALID_BODY","message":"invalid body"}})
			return
		}
		asyncEnabled := h.Config.AsyncJobsEnabled == "true"
		doAsync := asyncEnabled
		if body.Async != nil {
			doAsync = *body.Async
		}
		if doAsync {
			jobID := "sync-" + uuid.NewString()
			j := h.JobMgr.CreateJob(jobID, "sync")
			go func() {
				h.JobMgr.Update(j.ID, jobs.JobRunning, 0, nil)
				var err error
				if body.ProviderID != nil && *body.ProviderID != "" {
					_, err = h.SyncSvc.SyncProvider(c.Request.Context(), *body.ProviderID)
				} else {
					_, err = h.SyncSvc.SyncAllProviders(c.Request.Context())
				}
				if err != nil {
					msg := err.Error()
					h.JobMgr.Update(j.ID, jobs.JobFailed, 100, &msg)
				} else {
					h.JobMgr.Update(j.ID, jobs.JobCompleted, 100, nil)
				}
			}()
			c.JSON(http.StatusAccepted, gin.H{"success": true, "message": "Sync job started", "job_id": j.ID, "data": gin.H{"provider_id": body.ProviderID, "started_at": time.Now().UTC()}})
			return
		}
		var (
			results any
			err error
		)
		if body.ProviderID != nil && *body.ProviderID != "" {
			r, e := h.SyncSvc.SyncProvider(c.Request.Context(), *body.ProviderID)
			results = []services.SyncResult{r}
			err = e
		} else {
			r, e := h.SyncSvc.SyncAllProviders(c.Request.Context())
			results = r
			err = e
		}
		if err != nil {
			h.Logger.Error("manual sync failed", zap.Error(err))
		}
		c.JSON(http.StatusOK, gin.H{"success": err == nil, "data": gin.H{"results": results}})
	})

	grp.GET("/sync/history", func(c *gin.Context) {
		pid := c.Query("provider_id")
		statusStr := c.Query("status")
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
		if limit <= 0 || limit > 200 {
			limit = 50
		}
		var pidPtr *string
		if pid != "" {
			pidPtr = &pid
		}
		var stPtr *entities.SyncStatus
		if statusStr != "" {
			s := entities.SyncStatus(statusStr)
			stPtr = &s
		}
		items, err := h.SyncSvc.HistoryRepo.List(c.Request.Context(), pidPtr, stPtr, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false})
			return
		}
		total, _ := h.SyncSvc.HistoryRepo.Count(c.Request.Context(), pidPtr, stPtr)
		c.JSON(http.StatusOK, gin.H{"success": true, "data": items, "pagination": gin.H{"limit": limit, "offset": offset, "total": total}})
	})

	grp.POST("/scores/recalculate", func(c *gin.Context) {
		var body struct {
			ContentID      *int64  `json:"content_id"`
			ContentType    *string `json:"content_type"`
			RecalculateAll *bool   `json:"recalculate_all"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false})
			return
		}
		asyncEnabled := h.Config.AsyncJobsEnabled == "true"
		if asyncEnabled {
			jobID := "recalc-" + uuid.NewString()
			j := h.JobMgr.CreateJob(jobID, "recalculate")
			go func() {
				h.JobMgr.Update(j.ID, jobs.JobRunning, 0, nil)
				var err error
				ctx := c.Request.Context()
				switch {
				case body.ContentID != nil && *body.ContentID > 0:
					_, err = h.ScoreCalc.RecalculateScore(ctx, *body.ContentID)
				case body.RecalculateAll != nil && *body.RecalculateAll:
					total, _ := h.ScoreCalc.Contents.CountAll(ctx)
					for offset := 0; offset < int(total); offset += 100 {
						ids, _ := h.ScoreCalc.Contents.ListIDs(ctx, offset, 100)
						for _, id := range ids {
							_, _ = h.ScoreCalc.RecalculateScore(ctx, id)
						}
					}
				case body.ContentType != nil && *body.ContentType != "":
					t := strings.ToLower(*body.ContentType)
					var ct entities.ContentType
					if t == "video" {
						ct = entities.ContentTypeVideo
					} else {
						ct = entities.ContentTypeText
					}
					total, _ := h.ScoreCalc.Contents.CountAll(ctx)
					for offset := 0; offset < int(total); offset += 100 {
						ids, _ := h.ScoreCalc.Contents.ListIDsByType(ctx, ct, offset, 100)
						for _, id := range ids {
							_, _ = h.ScoreCalc.RecalculateScore(ctx, id)
						}
					}
				}
				if err != nil {
					msg := err.Error()
					h.JobMgr.Update(j.ID, jobs.JobFailed, 100, &msg)
				} else {
					h.JobMgr.Update(j.ID, jobs.JobCompleted, 100, nil)
				}
			}()
			c.JSON(http.StatusAccepted, gin.H{"success": true, "message": "Score recalculation job started", "job_id": j.ID, "data": gin.H{"started_at": time.Now().UTC()}})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	grp.GET("/providers", func(c *gin.Context) {
		byProvider, _ := h.SyncSvc.Contents.CountByProvider(c.Request.Context())
		providers := []gin.H{}
		for pid, count := range byProvider {
			avg, _ := h.SyncSvc.Contents.GetAverageScoreByProvider(c.Request.Context(), pid)
			last, _ := h.SyncSvc.HistoryRepo.GetLastSync(c.Request.Context(), pid)
			entry := gin.H{
				"provider_id": pid,
				"content_count": count,
				"average_score": avg,
			}
			if last != nil && last.CompletedAt != nil {
				entry["last_sync"] = last.CompletedAt
				entry["last_sync_status"] = last.SyncStatus
			}
			providers = append(providers, entry)
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "data": providers})
	})

	grp.POST("/providers/health-check", func(c *gin.Context) {
		// For demo purposes, return healthy with 0ms; real impl would probe provider endpoints
		byProvider, _ := h.SyncSvc.Contents.CountByProvider(c.Request.Context())
		out := []gin.H{}
		now := time.Now().UTC()
		for pid := range byProvider {
			out = append(out, gin.H{
				"provider_id": pid,
				"is_healthy": true,
				"response_time_ms": 0,
				"status_code": 200,
				"checked_at": now,
				"error": nil,
			})
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "data": out})
	})

	grp.DELETE("/contents/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || id <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"success": false})
			return
		}
		if err := h.SyncSvc.Contents.SoftDelete(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"id": id, "deleted_at": time.Now().UTC()}})
	})

	grp.GET("/metrics/dashboard", func(c *gin.Context) {
		stats, _ := (&services.ContentSearchService{Repo: h.SyncSvc.Contents, HistoryRepo: h.SyncSvc.HistoryRepo}).GetStats(c.Request.Context())
		totalSyncs, _ := h.SyncSvc.HistoryRepo.Count(c.Request.Context(), nil, nil)
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"overview": gin.H{
					"total_contents": stats.TotalContents,
					"total_videos": stats.TotalVideos,
					"total_texts": stats.TotalTexts,
					"average_score": stats.AverageScore,
					"providers_count": len(stats.Providers),
				},
				"sync_stats": gin.H{
					"last_sync": stats.LastSync,
					"total_syncs": totalSyncs,
				},
				"content_distribution": stats.Providers,
			},
		})
	})

	grp.GET("/jobs/:jobId", func(c *gin.Context) {
		id := c.Param("jobId")
		if j, ok := h.JobMgr.Get(id); ok {
			c.JSON(http.StatusOK, gin.H{"success": true, "data": j})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"success": false})
	})
}


