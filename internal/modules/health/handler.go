package health

import (
	"github.com/aolus-software/clean-go-fiber/internal/config"
	"github.com/aolus-software/clean-go-fiber/internal/libs/utils"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	cfg *config.Config
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{cfg: cfg}
}

// @Summary     Root
// @Description Returns service name and status
// @Tags        health
// @Produce     json
// @Success     200  {object} utils.Response{HomeResponse}
// @Router      / [get]
func (h *Handler) Root(c *fiber.Ctx) error {
	return utils.OK(c, &HomeResponse{
		Status:  "ok",
		Service: h.cfg.App.Name,
	})
}

// @Summary     Health check
// @Description Returns service health status
// @Tags        health
// @Produce     json
// @Success     200  {object} utils.Response{HealthResponse}
// @Router      /health [get]
func (h *Handler) Health(c *fiber.Ctx) error {
	return utils.OK(c, &HealthResponse{
		Status:  "ok",
		Service: h.cfg.App.Name,
	})
}
