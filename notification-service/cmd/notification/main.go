package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"os/signal"
	"syscall"

	"notification-service/internal/config"
	"notification-service/internal/controllers"
	grpcserver "notification-service/internal/grpc"
	"notification-service/internal/kafka"
	"notification-service/internal/repository"
	"notification-service/internal/storage"

	pb "semen_project/proto/gen/notificationpb"

	"google.golang.org/grpc"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		return
	}

	slog.Info("Configuration loaded", "app_name", cfg.AppName)

	dbPool, err := storage.ConnectToPostgres(ctx, cfg.PostgresDSN)
	if err != nil {
		slog.Error("Failed to connect to PostgreSQL", "error", err)
		return
	}
	defer dbPool.Close()

	slog.Info("PostgreSQL connection established")

	store := repository.NewStore(dbPool)
	handler := controllers.NewNotificationHandler(store)

	// Создаём TCP listener для gRPC.
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		slog.Error("Failed to listen on gRPC port", "error", err)
		return
	}

	// Создаём gRPC Server и регистрируем обработчик.
	grpcServer := grpc.NewServer()
	pb.RegisterNotificationServiceServer(
		grpcServer,
		grpcserver.NewServer(handler),
	)

	// Запускаем gRPC в отдельной горутине.
	grpcErrors := make(chan error, 1)

	go func() {
		slog.Info("gRPC server started", "port", 50051)
		grpcErrors <- grpcServer.Serve(listener)
	}()

	// Запускаем Kafka Consumer параллельно.
	kafkaErrors := make(chan error, 1)

	go func() {
		kafkaErrors <- kafka.StartConsumer(
			ctx,
			cfg.KafkaBroker,
			handler,
		)
	}()

	// Ожидаем остановки приложения или ошибки.
	select {
	case <-ctx.Done():
		slog.Info("Shutdown signal received")

	case err := <-grpcErrors:
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			slog.Error("gRPC server stopped", "error", err)
		}

	case err := <-kafkaErrors:
		if err != nil {
			slog.Error("Kafka consumer stopped", "error", err)
		}
	}

	// Останавливаем gRPC и Kafka.
	stop()
	grpcServer.GracefulStop()

	slog.Info("Notification service stopped")
}
