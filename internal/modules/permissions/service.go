package permissions

import (
	"context"
	"time"

	"github.com/aolus-software/clean-go-fiber/internal/errors"
	"github.com/aolus-software/clean-go-fiber/internal/libs/cache"
	"github.com/aolus-software/clean-go-fiber/internal/libs/pagination"
	"github.com/aolus-software/clean-go-fiber/internal/models"
	"github.com/aolus-software/clean-go-fiber/internal/repositories"
	"github.com/google/uuid"
	rediscache "github.com/redis/go-redis/v9"
)

const permissionCacheKey = "permission:id:%s"

type Service interface {
	List(ctx context.Context, q *pagination.QueryParamDatatable) (*pagination.PaginationResponse[PermissionResponse], error)
	GetByID(ctx context.Context, id string) (*PermissionResponse, error)
	Create(ctx context.Context, req *PermissionRequest) (*PermissionResponse, error)
	Update(ctx context.Context, id string, req *PermissionRequest) (*PermissionResponse, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo  repositories.PermissionRepository
	cache *cache.Cache
}

func NewService(repo repositories.PermissionRepository, c *cache.Cache) Service {
	return &service{repo: repo, cache: c}
}

func (s *service) List(ctx context.Context, q *pagination.QueryParamDatatable) (*pagination.PaginationResponse[PermissionResponse], error) {
	result, err := s.repo.FindAll(ctx, q)
	if err != nil {
		return nil, errors.NewInternalServerError("failed to list permissions")
	}
	items := make([]PermissionResponse, 0, len(result.Data))
	for _, p := range result.Data {
		items = append(items, toResponse(p))
	}
	return &pagination.PaginationResponse[PermissionResponse]{
		Total:       result.Total,
		PerPage:     result.PerPage,
		CurrentPage: result.CurrentPage,
		LastPage:    result.LastPage,
		Data:        items,
	}, nil
}

func (s *service) GetByID(ctx context.Context, id string) (*PermissionResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.NewBadRequestError("invalid permission id")
	}

	key := cache.Key(permissionCacheKey, id)
	var resp PermissionResponse
	if err := s.cache.Get(ctx, key, &resp); err == nil {
		return &resp, nil
	} else if err != rediscache.Nil {
		return nil, errors.NewInternalServerError("cache error")
	}

	p, err := s.repo.FindById(ctx, uid)
	if err != nil {
		return nil, errors.NewNotFoundError("permission not found")
	}
	resp = toResponse(*p)
	s.cache.Set(ctx, key, resp, 5*time.Minute)
	return &resp, nil
}

func (s *service) Create(ctx context.Context, req *PermissionRequest) (*PermissionResponse, error) {
	p, err := s.repo.Create(ctx, &models.Permission{Name: req.Name, Modules: req.Modules})
	if err != nil {
		return nil, errors.NewInternalServerError("failed to create permission")
	}
	resp := toResponse(*p)
	return &resp, nil
}

func (s *service) Update(ctx context.Context, id string, req *PermissionRequest) (*PermissionResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.NewBadRequestError("invalid permission id")
	}
	p, err := s.repo.FindById(ctx, uid)
	if err != nil {
		return nil, errors.NewNotFoundError("permission not found")
	}
	p.Name = req.Name
	p.Modules = req.Modules
	updated, err := s.repo.Update(ctx, p)
	if err != nil {
		return nil, errors.NewInternalServerError("failed to update permission")
	}
	s.cache.Delete(ctx, cache.Key(permissionCacheKey, id))
	resp := toResponse(*updated)
	return &resp, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return errors.NewBadRequestError("invalid permission id")
	}
	if _, err := s.repo.FindById(ctx, uid); err != nil {
		return errors.NewNotFoundError("permission not found")
	}
	if err := s.repo.Delete(ctx, uid); err != nil {
		return err
	}
	s.cache.Delete(ctx, cache.Key(permissionCacheKey, id))
	return nil
}

func toResponse(p models.Permission) PermissionResponse {
	return PermissionResponse{
		ID:        p.ID.String(),
		Name:      p.Name,
		Modules:   p.Modules,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
