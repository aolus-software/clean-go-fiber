package repositories

import (
	"context"
	"time"

	"github.com/aolus-software/clean-go-fiber/internal/libs/pagination"
	"github.com/aolus-software/clean-go-fiber/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserInformation struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Email       string   `json:"email"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}

func (u *UserInformation) HasRole(name string) bool {
	for _, r := range u.Roles {
		if r == name {
			return true
		}
	}
	return false
}

func (u *UserInformation) HasPermission(name string) bool {
	for _, p := range u.Permissions {
		if p == name {
			return true
		}
	}
	return false
}

type EmailVerificationRepository interface {
	WithTx(tx *gorm.DB) EmailVerificationRepository
	Create(ctx context.Context, ev *models.EmailVerification) error
	FindUsable(ctx context.Context, tokenHash string, now time.Time) (*models.EmailVerification, error)
	MarkAllUsed(ctx context.Context, userID uuid.UUID, now time.Time) error
}

type PasswordResetTokenRepository interface {
	WithTx(tx *gorm.DB) PasswordResetTokenRepository
	Create(ctx context.Context, t *models.PasswordResetToken) error
	FindUsable(ctx context.Context, tokenHash string, now time.Time) (*models.PasswordResetToken, error)
	MarkAllUsed(ctx context.Context, userID uuid.UUID, now time.Time) error
}

type RoleRepository interface {
	WithTx(tx *gorm.DB) RoleRepository
	FindAll(ctx context.Context, q *pagination.Query, tz *time.Location) ([]models.Role, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.Role, error)
	FindByName(ctx context.Context, name string, excludeID *uuid.UUID) (*models.Role, error)
	CountByIDs(ctx context.Context, ids []uuid.UUID) (int64, error)
	PermissionIDs(ctx context.Context, roleID uuid.UUID) ([]uuid.UUID, error)
	Create(ctx context.Context, role *models.Role) error
	Update(ctx context.Context, role *models.Role) error
	Delete(ctx context.Context, id uuid.UUID) error
	ReplacePermissions(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error
	SelectOptions(ctx context.Context) ([]models.Role, error)
}

type PermissionRepository interface {
	WithTx(tx *gorm.DB) PermissionRepository
	FindAll(ctx context.Context, q *pagination.Query, tz *time.Location) ([]models.Permission, int64, error)
	FindAllOrdered(ctx context.Context) ([]models.Permission, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.Permission, error)
	FindExistingNames(ctx context.Context, names []string, excludeID *uuid.UUID) ([]string, error)
	CountByIDs(ctx context.Context, ids []uuid.UUID) (int64, error)
	Create(ctx context.Context, permissions []models.Permission) error
	Update(ctx context.Context, permission *models.Permission) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type UserRepository interface {
	WithTx(tx *gorm.DB) UserRepository
	FindAll(ctx context.Context, q *pagination.Query, tz *time.Location) ([]models.User, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	FindLiveByEmail(ctx context.Context, email string, excludeID *uuid.UUID) (*models.User, error)
	FindUserInformation(ctx context.Context, id uuid.UUID) (*UserInformation, error)
	FindIDsByRoleIDs(ctx context.Context, roleIDs []uuid.UUID) ([]uuid.UUID, error)
	FindIDsByPermissionID(ctx context.Context, permissionID uuid.UUID) ([]uuid.UUID, error)
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	ReplaceRoles(ctx context.Context, userID uuid.UUID, roleIDs []uuid.UUID) error
}

func notFoundAsNil[T any](v *T, err error) (*T, error) {
	if err == nil {
		return v, nil
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return nil, err
}
