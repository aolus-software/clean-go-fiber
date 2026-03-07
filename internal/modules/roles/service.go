package roles

import (
	"context"

	"github.com/aolus-software/clean-go-fiber/internal/errors"
	"github.com/aolus-software/clean-go-fiber/internal/libs/pagination"
	"github.com/aolus-software/clean-go-fiber/internal/models"
	"github.com/aolus-software/clean-go-fiber/internal/repositories"
	"github.com/google/uuid"
)

type Service interface {
	List(ctx context.Context, q *pagination.QueryParamDatatable) (*pagination.PaginationResponse[RoleResponse], error)
	GetByID(ctx context.Context, id string) (*RoleResponse, error)
	Create(ctx context.Context, req *RoleRequest) (*RoleResponse, error)
	Update(ctx context.Context, id string, req *RoleRequest) (*RoleResponse, error)
	Delete(ctx context.Context, id string) error
	SyncPermissions(ctx context.Context, id string, req *SyncPermissionsRequest) (*RoleResponse, error)
}

type service struct {
	roleRepo       repositories.RoleRepository
	permissionRepo repositories.PermissionRepository
}

func NewService(roleRepo repositories.RoleRepository, permissionRepo repositories.PermissionRepository) Service {
	return &service{roleRepo: roleRepo, permissionRepo: permissionRepo}
}

func (s *service) List(ctx context.Context, q *pagination.QueryParamDatatable) (*pagination.PaginationResponse[RoleResponse], error) {
	result, err := s.roleRepo.FindAll(ctx, q)
	if err != nil {
		return nil, errors.NewInternalServerError("failed to list roles")
	}
	items := make([]RoleResponse, 0, len(result.Data))
	for _, r := range result.Data {
		items = append(items, toResponse(r))
	}
	return &pagination.PaginationResponse[RoleResponse]{
		Total:       result.Total,
		PerPage:     result.PerPage,
		CurrentPage: result.CurrentPage,
		LastPage:    result.LastPage,
		Data:        items,
	}, nil
}

func (s *service) GetByID(ctx context.Context, id string) (*RoleResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.NewBadRequestError("invalid role id")
	}
	role, err := s.roleRepo.FindById(ctx, uid)
	if err != nil {
		return nil, errors.NewNotFoundError("role not found")
	}
	resp := toResponse(*role)
	return &resp, nil
}

func (s *service) Create(ctx context.Context, req *RoleRequest) (*RoleResponse, error) {
	role, err := s.roleRepo.Create(ctx, &models.Role{Name: req.Name})
	if err != nil {
		return nil, errors.NewInternalServerError("failed to create role")
	}
	resp := toResponse(*role)
	return &resp, nil
}

func (s *service) Update(ctx context.Context, id string, req *RoleRequest) (*RoleResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.NewBadRequestError("invalid role id")
	}
	role, err := s.roleRepo.FindById(ctx, uid)
	if err != nil {
		return nil, errors.NewNotFoundError("role not found")
	}
	role.Name = req.Name
	updated, err := s.roleRepo.Update(ctx, role)
	if err != nil {
		return nil, errors.NewInternalServerError("failed to update role")
	}
	resp := toResponse(*updated)
	return &resp, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return errors.NewBadRequestError("invalid role id")
	}
	if _, err := s.roleRepo.FindById(ctx, uid); err != nil {
		return errors.NewNotFoundError("role not found")
	}
	return s.roleRepo.Delete(ctx, uid)
}

func (s *service) SyncPermissions(ctx context.Context, id string, req *SyncPermissionsRequest) (*RoleResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.NewBadRequestError("invalid role id")
	}
	role, err := s.roleRepo.FindById(ctx, uid)
	if err != nil {
		return nil, errors.NewNotFoundError("role not found")
	}
	permIDs := make([]uuid.UUID, 0, len(req.PermissionIDs))
	for _, pid := range req.PermissionIDs {
		puid, err := uuid.Parse(pid)
		if err != nil {
			return nil, errors.NewBadRequestError("invalid permission id: " + pid)
		}
		permIDs = append(permIDs, puid)
	}
	perms, err := s.permissionRepo.FindByIds(ctx, permIDs)
	if err != nil {
		return nil, errors.NewInternalServerError("failed to find permissions")
	}
	if err := s.roleRepo.SyncPermissions(ctx, role, perms); err != nil {
		return nil, errors.NewInternalServerError("failed to sync permissions")
	}
	role.Permissions = perms
	resp := toResponse(*role)
	return &resp, nil
}

func toResponse(r models.Role) RoleResponse {
	resp := RoleResponse{
		ID:        r.ID.String(),
		Name:      r.Name,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
	for _, p := range r.Permissions {
		resp.Permissions = append(resp.Permissions, PermissionItem{
			ID:        p.ID.String(),
			Name:      p.Name,
			Modules:   p.Modules,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		})
	}
	return resp
}
