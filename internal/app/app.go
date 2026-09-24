package app

import (
	"context"
	"fmt"
	"log/slog"

	"semen_project/internal/config"
	"semen_project/internal/controllers"
	"semen_project/internal/repository"
	"semen_project/internal/routes"
	"semen_project/internal/storage"

	"github.com/gin-gonic/gin"
)

func Run(cfg *config.Config) error {
	ctx := context.Background()
	slog.Info("initializing application", "app_name", cfg.AppName)

	slog.Info("connecting to PostgreSQL database...")

	dbPool, err := storage.ConnectToPg(ctx, cfg.PG, cfg.AppName)
	if err != nil {
		slog.Error("failed to connect to PostgreSQL", "error", err)
		return fmt.Errorf("database connection failed: %w", err)
	}
	defer dbPool.Close()

	slog.Info("PostgreSQL connection established")

	redisClient, err := storage.ConnectToRedis(ctx, "redis:6379")
	if err != nil {
		slog.Error("failed to connect to Redis", "error", err)
		return fmt.Errorf("redis connection failed: %w", err)
	}
	defer redisClient.Close()

	slog.Info("Redis connection established")

	store := repository.NewStore(dbPool, redisClient)

	handlers := controllers.NewHandlers(store, cfg.JWTSecret)

	router := gin.Default()
	routes.SetupRoutes(router, handlers, cfg.JWTSecret)
	if err := router.Run(fmt.Sprintf(":%d", cfg.PublicApiPort)); err != nil {
		return fmt.Errorf("server run failed: %w", err)
	}
	return nil
}
