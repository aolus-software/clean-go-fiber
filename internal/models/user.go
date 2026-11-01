package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	UserStatusActive    = "active"
	UserStatusInactive  = "inactive"
	UserStatusSuspended = "suspended"
	UserStatusBlocked   = "blocked"
)

var UserStatuses = []string{UserStatusActive, UserStatusInactive, UserStatusSuspended, UserStatusBlocked}

type User struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name            string         `gorm:"type:varchar(255);not null"`
	Email           string         `gorm:"type:varchar(255);not null"`
	Password        string         `gorm:"type:varchar(255);not null"`
	Status          string         `gorm:"type:user_status;not null;default:active"`
	Remark          *string        `gorm:"type:varchar(255)"`
	EmailVerifiedAt *time.Time     `gorm:"type:timestamp"`
	CreatedAt       time.Time      `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt       time.Time      `gorm:"type:timestamp;autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"type:timestamp"`
	Roles           []Role         `gorm:"many2many:user_roles;"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
