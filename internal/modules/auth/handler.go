package auth

import (
	"github.com/aolus-software/clean-go-fiber/internal/libs/utils"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// @Summary     Login
// @Tags        Authentication
// @Accept      json
// @Produce     json
// @Param       body  body      auth.LoginRequest  true  "Credentials"
// @Success     200   {object}  utils.Response{data=auth.LoginResponse}
// @Failure     400   {object}  utils.ErrorResponse
// @Failure     422   {object}  utils.ErrorResponse
// @Failure     429   {object}  utils.Response
// @Router      /auth/login [post]
func (h *Handler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}
	result, err := h.service.Login(c.Context(), &req)
	if err != nil {
		return err
	}
	return utils.OK(c, "Login successful", result)
}

// @Summary     Register a new account
// @Tags        Authentication
// @Accept      json
// @Produce     json
// @Param       body  body      auth.RegisterRequest  true  "Registration data"
// @Success     201   {object}  utils.Response
// @Failure     400   {object}  utils.ErrorResponse
// @Failure     422   {object}  utils.ErrorResponse
// @Router      /auth/register [post]
func (h *Handler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}
	if err := h.service.Register(c.Context(), &req); err != nil {
		return err
	}
	return utils.Created(c, "Registration successful, please verify your email", nil)
}

// @Summary     Resend the verification email
// @Description Always answers 200, whether or not the email is registered or already verified.
// @Tags        Authentication
// @Accept      json
// @Produce     json
// @Param       body  body      auth.EmailRequest  true  "Email"
// @Success     200   {object}  utils.Response
// @Failure     422   {object}  utils.ErrorResponse
// @Router      /auth/resend-verification [post]
func (h *Handler) ResendVerification(c *fiber.Ctx) error {
	var req EmailRequest
	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}
	if err := h.service.ResendVerification(c.Context(), &req); err != nil {
		return err
	}
	return utils.OK(c, "If the email is registered and unverified, a verification email has been sent", nil)
}

// @Summary     Verify an email address
// @Tags        Authentication
// @Accept      json
// @Produce     json
// @Param       body  body      auth.TokenRequest  true  "Verification token"
// @Success     200   {object}  utils.Response
// @Failure     400   {object}  utils.ErrorResponse
// @Failure     422   {object}  utils.ErrorResponse
// @Router      /auth/verify-email [post]
func (h *Handler) VerifyEmail(c *fiber.Ctx) error {
	var req TokenRequest
	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}
	if err := h.service.VerifyEmail(c.Context(), &req); err != nil {
		return err
	}
	return utils.OK(c, "Email verified", nil)
}

// @Summary     Request a password reset email
// @Description Always answers 200, whether or not the email is registered.
// @Tags        Authentication
// @Accept      json
// @Produce     json
// @Param       body  body      auth.EmailRequest  true  "Email"
// @Success     200   {object}  utils.Response
// @Failure     422   {object}  utils.ErrorResponse
// @Router      /auth/forgot-password [post]
func (h *Handler) ForgotPassword(c *fiber.Ctx) error {
	var req EmailRequest
	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}
	if err := h.service.ForgotPassword(c.Context(), &req); err != nil {
		return err
	}
	return utils.OK(c, "If the email is registered, a password reset email has been sent", nil)
}

// @Summary     Reset a password with a reset token
// @Tags        Authentication
// @Accept      json
// @Produce     json
// @Param       body  body      auth.ResetPasswordRequest  true  "Reset token and new password"
// @Success     200   {object}  utils.Response
// @Failure     400   {object}  utils.ErrorResponse
// @Failure     422   {object}  utils.ErrorResponse
// @Router      /auth/reset-password [post]
func (h *Handler) ResetPassword(c *fiber.Ctx) error {
	var req ResetPasswordRequest
	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}
	if err := h.service.ResetPassword(c.Context(), &req); err != nil {
		return err
	}
	return utils.OK(c, "Password has been reset", nil)
}

// @Summary     Exchange a refresh token for a new token pair
// @Tags        Authentication
// @Accept      json
// @Produce     json
// @Param       body  body      auth.RefreshRequest  true  "Refresh token"
// @Success     200   {object}  utils.Response{data=auth.RefreshResponse}
// @Failure     401   {object}  utils.Response
// @Failure     422   {object}  utils.ErrorResponse
// @Router      /auth/refresh [post]
func (h *Handler) Refresh(c *fiber.Ctx) error {
	var req RefreshRequest
	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}
	result, err := h.service.Refresh(c.Context(), &req)
	if err != nil {
		return err
	}
	return utils.OK(c, "Token refreshed", result)
}
