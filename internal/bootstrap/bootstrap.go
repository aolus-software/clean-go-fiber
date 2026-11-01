package bootstrap

import (
	"github.com/aolus-software/clean-go-fiber/internal/config"
	"github.com/aolus-software/clean-go-fiber/internal/libs/cache"
	"github.com/aolus-software/clean-go-fiber/internal/libs/txmanager"
	"github.com/aolus-software/clean-go-fiber/internal/middleware"
	"github.com/aolus-software/clean-go-fiber/internal/modules/auth"
	"github.com/aolus-software/clean-go-fiber/internal/modules/docs"
	"github.com/aolus-software/clean-go-fiber/internal/modules/health"
	"github.com/aolus-software/clean-go-fiber/internal/modules/permissions"
	"github.com/aolus-software/clean-go-fiber/internal/modules/profile"
	"github.com/aolus-software/clean-go-fiber/internal/modules/roles"
	"github.com/aolus-software/clean-go-fiber/internal/modules/selectoptions"
	"github.com/aolus-software/clean-go-fiber/internal/modules/users"
	"github.com/aolus-software/clean-go-fiber/internal/repositories"
	"github.com/gofiber/fiber/v2"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func RegisterRoutes(app *fiber.App, db *gorm.DB, redisClient *redis.Client, queueClient *asynq.Client, cfg *config.Config, log *zap.Logger) {
	c := cache.New(redisClient)
	tx := txmanager.New(db)

	userRepo := repositories.NewUserRepository(db)
	roleRepo := repositories.NewRoleRepository(db)
	permissionRepo := repositories.NewPermissionRepository(db)
	emailVerifyRepo := repositories.NewEmailVerificationRepository(db)
	resetTokenRepo := repositories.NewPasswordResetTokenRepository(db)

	authMiddleware := middleware.Auth(userRepo, c)
	accountMail := auth.NewAccountMail(emailVerifyRepo, resetTokenRepo, queueClient, tx, cfg)

	health.NewHandler(cfg, db, c, log).RegisterRoutes(app)
	docs.NewHandler(cfg).RegisterRoutes(app)

	auth.NewHandler(auth.NewService(userRepo, emailVerifyRepo, resetTokenRepo, accountMail, tx, c, cfg, log)).RegisterRoutes(app)
	profile.NewHandler(profile.NewService(userRepo, c)).RegisterRoutes(app, authMiddleware)

	settings := app.Group("/settings")
	users.NewHandler(users.NewService(userRepo, roleRepo, accountMail, tx, c, cfg)).RegisterRoutes(settings, authMiddleware)
	roles.NewHandler(roles.NewService(roleRepo, permissionRepo, userRepo, tx, c, cfg)).RegisterRoutes(settings, authMiddleware)
	permissions.NewHandler(permissions.NewService(permissionRepo, userRepo, c, cfg)).RegisterRoutes(settings, authMiddleware)
	selectoptions.NewHandler(selectoptions.NewService(permissionRepo, roleRepo)).RegisterRoutes(settings, authMiddleware)
}
