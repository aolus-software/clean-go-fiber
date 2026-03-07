package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func Performance(log *zap.Logger, threshold time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)

		if duration > threshold {
			log.Warn("slow request",
				zap.String("method", c.Method()),
				zap.String("path", c.Path()),
				zap.Duration("duration", duration),
			)
		}

		return err
	}
}
