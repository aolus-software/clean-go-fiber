package health

import (
	"context"
	"time"

	"github.com/aolus-software/clean-go-fiber/internal/config"
	"github.com/aolus-software/clean-go-fiber/internal/libs/cache"
	"github.com/aolus-software/clean-go-fiber/internal/libs/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Handler struct {
	cfg   *config.Config
	db    *gorm.DB
	cache *cache.Cache
	log   *zap.Logger
}

func NewHandler(cfg *config.Config, db *gorm.DB, c *cache.Cache, log *zap.Logger) *Handler {
	return &Handler{cfg: cfg, db: db, cache: c, log: log}
}

// @Summary     Root
// @Description Application name, environment, and server time.
// @Tags        General
// @Produce     json
// @Success     200  {object}  utils.Response{data=health.HomeResponse}
// @Router      / [get]
func (h *Handler) Root(c *fiber.Ctx) error {
	return utils.OK(c, "Success", HomeResponse{
		AppName: h.cfg.App.Name,
		AppEnv:  h.cfg.App.Env,
		Date:    time.Now().In(h.cfg.App.Timezone).Format(time.RFC3339),
	})
}

// @Summary     Health check
// @Description Probes PostgreSQL and Redis; 503 when either fails.
// @Tags        General
// @Produce     json
// @Success     200  {object}  utils.Response{data=health.HealthResponse}
// @Failure     503  {object}  utils.Response{data=health.HealthResponse}
// @Router      /health [get]
func (h *Handler) Health(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 3*time.Second)
	defer cancel()

	services := HealthServices{Database: "healthy", Redis: "healthy"}
	healthy := true

	sqlDB, err := h.db.DB()
	if err == nil {
		err = sqlDB.PingContext(ctx)
	}
	if err != nil {
		h.log.Error("health check: database probe failed", zap.Error(err))
		services.Database = "unhealthy"
		healthy = false
	}
	if err := h.cache.Ping(ctx); err != nil {
		h.log.Error("health check: redis probe failed", zap.Error(err))
		services.Redis = "unhealthy"
		healthy = false
	}

	payload := HealthResponse{Status: "healthy", Timestamp: time.Now().UTC().Format(time.RFC3339), Services: services}
	if !healthy {
		payload.Status = "unhealthy"
		return c.Status(fiber.StatusServiceUnavailable).JSON(utils.Response{
			Status: fiber.StatusServiceUnavailable, Success: false, Message: "Service unavailable", Data: payload,
		})
	}
	return utils.OK(c, "Service is healthy", payload)
}

// @Summary     Liveness probe
// @Tags        General
// @Produce     json
// @Success     200  {object}  utils.Response{data=health.LiveResponse}
// @Router      /live [get]
func (h *Handler) Live(c *fiber.Ctx) error {
	return utils.OK(c, "Alive", LiveResponse{Alive: true})
}
