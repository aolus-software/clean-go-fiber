package repositories

import (
	"context"
	"slices"
	"time"

	"github.com/aolus-software/clean-go-fiber/internal/libs/pagination"
	"github.com/aolus-software/clean-go-fiber/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var UserListOptions = pagination.Options{
	Sortable: map[string]string{
		"id":        "users.id",
		"name":      "users.name",
		"email":     "users.email",
		"status":    "users.status",
		"createdAt": "users.created_at",
		"updatedAt": "users.updated_at",
	},
	Filterable: map[string]pagination.Filter{
		"status":    {Column: "users.status", Kind: pagination.FilterEnum, Enum: models.UserStatuses},
		"name":      {Column: "users.name", Kind: pagination.FilterText},
		"email":     {Column: "users.email", Kind: pagination.FilterText},
		"roleId":    {Kind: pagination.FilterID, Expr: "users.id IN (SELECT user_id FROM user_roles WHERE role_id IN ?)"},
		"createdAt": {Column: "users.created_at", Kind: pagination.FilterDate},
		"updatedAt": {Column: "users.updated_at", Kind: pagination.FilterDate},
	},
	SearchColumns: []string{"users.name", "users.email", "users.status::text"},
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

func (r *userRepository) FindAll(ctx context.Context, q *pagination.Query, tz *time.Location) ([]models.User, int64, error) {
	query, err := q.Apply(r.db.WithContext(ctx).Model(&models.User{}), UserListOptions, tz)
	if err != nil {
		return nil, 0, err
	}

	var total int64
	if err := query.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var users []models.User
	if err := query.Preload("Roles", func(db *gorm.DB) *gorm.DB { return db.Order("roles.name") }).
		Offset(q.Offset()).Limit(q.PerPage).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Preload("Roles", func(db *gorm.DB) *gorm.DB { return db.Order("roles.name") }).
		First(&user, "id = ?", id).Error
	return notFoundAsNil(&user, err)
}

func (r *userRepository) FindLiveByEmail(ctx context.Context, email string, excludeID *uuid.UUID) (*models.User, error) {
	var user models.User
	query := r.db.WithContext(ctx).Where("LOWER(email) = LOWER(?)", email)
	if excludeID != nil {
		query = query.Where("id <> ?", *excludeID)
	}
	err := query.First(&user).Error
	return notFoundAsNil(&user, err)
}

func (r *userRepository) FindUserInformation(ctx context.Context, id uuid.UUID) (*UserInformation, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Preload("Roles.Permissions").
		Where("status = ?", models.UserStatusActive).
		First(&user, "id = ?", id).Error
	found, err := notFoundAsNil(&user, err)
	if err != nil || found == nil {
		return nil, err
	}

	info := &UserInformation{
		ID:          user.ID.String(),
		Name:        user.Name,
		Email:       user.Email,
		Roles:       []string{},
		Permissions: []string{},
	}
	for _, role := range user.Roles {
		info.Roles = append(info.Roles, role.Name)
		for _, perm := range role.Permissions {
			if !slices.Contains(info.Permissions, perm.Name) {
				info.Permissions = append(info.Permissions, perm.Name)
			}
		}
	}
	slices.Sort(info.Roles)
	slices.Sort(info.Permissions)
	return info, nil
}

func (r *userRepository) FindIDsByRoleIDs(ctx context.Context, roleIDs []uuid.UUID) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	if len(roleIDs) == 0 {
		return ids, nil
	}
	err := r.db.WithContext(ctx).Table("user_roles").Distinct("user_id").
		Where("role_id IN ?", roleIDs).Pluck("user_id", &ids).Error
	return ids, err
}

func (r *userRepository) FindIDsByPermissionID(ctx context.Context, permissionID uuid.UUID) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	err := r.db.WithContext(ctx).Table("user_roles").Distinct("user_roles.user_id").
		Joins("JOIN role_permissions ON role_permissions.role_id = user_roles.role_id").
		Where("role_permissions.permission_id = ?", permissionID).
		Pluck("user_roles.user_id", &ids).Error
	return ids, err
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Omit(clause.Associations).Create(user).Error
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Omit(clause.Associations).Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id).Error
}

func (r *userRepository) ReplaceRoles(ctx context.Context, userID uuid.UUID, roleIDs []uuid.UUID) error {
	db := r.db.WithContext(ctx)
	if err := db.Exec("DELETE FROM user_roles WHERE user_id = ?", userID).Error; err != nil {
		return err
	}
	if len(roleIDs) == 0 {
		return nil
	}
	rows := make([]map[string]any, 0, len(roleIDs))
	for _, roleID := range roleIDs {
		rows = append(rows, map[string]any{"user_id": userID, "role_id": roleID})
	}
	return db.Table("user_roles").Clauses(clause.OnConflict{DoNothing: true}).Create(rows).Error
}
