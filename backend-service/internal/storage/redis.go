package storage

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

func ConnectToRedis(ctx context.Context, addr string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		if err := client.Close(); err != nil {
			slog.Error("Failed to close Redis client", "error", err)
		}
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	return client, nil
}
