package txmanager

import (
	"context"

	"gorm.io/gorm"
)

type TxManager struct {
	db *gorm.DB
}

func New(db *gorm.DB) *TxManager {
	return &TxManager{db: db}
}

func (m *TxManager) WithinTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return m.db.WithContext(ctx).Transaction(fn)
}
