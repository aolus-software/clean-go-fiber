package selectoptions

import (
	"github.com/aolus-software/clean-go-fiber/internal/guards"
	"github.com/aolus-software/clean-go-fiber/internal/models"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) RegisterRoutes(settings fiber.Router, auth fiber.Handler) {
	r := settings.Group("/select-options", auth)
	r.Get("/permissions", guards.Role(models.SuperuserRole), h.Permissions)
	r.Get("/roles", guards.Role(models.SuperuserRole), h.Roles)
}
