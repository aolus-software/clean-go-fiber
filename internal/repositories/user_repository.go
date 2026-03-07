package repositories

import (
	"context"

	"github.com/aolus-software/clean-go-fiber/internal/libs/pagination"
	"github.com/aolus-software/clean-go-fiber/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var userDatatableOptions = pagination.DatatableOptions{
	AllowableFilter: []string{"name", "email", "status", "created_at", "updated_at"},
	AllowableSort:   []string{"id", "name", "email", "status", "created_at", "updated_at"},
	DefaultSort:     "created_at",
	DefaultSortDir:  "desc",
	DefaultPerPage:  10,
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) WithTx(tx *gorm.DB) UserRepository {
	return &userRepository{db: tx}
}

func (r *userRepository) FindAll(ctx context.Context, queryParams *pagination.QueryParamDatatable) (*pagination.PaginationResponse[models.User], error) {
	queryParams.ApplyDefaults(userDatatableOptions)
	if err := queryParams.Validate(userDatatableOptions); err != nil {
		return nil, err
	}

	query := r.db.WithContext(ctx).Model(&models.User{})

	if queryParams.Search != "" {
		query = query.Where("name ILIKE ? OR email ILIKE ?", "%"+queryParams.Search+"%", "%"+queryParams.Search+"%")
	}

	if len(queryParams.Filters) > 0 {
		if queryParams.Filters["name"] != "" {
			query = query.Where("name ILIKE ?", "%"+queryParams.Filters["name"]+"%")
		}
		if queryParams.Filters["email"] != "" {
			query = query.Where("email ILIKE ?", "%"+queryParams.Filters["email"]+"%")
		}
		if queryParams.Filters["status"] != "" {
			query = query.Where("status = ?", queryParams.Filters["status"])
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

	var users []models.User
	if err := query.Preload("Roles").Offset((queryParams.Page - 1) * queryParams.PerPage).Limit(queryParams.PerPage).Find(&users).Error; err != nil {
		return nil, err
	}

	return &pagination.PaginationResponse[models.User]{
		Total:       total,
		PerPage:     queryParams.PerPage,
		CurrentPage: queryParams.Page,
		LastPage:    pagination.CalcLastPage(total, queryParams.PerPage),
		Data:        users,
	}, nil
}

func (r *userRepository) FindById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Preload("Roles.Permissions").First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) (*models.User, error) {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id).Error
}

func (r *userRepository) SyncRoles(ctx context.Context, user *models.User, roles []models.Role) error {
	return r.db.WithContext(ctx).Model(user).Association("Roles").Replace(roles)
}
