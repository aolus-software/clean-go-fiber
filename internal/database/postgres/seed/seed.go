package seed

import (
	"fmt"

	"github.com/aolus-software/clean-go-fiber/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Run(db *gorm.DB) error {
	permissions := []models.Permission{
		{Name: "read", Modules: "users"},
		{Name: "create", Modules: "users"},
		{Name: "update", Modules: "users"},
		{Name: "delete", Modules: "users"},
		{Name: "read", Modules: "roles"},
		{Name: "create", Modules: "roles"},
		{Name: "update", Modules: "roles"},
		{Name: "delete", Modules: "roles"},
		{Name: "read", Modules: "permissions"},
		{Name: "create", Modules: "permissions"},
		{Name: "update", Modules: "permissions"},
		{Name: "delete", Modules: "permissions"},
	}

	for i := range permissions {
		result := db.Where(models.Permission{Name: permissions[i].Name, Modules: permissions[i].Modules}).
			FirstOrCreate(&permissions[i])
		if result.Error != nil {
			return fmt.Errorf("seed permission %s.%s: %w", permissions[i].Modules, permissions[i].Name, result.Error)
		}
	}

	superuserRole := models.Role{Name: "superuser"}
	if result := db.Where(models.Role{Name: "superuser"}).FirstOrCreate(&superuserRole); result.Error != nil {
		return fmt.Errorf("seed superuser role: %w", result.Error)
	}

	superadminRole := models.Role{Name: "superadmin"}
	if result := db.Where(models.Role{Name: "superadmin"}).FirstOrCreate(&superadminRole); result.Error != nil {
		return fmt.Errorf("seed superadmin role: %w", result.Error)
	}

	if err := db.Model(&superuserRole).Association("Permissions").Replace(permissions); err != nil {
		return fmt.Errorf("seed superuser role permissions: %w", err)
	}

	if err := db.Model(&superadminRole).Association("Permissions").Replace(permissions); err != nil {
		return fmt.Errorf("seed superadmin role permissions: %w", err)
	}

	hashedPw, err := bcrypt.GenerateFromPassword([]byte("Password123!"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("seed hash password: %w", err)
	}

	superuser := models.User{
		Name:     "Superuser",
		Email:    "superuser@example.com",
		Password: string(hashedPw),
		Status:   "active",
	}
	if result := db.Where(models.User{Email: "superuser@example.com"}).FirstOrCreate(&superuser); result.Error != nil {
		return fmt.Errorf("seed superuser user: %w", result.Error)
	}

	if err := db.Model(&superuser).Association("Roles").Replace([]models.Role{superuserRole}); err != nil {
		return fmt.Errorf("seed superuser user roles: %w", err)
	}

	return nil
}
