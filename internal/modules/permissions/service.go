package permissions

import (
	"context"
	"errors"
	"slices"
	"strings"

	"github.com/aolus-software/clean-go-fiber/internal/config"
	apperrors "github.com/aolus-software/clean-go-fiber/internal/errors"
	"github.com/aolus-software/clean-go-fiber/internal/libs/cache"
	"github.com/aolus-software/clean-go-fiber/internal/libs/pagination"
	"github.com/aolus-software/clean-go-fiber/internal/libs/utils"
	"github.com/aolus-software/clean-go-fiber/internal/models"
	"github.com/aolus-software/clean-go-fiber/internal/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	List(ctx context.Context, q *pagination.Query) (*utils.PaginatedData[PermissionList], error)
	Create(ctx context.Context, req *CreatePermissionRequest) error
	Detail(ctx context.Context, id uuid.UUID) (*PermissionList, error)
	Update(ctx context.Context, id uuid.UUID, req *UpdatePermissionRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type service struct {
	permissionRepo repositories.PermissionRepository
	userRepo       repositories.UserRepository
	cache          *cache.Cache
	cfg            *config.Config
}

func NewService(permissionRepo repositories.PermissionRepository, userRepo repositories.UserRepository, c *cache.Cache, cfg *config.Config) Service {
	return &service{permissionRepo: permissionRepo, userRepo: userRepo, cache: c, cfg: cfg}
}

func notFound() error {
	return apperrors.NewNotFoundError("Permission not found")
}

func namesTaken(names []string) error {
	return apperrors.NewUnprocessableEntityError("Permission already exists",
		apperrors.Field("name", "Permission already exists: "+strings.Join(names, ", ")))
}

func toList(p models.Permission) PermissionList {
	return PermissionList{ID: p.ID.String(), Name: p.Name, Group: p.Group, CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt}
}

func (s *service) List(ctx context.Context, q *pagination.Query) (*utils.PaginatedData[PermissionList], error) {
	perms, total, err := s.permissionRepo.FindAll(ctx, q, s.cfg.App.Timezone)
	if err != nil {
		return nil, err
	}
	items := make([]PermissionList, 0, len(perms))
	for _, p := range perms {
		items = append(items, toList(p))
	}
	return &utils.PaginatedData[PermissionList]{
		Data: items,
		Meta: utils.PaginationMeta{Page: q.Page, Limit: q.PerPage, TotalCount: total},
	}, nil
}

func (s *service) Create(ctx context.Context, req *CreatePermissionRequest) error {
	group := strings.TrimSpace(req.Group)
	names := make([]string, 0, len(req.Name))
	for _, action := range req.Name {
		full := group + " " + strings.TrimSpace(action)
		if len(full) > 255 {
			return apperrors.NewUnprocessableEntityError("Validation error", apperrors.Field("name", "group and action together must be at most 255 characters"))
		}
		if !slices.Contains(names, full) {
			names = append(names, full)
		}
	}
	existing, err := s.permissionRepo.FindExistingNames(ctx, names, nil)
	if err != nil {
		return err
	}
	if len(existing) > 0 {
		return namesTaken(existing)
	}

	rows := make([]models.Permission, 0, len(names))
	for _, n := range names {
		rows = append(rows, models.Permission{Name: n, Group: group})
	}
	err = s.permissionRepo.Create(ctx, rows)
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return namesTaken(names)
	}
	if err != nil {
		return err
	}
	return s.cache.ForgetSelectOptions(ctx)
}

func (s *service) Detail(ctx context.Context, id uuid.UUID) (*PermissionList, error) {
	p, err := s.permissionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, notFound()
	}
	result := toList(*p)
	return &result, nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, req *UpdatePermissionRequest) error {
	p, err := s.permissionRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if p == nil {
		return notFound()
	}
	name := strings.TrimSpace(req.Name)
	existing, err := s.permissionRepo.FindExistingNames(ctx, []string{name}, &p.ID)
	if err != nil {
		return err
	}
	if len(existing) > 0 {
		return namesTaken(existing)
	}

	p.Name = name
	p.Group = strings.TrimSpace(req.Group)
	if err := s.permissionRepo.Update(ctx, p); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return namesTaken([]string{name})
		}
		return err
	}
	if err := s.forgetHolders(ctx, p.ID); err != nil {
		return err
	}
	return s.cache.ForgetSelectOptions(ctx)
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	p, err := s.permissionRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if p == nil {
		return notFound()
	}
	holders, err := s.userRepo.FindIDsByPermissionID(ctx, p.ID)
	if err != nil {
		return err
	}
	if err := s.permissionRepo.Delete(ctx, p.ID); err != nil {
		return err
	}
	if err := s.cache.ForgetUsers(ctx, holders...); err != nil {
		return err
	}
	return s.cache.ForgetSelectOptions(ctx)
}

func (s *service) forgetHolders(ctx context.Context, permissionID uuid.UUID) error {
	holders, err := s.userRepo.FindIDsByPermissionID(ctx, permissionID)
	if err != nil {
		return err
	}
	return s.cache.ForgetUsers(ctx, holders...)
}
