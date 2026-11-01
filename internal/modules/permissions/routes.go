package permissions

import (
	"github.com/aolus-software/clean-go-fiber/internal/guards"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) RegisterRoutes(settings fiber.Router, auth fiber.Handler) {
	r := settings.Group("/permissions", auth)
	r.Get("", guards.Permission("permission list"), h.List)
	r.Post("", guards.Permission("permission create"), h.Create)
	r.Get("/:id", guards.Permission("permission detail"), h.Detail)
	r.Patch("/:id", guards.Permission("permission edit"), h.Update)
	r.Delete("/:id", guards.Permission("permission delete"), h.Delete)
}
