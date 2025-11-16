package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"search_engine/internal/api"
	"search_engine/internal/api/dto"
	"search_engine/internal/infrastructure/services"
)

func RegisterContentRoutes(router *gin.Engine, svc *services.ContentSearchService, defaultPageSize, maxPageSize int) {
	v1 := router.Group("/api/v1/contents")
	v1.GET("/search", func(c *gin.Context) {
		q := strings.TrimSpace(c.Query("q"))
		ct := strings.TrimSpace(c.Query("type"))
		sort := strings.TrimSpace(c.Query("sort"))
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
		req := dto.SearchRequest{
			Keyword: q, ContentType: ct, SortBy: sort, Page: page, PageSize: pageSize,
		}
		req.Normalize(1, defaultPageSize, maxPageSize)
		if req.PageSize < 1 || req.PageSize > maxPageSize {
			api.SendError(c, api.ErrInvalidParameter("page_size", "must be between 1 and 100"))
			return
		}
		items, total, err := svc.SearchContents(c.Request.Context(), req)
		if err != nil {
			api.SendError(c, api.ErrInvalidParameter("search_query", err.Error()))
			return
		}
		totalPages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
		c.JSON(http.StatusOK, dto.SearchResponse{
			Success: true,
			Data:    items,
			Pagination: dto.PaginationDTO{
				Page: req.Page, PageSize: req.PageSize, TotalItems: total, TotalPages: totalPages,
			},
		})
	})
	v1.GET("/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || id <= 0 {
			api.SendError(c, api.ErrInvalidParameter("id", "must be a valid positive integer"))
			return
		}
		item, err := svc.GetContentByID(c.Request.Context(), id)
		if err != nil {
			api.SendError(c, api.ErrInternal("Failed to retrieve content"))
			return
		}
		if item == nil {
			api.SendError(c, api.ErrContentNotFound(idStr))
			return
		}
		c.JSON(http.StatusOK, dto.APIContentResponse{Success: true, Data: item})
	})
	v1.GET("/stats", func(c *gin.Context) {
		stats, err := svc.GetStats(c.Request.Context())
		if err != nil {
			api.SendError(c, api.ErrInternal("Failed to retrieve statistics"))
			return
		}
		c.JSON(http.StatusOK, dto.StatsResponse{Success: true, Data: stats})
	})
}
