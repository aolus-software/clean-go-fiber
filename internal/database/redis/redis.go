package redis

import (
	"context"
	"time"

	"github.com/aolus-software/clean-go-fiber/internal/config"
	"github.com/redis/go-redis/v9"
)

func Connect() (*redis.Client, error) {
	cfg := config.C.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, err
	}
	return client, nil
}
