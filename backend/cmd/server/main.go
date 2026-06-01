package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/YASSERRMD/sql-sage/backend/internal/api"
	"github.com/YASSERRMD/sql-sage/backend/internal/auth"
	"github.com/YASSERRMD/sql-sage/backend/internal/config"
	"github.com/YASSERRMD/sql-sage/backend/internal/database"
	"github.com/YASSERRMD/sql-sage/backend/internal/middleware"
	"github.com/YASSERRMD/sql-sage/backend/internal/repositories"
	"github.com/YASSERRMD/sql-sage/backend/internal/services"
	"github.com/YASSERRMD/sql-sage/backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config: %v\n", err)
		os.Exit(1)
	}
	log := logger.New(cfg.LogLevel)

	db, err := database.Open(cfg, log)
	if err != nil {
		log.Error("db open", "err", err)
		os.Exit(1)
	}
	if err := database.AutoMigrate(db); err != nil {
		log.Error("auto migrate", "err", err)
		os.Exit(1)
	}

	jwtSvc := auth.NewService(cfg)
	userRepo := repositories.NewUserRepository(db)
	rtRepo := repositories.NewRefreshTokenRepository(db)
	authSvc := services.NewAuthService(userRepo, rtRepo, jwtSvc)
	authH := api.NewAuthHandler(authSvc)

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery(), requestLogger(log))
	r.Use(middleware.NewRateLimiter(cfg.RateLimitPerMin).Middleware())

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/readyz", readyz(db))

	v1 := r.Group("/api/v1")
	authH.Register(v1, middleware.AuthRequired(jwtSvc))

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.HTTPPort),
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Info("server starting", "port", cfg.HTTPPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("listen", "err", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}

func readyz(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sqlDB, err := db.DB()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "down"})
			return
		}
		if err := sqlDB.PingContext(c.Request.Context()); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "down"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	}
}

func requestLogger(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		log.Info("request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency_ms", time.Since(start).Milliseconds(),
		)
	}
}
