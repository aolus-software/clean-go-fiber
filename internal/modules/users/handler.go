package users

import (
	"github.com/aolus-software/clean-go-fiber/internal/libs/pagination"
	"github.com/aolus-software/clean-go-fiber/internal/libs/utils"
	"github.com/aolus-software/clean-go-fiber/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// @Summary     List users
// @Description Paginated list of live users. Requires 'user list' permission. Filters: filter[status] (comma list of active,inactive,suspended,blocked), filter[name], filter[email], filter[roleId] (comma list of uuids), filter[createdAt], filter[updatedAt] (YYYY-MM-DD or start,end).
// @Tags        Settings/Users
// @Produce     json
// @Security    BearerAuth
// @Param       page           query  int     false  "Page (>= 1)"
// @Param       perPage        query  int     false  "Page size (1-100)"
// @Param       search         query  string  false  "Search name, email, status"
// @Param       sort           query  string  false  "id, name, email, status, createdAt, updatedAt"
// @Param       sortDirection  query  string  false  "asc or desc"
// @Success     200  {object}  utils.Response{data=utils.PaginatedData[users.UserList]}
// @Failure     400  {object}  utils.ErrorResponse
// @Failure     401  {object}  utils.Response
// @Failure     403  {object}  utils.Response
// @Failure     422  {object}  utils.ErrorResponse
// @Router      /settings/users [get]
func (h *Handler) List(c *fiber.Ctx) error {
	q, err := pagination.Parse(c)
	if err != nil {
		return err
	}
	result, err := h.service.List(c.Context(), q)
	if err != nil {
		return err
	}
	return utils.OK(c, "Users retrieved", result)
}

// @Summary     Create user
// @Description Requires 'user create' permission. Only a superuser may assign the superuser role.
// @Tags        Settings/Users
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       body  body      users.CreateUserRequest  true  "User"
// @Success     201   {object}  utils.Response
// @Failure     400   {object}  utils.ErrorResponse
// @Failure     401   {object}  utils.Response
// @Failure     403   {object}  utils.Response
// @Failure     422   {object}  utils.ErrorResponse
// @Router      /settings/users [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	actor, err := middleware.CurrentUser(c)
	if err != nil {
		return err
	}
	var req CreateUserRequest
	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}
	if err := h.service.Create(c.Context(), actor, &req); err != nil {
		return err
	}
	return utils.Created(c, "User created", nil)
}

// @Summary     User detail
// @Description Requires 'user detail' permission.
// @Tags        Settings/Users
// @Produce     json
// @Security    BearerAuth
// @Param       id   path      string  true  "User ID (uuid)"
// @Success     200  {object}  utils.Response{data=users.UserDetail}
// @Failure     401  {object}  utils.Response
// @Failure     403  {object}  utils.Response
// @Failure     404  {object}  utils.Response
// @Failure     422  {object}  utils.ErrorResponse
// @Router      /settings/users/{id} [get]
func (h *Handler) Detail(c *fiber.Ctx) error {
	id, err := utils.ParamUUID(c, "id")
	if err != nil {
		return err
	}
	result, err := h.service.Detail(c.Context(), id)
	if err != nil {
		return err
	}
	return utils.OK(c, "User retrieved", result)
}

// @Summary     Update user
// @Description Full replace of name, email, status, remark, roleIds. Requires 'user edit' permission. Only a superuser may modify a superuser account or assign the superuser role.
// @Tags        Settings/Users
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id    path      string                   true  "User ID (uuid)"
// @Param       body  body      users.UpdateUserRequest  true  "User"
// @Success     200   {object}  utils.Response
// @Failure     400   {object}  utils.ErrorResponse
// @Failure     401   {object}  utils.Response
// @Failure     403   {object}  utils.Response
// @Failure     404   {object}  utils.Response
// @Failure     422   {object}  utils.ErrorResponse
// @Router      /settings/users/{id} [patch]
func (h *Handler) Update(c *fiber.Ctx) error {
	actor, err := middleware.CurrentUser(c)
	if err != nil {
		return err
	}
	id, err := utils.ParamUUID(c, "id")
	if err != nil {
		return err
	}
	var req UpdateUserRequest
	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}
	if err := h.service.Update(c.Context(), actor, id, &req); err != nil {
		return err
	}
	return utils.OK(c, "User updated", nil)
}

// @Summary     Reset a user's password
// @Description Requires the 'superuser' role.
// @Tags        Settings/Users
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id    path      string                      true  "User ID (uuid)"
// @Param       body  body      users.ResetPasswordRequest  true  "New password"
// @Success     200   {object}  utils.Response
// @Failure     400   {object}  utils.ErrorResponse
// @Failure     401   {object}  utils.Response
// @Failure     403   {object}  utils.Response
// @Failure     404   {object}  utils.Response
// @Failure     422   {object}  utils.ErrorResponse
// @Router      /settings/users/{id}/reset-password [post]
func (h *Handler) ResetPassword(c *fiber.Ctx) error {
	id, err := utils.ParamUUID(c, "id")
	if err != nil {
		return err
	}
	var req ResetPasswordRequest
	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}
	if err := h.service.ResetPassword(c.Context(), id, &req); err != nil {
		return err
	}
	return utils.OK(c, "Password reset", nil)
}

// @Summary     Send a verification email to a user
// @Description Requires 'user create' permission.
// @Tags        Settings/Users
// @Produce     json
// @Security    BearerAuth
// @Param       id   path      string  true  "User ID (uuid)"
// @Success     200  {object}  utils.Response
// @Failure     400  {object}  utils.ErrorResponse
// @Failure     401  {object}  utils.Response
// @Failure     403  {object}  utils.Response
// @Failure     404  {object}  utils.Response
// @Failure     422  {object}  utils.ErrorResponse
// @Router      /settings/users/{id}/send-verification-email [post]
func (h *Handler) SendVerificationEmail(c *fiber.Ctx) error {
	id, err := utils.ParamUUID(c, "id")
	if err != nil {
		return err
	}
	if err := h.service.SendVerificationEmail(c.Context(), id); err != nil {
		return err
	}
	return utils.OK(c, "Verification email queued", nil)
}

// @Summary     Send a password reset email to a user
// @Description Requires 'user create' permission.
// @Tags        Settings/Users
// @Produce     json
// @Security    BearerAuth
// @Param       id   path      string  true  "User ID (uuid)"
// @Success     200  {object}  utils.Response
// @Failure     401  {object}  utils.Response
// @Failure     403  {object}  utils.Response
// @Failure     404  {object}  utils.Response
// @Failure     422  {object}  utils.ErrorResponse
// @Router      /settings/users/{id}/send-reset-password-email [post]
func (h *Handler) SendResetPasswordEmail(c *fiber.Ctx) error {
	id, err := utils.ParamUUID(c, "id")
	if err != nil {
		return err
	}
	if err := h.service.SendResetPasswordEmail(c.Context(), id); err != nil {
		return err
	}
	return utils.OK(c, "Password reset email queued", nil)
}

// @Summary     Delete user
// @Description Soft delete. Requires 'user delete' permission. Only a superuser may delete a superuser account.
// @Tags        Settings/Users
// @Produce     json
// @Security    BearerAuth
// @Param       id   path      string  true  "User ID (uuid)"
// @Success     200  {object}  utils.Response
// @Failure     401  {object}  utils.Response
// @Failure     403  {object}  utils.Response
// @Failure     404  {object}  utils.Response
// @Failure     422  {object}  utils.ErrorResponse
// @Router      /settings/users/{id} [delete]
func (h *Handler) Delete(c *fiber.Ctx) error {
	actor, err := middleware.CurrentUser(c)
	if err != nil {
		return err
	}
	id, err := utils.ParamUUID(c, "id")
	if err != nil {
		return err
	}
	if err := h.service.Delete(c.Context(), actor, id); err != nil {
		return err
	}
	return utils.OK(c, "User deleted", nil)
}
