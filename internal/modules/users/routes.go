package users

import (
	"github.com/aolus-software/clean-go-fiber/internal/guards"
	"github.com/aolus-software/clean-go-fiber/internal/models"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) RegisterRoutes(settings fiber.Router, auth fiber.Handler) {
	r := settings.Group("/users", auth)
	r.Get("", guards.Permission("user list"), h.List)
	r.Post("", guards.Permission("user create"), h.Create)
	r.Get("/:id", guards.Permission("user detail"), h.Detail)
	r.Patch("/:id", guards.Permission("user edit"), h.Update)
	r.Post("/:id/reset-password", guards.Role(models.SuperuserRole), h.ResetPassword)
	r.Post("/:id/send-verification-email", guards.Permission("user create"), h.SendVerificationEmail)
	r.Post("/:id/send-reset-password-email", guards.Permission("user create"), h.SendResetPasswordEmail)
	r.Delete("/:id", guards.Permission("user delete"), h.Delete)
}
