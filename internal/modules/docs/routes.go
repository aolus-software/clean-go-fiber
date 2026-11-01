package docs

import "github.com/gofiber/fiber/v2"

func (h *Handler) RegisterRoutes(router fiber.Router) {
	if !h.cfg.App.EnableAPIDocs {
		return
	}
	router.Get("/docs", h.UI)
	router.Get("/docs/openapi.json", h.Spec)
}
