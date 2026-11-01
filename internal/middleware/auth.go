package middleware

import (
	"errors"
	"strings"

	"github.com/aolus-software/clean-go-fiber/internal/config"
	apperrors "github.com/aolus-software/clean-go-fiber/internal/errors"
	"github.com/aolus-software/clean-go-fiber/internal/libs/cache"
	"github.com/aolus-software/clean-go-fiber/internal/libs/utils"
	"github.com/aolus-software/clean-go-fiber/internal/repositories"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const userLocalsKey = "user"

func Auth(userRepo repositories.UserRepository, c *cache.Cache) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		header := ctx.Get(fiber.HeaderAuthorization)
		if !strings.HasPrefix(header, "Bearer ") {
			return apperrors.NewUnauthorizedError("Unauthorized")
		}
		claims, err := utils.ParseToken(strings.TrimPrefix(header, "Bearer "), utils.TokenTypeAccess, config.C.Jwt.Secret)
		if err != nil {
			return apperrors.NewUnauthorizedError("Unauthorized")
		}
		userID, err := uuid.Parse(claims.Subject)
		if err != nil {
			return apperrors.NewUnauthorizedError("Unauthorized")
		}

		key := cache.Key(cache.UserInformationKey, userID.String())
		var info repositories.UserInformation
		err = c.Get(ctx.Context(), key, &info)
		if err == nil {
			ctx.Locals(userLocalsKey, &info)
			return ctx.Next()
		}
		if !errors.Is(err, redis.Nil) {
			return err
		}

		loaded, err := userRepo.FindUserInformation(ctx.Context(), userID)
		if err != nil {
			return err
		}
		if loaded == nil {
			return apperrors.NewUnauthorizedError("Unauthorized")
		}
		if err := c.Set(ctx.Context(), key, loaded, cache.UserInformationTTL); err != nil {
			return err
		}
		ctx.Locals(userLocalsKey, loaded)
		return ctx.Next()
	}
}

func CurrentUser(ctx *fiber.Ctx) (*repositories.UserInformation, error) {
	user, ok := ctx.Locals(userLocalsKey).(*repositories.UserInformation)
	if !ok || user == nil {
		return nil, apperrors.NewUnauthorizedError("Unauthorized")
	}
	return user, nil
}
