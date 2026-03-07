package auth

import (
	"github.com/aolus-software/clean-go-fiber/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) RegisterRoutes(router fiber.Router) {
	auth := router.Group("/auth")

	auth.Post("/register", h.Register)
	auth.Post("/login", h.Login)
	auth.Post("/refresh-token", h.RefreshToken)
	auth.Post("/forgot-password", h.ForgotPassword)
	auth.Post("/reset-password", h.ResetPassword)
	auth.Post("/verify-email", h.VerifyEmail)
	auth.Post("/resend-verification", h.ResendVerification)

	auth.Get("/me", middleware.Auth(), h.Me)
}
