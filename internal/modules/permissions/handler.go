package permissions

import (
	"github.com/aolus-software/clean-go-fiber/internal/libs/pagination"
	"github.com/aolus-software/clean-go-fiber/internal/libs/utils"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// @Summary     List permissions
// @Description Requires 'permission list' permission. Filters: filter[name], filter[group], filter[createdAt], filter[updatedAt] (YYYY-MM-DD or start,end).
// @Tags        Settings/Permissions
// @Produce     json
// @Security    BearerAuth
// @Param       page           query  int     false  "Page (>= 1)"
// @Param       perPage        query  int     false  "Page size (1-100)"
// @Param       search         query  string  false  "Search name, group"
// @Param       sort           query  string  false  "id, name, group, createdAt, updatedAt"
// @Param       sortDirection  query  string  false  "asc or desc"
// @Success     200  {object}  utils.Response{data=utils.PaginatedData[permissions.PermissionList]}
// @Failure     400  {object}  utils.ErrorResponse
// @Failure     401  {object}  utils.Response
// @Failure     403  {object}  utils.Response
// @Failure     422  {object}  utils.ErrorResponse
// @Router      /settings/permissions [get]
func (h *Handler) List(c *fiber.Ctx) error {
	q, err := pagination.Parse(c)
	if err != nil {
		return err
	}
	result, err := h.service.List(c.Context(), q)
	if err != nil {
		return err
	}
	return utils.OK(c, "Permissions retrieved", result)
}

// @Summary     Create permissions
// @Description Creates one permission per entry in name, stored as "<group> <name>". Requires 'permission create' permission.
// @Tags        Settings/Permissions
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       body  body      permissions.CreatePermissionRequest  true  "Group and actions"
// @Success     201   {object}  utils.Response
// @Failure     400   {object}  utils.ErrorResponse
// @Failure     401   {object}  utils.Response
// @Failure     403   {object}  utils.Response
// @Failure     422   {object}  utils.ErrorResponse
// @Router      /settings/permissions [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	var req CreatePermissionRequest
	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}
	if err := h.service.Create(c.Context(), &req); err != nil {
		return err
	}
	return utils.Created(c, "Permissions created", nil)
}

// @Summary     Permission detail
// @Description Requires 'permission detail' permission.
// @Tags        Settings/Permissions
// @Produce     json
// @Security    BearerAuth
// @Param       id   path      string  true  "Permission ID (uuid)"
// @Success     200  {object}  utils.Response{data=permissions.PermissionList}
// @Failure     401  {object}  utils.Response
// @Failure     403  {object}  utils.Response
// @Failure     404  {object}  utils.Response
// @Failure     422  {object}  utils.ErrorResponse
// @Router      /settings/permissions/{id} [get]
func (h *Handler) Detail(c *fiber.Ctx) error {
	id, err := utils.ParamUUID(c, "id")
	if err != nil {
		return err
	}
	result, err := h.service.Detail(c.Context(), id)
	if err != nil {
		return err
	}
	return utils.OK(c, "Permission retrieved", result)
}

// @Summary     Update permission
// @Description name is stored verbatim, so send the full "<group> <action>" string. Requires 'permission edit' permission.
// @Tags        Settings/Permissions
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id    path      string                               true  "Permission ID (uuid)"
// @Param       body  body      permissions.UpdatePermissionRequest  true  "Permission"
// @Success     200   {object}  utils.Response
// @Failure     400   {object}  utils.ErrorResponse
// @Failure     401   {object}  utils.Response
// @Failure     403   {object}  utils.Response
// @Failure     404   {object}  utils.Response
// @Failure     422   {object}  utils.ErrorResponse
// @Router      /settings/permissions/{id} [patch]
func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := utils.ParamUUID(c, "id")
	if err != nil {
		return err
	}
	var req UpdatePermissionRequest
	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}
	if err := h.service.Update(c.Context(), id, &req); err != nil {
		return err
	}
	return utils.OK(c, "Permission updated", nil)
}

// @Summary     Delete permission
// @Description Hard delete; role grants cascade. Requires 'permission delete' permission.
// @Tags        Settings/Permissions
// @Produce     json
// @Security    BearerAuth
// @Param       id   path      string  true  "Permission ID (uuid)"
// @Success     200  {object}  utils.Response
// @Failure     401  {object}  utils.Response
// @Failure     403  {object}  utils.Response
// @Failure     404  {object}  utils.Response
// @Failure     422  {object}  utils.ErrorResponse
// @Router      /settings/permissions/{id} [delete]
func (h *Handler) Delete(c *fiber.Ctx) error {
	id, err := utils.ParamUUID(c, "id")
	if err != nil {
		return err
	}
	if err := h.service.Delete(c.Context(), id); err != nil {
		return err
	}
	return utils.OK(c, "Permission deleted", nil)
}
