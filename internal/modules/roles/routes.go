package roles

import (
	"github.com/aolus-software/clean-go-fiber/internal/guards"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) RegisterRoutes(settings fiber.Router, auth fiber.Handler) {
	r := settings.Group("/roles", auth)
	r.Get("", guards.Permission("role list"), h.List)
	r.Post("", guards.Permission("role create"), h.Create)
	r.Get("/:id", guards.Permission("role detail"), h.Detail)
	r.Patch("/:id", guards.Permission("role edit"), h.Update)
	r.Delete("/:id", guards.Permission("role delete"), h.Delete)
}
