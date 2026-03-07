package guards

import (
	"github.com/aolus-software/clean-go-fiber/internal/errors"
	"github.com/aolus-software/clean-go-fiber/internal/libs/utils"
	"github.com/aolus-software/clean-go-fiber/internal/repositories"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func RoleGuard(userRepo repositories.UserRepository, roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, ok := c.Locals("claims").(*utils.Claims)
		if !ok || claims == nil {
			return errors.NewUnauthorizedError("unauthorized")
		}

		id, err := uuid.Parse(claims.UserID)
		if err != nil {
			return errors.NewUnauthorizedError("unauthorized")
		}

		user, err := userRepo.FindById(c.Context(), id)
		if err != nil {
			return errors.NewUnauthorizedError("unauthorized")
		}

		for _, role := range user.Roles {
			if role.Name == "superuser" {
				return c.Next()
			}
		}

		for _, required := range roles {
			for _, role := range user.Roles {
				if role.Name == required {
					return c.Next()
				}
			}
		}

		return errors.NewForbiddenError("insufficient role")
	}
}
