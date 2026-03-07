package users

import (
	"context"
	"time"

	"github.com/aolus-software/clean-go-fiber/internal/errors"
	"github.com/aolus-software/clean-go-fiber/internal/libs/cache"
	"github.com/aolus-software/clean-go-fiber/internal/libs/pagination"
	"github.com/aolus-software/clean-go-fiber/internal/libs/utils"
	"github.com/aolus-software/clean-go-fiber/internal/models"
	"github.com/aolus-software/clean-go-fiber/internal/repositories"
	"github.com/google/uuid"
	rediscache "github.com/redis/go-redis/v9"
)

type Service interface {
	List(ctx context.Context, q *pagination.QueryParamDatatable) (*pagination.PaginationResponse[UserResponse], error)
	GetByID(ctx context.Context, id string) (*UserResponse, error)
	Create(ctx context.Context, req *CreateUserRequest) (*UserResponse, error)
	Update(ctx context.Context, id string, req *UpdateUserRequest) (*UserResponse, error)
	Delete(ctx context.Context, id string) error
	SyncRoles(ctx context.Context, id string, req *SyncRolesRequest) (*UserResponse, error)
}

type service struct {
	userRepo repositories.UserRepository
	roleRepo repositories.RoleRepository
	cache    *cache.Cache
}

func NewService(userRepo repositories.UserRepository, roleRepo repositories.RoleRepository, c *cache.Cache) Service {
	return &service{userRepo: userRepo, roleRepo: roleRepo, cache: c}
}

func (s *service) List(ctx context.Context, q *pagination.QueryParamDatatable) (*pagination.PaginationResponse[UserResponse], error) {
	result, err := s.userRepo.FindAll(ctx, q)
	if err != nil {
		return nil, errors.NewInternalServerError("failed to list users")
	}
	items := make([]UserResponse, 0, len(result.Data))
	for _, u := range result.Data {
		items = append(items, toResponse(u))
	}
	return &pagination.PaginationResponse[UserResponse]{
		Total:       result.Total,
		PerPage:     result.PerPage,
		CurrentPage: result.CurrentPage,
		LastPage:    result.LastPage,
		Data:        items,
	}, nil
}

func (s *service) GetByID(ctx context.Context, id string) (*UserResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.NewBadRequestError("invalid user id")
	}

	key := cache.Key(cache.UserByID, id)
	var resp UserResponse
	if err := s.cache.Get(ctx, key, &resp); err == nil {
		return &resp, nil
	} else if err != rediscache.Nil {
		return nil, errors.NewInternalServerError("cache error")
	}

	user, err := s.userRepo.FindById(ctx, uid)
	if err != nil {
		return nil, errors.NewNotFoundError("user not found")
	}
	resp = toResponse(*user)
	s.cache.Set(ctx, key, resp, 5*time.Minute)
	return &resp, nil
}

func (s *service) Create(ctx context.Context, req *CreateUserRequest) (*UserResponse, error) {
	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.NewInternalServerError("failed to process password")
	}
	user, err := s.userRepo.Create(ctx, &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashed,
		Status:   req.Status,
	})
	if err != nil {
		return nil, errors.NewInternalServerError("failed to create user")
	}
	resp := toResponse(*user)
	return &resp, nil
}

func (s *service) Update(ctx context.Context, id string, req *UpdateUserRequest) (*UserResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.NewBadRequestError("invalid user id")
	}
	user, err := s.userRepo.FindById(ctx, uid)
	if err != nil {
		return nil, errors.NewNotFoundError("user not found")
	}
	user.Name = req.Name
	user.Email = req.Email
	user.Status = req.Status
	updated, err := s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, errors.NewInternalServerError("failed to update user")
	}
	s.cache.Delete(ctx, cache.Key(cache.UserByID, id))
	resp := toResponse(*updated)
	return &resp, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return errors.NewBadRequestError("invalid user id")
	}
	if _, err := s.userRepo.FindById(ctx, uid); err != nil {
		return errors.NewNotFoundError("user not found")
	}
	if err := s.userRepo.Delete(ctx, uid); err != nil {
		return err
	}
	s.cache.Delete(ctx, cache.Key(cache.UserByID, id))
	return nil
}

func (s *service) SyncRoles(ctx context.Context, id string, req *SyncRolesRequest) (*UserResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.NewBadRequestError("invalid user id")
	}
	user, err := s.userRepo.FindById(ctx, uid)
	if err != nil {
		return nil, errors.NewNotFoundError("user not found")
	}
	roleIDs := make([]uuid.UUID, 0, len(req.RoleIDs))
	for _, rid := range req.RoleIDs {
		ruid, err := uuid.Parse(rid)
		if err != nil {
			return nil, errors.NewBadRequestError("invalid role id: " + rid)
		}
		roleIDs = append(roleIDs, ruid)
	}
	roles := make([]models.Role, 0, len(roleIDs))
	for _, ruid := range roleIDs {
		role, err := s.roleRepo.FindById(ctx, ruid)
		if err != nil {
			return nil, errors.NewNotFoundError("role not found: " + ruid.String())
		}
		roles = append(roles, *role)
	}
	if err := s.userRepo.SyncRoles(ctx, user, roles); err != nil {
		return nil, errors.NewInternalServerError("failed to sync roles")
	}
	s.cache.Delete(ctx, cache.Key(cache.UserByID, id))
	user.Roles = roles
	resp := toResponse(*user)
	return &resp, nil
}

func toResponse(u models.User) UserResponse {
	resp := UserResponse{
		ID:              u.ID.String(),
		Name:            u.Name,
		Email:           u.Email,
		Status:          u.Status,
		EmailVerifiedAt: u.EmailVerifiedAt,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
	}
	for _, r := range u.Roles {
		resp.Roles = append(resp.Roles, RoleItem{ID: r.ID.String(), Name: r.Name})
	}
	return resp
}
