package app

import (
	"context"
	"fmt"
	"log/slog"

	"semen_project/internal/config"
	"semen_project/internal/controllers"
	"semen_project/internal/kafka"
	"semen_project/internal/repository"
	"semen_project/internal/routes"
	"semen_project/internal/storage"

	pb "semen_project/proto/gen/notificationpb"

	"github.com/gin-gonic/gin"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	redisClient, err := storage.ConnectToRedis(ctx, cfg.RedisAddr)
	if err != nil {
		slog.Error("failed to connect to Redis", "error", err)
		return fmt.Errorf("redis connection failed: %w", err)
	}
	defer func() {
		if err := redisClient.Close(); err != nil {
			slog.Error("Failed to close Redis client", "error", err)
		}
	}()
	slog.Info("Redis connection established")

	store := repository.NewStore(dbPool, redisClient)

	producer := kafka.NewMessageProducer(cfg.KafkaBroker)

	defer func() {
		if err := producer.Close(); err != nil {
			slog.Error("Failed to close Kafka producer", "error", err)
		}
	}()
	grpcConn, err := grpc.NewClient(
		cfg.NotificationGRPCAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("create gRPC client: %w", err)
	}
	defer grpcConn.Close()

	notificationClient := pb.NewNotificationServiceClient(grpcConn)
	handlers := controllers.NewHandlers(
		store,
		cfg.JWTSecret,
		producer,
		notificationClient,
	)
	router := gin.Default()
	routes.SetupRoutes(router, handlers, cfg.JWTSecret)
	if err := router.Run(fmt.Sprintf(":%d", cfg.PublicApiPort)); err != nil {
		return fmt.Errorf("server run failed: %w", err)
	}
	return nil
}
