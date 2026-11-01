package utils

import (
	apperrors "github.com/aolus-software/clean-go-fiber/internal/errors"
	"github.com/aolus-software/clean-go-fiber/internal/libs/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func ParseBody(c *fiber.Ctx, dst any) error {
	if err := c.BodyParser(dst); err != nil {
		return apperrors.NewBadRequestError("Invalid request body", apperrors.Field("body", "request body could not be parsed"))
	}
	return validator.Validate(dst)
}

func ParamUUID(c *fiber.Ctx, name string) (uuid.UUID, error) {
	id, err := uuid.Parse(c.Params(name))
	if err != nil {
		return uuid.Nil, apperrors.NewUnprocessableEntityError("Validation error", apperrors.Field(name, "must be a valid uuid"))
	}
	return id, nil
}
