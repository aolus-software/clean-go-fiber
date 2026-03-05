package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID              uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v7()"`
	Name            string         `json:"name" gorm:"type:varchar(255)"`
	Email           string         `json:"email" gorm:"type:varchar(255);uniqueIndex"`
	Password        string         `json:"-" gorm:"type:varchar(255)"`
	Status          string         `json:"status" gorm:"type:varchar(255)"`
	EmailVerifiedAt *time.Time     `json:"email_verified_at,omitempty" gorm:"type:timestamp;default:null"`
	CreatedAt       time.Time      `json:"created_at" gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt       time.Time      `json:"updated_at" gorm:"type:timestamp;autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"type:timestamp;default:null"`
	Roles           []Role         `json:"roles" gorm:"many2many:user_roles;"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == (uuid.UUID{}) {
		u.ID = uuid.New()
	}
	return nil
}
