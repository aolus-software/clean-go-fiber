package repositories

import (
	"context"
	"time"

	"github.com/aolus-software/clean-go-fiber/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type passwordResetTokenRepository struct {
	db *gorm.DB
}

func NewPasswordResetTokenRepository(db *gorm.DB) PasswordResetTokenRepository {
	return &passwordResetTokenRepository{db: db}
}

func (r *passwordResetTokenRepository) WithTx(tx *gorm.DB) PasswordResetTokenRepository {
	return &passwordResetTokenRepository{db: tx}
}

func (r *passwordResetTokenRepository) Create(ctx context.Context, t *models.PasswordResetToken) error {
	return r.db.WithContext(ctx).Create(t).Error
}

func (r *passwordResetTokenRepository) FindUsable(ctx context.Context, tokenHash string, now time.Time) (*models.PasswordResetToken, error) {
	var t models.PasswordResetToken
	err := r.db.WithContext(ctx).
		Where("token = ? AND used_at IS NULL AND expired_at > ?", tokenHash, now).
		First(&t).Error
	return notFoundAsNil(&t, err)
}

func (r *passwordResetTokenRepository) MarkAllUsed(ctx context.Context, userID uuid.UUID, now time.Time) error {
	return r.db.WithContext(ctx).Model(&models.PasswordResetToken{}).
		Where("user_id = ? AND used_at IS NULL", userID).
		Update("used_at", now).Error
}
