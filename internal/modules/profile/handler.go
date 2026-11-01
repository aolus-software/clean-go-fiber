package profile

import (
	apperrors "github.com/aolus-software/clean-go-fiber/internal/errors"
	"github.com/aolus-software/clean-go-fiber/internal/libs/utils"
	"github.com/aolus-software/clean-go-fiber/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// @Summary     Get own profile
// @Description Returns the caller's identity, roles, and permissions.
// @Tags        Profile
// @Produce     json
// @Security    BearerAuth
// @Success     200  {object}  utils.Response{data=repositories.UserInformation}
// @Failure     401  {object}  utils.Response
// @Router      /profile [get]
func (h *Handler) Get(c *fiber.Ctx) error {
	user, err := middleware.CurrentUser(c)
	if err != nil {
		return err
	}
	return utils.OK(c, "Profile retrieved", user)
}

// @Summary     Update own profile
// @Tags        Profile
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       body  body      profile.UpdateProfileRequest  true  "Profile"
// @Success     200   {object}  utils.Response{data=repositories.UserInformation}
// @Failure     400   {object}  utils.ErrorResponse
// @Failure     401   {object}  utils.Response
// @Failure     422   {object}  utils.ErrorResponse
// @Router      /profile [patch]
func (h *Handler) Update(c *fiber.Ctx) error {
	user, err := middleware.CurrentUser(c)
	if err != nil {
		return err
	}
	var req UpdateProfileRequest
	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}
	userID, err := uuid.Parse(user.ID)
	if err != nil {
		return apperrors.NewUnauthorizedError("Unauthorized")
	}
	result, err := h.service.Update(c.Context(), userID, &req)
	if err != nil {
		return err
	}
	return utils.OK(c, "Profile updated", result)
}

// @Summary     Change own password
// @Description Checks the current password, then stores the new one and clears the cached user information.
// @Tags        Profile
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       body  body      profile.UpdatePasswordRequest  true  "Passwords"
// @Success     200   {object}  utils.Response
// @Failure     400   {object}  utils.ErrorResponse
// @Failure     401   {object}  utils.Response
// @Failure     422   {object}  utils.ErrorResponse
// @Router      /profile/password [patch]
func (h *Handler) UpdatePassword(c *fiber.Ctx) error {
	user, err := middleware.CurrentUser(c)
	if err != nil {
		return err
	}
	var req UpdatePasswordRequest
	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}
	userID, err := uuid.Parse(user.ID)
	if err != nil {
		return apperrors.NewUnauthorizedError("Unauthorized")
	}
	if err := h.service.UpdatePassword(c.Context(), userID, &req); err != nil {
		return err
	}
	return utils.OK(c, "Password updated", nil)
}
