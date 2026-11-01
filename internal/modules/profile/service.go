package profile

import (
	"context"
	"errors"
	"strings"

	apperrors "github.com/aolus-software/clean-go-fiber/internal/errors"
	"github.com/aolus-software/clean-go-fiber/internal/libs/cache"
	"github.com/aolus-software/clean-go-fiber/internal/libs/utils"
	"github.com/aolus-software/clean-go-fiber/internal/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	Update(ctx context.Context, userID uuid.UUID, req *UpdateProfileRequest) (*repositories.UserInformation, error)
	UpdatePassword(ctx context.Context, userID uuid.UUID, req *UpdatePasswordRequest) error
}

type service struct {
	userRepo repositories.UserRepository
	cache    *cache.Cache
}

func NewService(userRepo repositories.UserRepository, c *cache.Cache) Service {
	return &service{userRepo: userRepo, cache: c}
}

func emailTaken() error {
	return apperrors.NewUnprocessableEntityError("Validation error", apperrors.Field("email", "Email is already registered"))
}

func (s *service) Update(ctx context.Context, userID uuid.UUID, req *UpdateProfileRequest) (*repositories.UserInformation, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperrors.NewNotFoundError("User not found")
	}

	email := strings.TrimSpace(req.Email)
	existing, err := s.userRepo.FindLiveByEmail(ctx, email, &user.ID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, emailTaken()
	}

	user.Name = req.Name
	user.Email = email
	if err := s.userRepo.Update(ctx, user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, emailTaken()
		}
		return nil, err
	}

	info, err := s.userRepo.FindUserInformation(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, apperrors.NewUnauthorizedError("Unauthorized")
	}
	if err := s.cache.Set(ctx, cache.Key(cache.UserInformationKey, user.ID.String()), info, cache.UserInformationTTL); err != nil {
		return nil, err
	}
	return info, nil
}

func (s *service) UpdatePassword(ctx context.Context, userID uuid.UUID, req *UpdatePasswordRequest) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return apperrors.NewUnauthorizedError("Unauthorized")
	}
	if !utils.CheckPasswordHash(req.CurrentPassword, user.Password) {
		return apperrors.NewUnprocessableEntityError("Validation error", apperrors.Field("currentPassword", "Current password is incorrect"))
	}

	hashed, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}
	user.Password = hashed
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}
	return s.cache.RevokeSessions(ctx, user.ID)
}
