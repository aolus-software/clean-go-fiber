package auth

import (
	"context"
	"time"

	"github.com/aolus-software/clean-go-fiber/internal/config"
	"github.com/aolus-software/clean-go-fiber/internal/libs/mailer"
	"github.com/aolus-software/clean-go-fiber/internal/libs/txmanager"
	"github.com/aolus-software/clean-go-fiber/internal/libs/utils"
	"github.com/aolus-software/clean-go-fiber/internal/models"
	"github.com/aolus-software/clean-go-fiber/internal/queue/tasks"
	"github.com/aolus-software/clean-go-fiber/internal/repositories"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

type AccountMail struct {
	emailVerifyRepo repositories.EmailVerificationRepository
	resetTokenRepo  repositories.PasswordResetTokenRepository
	queueClient     *asynq.Client
	tx              *txmanager.TxManager
	cfg             *config.Config
}

func NewAccountMail(
	emailVerifyRepo repositories.EmailVerificationRepository,
	resetTokenRepo repositories.PasswordResetTokenRepository,
	queueClient *asynq.Client,
	tx *txmanager.TxManager,
	cfg *config.Config,
) *AccountMail {
	return &AccountMail{
		emailVerifyRepo: emailVerifyRepo,
		resetTokenRepo:  resetTokenRepo,
		queueClient:     queueClient,
		tx:              tx,
		cfg:             cfg,
	}
}

func (m *AccountMail) IssueVerificationToken(ctx context.Context, tx *gorm.DB, userID uuid.UUID) (string, error) {
	raw, err := utils.GenerateSecureToken()
	if err != nil {
		return "", err
	}
	repo := m.emailVerifyRepo.WithTx(tx)
	if err := repo.MarkAllUsed(ctx, userID, time.Now().UTC()); err != nil {
		return "", err
	}
	err = repo.Create(ctx, &models.EmailVerification{
		UserID:    userID,
		Token:     utils.HashToken(raw),
		ExpiredAt: utils.VerificationTokenExpiry(),
	})
	return raw, err
}

func (m *AccountMail) IssuePasswordResetToken(ctx context.Context, tx *gorm.DB, userID uuid.UUID) (string, error) {
	raw, err := utils.GenerateSecureToken()
	if err != nil {
		return "", err
	}
	repo := m.resetTokenRepo.WithTx(tx)
	if err := repo.MarkAllUsed(ctx, userID, time.Now().UTC()); err != nil {
		return "", err
	}
	err = repo.Create(ctx, &models.PasswordResetToken{
		UserID:    userID,
		Token:     utils.HashToken(raw),
		ExpiredAt: utils.PasswordResetTokenExpiry(),
	})
	return raw, err
}

func (m *AccountMail) EnqueueVerification(ctx context.Context, user *models.User, rawToken string) error {
	subject, body := mailer.VerificationEmail(m.cfg.App.Name, user.Name, mailer.VerificationLink(m.cfg.App.ClientURL, rawToken))
	return tasks.EnqueueEmail(ctx, m.queueClient, tasks.EmailSendPayload{To: user.Email, Subject: subject, Body: body})
}

func (m *AccountMail) EnqueuePasswordReset(ctx context.Context, user *models.User, rawToken string) error {
	subject, body := mailer.ResetPasswordEmail(m.cfg.App.Name, user.Name, mailer.ResetPasswordLink(m.cfg.App.ClientURL, rawToken))
	return tasks.EnqueueEmail(ctx, m.queueClient, tasks.EmailSendPayload{To: user.Email, Subject: subject, Body: body})
}

func (m *AccountMail) SendVerification(ctx context.Context, user *models.User) error {
	var raw string
	err := m.tx.WithinTransaction(ctx, func(tx *gorm.DB) error {
		var err error
		raw, err = m.IssueVerificationToken(ctx, tx, user.ID)
		return err
	})
	if err != nil {
		return err
	}
	return m.EnqueueVerification(ctx, user, raw)
}

func (m *AccountMail) SendPasswordReset(ctx context.Context, user *models.User) error {
	var raw string
	err := m.tx.WithinTransaction(ctx, func(tx *gorm.DB) error {
		var err error
		raw, err = m.IssuePasswordResetToken(ctx, tx, user.ID)
		return err
	})
	if err != nil {
		return err
	}
	return m.EnqueuePasswordReset(ctx, user, raw)
}
