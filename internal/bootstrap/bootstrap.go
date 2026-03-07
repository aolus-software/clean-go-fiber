package bootstrap

import (
	"github.com/aolus-software/clean-go-fiber/internal/config"
	"github.com/aolus-software/clean-go-fiber/internal/modules/auth"
	"github.com/aolus-software/clean-go-fiber/internal/modules/docs"
	"github.com/aolus-software/clean-go-fiber/internal/modules/health"
	"github.com/aolus-software/clean-go-fiber/internal/modules/profile"
	"github.com/aolus-software/clean-go-fiber/internal/modules/settings"
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
	roleRepo := repositories.NewRoleRepository(db)
	permissionRepo := repositories.NewPermissionRepository(db)

	authService := auth.NewService(userRepo, emailVerifyRepo, forgotPasswordRepo, cfg)
	auth.NewHandler(authService).RegisterRoutes(app)

	profileService := profile.NewService(userRepo)
	profile.NewHandler(profileService).RegisterRoutes(app)

	settingsService := settings.NewService(roleRepo, permissionRepo, userRepo)
	settings.NewHandler(settingsService).RegisterRoutes(app, userRepo)
}
