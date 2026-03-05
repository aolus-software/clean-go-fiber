package main

import (
	"log"
	"time"

	"github.com/aolus-software/clean-go-fiber/internal/config"
	"github.com/aolus-software/clean-go-fiber/internal/database/postgres"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func main() {
	config.Load()

	// Logger
	logger, _ := zap.NewProduction()
	if config.C.App.Env == "development" {
		logger, _ = zap.NewDevelopment()
	}
	defer logger.Sync()

	// Connect to database
	if _, err := postgres.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	logger.Info("Connected to PostgreSQL")

	// Run database migrations
	if err := postgres.Migrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	logger.Info("Database migrations applied")

	app := fiber.New(fiber.Config{
		BodyLimit:    100 * 1024,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		},
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "service": config.C.App.Name})
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "service": config.C.App.Name})
	})

	logger.Info("Service running on port" + config.C.App.Port)
	log.Fatal(app.Listen(":" + config.C.App.Port))
}
