package roles

import (
	"github.com/aolus-software/clean-go-fiber/internal/guards"
	"github.com/aolus-software/clean-go-fiber/internal/middleware"
	"github.com/aolus-software/clean-go-fiber/internal/repositories"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) RegisterRoutes(router fiber.Router, userRepo repositories.UserRepository) {
	r := router.Group("/roles", middleware.Auth())
	r.Get("/", guards.PermissionGuard(userRepo, "read"), h.List)
	r.Post("/", guards.PermissionGuard(userRepo, "create"), h.Create)
	r.Get("/:id", guards.PermissionGuard(userRepo, "read"), h.Get)
	r.Patch("/:id", guards.PermissionGuard(userRepo, "update"), h.Update)
	r.Delete("/:id", guards.PermissionGuard(userRepo, "delete"), h.Delete)
	r.Post("/:id/permissions", guards.PermissionGuard(userRepo, "update"), h.SyncPermissions)
}
