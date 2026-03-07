// @title           Clean Go Fiber API
// @version         1.0
// @description     Production-ready Go REST API boilerplate
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization

package main

import (
	"errors"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	apperrors "github.com/aolus-software/clean-go-fiber/internal/errors"

	"github.com/aolus-software/clean-go-fiber/internal/bootstrap"
	"github.com/aolus-software/clean-go-fiber/internal/config"
	"github.com/aolus-software/clean-go-fiber/internal/database/postgres"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/zap"
)

func main() {
	config.Load()

	logger, _ := zap.NewProduction()
	if config.C.App.Env == "development" {
		logger, _ = zap.NewDevelopment()
	}
	defer logger.Sync()

	db, err := postgres.Connect()
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	logger.Info("Connected to PostgreSQL")

	if err := postgres.Migrate(); err != nil {
		logger.Fatal("Failed to run migrations", zap.Error(err))
	}
	logger.Info("Database migrations applied")

	app := fiber.New(fiber.Config{
		BodyLimit:    100 * 1024,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			var appErr *apperrors.AppError
			if errors.As(err, &appErr) {
				code = appErr.Code
			} else if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		},
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: strings.Join(config.C.Cors.Host, ","),
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: strings.Join(config.C.Cors.Method, ","),
	}))

	bootstrap.RegisterRoutes(app, db, config.C)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Info("Service running", zap.String("port", config.C.App.Port))
		if err := app.Listen(":" + config.C.App.Port); err != nil {
			logger.Fatal("Server error", zap.Error(err))
		}
	}()

	<-quit
	logger.Info("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}
	logger.Info("Server stopped")
}
