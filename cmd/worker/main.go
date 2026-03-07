package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/aolus-software/clean-go-fiber/internal/config"
	"github.com/aolus-software/clean-go-fiber/internal/libs/mailer"
	"github.com/aolus-software/clean-go-fiber/internal/worker"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

func main() {
	config.Load()

	log, _ := zap.NewProduction()
	if config.C.App.Env == "development" {
		log, _ = zap.NewDevelopment()
	}
	defer log.Sync()

	m := mailer.New(&config.C.Mail)
	srv := worker.NewServer(log)
	mux := asynq.NewServeMux()
	worker.RegisterHandlers(mux, m, log)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Info("Worker starting")
		if err := srv.Run(mux); err != nil {
			log.Fatal("Worker error", zap.Error(err))
		}
	}()

	<-quit
	log.Info("Shutting down worker...")
	srv.Shutdown()
	log.Info("Worker stopped")
}
