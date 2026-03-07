package profile

import (
	"github.com/aolus-software/clean-go-fiber/internal/libs/utils"
	"github.com/aolus-software/clean-go-fiber/internal/libs/validator"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// @Summary     Get current user profile
// @Tags        profile
// @Produce     json
// @Security    BearerAuth
// @Success     200  {object}  utils.Response{data=profile.ProfileResponse}
// @Failure     401  {object}  utils.Response
// @Router      /profile/me [get]
func (h *Handler) GetProfile(c *fiber.Ctx) error {
	claims := c.Locals("claims").(*utils.Claims)
	result, err := h.service.GetProfile(c.Context(), claims.UserID)
	if err != nil {
		return err
	}
	return utils.OK(c, result)
}

// @Summary     Update current user profile
// @Tags        profile
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       body  body      profile.UpdateProfileRequest  true  "Profile data"
// @Success     200   {object}  utils.Response{data=profile.ProfileResponse}
// @Failure     400   {object}  utils.Response
// @Failure     401   {object}  utils.Response
// @Failure     422   {object}  utils.Response
// @Router      /profile/me [patch]
func (h *Handler) UpdateProfile(c *fiber.Ctx) error {
	claims := c.Locals("claims").(*utils.Claims)
	var req UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}
	if err := validator.Validate(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"errors":  err,
		})
	}
	result, err := h.service.UpdateProfile(c.Context(), claims.UserID, &req)
	if err != nil {
		return err
	}
	return utils.OK(c, result)
}

// @Summary     Change current user password
// @Tags        profile
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       body  body      profile.ChangePasswordRequest  true  "Password change data"
// @Success     200   {object}  utils.Response
// @Failure     400   {object}  utils.Response
// @Failure     401   {object}  utils.Response
// @Failure     422   {object}  utils.Response
// @Router      /profile/me/password [patch]
func (h *Handler) ChangePassword(c *fiber.Ctx) error {
	claims := c.Locals("claims").(*utils.Claims)
	var req ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}
	if err := validator.Validate(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"errors":  err,
		})
	}
	if err := h.service.ChangePassword(c.Context(), claims.UserID, &req); err != nil {
		return err
	}
	return utils.OK(c, fiber.Map{"message": "password changed"})
}
