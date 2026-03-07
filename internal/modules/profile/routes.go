package profile

import (
	"github.com/aolus-software/clean-go-fiber/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) RegisterRoutes(router fiber.Router) {
	p := router.Group("/profile", middleware.Auth())
	p.Get("/me", h.GetProfile)
	p.Patch("/me", h.UpdateProfile)
	p.Patch("/me/password", h.ChangePassword)
}
