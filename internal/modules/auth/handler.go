package auth

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

// @Summary     Register a new user
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body  body      auth.RegisterRequest  true  "Registration data"
// @Success     201   {object}  utils.Response{data=auth.UserResponse}
// @Failure     400   {object}  utils.Response
// @Failure     409   {object}  utils.Response
// @Failure     422   {object}  utils.Response
// @Router      /auth/register [post]
func (h *Handler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}
	if err := validator.Validate(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"errors":  err,
		})
	}
	result, err := h.service.Register(c.Context(), &req)
	if err != nil {
		return err
	}
	return utils.Created(c, result)
}

// @Summary     Login
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body  body      auth.LoginRequest  true  "Login credentials"
// @Success     200   {object}  utils.Response{data=auth.TokenResponse}
// @Failure     400   {object}  utils.Response
// @Failure     401   {object}  utils.Response
// @Failure     422   {object}  utils.Response
// @Router      /auth/login [post]
func (h *Handler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}
	if err := validator.Validate(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"errors":  err,
		})
	}
	result, err := h.service.Login(c.Context(), &req)
	if err != nil {
		return err
	}
	return utils.OK(c, result)
}

// @Summary     Refresh access token
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body  body      auth.RefreshTokenRequest  true  "Refresh token"
// @Success     200   {object}  utils.Response{data=auth.TokenResponse}
// @Failure     400   {object}  utils.Response
// @Failure     401   {object}  utils.Response
// @Failure     422   {object}  utils.Response
// @Router      /auth/refresh-token [post]
func (h *Handler) RefreshToken(c *fiber.Ctx) error {
	var req RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}
	if err := validator.Validate(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"errors":  err,
		})
	}
	result, err := h.service.RefreshToken(c.Context(), &req)
	if err != nil {
		return err
	}
	return utils.OK(c, result)
}

// @Summary     Request password reset email
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body  body      auth.ForgotPasswordRequest  true  "Email address"
// @Success     200   {object}  utils.Response
// @Failure     400   {object}  utils.Response
// @Failure     422   {object}  utils.Response
// @Router      /auth/forgot-password [post]
func (h *Handler) ForgotPassword(c *fiber.Ctx) error {
	var req ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}
	if err := validator.Validate(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"errors":  err,
		})
	}
	if err := h.service.ForgotPassword(c.Context(), &req); err != nil {
		return err
	}
	return utils.OK(c, fiber.Map{"message": "if the email exists, a reset link has been sent"})
}

// @Summary     Reset password using token
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body  body      auth.ResetPasswordRequest  true  "Reset token and new password"
// @Success     200   {object}  utils.Response
// @Failure     400   {object}  utils.Response
// @Failure     422   {object}  utils.Response
// @Router      /auth/reset-password [post]
func (h *Handler) ResetPassword(c *fiber.Ctx) error {
	var req ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}
	if err := validator.Validate(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"errors":  err,
		})
	}
	if err := h.service.ResetPassword(c.Context(), &req); err != nil {
		return err
	}
	return utils.OK(c, fiber.Map{"message": "password has been reset"})
}

// @Summary     Verify email address
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body  body      auth.VerifyEmailRequest  true  "Verification token"
// @Success     200   {object}  utils.Response
// @Failure     400   {object}  utils.Response
// @Failure     422   {object}  utils.Response
// @Router      /auth/verify-email [post]
func (h *Handler) VerifyEmail(c *fiber.Ctx) error {
	var req VerifyEmailRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}
	if err := validator.Validate(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"errors":  err,
		})
	}
	if err := h.service.VerifyEmail(c.Context(), &req); err != nil {
		return err
	}
	return utils.OK(c, fiber.Map{"message": "email verified"})
}

// @Summary     Resend email verification
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body  body      auth.ResendVerificationRequest  true  "Email address"
// @Success     200   {object}  utils.Response
// @Failure     400   {object}  utils.Response
// @Failure     422   {object}  utils.Response
// @Router      /auth/resend-verification [post]
func (h *Handler) ResendVerification(c *fiber.Ctx) error {
	var req ResendVerificationRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}
	if err := validator.Validate(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"errors":  err,
		})
	}
	if err := h.service.ResendVerification(c.Context(), &req); err != nil {
		return err
	}
	return utils.OK(c, fiber.Map{"message": "if the email exists and is unverified, a new verification email has been sent"})
}

// @Summary     Get current user
// @Tags        auth
// @Produce     json
// @Security    BearerAuth
// @Success     200  {object}  utils.Response{data=auth.UserResponse}
// @Failure     401  {object}  utils.Response
// @Router      /auth/me [get]
func (h *Handler) Me(c *fiber.Ctx) error {
	claims := c.Locals("claims").(*utils.Claims)
	result, err := h.service.Me(c.Context(), claims.UserID)
	if err != nil {
		return err
	}
	return utils.OK(c, result)
}
