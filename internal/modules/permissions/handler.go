package permissions

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

// @Summary     List permissions
// @Tags        permissions
// @Produce     json
// @Security    BearerAuth
// @Param       page            query  int     false  "Page"
// @Param       per_page        query  int     false  "Per page"
// @Param       search          query  string  false  "Search"
// @Param       order_by        query  string  false  "Order by"
// @Param       sort_direction  query  string  false  "Sort direction"
// @Success     200  {object}  utils.Response
// @Failure     401  {object}  utils.Response
// @Router      /permissions [get]
func (h *Handler) List(c *fiber.Ctx) error {
	q := pagination.ParseFromFiber(c)
	result, err := h.service.List(c.Context(), q)
	if err != nil {
		return err
	}
	return utils.OK(c, result)
}

// @Summary     Get permission by ID
// @Tags        permissions
// @Produce     json
// @Security    BearerAuth
// @Param       id   path      string  true  "Permission ID"
// @Success     200  {object}  utils.Response{data=permissions.PermissionResponse}
// @Failure     401  {object}  utils.Response
// @Failure     404  {object}  utils.Response
// @Router      /permissions/{id} [get]
func (h *Handler) Get(c *fiber.Ctx) error {
	result, err := h.service.GetByID(c.Context(), c.Params("id"))
	if err != nil {
		return err
	}
	return utils.OK(c, result)
}

// @Summary     Create permission
// @Tags        permissions
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       body  body      permissions.PermissionRequest  true  "Permission data"
// @Success     201   {object}  utils.Response{data=permissions.PermissionResponse}
// @Failure     401   {object}  utils.Response
// @Failure     422   {object}  utils.Response
// @Router      /permissions [post]
func (h *Handler) Create(c *fiber.Ctx) error {
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
	result, err := h.service.Create(c.Context(), &req)
	if err != nil {
		return err
	}
	return utils.Created(c, result)
}

// @Summary     Update permission
// @Tags        permissions
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id    path      string                         true  "Permission ID"
// @Param       body  body      permissions.PermissionRequest  true  "Permission data"
// @Success     200   {object}  utils.Response{data=permissions.PermissionResponse}
// @Failure     401   {object}  utils.Response
// @Failure     404   {object}  utils.Response
// @Failure     422   {object}  utils.Response
// @Router      /permissions/{id} [patch]
func (h *Handler) Update(c *fiber.Ctx) error {
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
	result, err := h.service.Update(c.Context(), c.Params("id"), &req)
	if err != nil {
		return err
	}
	return utils.OK(c, result)
}

// @Summary     Delete permission
// @Tags        permissions
// @Produce     json
// @Security    BearerAuth
// @Param       id   path  string  true  "Permission ID"
// @Success     204
// @Failure     401  {object}  utils.Response
// @Failure     404  {object}  utils.Response
// @Router      /permissions/{id} [delete]
func (h *Handler) Delete(c *fiber.Ctx) error {
	if err := h.service.Delete(c.Context(), c.Params("id")); err != nil {
		return err
	}
	return utils.NoContent(c)
}
