package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v7()"`
	Name        string         `json:"name" gorm:"type:varchar(255)"`
	CreatedAt   time.Time      `json:"created_at" gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"type:timestamp;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"type:timestamp;default:null"`
	Permissions []Permission   `json:"permissions" gorm:"many2many:role_permissions;"`
}
