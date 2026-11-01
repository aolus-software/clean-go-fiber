package auth

import "github.com/gofiber/fiber/v2"

func (h *Handler) RegisterRoutes(router fiber.Router) {
	r := router.Group("/auth")
	r.Post("/login", h.Login)
	r.Post("/register", h.Register)
	r.Post("/resend-verification", h.ResendVerification)
	r.Post("/verify-email", h.VerifyEmail)
	r.Post("/forgot-password", h.ForgotPassword)
	r.Post("/reset-password", h.ResetPassword)
	r.Post("/refresh", h.Refresh)
}
