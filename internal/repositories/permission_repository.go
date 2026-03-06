package repositories

import (
	"context"

	"github.com/aolus-software/clean-go-fiber/internal/libs/pagination"
	"github.com/aolus-software/clean-go-fiber/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var permissionDatatableOptions = pagination.DatatableOptions{
	AllowableFilter: []string{"name", "modules", "created_at", "updated_at"},
	AllowableSort:   []string{"id", "name", "modules", "created_at", "updated_at"},
	DefaultSort:     "created_at",
	DefaultSortDir:  "desc",
	DefaultPerPage:  10,
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

func (r *permissionRepository) FindAll(ctx context.Context, queryParams *pagination.QueryParamDatatable) (*pagination.PaginationResponse[models.Permission], error) {
	queryParams.ApplyDefaults(permissionDatatableOptions)
	if err := queryParams.Validate(permissionDatatableOptions); err != nil {
		return nil, err
	}

	query := r.db.WithContext(ctx).Model(&models.Permission{})

	if queryParams.Search != "" {
		query = query.Where("name ILIKE ? OR modules ILIKE ?", "%"+queryParams.Search+"%", "%"+queryParams.Search+"%")
	}

	if len(queryParams.Filters) > 0 {
		if queryParams.Filters["name"] != "" {
			query = query.Where("name ILIKE ?", "%"+queryParams.Filters["name"]+"%")
		}
		if queryParams.Filters["modules"] != "" {
			query = query.Where("modules ILIKE ?", "%"+queryParams.Filters["modules"]+"%")
		}
		if queryParams.Filters["created_at"] != "" {
			query = query.Where("created_at::date = ?", queryParams.Filters["created_at"])
		}
		if queryParams.Filters["updated_at"] != "" {
			query = query.Where("updated_at::date = ?", queryParams.Filters["updated_at"])
		}
	}

	if queryParams.OrderBy != "" {
		query = query.Order(queryParams.OrderBy + " " + queryParams.SortDirection)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	var permissions []models.Permission
	if err := query.Offset((queryParams.Page - 1) * queryParams.PerPage).Limit(queryParams.PerPage).Find(&permissions).Error; err != nil {
		return nil, err
	}

	return &pagination.PaginationResponse[models.Permission]{
		Total:       total,
		PerPage:     queryParams.PerPage,
		CurrentPage: queryParams.Page,
		Data:        permissions,
	}, nil
}

func (r *permissionRepository) FindById(ctx context.Context, id uuid.UUID) (*models.Permission, error) {
	var permission models.Permission
	if err := r.db.WithContext(ctx).First(&permission, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepository) FindByIds(ctx context.Context, ids []uuid.UUID) ([]models.Permission, error) {
	var permissions []models.Permission
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r *permissionRepository) Create(ctx context.Context, permission *models.Permission) (*models.Permission, error) {
	if err := r.db.WithContext(ctx).Create(permission).Error; err != nil {
		return nil, err
	}
	return permission, nil
}

func (r *permissionRepository) Update(ctx context.Context, permission *models.Permission) (*models.Permission, error) {
	if err := r.db.WithContext(ctx).Save(permission).Error; err != nil {
		return nil, err
	}
	return permission, nil
}

func (r *permissionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Permission{}, "id = ?", id).Error
}
