package guards

import (
	"github.com/aolus-software/clean-go-fiber/internal/errors"
	"github.com/aolus-software/clean-go-fiber/internal/middleware"
	"github.com/aolus-software/clean-go-fiber/internal/models"
	"github.com/gofiber/fiber/v2"
)

func Role(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := middleware.CurrentUser(c)
		if err != nil {
			return err
		}
		if user.HasRole(models.SuperuserRole) {
			return c.Next()
		}
		for _, required := range roles {
			if !user.HasRole(required) {
				return errors.NewForbiddenError("Forbidden")
			}
		}
		return c.Next()
	}
}
