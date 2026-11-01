package selectoptions

import (
	"github.com/aolus-software/clean-go-fiber/internal/libs/utils"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// @Summary     Permission options
// @Description Every permission grouped by group. Requires the 'superuser' role.
// @Tags        Settings/Select Options
// @Produce     json
// @Security    BearerAuth
// @Success     200  {object}  utils.Response{data=[]selectoptions.PermissionGroupOption}
// @Failure     401  {object}  utils.Response
// @Failure     403  {object}  utils.Response
// @Router      /settings/select-options/permissions [get]
func (h *Handler) Permissions(c *fiber.Ctx) error {
	result, err := h.service.Permissions(c.Context())
	if err != nil {
		return err
	}
	return utils.OK(c, "Permission options retrieved", result)
}

// @Summary     Role options
// @Description Every role except superuser, ordered by name. Requires the 'superuser' role.
// @Tags        Settings/Select Options
// @Produce     json
// @Security    BearerAuth
// @Success     200  {object}  utils.Response{data=[]selectoptions.RoleOption}
// @Failure     401  {object}  utils.Response
// @Failure     403  {object}  utils.Response
// @Router      /settings/select-options/roles [get]
func (h *Handler) Roles(c *fiber.Ctx) error {
	result, err := h.service.Roles(c.Context())
	if err != nil {
		return err
	}
	return utils.OK(c, "Role options retrieved", result)
}
