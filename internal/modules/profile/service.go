package profile

import (
	"context"

	"github.com/aolus-software/clean-go-fiber/internal/errors"
	"github.com/aolus-software/clean-go-fiber/internal/libs/utils"
	"github.com/aolus-software/clean-go-fiber/internal/repositories"
	"github.com/google/uuid"
)

type Service interface {
	GetProfile(ctx context.Context, userID string) (*ProfileResponse, error)
	UpdateProfile(ctx context.Context, userID string, req *UpdateProfileRequest) (*ProfileResponse, error)
	ChangePassword(ctx context.Context, userID string, req *ChangePasswordRequest) error
}

type service struct {
	userRepo repositories.UserRepository
}

func NewService(userRepo repositories.UserRepository) Service {
	return &service{userRepo: userRepo}
}

func (s *service) GetProfile(ctx context.Context, userID string) (*ProfileResponse, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.NewBadRequestError("invalid user id")
	}
	user, err := s.userRepo.FindById(ctx, id)
	if err != nil {
		return nil, errors.NewNotFoundError("user not found")
	}
	resp := &ProfileResponse{
		ID:              user.ID.String(),
		Name:            user.Name,
		Email:           user.Email,
		Status:          user.Status,
		EmailVerifiedAt: user.EmailVerifiedAt,
	}
	for _, r := range user.Roles {
		resp.Roles = append(resp.Roles, RoleItem{ID: r.ID.String(), Name: r.Name})
	}
	return resp, nil
}

func (s *service) UpdateProfile(ctx context.Context, userID string, req *UpdateProfileRequest) (*ProfileResponse, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.NewBadRequestError("invalid user id")
	}
	user, err := s.userRepo.FindById(ctx, id)
	if err != nil {
		return nil, errors.NewNotFoundError("user not found")
	}
	user.Name = req.Name
	updated, err := s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, errors.NewInternalServerError("failed to update profile")
	}
	resp := &ProfileResponse{
		ID:              updated.ID.String(),
		Name:            updated.Name,
		Email:           updated.Email,
		Status:          updated.Status,
		EmailVerifiedAt: updated.EmailVerifiedAt,
	}
	for _, r := range updated.Roles {
		resp.Roles = append(resp.Roles, RoleItem{ID: r.ID.String(), Name: r.Name})
	}
	return resp, nil
}

func (s *service) ChangePassword(ctx context.Context, userID string, req *ChangePasswordRequest) error {
	id, err := uuid.Parse(userID)
	if err != nil {
		return errors.NewBadRequestError("invalid user id")
	}
	user, err := s.userRepo.FindById(ctx, id)
	if err != nil {
		return errors.NewNotFoundError("user not found")
	}
	if !utils.CheckPasswordHash(req.CurrentPassword, user.Password) {
		return errors.NewBadRequestError("current password is incorrect")
	}
	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		return errors.NewInternalServerError("failed to process password")
	}
	user.Password = hashed
	if _, err := s.userRepo.Update(ctx, user); err != nil {
		return errors.NewInternalServerError("failed to update password")
	}
	return nil
}
