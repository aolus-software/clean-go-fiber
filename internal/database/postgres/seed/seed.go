package seed

import (
	"fmt"
	"time"

	"github.com/aolus-software/clean-go-fiber/internal/libs/utils"
	"github.com/aolus-software/clean-go-fiber/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const DefaultPassword = "S3crEtP4ssw0rd!"

var (
	permissionGroups  = []string{"user", "role", "permission"}
	permissionActions = []string{"list", "create", "detail", "edit", "delete"}
	seededRoles       = []string{models.SuperuserRole, "admin"}
)

func Run(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		names, err := seedPermissions(tx)
		if err != nil {
			return err
		}
		roles, err := seedRoles(tx, names)
		if err != nil {
			return err
		}
		return seedUsers(tx, roles)
	})
}

func seedPermissions(tx *gorm.DB) ([]string, error) {
	names := make([]string, 0, len(permissionGroups)*len(permissionActions))
	rows := make([]models.Permission, 0, cap(names))
	for _, group := range permissionGroups {
		for _, action := range permissionActions {
			name := group + " " + action
			names = append(names, name)
			rows = append(rows, models.Permission{Name: name, Group: group})
		}
	}
	if err := tx.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "name"}}, DoNothing: true}).Create(&rows).Error; err != nil {
		return nil, fmt.Errorf("seed permissions: %w", err)
	}
	return names, nil
}

func seedRoles(tx *gorm.DB, permissionNames []string) (map[string]models.Role, error) {
	var permissions []models.Permission
	if err := tx.Where("name IN ?", permissionNames).Find(&permissions).Error; err != nil {
		return nil, fmt.Errorf("seed roles: load permissions: %w", err)
	}

	roles := make(map[string]models.Role, len(seededRoles))
	for _, name := range seededRoles {
		role := models.Role{Name: name}
		if err := tx.Where(models.Role{Name: name}).FirstOrCreate(&role).Error; err != nil {
			return nil, fmt.Errorf("seed role %s: %w", name, err)
		}
		grants := make([]map[string]any, 0, len(permissions))
		for _, p := range permissions {
			grants = append(grants, map[string]any{"role_id": role.ID, "permission_id": p.ID})
		}
		if len(grants) > 0 {
			if err := tx.Table("role_permissions").Clauses(clause.OnConflict{DoNothing: true}).Create(grants).Error; err != nil {
				return nil, fmt.Errorf("seed role %s grants: %w", name, err)
			}
		}
		roles[name] = role
	}
	return roles, nil
}

func seedUsers(tx *gorm.DB, roles map[string]models.Role) error {
	hashed, err := utils.HashPassword(DefaultPassword)
	if err != nil {
		return fmt.Errorf("seed users: hash password: %w", err)
	}
	now := time.Now().UTC()

	for _, roleName := range seededRoles {
		email := roleName + "@example.com"
		var user models.User
		err := tx.Where("email = ?", email).First(&user).Error
		if err == nil {
			continue
		}
		if err != gorm.ErrRecordNotFound {
			return fmt.Errorf("seed user %s: %w", email, err)
		}

		user = models.User{
			Name:            roleName,
			Email:           email,
			Password:        hashed,
			Status:          models.UserStatusActive,
			EmailVerifiedAt: &now,
		}
		if err := tx.Omit(clause.Associations).Create(&user).Error; err != nil {
			return fmt.Errorf("seed user %s: %w", email, err)
		}
		assignment := map[string]any{"user_id": user.ID, "role_id": roles[roleName].ID}
		if err := tx.Table("user_roles").Clauses(clause.OnConflict{DoNothing: true}).Create(assignment).Error; err != nil {
			return fmt.Errorf("seed user %s role: %w", email, err)
		}
	}
	return nil
}
