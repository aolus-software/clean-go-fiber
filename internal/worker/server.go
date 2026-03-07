package worker

import (
	"context"
	"fmt"

	"github.com/aolus-software/clean-go-fiber/internal/config"
	"github.com/aolus-software/clean-go-fiber/internal/libs/mailer"
	"github.com/aolus-software/clean-go-fiber/internal/queue/tasks"
	"github.com/aolus-software/clean-go-fiber/internal/worker/handlers"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

func NewServer(log *zap.Logger) *asynq.Server {
	cfg := config.C.Redis
	return asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
			Password: cfg.Password,
		},
		asynq.Config{
			Concurrency: 10,
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Error("task failed", zap.String("type", task.Type()), zap.Error(err))
			}),
		},
	)
}

func RegisterHandlers(mux *asynq.ServeMux, m *mailer.Mailer, log *zap.Logger) {
	emailHandler := handlers.NewEmailHandler(m, log)
	mux.HandleFunc(tasks.TypeEmailSend, emailHandler.HandleEmailSend)
}
