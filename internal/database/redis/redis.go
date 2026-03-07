package redis

import (
	"fmt"

	"github.com/aolus-software/clean-go-fiber/internal/config"
	"github.com/redis/go-redis/v9"
)

func Connect() (*redis.Client, error) {
	cfg := config.C.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       0,
	})
	return client, nil
}
