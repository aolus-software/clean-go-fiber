// @title           Clean Go Fiber API
// @version         1.0
// @description     Production-ready Go REST API boilerplate
// @host            localhost:3000
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization

package main

import (
	"errors"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/aolus-software/clean-go-fiber/internal/bootstrap"
	"github.com/aolus-software/clean-go-fiber/internal/config"
	"github.com/aolus-software/clean-go-fiber/internal/database/postgres"
	redisdb "github.com/aolus-software/clean-go-fiber/internal/database/redis"
	apperrors "github.com/aolus-software/clean-go-fiber/internal/errors"
	"github.com/aolus-software/clean-go-fiber/internal/libs/utils"
	"github.com/aolus-software/clean-go-fiber/internal/logging"
	"github.com/aolus-software/clean-go-fiber/internal/middleware"
	"github.com/aolus-software/clean-go-fiber/internal/queue"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"go.uber.org/zap"
)

func main() {
	config.Load()

	logger, err := logging.New(config.C.App)
	if err != nil {
		panic(err)
	}
	defer func() { _ = logger.Sync() }()

	db, err := postgres.Connect()
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	logger.Info("Connected to PostgreSQL")

	if err := postgres.Migrate(); err != nil {
		logger.Fatal("Failed to run migrations", zap.Error(err))
	}
	logger.Info("Database migrations applied")

	redisClient, err := redisdb.Connect()
	if err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
	}
	logger.Info("Connected to Redis")

	queueClient := queue.NewClient()
	defer func() { _ = queueClient.Close() }()

	app := fiber.New(fiber.Config{
		AppName:      config.C.App.Name,
		BodyLimit:    100 * 1024,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		ErrorHandler: errorHandler(logger),
	})

	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(helmet.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: config.C.Cors.AllowedOrigins,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,PATCH,DELETE,HEAD,OPTIONS",
	}))
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return utils.Error(c, fiber.StatusTooManyRequests, "Too many requests", nil)
		},
	}))
	app.Use(middleware.Logger(logger))
	app.Use(middleware.Performance(logger, 1*time.Second))

	bootstrap.RegisterRoutes(app, db, redisClient, queueClient, config.C, logger)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		port := strconv.Itoa(config.C.App.Port)
		logger.Info("Service running", zap.String("port", port))
		if err := app.Listen(":" + port); err != nil {
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

func errorHandler(logger *zap.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			return utils.Error(c, appErr.Code, appErr.Message, appErr.Errors)
		}

		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			switch fiberErr.Code {
			case fiber.StatusNotFound:
				return utils.Error(c, fiber.StatusNotFound, "Route not found", nil)
			case fiber.StatusBadRequest, fiber.StatusUnprocessableEntity:
				return utils.Error(c, fiber.StatusBadRequest, "Invalid request body", []apperrors.FieldError{apperrors.Field("body", "request body could not be parsed")})
			}
			if fiberErr.Code < fiber.StatusInternalServerError {
				return utils.Error(c, fiberErr.Code, fiberErr.Message, nil)
			}
		}

		logger.Error("unhandled error",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Error(err),
		)
		return utils.Error(c, fiber.StatusInternalServerError, "Internal server error", nil)
	}
}
