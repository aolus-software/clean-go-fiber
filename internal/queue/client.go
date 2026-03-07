package queue

import (
	"fmt"

	"github.com/aolus-software/clean-go-fiber/internal/config"
	"github.com/hibiken/asynq"
)

func NewClient() *asynq.Client {
	cfg := config.C.Redis
	return asynq.NewClient(asynq.RedisClientOpt{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
	})
}
