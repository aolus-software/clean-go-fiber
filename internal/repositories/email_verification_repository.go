package repositories

import (
	"context"

	"github.com/aolus-software/clean-go-fiber/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type emailVerificationRepository struct {
	db *gorm.DB
}

func NewEmailVerificationRepository(db *gorm.DB) EmailVerificationRepository {
	return &emailVerificationRepository{db: db}
}

func (r *emailVerificationRepository) Create(ctx context.Context, ev *models.EmailVerification) (*models.EmailVerification, error) {
	if err := r.db.WithContext(ctx).Create(ev).Error; err != nil {
		return nil, err
	}
	return ev, nil
}

func (r *emailVerificationRepository) FindByToken(ctx context.Context, token string) (*models.EmailVerification, error) {
	var ev models.EmailVerification
	if err := r.db.WithContext(ctx).Where("token = ?", token).First(&ev).Error; err != nil {
		return nil, err
	}
	return &ev, nil
}

func (r *emailVerificationRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.EmailVerification{}).Error
}
