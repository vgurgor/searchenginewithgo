package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	"search_engine/internal/api/handlers"
	"search_engine/internal/config"
	"search_engine/internal/infrastructure/cache"
	"search_engine/internal/infrastructure/database"
	"search_engine/internal/middleware"
	"search_engine/pkg/logger"
	_ "search_engine/docs"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()
	log := logger.NewLogger(cfg.LogLevel)
	defer func(l *zap.Logger) { _ = l.Sync() }(log)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbPool, err := database.ConnectPostgres(ctx, cfg.DatabaseURL)
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
	router.Use(middleware.ErrorHandler(log))

	healthHandler := handlers.NewHealthHandler(dbPool, redisClient)
	router.GET("/health", healthHandler)

	// Swagger UI served by gin-swagger at /swagger/index.html
	router.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler))

	addr := ":" + cfg.APIPort
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	log.Info("API starting", zap.String("addr", addr))
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("server error", zap.Error(err))
	}
}


