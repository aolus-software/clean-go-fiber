package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Permission struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name      string    `gorm:"type:varchar(255);not null;unique"`
	Group     string    `gorm:"column:group;type:varchar(100);not null"`
	CreatedAt time.Time `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt time.Time `gorm:"type:timestamp;autoUpdateTime"`
}

func (p *Permission) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
