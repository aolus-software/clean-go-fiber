package auth

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/aolus-software/clean-go-fiber/internal/config"
	apperrors "github.com/aolus-software/clean-go-fiber/internal/errors"
	"github.com/aolus-software/clean-go-fiber/internal/libs/cache"
	"github.com/aolus-software/clean-go-fiber/internal/libs/txmanager"
	"github.com/aolus-software/clean-go-fiber/internal/libs/utils"
	"github.com/aolus-software/clean-go-fiber/internal/models"
	"github.com/aolus-software/clean-go-fiber/internal/repositories"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service interface {
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	Register(ctx context.Context, req *RegisterRequest) error
	ResendVerification(ctx context.Context, req *EmailRequest) error
	VerifyEmail(ctx context.Context, req *TokenRequest) error
	ForgotPassword(ctx context.Context, req *EmailRequest) error
	ResetPassword(ctx context.Context, req *ResetPasswordRequest) error
	Refresh(ctx context.Context, req *RefreshRequest) (*RefreshResponse, error)
}

type service struct {
	userRepo        repositories.UserRepository
	emailVerifyRepo repositories.EmailVerificationRepository
	resetTokenRepo  repositories.PasswordResetTokenRepository
	accountMail     *AccountMail
	tx              *txmanager.TxManager
	cache           *cache.Cache
	cfg             *config.Config
	log             *zap.Logger
}

func NewService(
	userRepo repositories.UserRepository,
	emailVerifyRepo repositories.EmailVerificationRepository,
	resetTokenRepo repositories.PasswordResetTokenRepository,
	accountMail *AccountMail,
	tx *txmanager.TxManager,
	c *cache.Cache,
	cfg *config.Config,
	log *zap.Logger,
) Service {
	return &service{
		userRepo:        userRepo,
		emailVerifyRepo: emailVerifyRepo,
		resetTokenRepo:  resetTokenRepo,
		accountMail:     accountMail,
		tx:              tx,
		cache:           c,
		cfg:             cfg,
		log:             log,
	}
}

var dummyPasswordHash = sync.OnceValue(func() string {
	hash, err := utils.HashPassword("timing-equaliser-not-a-real-password")
	if err != nil {
		return ""
	}
	return hash
})

func invalidCredentials() error {
	return apperrors.NewBadRequestError("Invalid credentials", apperrors.Field("email", "Invalid email or password"))
}

func invalidToken() error {
	return apperrors.NewBadRequestError("Invalid token", apperrors.Field("token", "Invalid or expired token"))
}

func (s *service) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	user, err := s.userRepo.FindLiveByEmail(ctx, strings.TrimSpace(req.Email), nil)
	if err != nil {
		return nil, err
	}

	hash := dummyPasswordHash()
	if user != nil {
		hash = user.Password
	}
	passwordOK := utils.CheckPasswordHash(req.Password, hash)
	if user == nil || !passwordOK {
		return nil, invalidCredentials()
	}
	if user.EmailVerifiedAt == nil {
		return nil, apperrors.NewBadRequestError("Email not verified", apperrors.Field("email", "Email address is not verified"))
	}
	if user.Status != models.UserStatusActive {
		return nil, apperrors.NewBadRequestError("Account inactive", apperrors.Field("email", "Account is not active"))
	}

	info, err := s.userRepo.FindUserInformation(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, invalidCredentials()
	}
	if err := s.cache.Set(ctx, cache.Key(cache.UserInformationKey, user.ID.String()), info, cache.UserInformationTTL); err != nil {
		return nil, err
	}

	access, refresh, err := s.issueTokens(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	return &LoginResponse{User: info, AccessToken: access, RefreshToken: refresh}, nil
}

func (s *service) Register(ctx context.Context, req *RegisterRequest) error {
	email := strings.TrimSpace(req.Email)
	existing, err := s.userRepo.FindLiveByEmail(ctx, email, nil)
	if err != nil {
		return err
	}
	if existing != nil {
		return emailTaken()
	}

	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := &models.User{Name: req.Name, Email: email, Password: hashed, Status: models.UserStatusActive}
	var rawToken string
	err = s.tx.WithinTransaction(ctx, func(tx *gorm.DB) error {
		if err := s.userRepo.WithTx(tx).Create(ctx, user); err != nil {
			return err
		}
		rawToken, err = s.accountMail.IssueVerificationToken(ctx, tx, user.ID)
		return err
	})
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return emailTaken()
	}
	if err != nil {
		return err
	}

	if err := s.accountMail.EnqueueVerification(ctx, user, rawToken); err != nil {
		s.log.Error("failed to enqueue verification email", zap.String("userId", user.ID.String()), zap.Error(err))
	}
	return nil
}

func emailTaken() error {
	return apperrors.NewUnprocessableEntityError("Validation error", apperrors.Field("email", "Email is already registered"))
}

func (s *service) ResendVerification(ctx context.Context, req *EmailRequest) error {
	user, err := s.userRepo.FindLiveByEmail(ctx, strings.TrimSpace(req.Email), nil)
	if err != nil {
		return err
	}
	if user == nil || user.EmailVerifiedAt != nil {
		return nil
	}
	return s.accountMail.SendVerification(ctx, user)
}

func (s *service) VerifyEmail(ctx context.Context, req *TokenRequest) error {
	now := time.Now().UTC()
	record, err := s.emailVerifyRepo.FindUsable(ctx, utils.HashToken(req.Token), now)
	if err != nil {
		return err
	}
	if record == nil {
		return invalidToken()
	}
	user, err := s.userRepo.FindByID(ctx, record.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		return invalidToken()
	}

	err = s.tx.WithinTransaction(ctx, func(tx *gorm.DB) error {
		user.EmailVerifiedAt = &now
		if err := s.userRepo.WithTx(tx).Update(ctx, user); err != nil {
			return err
		}
		return s.emailVerifyRepo.WithTx(tx).MarkAllUsed(ctx, user.ID, now)
	})
	if err != nil {
		return err
	}
	return s.cache.ForgetUsers(ctx, user.ID)
}

func (s *service) ForgotPassword(ctx context.Context, req *EmailRequest) error {
	user, err := s.userRepo.FindLiveByEmail(ctx, strings.TrimSpace(req.Email), nil)
	if err != nil {
		return err
	}
	if user == nil {
		return nil
	}
	return s.accountMail.SendPasswordReset(ctx, user)
}

func (s *service) ResetPassword(ctx context.Context, req *ResetPasswordRequest) error {
	now := time.Now().UTC()
	record, err := s.resetTokenRepo.FindUsable(ctx, utils.HashToken(req.Token), now)
	if err != nil {
		return err
	}
	if record == nil {
		return invalidToken()
	}
	user, err := s.userRepo.FindByID(ctx, record.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		return invalidToken()
	}

	hashed, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	err = s.tx.WithinTransaction(ctx, func(tx *gorm.DB) error {
		user.Password = hashed
		if err := s.userRepo.WithTx(tx).Update(ctx, user); err != nil {
			return err
		}
		return s.resetTokenRepo.WithTx(tx).MarkAllUsed(ctx, user.ID, now)
	})
	if err != nil {
		return err
	}
	return s.cache.RevokeSessions(ctx, user.ID)
}

func (s *service) Refresh(ctx context.Context, req *RefreshRequest) (*RefreshResponse, error) {
	claims, err := utils.ParseToken(req.RefreshToken, utils.TokenTypeRefresh, s.cfg.Jwt.RefreshSecret)
	if err != nil {
		return nil, apperrors.NewUnauthorizedError("Invalid or expired refresh token")
	}
	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, apperrors.NewUnauthorizedError("Invalid or expired refresh token")
	}
	storedID, err := s.cache.RefreshTokenID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if claims.ID == "" || storedID != claims.ID {
		return nil, apperrors.NewUnauthorizedError("Invalid or expired refresh token")
	}
	info, err := s.userRepo.FindUserInformation(ctx, userID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		if err := s.cache.RevokeSessions(ctx, userID); err != nil {
			return nil, err
		}
		return nil, apperrors.NewUnauthorizedError("Invalid or expired refresh token")
	}
	if err := s.cache.ForgetUsers(ctx, userID); err != nil {
		return nil, err
	}

	access, refresh, err := s.issueTokens(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &RefreshResponse{AccessToken: access, RefreshToken: refresh}, nil
}

/*
The refresh token's jti is stored per user in Redis, so issuing a pair rotates out the previous
refresh token and deleting the key revokes it.
*/
func (s *service) issueTokens(ctx context.Context, userID uuid.UUID) (string, string, error) {
	access, err := utils.SignToken(userID.String(), utils.TokenTypeAccess, "", s.cfg.Jwt.Secret, config.AccessTokenLifetime)
	if err != nil {
		return "", "", err
	}
	tokenID, err := utils.GenerateSecureToken()
	if err != nil {
		return "", "", err
	}
	refresh, err := utils.SignToken(userID.String(), utils.TokenTypeRefresh, tokenID, s.cfg.Jwt.RefreshSecret, config.RefreshTokenLifetime)
	if err != nil {
		return "", "", err
	}
	if err := s.cache.Set(ctx, cache.Key(cache.RefreshTokenKey, userID.String()), tokenID, config.RefreshTokenLifetime); err != nil {
		return "", "", err
	}
	return access, refresh, nil
}
