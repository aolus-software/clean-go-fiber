package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/aolus-software/clean-go-fiber/internal/config"
	"github.com/aolus-software/clean-go-fiber/internal/errors"
	"github.com/aolus-software/clean-go-fiber/internal/models"
	"github.com/aolus-software/clean-go-fiber/internal/libs/utils"
	"github.com/aolus-software/clean-go-fiber/internal/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	Register(ctx context.Context, req *RegisterRequest) (*TokenResponse, error)
	Login(ctx context.Context, req *LoginRequest) (*TokenResponse, error)
	RefreshToken(ctx context.Context, req *RefreshTokenRequest) (*TokenResponse, error)
	ForgotPassword(ctx context.Context, req *ForgotPasswordRequest) error
	ResetPassword(ctx context.Context, req *ResetPasswordRequest) error
	VerifyEmail(ctx context.Context, req *VerifyEmailRequest) error
	Me(ctx context.Context, userID string) (*UserResponse, error)
}

type service struct {
	userRepo          repositories.UserRepository
	emailVerifyRepo   repositories.EmailVerificationRepository
	forgotPasswordRepo repositories.ForgotPasswordRepository
	cfg               *config.Config
}

func NewService(
	userRepo repositories.UserRepository,
	emailVerifyRepo repositories.EmailVerificationRepository,
	forgotPasswordRepo repositories.ForgotPasswordRepository,
	cfg *config.Config,
) Service {
	return &service{
		userRepo:          userRepo,
		emailVerifyRepo:   emailVerifyRepo,
		forgotPasswordRepo: forgotPasswordRepo,
		cfg:               cfg,
	}
}

func (s *service) Register(ctx context.Context, req *RegisterRequest) (*TokenResponse, error) {
	existing, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.NewInternalServerError("failed to check existing user")
	}
	if existing != nil {
		return nil, errors.NewConflictError("email already in use")
	}

	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.NewInternalServerError("failed to process password")
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashed,
		Status:   "active",
	}
	created, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, errors.NewInternalServerError("failed to create user")
	}

	return s.issueTokenPair(created)
}

func (s *service) Login(ctx context.Context, req *LoginRequest) (*TokenResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.NewUnauthorizedError("invalid credentials")
	}
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.NewUnauthorizedError("invalid credentials")
	}
	return s.issueTokenPair(user)
}

func (s *service) RefreshToken(ctx context.Context, req *RefreshTokenRequest) (*TokenResponse, error) {
	claims, err := utils.ParseToken(req.RefreshToken, s.cfg.Jwt.RefreshSecretKey)
	if err != nil {
		return nil, errors.NewUnauthorizedError("invalid or expired refresh token")
	}

	user, err := s.userRepo.FindByEmail(ctx, claims.Email)
	if err != nil {
		return nil, errors.NewUnauthorizedError("user not found")
	}

	return s.issueTokenPair(user)
}

func (s *service) ForgotPassword(ctx context.Context, req *ForgotPasswordRequest) error {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil
	}

	if err := s.forgotPasswordRepo.DeleteByUserID(ctx, user.ID); err != nil {
		return errors.NewInternalServerError("failed to process request")
	}

	token, err := generateSecureToken()
	if err != nil {
		return errors.NewInternalServerError("failed to generate token")
	}

	fp := &models.ForgotPassword{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}
	if _, err := s.forgotPasswordRepo.Create(ctx, fp); err != nil {
		return errors.NewInternalServerError("failed to create reset token")
	}

	return nil
}

func (s *service) ResetPassword(ctx context.Context, req *ResetPasswordRequest) error {
	fp, err := s.forgotPasswordRepo.FindByToken(ctx, req.Token)
	if err != nil {
		return errors.NewBadRequestError("invalid or expired token")
	}
	if time.Now().After(fp.ExpiresAt) {
		return errors.NewBadRequestError("token has expired")
	}

	user, err := s.userRepo.FindById(ctx, fp.UserID)
	if err != nil {
		return errors.NewNotFoundError("user not found")
	}

	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		return errors.NewInternalServerError("failed to process password")
	}
	user.Password = hashed

	if _, err := s.userRepo.Update(ctx, user); err != nil {
		return errors.NewInternalServerError("failed to update password")
	}

	return s.forgotPasswordRepo.DeleteByUserID(ctx, user.ID)
}

func (s *service) VerifyEmail(ctx context.Context, req *VerifyEmailRequest) error {
	ev, err := s.emailVerifyRepo.FindByToken(ctx, req.Token)
	if err != nil {
		return errors.NewBadRequestError("invalid or expired token")
	}
	if time.Now().After(ev.ExpiresIn) {
		return errors.NewBadRequestError("token has expired")
	}

	user, err := s.userRepo.FindById(ctx, ev.UserID)
	if err != nil {
		return errors.NewNotFoundError("user not found")
	}

	now := time.Now()
	user.EmailVerifiedAt = &now
	if _, err := s.userRepo.Update(ctx, user); err != nil {
		return errors.NewInternalServerError("failed to verify email")
	}

	return s.emailVerifyRepo.DeleteByUserID(ctx, user.ID)
}

func (s *service) Me(ctx context.Context, userID string) (*UserResponse, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.NewBadRequestError("invalid user id")
	}
	user, err := s.userRepo.FindById(ctx, id)
	if err != nil {
		return nil, errors.NewNotFoundError("user not found")
	}
	return &UserResponse{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *service) issueTokenPair(user *models.User) (*TokenResponse, error) {
	userID := user.ID.String()
	accessToken, err := utils.GenerateAccessToken(userID, user.Email, s.cfg.Jwt.SecretKey, s.cfg.Jwt.ExpiresIn)
	if err != nil {
		return nil, errors.NewInternalServerError("failed to generate access token")
	}
	refreshToken, err := utils.GenerateRefreshToken(userID, user.Email, s.cfg.Jwt.RefreshSecretKey, s.cfg.Jwt.RefreshExpiresIn)
	if err != nil {
		return nil, errors.NewInternalServerError("failed to generate refresh token")
	}
	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
	}, nil
}

func generateSecureToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
