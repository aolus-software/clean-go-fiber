package docs

import "github.com/gofiber/fiber/v2"

func (h *Handler) RegisterRoutes(app *fiber.App) {
	app.Get("/openapi.json", h.Spec)
	app.Get("/docs", h.UI)
}
