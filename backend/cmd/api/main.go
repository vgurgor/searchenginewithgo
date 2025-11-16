package main

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	_ "search_engine/docs"
	"search_engine/internal/api"
	"search_engine/internal/api/handlers"
	"search_engine/internal/config"
	"search_engine/internal/domain/scoring"
	"search_engine/internal/infrastructure/cache"
	"search_engine/internal/infrastructure/database"
	"search_engine/internal/infrastructure/jobs"
	infraproviders "search_engine/internal/infrastructure/providers"
	"search_engine/internal/infrastructure/ratelimiter"
	"search_engine/internal/infrastructure/repository/postgres"
	"search_engine/internal/infrastructure/services"
	"search_engine/internal/middleware"
	pubmw "search_engine/middleware"
	"search_engine/pkg/logger"
)

var appStart = time.Now().UTC()

func main() {
	_ = godotenv.Load()

	cfg := config.Load()
	log := logger.NewLogger(cfg.LogLevel)
	defer func(l *zap.Logger) { _ = l.Sync() }(log)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	dbPool, err := database.ConnectPostgres(ctx, cfg.DatabaseURL)
	cancel()
	if err != nil {
		log.Fatal("postgres connection error", zap.Error(err))
	}
	defer dbPool.Close()

	redisClient, err := cache.NewRedisClient(cfg.RedisURL)
	if err != nil {
		log.Fatal("redis connection error", zap.Error(err))
	}
	defer func() { _ = redisClient.Close() }()

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.RequestLogger(log))
	router.Use(api.ErrorHandler(log))
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173", "*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-API-Key"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))

	// Public rate limit middleware (per IP)
	publicLimit, _ := strconv.Atoi(cfg.PublicRateLimit)
	publicWindow, _ := time.ParseDuration(cfg.PublicRateLimitWindow)
	publicRL := pubmw.PublicRateLimit(redisClient, publicLimit > 0, publicLimit, publicWindow)
	router.Use(publicRL)

	version := os.Getenv("APP_VERSION")
	if version == "" {
		version = "1.0.0"
	}
	healthHandler := handlers.NewHealthHandler(dbPool, redisClient, appStart, version, log)
	router.GET("/health", healthHandler)

	// Detailed metrics endpoint (for monitoring)
	metricsHandler := handlers.DetailedMetricsHandler(dbPool, redisClient, appStart, log)
	router.GET("/api/v1/admin/metrics/system", metricsHandler)

	// Mock providers
	router.GET("/mock/provider1/contents", handlers.MockProvider1Handler)
	router.GET("/mock/provider2/feed", handlers.MockProvider2Handler)

	// Swagger UI served by gin-swagger at /swagger/index.html
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Provider factory and service wiring (for use in future endpoints/jobs)
	providerTimeout, _ := time.ParseDuration(cfg.ProviderTimeout)
	factory := infraproviders.NewProviderFactory()
	jsonProvider := infraproviders.NewJSONProvider(cfg.Provider1BaseURL, providerTimeout)
	xmlProvider := infraproviders.NewXMLProvider(cfg.Provider2BaseURL, providerTimeout)
	factory.RegisterProvider(jsonProvider)
	factory.RegisterProvider(xmlProvider)
	rateLimiter := ratelimiter.NewRedisLimiter(redisClient, cfg.RateLimitEnabled == "true")
	providerSvc := &services.ProviderService{
		Factory: factory,
		Limiter: rateLimiter,
		Logger:  log,
		Timeout: providerTimeout,
	}

	// Scoring services wiring
	// Parse scoring configs
	videoMul, _ := strconv.ParseFloat(cfg.VideoTypeMultiplier, 64)
	textMul, _ := strconv.ParseFloat(cfg.TextTypeMultiplier, 64)
	fresh1w, _ := strconv.ParseFloat(cfg.Freshness1Week, 64)
	fresh1m, _ := strconv.ParseFloat(cfg.Freshness1Month, 64)
	fresh3m, _ := strconv.ParseFloat(cfg.Freshness3Months, 64)
	engine := &scoring.ScoringEngine{
		VideoTypeMultiplier: videoMul,
		TextTypeMultiplier:  textMul,
		Freshness: scoring.FreshnessConfig{
			WithinOneWeekScore:     fresh1w,
			WithinOneMonthScore:    fresh1m,
			WithinThreeMonthsScore: fresh3m,
		},
	}
	scoreCalc := &services.ScoreCalculatorService{
		Contents: postgres.NewContentRepository(dbPool),
		Metrics:  postgres.NewContentMetricsRepository(dbPool),
		Engine:   engine,
		Logger:   log,
	}
	// Optional background job
	if cfg.ScoreRecalcEnabled == "true" {
		recalcEvery, _ := time.ParseDuration(cfg.ScoreRecalcInterval)
		batchSize, _ := strconv.Atoi(cfg.ScoreBatchSize)
		job := jobs.NewScoreRecalculationJob(log, postgres.NewContentRepository(dbPool), scoreCalc, batchSize, recalcEvery)
		job.Start()
		defer job.Stop()
	}

	// Content Sync service and job
	thPercent, _ := strconv.Atoi(cfg.MetricsChangeThresholdPercent)
	thAbsViews, _ := strconv.Atoi(cfg.MetricsChangeThresholdAbsViews)
	thAbsLikes, _ := strconv.Atoi(cfg.MetricsChangeThresholdAbsLikes)
	thAbsReac, _ := strconv.Atoi(cfg.MetricsChangeThresholdAbsReactions)
	syncSvc := &services.ContentSyncService{
		Logger:         log,
		Factory:        factory,
		ProviderClient: providerSvc,
		Contents:       postgres.NewContentRepository(dbPool),
		Metrics:        postgres.NewContentMetricsRepository(dbPool),
		ScoreCalc:      scoreCalc,
		HistoryRepo:    postgres.NewSyncHistoryRepository(dbPool),
		Thresholds:     services.MetricsThresholds{Percent: thPercent, AbsViews: thAbsViews, AbsLikes: thAbsLikes, AbsReactions: thAbsReac},
	}
	if cfg.ContentSyncEnabled == "true" {
		syncEvery, _ := time.ParseDuration(cfg.ContentSyncInterval)
		retryCnt, _ := strconv.Atoi(cfg.ContentSyncRetryCount)
		retryDelay, _ := time.ParseDuration(cfg.ContentSyncRetryDelay)
		sjob := jobs.NewContentSyncJob(log, syncSvc, syncEvery, true, retryCnt, retryDelay)
		sjob.Start()
		defer sjob.Stop()
	}

	// Admin API (secured)
	jobMgr := jobs.NewJobManager()
	adminHandlers := &handlers.AdminHandlers{
		Logger:    log,
		Config:    cfg,
		SyncSvc:   syncSvc,
		ScoreCalc: scoreCalc,
		JobMgr:    jobMgr,
	}
	handlers.RegisterAdminRoutes(router, adminHandlers)

	// Content search endpoints
	defPage, _ := strconv.Atoi(cfg.DefaultPageSize)
	maxPage, _ := strconv.Atoi(cfg.MaxPageSize)
	searchCacheTTL, _ := time.ParseDuration(cfg.SearchCacheTTL)
	searchSvc := &services.ContentSearchService{
		Repo:            postgres.NewContentRepository(dbPool),
		HistoryRepo:     postgres.NewSyncHistoryRepository(dbPool),
		DefaultPageSize: defPage,
		MaxPageSize:     maxPage,
		CacheClient:     redisClient,
		CacheEnabled:    cfg.SearchCacheEnabled == "true",
		CacheTTL:        searchCacheTTL,
	}
	handlers.RegisterContentRoutes(router, searchSvc, defPage, maxPage)

	addr := ":" + cfg.APIPort
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	srv := &http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	log.Info("API starting", zap.String("addr", addr))
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("server error", zap.Error(err))
	}
}
