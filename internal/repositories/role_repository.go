package repositories

import (
	"context"
	"time"

	"github.com/aolus-software/clean-go-fiber/internal/libs/pagination"
	"github.com/aolus-software/clean-go-fiber/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var RoleListOptions = pagination.Options{
	Sortable: map[string]string{
		"id":        "id",
		"name":      "name",
		"createdAt": "created_at",
		"updatedAt": "updated_at",
	},
	Filterable: map[string]pagination.Filter{
		"name":      {Column: "name", Kind: pagination.FilterText},
		"createdAt": {Column: "created_at", Kind: pagination.FilterDate},
		"updatedAt": {Column: "updated_at", Kind: pagination.FilterDate},
	},
	SearchColumns: []string{"name"},
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

func (r *roleRepository) FindAll(ctx context.Context, q *pagination.Query, tz *time.Location) ([]models.Role, int64, error) {
	query, err := q.Apply(r.db.WithContext(ctx).Model(&models.Role{}), RoleListOptions, tz)
	if err != nil {
		return nil, 0, err
	}
	var total int64
	if err := query.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var roles []models.Role
	if err := query.Offset(q.Offset()).Limit(q.PerPage).Find(&roles).Error; err != nil {
		return nil, 0, err
	}
	return roles, total, nil
}

func (r *roleRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	var role models.Role
	err := r.db.WithContext(ctx).First(&role, "id = ?", id).Error
	return notFoundAsNil(&role, err)
}

func (r *roleRepository) FindByName(ctx context.Context, name string, excludeID *uuid.UUID) (*models.Role, error) {
	var role models.Role
	query := r.db.WithContext(ctx).Where("name = ?", name)
	if excludeID != nil {
		query = query.Where("id <> ?", *excludeID)
	}
	err := query.First(&role).Error
	return notFoundAsNil(&role, err)
}

func (r *roleRepository) CountByIDs(ctx context.Context, ids []uuid.UUID) (int64, error) {
	var count int64
	if len(ids) == 0 {
		return 0, nil
	}
	err := r.db.WithContext(ctx).Model(&models.Role{}).Where("id IN ?", ids).Count(&count).Error
	return count, err
}

func (r *roleRepository) PermissionIDs(ctx context.Context, roleID uuid.UUID) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	err := r.db.WithContext(ctx).Table("role_permissions").
		Where("role_id = ?", roleID).Pluck("permission_id", &ids).Error
	return ids, err
}

func (r *roleRepository) Create(ctx context.Context, role *models.Role) error {
	return r.db.WithContext(ctx).Omit(clause.Associations).Create(role).Error
}

func (r *roleRepository) Update(ctx context.Context, role *models.Role) error {
	return r.db.WithContext(ctx).Omit(clause.Associations).Save(role).Error
}

func (r *roleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Role{}, "id = ?", id).Error
}

func (r *roleRepository) ReplacePermissions(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error {
	db := r.db.WithContext(ctx)
	if err := db.Exec("DELETE FROM role_permissions WHERE role_id = ?", roleID).Error; err != nil {
		return err
	}
	if len(permissionIDs) == 0 {
		return nil
	}
	rows := make([]map[string]any, 0, len(permissionIDs))
	for _, permissionID := range permissionIDs {
		rows = append(rows, map[string]any{"role_id": roleID, "permission_id": permissionID})
	}
	return db.Table("role_permissions").Clauses(clause.OnConflict{DoNothing: true}).Create(rows).Error
}

func (r *roleRepository) SelectOptions(ctx context.Context) ([]models.Role, error) {
	var roles []models.Role
	err := r.db.WithContext(ctx).Where("name <> ?", models.SuperuserRole).Order("name").Find(&roles).Error
	return roles, err
}
