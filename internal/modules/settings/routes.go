package settings

import (
	"github.com/aolus-software/clean-go-fiber/internal/guards"
	"github.com/aolus-software/clean-go-fiber/internal/middleware"
	"github.com/aolus-software/clean-go-fiber/internal/repositories"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) RegisterRoutes(router fiber.Router, userRepo repositories.UserRepository) {
	s := router.Group("/settings", middleware.Auth())

	roles := s.Group("/roles")
	roles.Get("/", guards.PermissionGuard(userRepo, "read"), h.ListRoles)
	roles.Post("/", guards.PermissionGuard(userRepo, "create"), h.CreateRole)
	roles.Get("/:id", guards.PermissionGuard(userRepo, "read"), h.GetRole)
	roles.Patch("/:id", guards.PermissionGuard(userRepo, "update"), h.UpdateRole)
	roles.Delete("/:id", guards.PermissionGuard(userRepo, "delete"), h.DeleteRole)
	roles.Post("/:id/permissions", guards.PermissionGuard(userRepo, "update"), h.SyncRolePermissions)

	perms := s.Group("/permissions")
	perms.Get("/", guards.PermissionGuard(userRepo, "read"), h.ListPermissions)
	perms.Post("/", guards.PermissionGuard(userRepo, "create"), h.CreatePermission)
	perms.Get("/:id", guards.PermissionGuard(userRepo, "read"), h.GetPermission)
	perms.Patch("/:id", guards.PermissionGuard(userRepo, "update"), h.UpdatePermission)
	perms.Delete("/:id", guards.PermissionGuard(userRepo, "delete"), h.DeletePermission)

	users := s.Group("/users")
	users.Post("/:id/roles", guards.PermissionGuard(userRepo, "update"), h.SyncUserRoles)
}
