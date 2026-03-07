package settings

import (
	"github.com/aolus-software/clean-go-fiber/internal/libs/pagination"
	"github.com/aolus-software/clean-go-fiber/internal/libs/utils"
	"github.com/aolus-software/clean-go-fiber/internal/libs/validator"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// @Summary     List roles
// @Tags        settings
// @Produce     json
// @Security    BearerAuth
// @Param       page            query  int     false  "Page"
// @Param       per_page        query  int     false  "Per page"
// @Param       search          query  string  false  "Search"
// @Param       order_by        query  string  false  "Order by"
// @Param       sort_direction  query  string  false  "Sort direction"
// @Success     200  {object}  utils.Response
// @Failure     401  {object}  utils.Response
// @Router      /settings/roles [get]
func (h *Handler) ListRoles(c *fiber.Ctx) error {
	q := pagination.ParseFromFiber(c)
	result, err := h.service.ListRoles(c.Context(), q)
	if err != nil {
		return err
	}
	return utils.OK(c, result)
}

// @Summary     Get role by ID
// @Tags        settings
// @Produce     json
// @Security    BearerAuth
// @Param       id   path      string  true  "Role ID"
// @Success     200  {object}  utils.Response{data=settings.RoleResponse}
// @Failure     401  {object}  utils.Response
// @Failure     404  {object}  utils.Response
// @Router      /settings/roles/{id} [get]
func (h *Handler) GetRole(c *fiber.Ctx) error {
	result, err := h.service.GetRole(c.Context(), c.Params("id"))
	if err != nil {
		return err
	}
	return utils.OK(c, result)
}

// @Summary     Create role
// @Tags        settings
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       body  body      settings.RoleRequest  true  "Role data"
// @Success     201   {object}  utils.Response{data=settings.RoleResponse}
// @Failure     401   {object}  utils.Response
// @Failure     422   {object}  utils.Response
// @Router      /settings/roles [post]
func (h *Handler) CreateRole(c *fiber.Ctx) error {
	var req RoleRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}
	if err := validator.Validate(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"errors":  err,
		})
	}
	result, err := h.service.CreateRole(c.Context(), &req)
	if err != nil {
		return err
	}
	return utils.Created(c, result)
}

// @Summary     Update role
// @Tags        settings
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id    path      string               true  "Role ID"
// @Param       body  body      settings.RoleRequest  true  "Role data"
// @Success     200   {object}  utils.Response{data=settings.RoleResponse}
// @Failure     401   {object}  utils.Response
// @Failure     404   {object}  utils.Response
// @Failure     422   {object}  utils.Response
// @Router      /settings/roles/{id} [patch]
func (h *Handler) UpdateRole(c *fiber.Ctx) error {
	var req RoleRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}
	if err := validator.Validate(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"errors":  err,
		})
	}
	result, err := h.service.UpdateRole(c.Context(), c.Params("id"), &req)
	if err != nil {
		return err
	}
	return utils.OK(c, result)
}

// @Summary     Delete role
// @Tags        settings
// @Produce     json
// @Security    BearerAuth
// @Param       id   path  string  true  "Role ID"
// @Success     204
// @Failure     401  {object}  utils.Response
// @Failure     404  {object}  utils.Response
// @Router      /settings/roles/{id} [delete]
func (h *Handler) DeleteRole(c *fiber.Ctx) error {
	if err := h.service.DeleteRole(c.Context(), c.Params("id")); err != nil {
		return err
	}
	return utils.NoContent(c)
}

// @Summary     Sync role permissions
// @Tags        settings
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id    path      string                          true  "Role ID"
// @Param       body  body      settings.SyncPermissionsRequest  true  "Permission IDs"
// @Success     200   {object}  utils.Response{data=settings.RoleResponse}
// @Failure     401   {object}  utils.Response
// @Failure     404   {object}  utils.Response
// @Failure     422   {object}  utils.Response
// @Router      /settings/roles/{id}/permissions [post]
func (h *Handler) SyncRolePermissions(c *fiber.Ctx) error {
	var req SyncPermissionsRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}
	if err := validator.Validate(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"errors":  err,
		})
	}
	result, err := h.service.SyncRolePermissions(c.Context(), c.Params("id"), &req)
	if err != nil {
		return err
	}
	return utils.OK(c, result)
}

// @Summary     List permissions
// @Tags        settings
// @Produce     json
// @Security    BearerAuth
// @Param       page            query  int     false  "Page"
// @Param       per_page        query  int     false  "Per page"
// @Param       search          query  string  false  "Search"
// @Param       order_by        query  string  false  "Order by"
// @Param       sort_direction  query  string  false  "Sort direction"
// @Success     200  {object}  utils.Response
// @Failure     401  {object}  utils.Response
// @Router      /settings/permissions [get]
func (h *Handler) ListPermissions(c *fiber.Ctx) error {
	q := pagination.ParseFromFiber(c)
	result, err := h.service.ListPermissions(c.Context(), q)
	if err != nil {
		return err
	}
	return utils.OK(c, result)
}

// @Summary     Get permission by ID
// @Tags        settings
// @Produce     json
// @Security    BearerAuth
// @Param       id   path      string  true  "Permission ID"
// @Success     200  {object}  utils.Response{data=settings.PermissionResponse}
// @Failure     401  {object}  utils.Response
// @Failure     404  {object}  utils.Response
// @Router      /settings/permissions/{id} [get]
func (h *Handler) GetPermission(c *fiber.Ctx) error {
	result, err := h.service.GetPermission(c.Context(), c.Params("id"))
	if err != nil {
		return err
	}
	return utils.OK(c, result)
}

// @Summary     Create permission
// @Tags        settings
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       body  body      settings.PermissionRequest  true  "Permission data"
// @Success     201   {object}  utils.Response{data=settings.PermissionResponse}
// @Failure     401   {object}  utils.Response
// @Failure     422   {object}  utils.Response
// @Router      /settings/permissions [post]
func (h *Handler) CreatePermission(c *fiber.Ctx) error {
	var req PermissionRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}
	if err := validator.Validate(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"errors":  err,
		})
	}
	result, err := h.service.CreatePermission(c.Context(), &req)
	if err != nil {
		return err
	}
	return utils.Created(c, result)
}

// @Summary     Update permission
// @Tags        settings
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id    path      string                      true  "Permission ID"
// @Param       body  body      settings.PermissionRequest  true  "Permission data"
// @Success     200   {object}  utils.Response{data=settings.PermissionResponse}
// @Failure     401   {object}  utils.Response
// @Failure     404   {object}  utils.Response
// @Failure     422   {object}  utils.Response
// @Router      /settings/permissions/{id} [patch]
func (h *Handler) UpdatePermission(c *fiber.Ctx) error {
	var req PermissionRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}
	if err := validator.Validate(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"errors":  err,
		})
	}
	result, err := h.service.UpdatePermission(c.Context(), c.Params("id"), &req)
	if err != nil {
		return err
	}
	return utils.OK(c, result)
}

// @Summary     Delete permission
// @Tags        settings
// @Produce     json
// @Security    BearerAuth
// @Param       id   path  string  true  "Permission ID"
// @Success     204
// @Failure     401  {object}  utils.Response
// @Failure     404  {object}  utils.Response
// @Router      /settings/permissions/{id} [delete]
func (h *Handler) DeletePermission(c *fiber.Ctx) error {
	if err := h.service.DeletePermission(c.Context(), c.Params("id")); err != nil {
		return err
	}
	return utils.NoContent(c)
}

// @Summary     Assign roles to user
// @Tags        settings
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id    path      string                    true  "User ID"
// @Param       body  body      settings.SyncRolesRequest  true  "Role IDs"
// @Success     200   {object}  utils.Response
// @Failure     401   {object}  utils.Response
// @Failure     404   {object}  utils.Response
// @Failure     422   {object}  utils.Response
// @Router      /settings/users/{id}/roles [post]
func (h *Handler) SyncUserRoles(c *fiber.Ctx) error {
	var req SyncRolesRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}
	if err := validator.Validate(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"errors":  err,
		})
	}
	if err := h.service.SyncUserRoles(c.Context(), c.Params("id"), &req); err != nil {
		return err
	}
	return utils.OK(c, fiber.Map{"message": "roles assigned"})
}
