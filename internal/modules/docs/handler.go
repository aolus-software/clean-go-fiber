package docs

import (
	"html"

	"github.com/aolus-software/clean-go-fiber/internal/config"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	cfg *config.Config
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{cfg: cfg}
}

// @Summary     OpenAPI specification
// @Description Returns the raw OpenAPI JSON specification
// @Tags        General
// @Produce     json
// @Success     200  {object}  map[string]interface{}
// @Router      /docs/openapi.json [get]
func (h *Handler) Spec(c *fiber.Ctx) error {
	return c.SendFile("./docs/swagger.json")
}

// @Summary     API documentation UI
// @Description Serves the Scalar API reference UI
// @Tags        General
// @Produce     html
// @Success     200  {string}  string  "HTML page"
// @Router      /docs [get]
func (h *Handler) UI(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/html")
	return c.SendString(`<!DOCTYPE html>
<html>
<head><title>` + html.EscapeString(h.cfg.App.Name) + ` API Docs</title></head>
<body>
<script id="api-reference" data-url="/docs/openapi.json"></script>
<script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference@1.25.0"></script>
</body>
</html>`)
}
