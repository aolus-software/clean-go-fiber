package repositories

import (
	"context"
	"time"

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

func (r *emailVerificationRepository) WithTx(tx *gorm.DB) EmailVerificationRepository {
	return &emailVerificationRepository{db: tx}
}

func (r *emailVerificationRepository) Create(ctx context.Context, ev *models.EmailVerification) error {
	return r.db.WithContext(ctx).Create(ev).Error
}

func (r *emailVerificationRepository) FindUsable(ctx context.Context, tokenHash string, now time.Time) (*models.EmailVerification, error) {
	var ev models.EmailVerification
	err := r.db.WithContext(ctx).
		Where("token = ? AND used_at IS NULL AND expired_at > ?", tokenHash, now).
		First(&ev).Error
	return notFoundAsNil(&ev, err)
}

func (r *emailVerificationRepository) MarkAllUsed(ctx context.Context, userID uuid.UUID, now time.Time) error {
	return r.db.WithContext(ctx).Model(&models.EmailVerification{}).
		Where("user_id = ? AND used_at IS NULL", userID).
		Update("used_at", now).Error
}
