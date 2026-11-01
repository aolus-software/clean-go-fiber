package queue

import (
	"github.com/aolus-software/clean-go-fiber/internal/config"
	"github.com/hibiken/asynq"
)

func RedisOpt() asynq.RedisClientOpt {
	cfg := config.C.Redis
	return asynq.RedisClientOpt{Addr: cfg.Addr(), Password: cfg.Password, DB: cfg.DB}
}

func NewClient() *asynq.Client {
	return asynq.NewClient(RedisOpt())
}
