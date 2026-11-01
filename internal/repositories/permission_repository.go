package repositories

import (
	"context"
	"time"

	"github.com/aolus-software/clean-go-fiber/internal/libs/pagination"
	"github.com/aolus-software/clean-go-fiber/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var PermissionListOptions = pagination.Options{
	Sortable: map[string]string{
		"id":        "id",
		"name":      "name",
		"group":     `"group"`,
		"createdAt": "created_at",
		"updatedAt": "updated_at",
	},
	Filterable: map[string]pagination.Filter{
		"name":      {Column: "name", Kind: pagination.FilterText},
		"group":     {Column: `"group"`, Kind: pagination.FilterText},
		"createdAt": {Column: "created_at", Kind: pagination.FilterDate},
		"updatedAt": {Column: "updated_at", Kind: pagination.FilterDate},
	},
	SearchColumns: []string{"name", `"group"`},
}

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) WithTx(tx *gorm.DB) PermissionRepository {
	return &permissionRepository{db: tx}
}

func (r *permissionRepository) FindAll(ctx context.Context, q *pagination.Query, tz *time.Location) ([]models.Permission, int64, error) {
	query, err := q.Apply(r.db.WithContext(ctx).Model(&models.Permission{}), PermissionListOptions, tz)
	if err != nil {
		return nil, 0, err
	}
	var total int64
	if err := query.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var permissions []models.Permission
	if err := query.Offset(q.Offset()).Limit(q.PerPage).Find(&permissions).Error; err != nil {
		return nil, 0, err
	}
	return permissions, total, nil
}

func (r *permissionRepository) FindAllOrdered(ctx context.Context) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.WithContext(ctx).Order(`"group"`).Order("name").Find(&permissions).Error
	return permissions, err
}

func (r *permissionRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Permission, error) {
	var permission models.Permission
	err := r.db.WithContext(ctx).First(&permission, "id = ?", id).Error
	return notFoundAsNil(&permission, err)
}

func (r *permissionRepository) FindExistingNames(ctx context.Context, names []string, excludeID *uuid.UUID) ([]string, error) {
	var existing []string
	if len(names) == 0 {
		return existing, nil
	}
	query := r.db.WithContext(ctx).Model(&models.Permission{}).Where("name IN ?", names)
	if excludeID != nil {
		query = query.Where("id <> ?", *excludeID)
	}
	err := query.Pluck("name", &existing).Error
	return existing, err
}

func (r *permissionRepository) CountByIDs(ctx context.Context, ids []uuid.UUID) (int64, error) {
	var count int64
	if len(ids) == 0 {
		return 0, nil
	}
	err := r.db.WithContext(ctx).Model(&models.Permission{}).Where("id IN ?", ids).Count(&count).Error
	return count, err
}

func (r *permissionRepository) Create(ctx context.Context, permissions []models.Permission) error {
	return r.db.WithContext(ctx).Create(&permissions).Error
}

func (r *permissionRepository) Update(ctx context.Context, permission *models.Permission) error {
	return r.db.WithContext(ctx).Save(permission).Error
}

func (r *permissionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Permission{}, "id = ?", id).Error
}
