package main

import (
	"log"
	"time"

	"github.com/aolus-software/clean-go-fiber/internal/config"
	"github.com/gofiber/fiber/v2"
)

func main() {
	config.Load()

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

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "service": config.C.App.Name})
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "service": config.C.App.Name})
	})

	log.Fatal(app.Listen(":" + config.C.App.Port))
}