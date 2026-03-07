package repositories

import (
	"context"

	"github.com/aolus-software/clean-go-fiber/internal/libs/pagination"
	"github.com/aolus-software/clean-go-fiber/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EmailVerificationRepository interface {
	Create(ctx context.Context, ev *models.EmailVerification) (*models.EmailVerification, error)
	FindByToken(ctx context.Context, token string) (*models.EmailVerification, error)
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}

type ForgotPasswordRepository interface {
	Create(ctx context.Context, fp *models.ForgotPassword) (*models.ForgotPassword, error)
	FindByToken(ctx context.Context, token string) (*models.ForgotPassword, error)
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}

type RoleRepository interface {
	WithTx(tx *gorm.DB) RoleRepository
	FindAll(ctx context.Context, queryParams *pagination.QueryParamDatatable) (*pagination.PaginationResponse[models.Role], error)
	FindById(ctx context.Context, id uuid.UUID) (*models.Role, error)
	Create(ctx context.Context, role *models.Role) (*models.Role, error)
	Update(ctx context.Context, role *models.Role) (*models.Role, error)
	Delete(ctx context.Context, id uuid.UUID) error
	SyncPermissions(ctx context.Context, role *models.Role, permissions []models.Permission) error
}

type PermissionRepository interface {
	WithTx(tx *gorm.DB) PermissionRepository
	FindAll(ctx context.Context, queryParams *pagination.QueryParamDatatable) (*pagination.PaginationResponse[models.Permission], error)
	FindById(ctx context.Context, id uuid.UUID) (*models.Permission, error)
	FindByIds(ctx context.Context, ids []uuid.UUID) ([]models.Permission, error)
	Create(ctx context.Context, permission *models.Permission) (*models.Permission, error)
	Update(ctx context.Context, permission *models.Permission) (*models.Permission, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type UserRepository interface {
	WithTx(tx *gorm.DB) UserRepository
	FindAll(ctx context.Context, queryParams *pagination.QueryParamDatatable) (*pagination.PaginationResponse[models.User], error)
	FindById(ctx context.Context, id uuid.UUID) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, user *models.User) (*models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	SyncRoles(ctx context.Context, user *models.User, roles []models.Role) error
}
