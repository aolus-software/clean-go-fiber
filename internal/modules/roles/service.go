package roles

import (
	"context"
	"errors"
	"slices"
	"strings"

	"github.com/aolus-software/clean-go-fiber/internal/config"
	apperrors "github.com/aolus-software/clean-go-fiber/internal/errors"
	"github.com/aolus-software/clean-go-fiber/internal/libs/cache"
	"github.com/aolus-software/clean-go-fiber/internal/libs/pagination"
	"github.com/aolus-software/clean-go-fiber/internal/libs/txmanager"
	"github.com/aolus-software/clean-go-fiber/internal/libs/utils"
	"github.com/aolus-software/clean-go-fiber/internal/models"
	"github.com/aolus-software/clean-go-fiber/internal/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	List(ctx context.Context, q *pagination.Query) (*utils.PaginatedData[RoleList], error)
	Create(ctx context.Context, req *RoleRequest) error
	Detail(ctx context.Context, id uuid.UUID) (*RoleDetail, error)
	Update(ctx context.Context, id uuid.UUID, req *RoleRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type service struct {
	roleRepo       repositories.RoleRepository
	permissionRepo repositories.PermissionRepository
	userRepo       repositories.UserRepository
	tx             *txmanager.TxManager
	cache          *cache.Cache
	cfg            *config.Config
}

func NewService(
	roleRepo repositories.RoleRepository,
	permissionRepo repositories.PermissionRepository,
	userRepo repositories.UserRepository,
	tx *txmanager.TxManager,
	c *cache.Cache,
	cfg *config.Config,
) Service {
	return &service{roleRepo: roleRepo, permissionRepo: permissionRepo, userRepo: userRepo, tx: tx, cache: c, cfg: cfg}
}

func notFound() error {
	return apperrors.NewNotFoundError("Role not found")
}

func nameTaken(name string) error {
	return apperrors.NewUnprocessableEntityError("Role name already exists", apperrors.Field("name", "Role "+name+" already exists"))
}

func (s *service) List(ctx context.Context, q *pagination.Query) (*utils.PaginatedData[RoleList], error) {
	roles, total, err := s.roleRepo.FindAll(ctx, q, s.cfg.App.Timezone)
	if err != nil {
		return nil, err
	}
	items := make([]RoleList, 0, len(roles))
	for _, r := range roles {
		items = append(items, RoleList{ID: r.ID.String(), Name: r.Name, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt})
	}
	return &utils.PaginatedData[RoleList]{
		Data: items,
		Meta: utils.PaginationMeta{Page: q.Page, Limit: q.PerPage, TotalCount: total},
	}, nil
}

func (s *service) Create(ctx context.Context, req *RoleRequest) error {
	name := strings.TrimSpace(req.Name)
	existing, err := s.roleRepo.FindByName(ctx, name, nil)
	if err != nil {
		return err
	}
	if existing != nil {
		return nameTaken(name)
	}
	permissionIDs, err := s.resolvePermissionIDs(ctx, req.PermissionIDs)
	if err != nil {
		return err
	}

	role := &models.Role{Name: name}
	err = s.tx.WithinTransaction(ctx, func(tx *gorm.DB) error {
		repo := s.roleRepo.WithTx(tx)
		if err := repo.Create(ctx, role); err != nil {
			return err
		}
		return repo.ReplacePermissions(ctx, role.ID, permissionIDs)
	})
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return nameTaken(name)
	}
	return err
}

func (s *service) Detail(ctx context.Context, id uuid.UUID) (*RoleDetail, error) {
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, notFound()
	}
	assigned, err := s.roleRepo.PermissionIDs(ctx, role.ID)
	if err != nil {
		return nil, err
	}
	all, err := s.permissionRepo.FindAllOrdered(ctx)
	if err != nil {
		return nil, err
	}

	groups := make([]PermissionGroup, 0)
	for _, p := range all {
		if len(groups) == 0 || groups[len(groups)-1].Group != p.Group {
			groups = append(groups, PermissionGroup{Group: p.Group, Names: []PermissionAssignment{}})
		}
		last := &groups[len(groups)-1]
		last.Names = append(last.Names, PermissionAssignment{
			ID:         p.ID.String(),
			Name:       p.Name,
			IsAssigned: slices.Contains(assigned, p.ID),
		})
	}

	return &RoleDetail{
		ID:          role.ID.String(),
		Name:        role.Name,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
		Permissions: groups,
	}, nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, req *RoleRequest) error {
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if role == nil {
		return notFound()
	}
	name := strings.TrimSpace(req.Name)
	if role.Name == models.SuperuserRole && name != models.SuperuserRole {
		return apperrors.NewUnprocessableEntityError("Validation error", apperrors.Field("name", "The superuser role cannot be renamed"))
	}
	existing, err := s.roleRepo.FindByName(ctx, name, &role.ID)
	if err != nil {
		return err
	}
	if existing != nil {
		return nameTaken(name)
	}
	permissionIDs, err := s.resolvePermissionIDs(ctx, req.PermissionIDs)
	if err != nil {
		return err
	}

	role.Name = name
	err = s.tx.WithinTransaction(ctx, func(tx *gorm.DB) error {
		repo := s.roleRepo.WithTx(tx)
		if err := repo.Update(ctx, role); err != nil {
			return err
		}
		return repo.ReplacePermissions(ctx, role.ID, permissionIDs)
	})
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return nameTaken(name)
	}
	if err != nil {
		return err
	}
	return s.forgetHolders(ctx, role.ID)
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if role == nil {
		return notFound()
	}
	if role.Name == models.SuperuserRole {
		return apperrors.NewUnprocessableEntityError("Validation error", apperrors.Field("id", "The superuser role cannot be deleted"))
	}
	holders, err := s.userRepo.FindIDsByRoleIDs(ctx, []uuid.UUID{role.ID})
	if err != nil {
		return err
	}
	if err := s.roleRepo.Delete(ctx, role.ID); err != nil {
		return err
	}
	return s.cache.ForgetUsers(ctx, holders...)
}

func (s *service) forgetHolders(ctx context.Context, roleID uuid.UUID) error {
	holders, err := s.userRepo.FindIDsByRoleIDs(ctx, []uuid.UUID{roleID})
	if err != nil {
		return err
	}
	return s.cache.ForgetUsers(ctx, holders...)
}

func (s *service) resolvePermissionIDs(ctx context.Context, raw []string) ([]uuid.UUID, error) {
	ids := make([]uuid.UUID, 0, len(raw))
	for _, r := range raw {
		id, err := uuid.Parse(r)
		if err != nil {
			return nil, apperrors.NewUnprocessableEntityError("Validation error", apperrors.Field("permissionIds", "must contain valid uuids"))
		}
		if !slices.Contains(ids, id) {
			ids = append(ids, id)
		}
	}
	count, err := s.permissionRepo.CountByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	if int(count) != len(ids) {
		return nil, apperrors.NewUnprocessableEntityError("Validation error", apperrors.Field("permissionIds", "contains permissions that do not exist"))
	}
	return ids, nil
}
