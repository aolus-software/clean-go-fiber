package health

import "github.com/gofiber/fiber/v2"

func (h *Handler) RegisterRoutes(app *fiber.App) {
	app.Get("/", h.Root)
	app.Get("/health", h.Health)
}
