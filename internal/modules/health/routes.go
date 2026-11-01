package health

import "github.com/gofiber/fiber/v2"

func (h *Handler) RegisterRoutes(router fiber.Router) {
	router.Get("/", h.Root)
	router.Get("/health", h.Health)
	router.Get("/live", h.Live)
}
