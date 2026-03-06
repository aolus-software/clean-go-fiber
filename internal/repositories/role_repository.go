package repositories

import (
	"context"

	"github.com/aolus-software/clean-go-fiber/internal/libs/pagination"
	"github.com/aolus-software/clean-go-fiber/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var roleDatatableOptions = pagination.DatatableOptions{
	AllowableFilter: []string{"name", "created_at", "updated_at"},
	AllowableSort:   []string{"id", "name", "created_at", "updated_at"},
	DefaultSort:     "created_at",
	DefaultSortDir:  "desc",
	DefaultPerPage:  10,
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) WithTx(tx *gorm.DB) RoleRepository {
	return &roleRepository{db: tx}
}

func (r *roleRepository) FindAll(ctx context.Context, queryParams *pagination.QueryParamDatatable) (*pagination.PaginationResponse[models.Role], error) {
	queryParams.ApplyDefaults(roleDatatableOptions)
	if err := queryParams.Validate(roleDatatableOptions); err != nil {
		return nil, err
	}

	query := r.db.WithContext(ctx).Model(&models.Role{})

	if queryParams.Search != "" {
		query = query.Where("name ILIKE ?", "%"+queryParams.Search+"%")
	}

	if len(queryParams.Filters) > 0 {
		if queryParams.Filters["name"] != "" {
			query = query.Where("name ILIKE ?", "%"+queryParams.Filters["name"]+"%")
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

	var roles []models.Role
	if err := query.Offset((queryParams.Page - 1) * queryParams.PerPage).Limit(queryParams.PerPage).Find(&roles).Error; err != nil {
		return nil, err
	}

	return &pagination.PaginationResponse[models.Role]{
		Total:       total,
		PerPage:     queryParams.PerPage,
		CurrentPage: queryParams.Page,
		Data:        roles,
	}, nil
}

func (r *roleRepository) FindById(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	var role models.Role
	if err := r.db.WithContext(ctx).Preload("Permissions").First(&role, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) SyncPermissions(ctx context.Context, role *models.Role, permissions []models.Permission) error {
	return r.db.WithContext(ctx).Model(role).Association("Permissions").Replace(permissions)
}

func (r *roleRepository) Create(ctx context.Context, role *models.Role) (*models.Role, error) {
	if err := r.db.WithContext(ctx).Create(role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

func (r *roleRepository) Update(ctx context.Context, role *models.Role) (*models.Role, error) {
	if err := r.db.WithContext(ctx).Save(role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

func (r *roleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Role{}, "id = ?", id).Error
}
