package bootstrap

import (
	"github.com/aolus-software/clean-go-fiber/internal/config"
	"github.com/aolus-software/clean-go-fiber/internal/modules/auth"
	"github.com/aolus-software/clean-go-fiber/internal/modules/docs"
	"github.com/aolus-software/clean-go-fiber/internal/modules/health"
	"github.com/aolus-software/clean-go-fiber/internal/repositories"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRoutes(app *fiber.App, db *gorm.DB, cfg *config.Config) {
	health.NewHandler(cfg).RegisterRoutes(app)
	docs.NewHandler(cfg).RegisterRoutes(app)

	userRepo := repositories.NewUserRepository(db)
	emailVerifyRepo := repositories.NewEmailVerificationRepository(db)
	forgotPasswordRepo := repositories.NewForgotPasswordRepository(db)

	authService := auth.NewService(userRepo, emailVerifyRepo, forgotPasswordRepo, cfg)
	authHandler := auth.NewHandler(authService)

	api := app.Group("/")
	authHandler.RegisterRoutes(api)
}
