package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"search_engine/internal/api"
	"search_engine/internal/config"
	"search_engine/internal/domain/entities"
	"search_engine/internal/infrastructure/jobs"
	"search_engine/internal/infrastructure/services"
	"search_engine/internal/middleware"
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
			api.SendError(c, api.ErrInvalidParameter("body", "invalid JSON format"))
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
			providerID := ""
			if body.ProviderID != nil {
				providerID = strings.TrimSpace(*body.ProviderID)
			}
			go func(providerID string) {
				ctx, cancel := jobContext(h.Config.JobTimeout)
				defer cancel()
				h.JobMgr.Update(j.ID, jobs.JobRunning, 0, nil)
				var err error
				if providerID != "" {
					_, err = h.SyncSvc.SyncProvider(ctx, providerID)
				} else {
					_, err = h.SyncSvc.SyncAllProviders(ctx)
				}
				if err != nil {
					msg := err.Error()
					h.JobMgr.Update(j.ID, jobs.JobFailed, 100, &msg)
					h.Logger.Error("async sync failed", zap.Error(err))
				} else {
					h.JobMgr.Update(j.ID, jobs.JobCompleted, 100, nil)
				}
			}(providerID)
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
		performRecalc := func(ctx context.Context) error {
			switch {
			case body.ContentID != nil && *body.ContentID > 0:
				_, err := h.ScoreCalc.RecalculateScore(ctx, *body.ContentID)
				return err
			case body.RecalculateAll != nil && *body.RecalculateAll:
				return recalcAllContents(ctx, h.ScoreCalc, 100)
			case body.ContentType != nil && *body.ContentType != "":
				t := strings.ToLower(*body.ContentType)
				var ct entities.ContentType
				if t == "video" {
					ct = entities.ContentTypeVideo
				} else if t == "text" {
					ct = entities.ContentTypeText
				} else {
					return errors.New("invalid content_type")
				}
				return recalcByType(ctx, h.ScoreCalc, ct, 100)
			default:
				return errors.New("no recalculation scope provided")
			}
		}
		if asyncEnabled {
			jobID := "recalc-" + uuid.NewString()
			j := h.JobMgr.CreateJob(jobID, "recalculate")
			go func() {
				ctx, cancel := jobContext(h.Config.JobTimeout)
				defer cancel()
				h.JobMgr.Update(j.ID, jobs.JobRunning, 0, nil)
				err := performRecalc(ctx)
				if err != nil {
					msg := err.Error()
					h.JobMgr.Update(j.ID, jobs.JobFailed, 100, &msg)
					h.Logger.Error("async score recalculation failed", zap.Error(err))
				} else {
					h.JobMgr.Update(j.ID, jobs.JobCompleted, 100, nil)
				}
			}()
			c.JSON(http.StatusAccepted, gin.H{"success": true, "message": "Score recalculation job started", "job_id": j.ID, "data": gin.H{"started_at": time.Now().UTC()}})
			return
		}
		if err := performRecalc(c.Request.Context()); err != nil {
			h.Logger.Error("score recalculation failed", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
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
			api.SendError(c, api.ErrInvalidParameter("id", "must be a valid positive integer"))
			return
		}
		if err := h.SyncSvc.Contents.SoftDelete(c.Request.Context(), id); err != nil {
			api.SendError(c, api.ErrInternal("Failed to delete content"))
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
		api.SendError(c, api.NewError(api.ErrCodeNotFound, "Job not found").WithDetails("job_id", id))
	})
}

func jobContext(timeout string) (context.Context, context.CancelFunc) {
	if d, err := time.ParseDuration(timeout); err == nil && d > 0 {
		return context.WithTimeout(context.Background(), d)
	}
	return context.WithCancel(context.Background())
}

func recalcAllContents(ctx context.Context, svc *services.ScoreCalculatorService, batch int) error {
	total, err := svc.Contents.CountAll(ctx)
	if err != nil {
		return err
	}
	for offset := 0; offset < int(total); offset += batch {
		ids, err := svc.Contents.ListIDs(ctx, offset, batch)
		if err != nil {
			return err
		}
		if len(ids) == 0 {
			break
		}
		for _, id := range ids {
			if _, err := svc.RecalculateScore(ctx, id); err != nil {
				return err
			}
		}
	}
	return nil
}

func recalcByType(ctx context.Context, svc *services.ScoreCalculatorService, ct entities.ContentType, batch int) error {
	offset := 0
	for {
		ids, err := svc.Contents.ListIDsByType(ctx, ct, offset, batch)
		if err != nil {
			return err
		}
		if len(ids) == 0 {
			break
		}
		for _, id := range ids {
			if _, err := svc.RecalculateScore(ctx, id); err != nil {
				return err
			}
		}
		offset += batch
	}
	return nil
}


