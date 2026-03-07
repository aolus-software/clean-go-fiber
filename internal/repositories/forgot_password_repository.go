package repositories

import (
	"context"

	"github.com/aolus-software/clean-go-fiber/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type forgotPasswordRepository struct {
	db *gorm.DB
}

func NewForgotPasswordRepository(db *gorm.DB) ForgotPasswordRepository {
	return &forgotPasswordRepository{db: db}
}

func (r *forgotPasswordRepository) Create(ctx context.Context, fp *models.ForgotPassword) (*models.ForgotPassword, error) {
	if err := r.db.WithContext(ctx).Create(fp).Error; err != nil {
		return nil, err
	}
	return fp, nil
}

func (r *forgotPasswordRepository) FindByToken(ctx context.Context, token string) (*models.ForgotPassword, error) {
	var fp models.ForgotPassword
	if err := r.db.WithContext(ctx).Where("token = ?", token).First(&fp).Error; err != nil {
		return nil, err
	}
	return &fp, nil
}

func (r *forgotPasswordRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.ForgotPassword{}).Error
}
