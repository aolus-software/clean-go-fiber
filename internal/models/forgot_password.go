package models

import (
	"time"

	"github.com/google/uuid"
)

type ForgotPassword struct {
	Id        uuid.UUID `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v7();"`
	UserId    uuid.UUID `json:"userId" gorm:"type:uuid;not null;"`
	Token     string    `json:"token" gorm:"type:varchar(255);not null;"`
	ExpiresAt time.Time `json:"expiresAt" gorm:"type:timestamp;not null;"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;autoCreateTime;"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp;autoUpdateTime;"`
}
