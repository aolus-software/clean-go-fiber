package profile

import "github.com/gofiber/fiber/v2"

func (h *Handler) RegisterRoutes(router fiber.Router, auth fiber.Handler) {
	r := router.Group("/profile", auth)
	r.Get("", h.Get)
	r.Patch("", h.Update)
	r.Patch("/password", h.UpdatePassword)
}
