package roles

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

// @Summary     List roles
// @Description Requires 'role list' permission. Filters: filter[name], filter[createdAt], filter[updatedAt] (YYYY-MM-DD or start,end).
// @Tags        Settings/Roles
// @Produce     json
// @Security    BearerAuth
// @Param       page           query  int     false  "Page (>= 1)"
// @Param       perPage        query  int     false  "Page size (1-100)"
// @Param       search         query  string  false  "Search name"
// @Param       sort           query  string  false  "id, name, createdAt, updatedAt"
// @Param       sortDirection  query  string  false  "asc or desc"
// @Success     200  {object}  utils.Response{data=utils.PaginatedData[roles.RoleList]}
// @Failure     400  {object}  utils.ErrorResponse
// @Failure     401  {object}  utils.Response
// @Failure     403  {object}  utils.Response
// @Failure     422  {object}  utils.ErrorResponse
// @Router      /settings/roles [get]
func (h *Handler) List(c *fiber.Ctx) error {
	q, err := pagination.Parse(c)
	if err != nil {
		return err
	}
	result, err := h.service.List(c.Context(), q)
	if err != nil {
		return err
	}
	return utils.OK(c, "Roles retrieved", result)
}

// @Summary     Create role
// @Description Requires 'role create' permission.
// @Tags        Settings/Roles
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       body  body      roles.RoleRequest  true  "Role"
// @Success     201   {object}  utils.Response
// @Failure     400   {object}  utils.ErrorResponse
// @Failure     401   {object}  utils.Response
// @Failure     403   {object}  utils.Response
// @Failure     422   {object}  utils.ErrorResponse
// @Router      /settings/roles [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	var req RoleRequest
	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}
	if err := h.service.Create(c.Context(), &req); err != nil {
		return err
	}
	return utils.Created(c, "Role created", nil)
}

// @Summary     Role detail
// @Description Every permission, grouped, with isAssigned for this role. Requires 'role detail' permission.
// @Tags        Settings/Roles
// @Produce     json
// @Security    BearerAuth
// @Param       id   path      string  true  "Role ID (uuid)"
// @Success     200  {object}  utils.Response{data=roles.RoleDetail}
// @Failure     401  {object}  utils.Response
// @Failure     403  {object}  utils.Response
// @Failure     404  {object}  utils.Response
// @Failure     422  {object}  utils.ErrorResponse
// @Router      /settings/roles/{id} [get]
func (h *Handler) Detail(c *fiber.Ctx) error {
	id, err := utils.ParamUUID(c, "id")
	if err != nil {
		return err
	}
	result, err := h.service.Detail(c.Context(), id)
	if err != nil {
		return err
	}
	return utils.OK(c, "Role retrieved", result)
}

// @Summary     Update role
// @Description Replace a role's name and permissions. Requires 'role edit' permission. The superuser role cannot be renamed.
// @Tags        Settings/Roles
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id    path      string             true  "Role ID (uuid)"
// @Param       body  body      roles.RoleRequest  true  "Role"
// @Success     200   {object}  utils.Response
// @Failure     400   {object}  utils.ErrorResponse
// @Failure     401   {object}  utils.Response
// @Failure     403   {object}  utils.Response
// @Failure     404   {object}  utils.Response
// @Failure     422   {object}  utils.ErrorResponse
// @Router      /settings/roles/{id} [patch]
func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := utils.ParamUUID(c, "id")
	if err != nil {
		return err
	}
	var req RoleRequest
	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}
	if err := h.service.Update(c.Context(), id, &req); err != nil {
		return err
	}
	return utils.OK(c, "Role updated", nil)
}

// @Summary     Delete role
// @Description Hard delete; assignments cascade. Requires 'role delete' permission. The superuser role cannot be deleted.
// @Tags        Settings/Roles
// @Produce     json
// @Security    BearerAuth
// @Param       id   path      string  true  "Role ID (uuid)"
// @Success     200  {object}  utils.Response
// @Failure     401  {object}  utils.Response
// @Failure     403  {object}  utils.Response
// @Failure     404  {object}  utils.Response
// @Failure     422  {object}  utils.ErrorResponse
// @Router      /settings/roles/{id} [delete]
func (h *Handler) Delete(c *fiber.Ctx) error {
	id, err := utils.ParamUUID(c, "id")
	if err != nil {
		return err
	}
	if err := h.service.Delete(c.Context(), id); err != nil {
		return err
	}
	return utils.OK(c, "Role deleted", nil)
}
