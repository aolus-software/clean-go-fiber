package middleware

import (
	"strings"

	"github.com/aolus-software/clean-go-fiber/internal/config"
	"github.com/aolus-software/clean-go-fiber/internal/errors"
	"github.com/aolus-software/clean-go-fiber/internal/libs/utils"
	"github.com/gofiber/fiber/v2"
)

func Auth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			return errors.NewUnauthorizedError("missing or invalid authorization header")
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := utils.ParseToken(tokenStr, config.C.Jwt.SecretKey)
		if err != nil {
			return errors.NewUnauthorizedError("invalid or expired token")
		}
		c.Locals("claims", claims)
		return c.Next()
	}
}
