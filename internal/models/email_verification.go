package models

import (
	"time"

	"github.com/google/uuid"
)

type EmailVerification struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null"`
	Token     string     `gorm:"type:varchar(255);not null;unique"`
	ExpiredAt time.Time  `gorm:"type:timestamp;not null"`
	UsedAt    *time.Time `gorm:"type:timestamp"`
	CreatedAt time.Time  `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt time.Time  `gorm:"type:timestamp;autoUpdateTime"`
}

func (EmailVerification) TableName() string {
	return "email_verifications"
}
