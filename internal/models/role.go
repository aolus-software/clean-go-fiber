package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const SuperuserRole = "superuser"

type Role struct {
	ID          uuid.UUID    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string       `gorm:"type:varchar(100);not null;unique"`
	CreatedAt   time.Time    `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt   time.Time    `gorm:"type:timestamp;autoUpdateTime"`
	Permissions []Permission `gorm:"many2many:role_permissions;"`
}

func (r *Role) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}
