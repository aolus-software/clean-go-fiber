package settings

import (
	"context"

	"github.com/aolus-software/clean-go-fiber/internal/errors"
	"github.com/aolus-software/clean-go-fiber/internal/libs/pagination"
	"github.com/aolus-software/clean-go-fiber/internal/models"
	"github.com/aolus-software/clean-go-fiber/internal/repositories"
	"github.com/google/uuid"
)

type Service interface {
	ListRoles(ctx context.Context, q *pagination.QueryParamDatatable) (*pagination.PaginationResponse[RoleResponse], error)
	GetRole(ctx context.Context, id string) (*RoleResponse, error)
	CreateRole(ctx context.Context, req *RoleRequest) (*RoleResponse, error)
	UpdateRole(ctx context.Context, id string, req *RoleRequest) (*RoleResponse, error)
	DeleteRole(ctx context.Context, id string) error
	SyncRolePermissions(ctx context.Context, id string, req *SyncPermissionsRequest) (*RoleResponse, error)

	ListPermissions(ctx context.Context, q *pagination.QueryParamDatatable) (*pagination.PaginationResponse[PermissionResponse], error)
	GetPermission(ctx context.Context, id string) (*PermissionResponse, error)
	CreatePermission(ctx context.Context, req *PermissionRequest) (*PermissionResponse, error)
	UpdatePermission(ctx context.Context, id string, req *PermissionRequest) (*PermissionResponse, error)
	DeletePermission(ctx context.Context, id string) error

	SyncUserRoles(ctx context.Context, userID string, req *SyncRolesRequest) error
}

type service struct {
	roleRepo       repositories.RoleRepository
	permissionRepo repositories.PermissionRepository
	userRepo       repositories.UserRepository
}

func NewService(
	roleRepo repositories.RoleRepository,
	permissionRepo repositories.PermissionRepository,
	userRepo repositories.UserRepository,
) Service {
	return &service{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		userRepo:       userRepo,
	}
}

func (s *service) ListRoles(ctx context.Context, q *pagination.QueryParamDatatable) (*pagination.PaginationResponse[RoleResponse], error) {
	result, err := s.roleRepo.FindAll(ctx, q)
	if err != nil {
		return nil, errors.NewInternalServerError("failed to list roles")
	}
	items := make([]RoleResponse, 0, len(result.Data))
	for _, r := range result.Data {
		items = append(items, toRoleResponse(r))
	}
	return &pagination.PaginationResponse[RoleResponse]{
		Total:       result.Total,
		PerPage:     result.PerPage,
		CurrentPage: result.CurrentPage,
		LastPage:    result.LastPage,
		Data:        items,
	}, nil
}

func (s *service) GetRole(ctx context.Context, id string) (*RoleResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.NewBadRequestError("invalid role id")
	}
	role, err := s.roleRepo.FindById(ctx, uid)
	if err != nil {
		return nil, errors.NewNotFoundError("role not found")
	}
	resp := toRoleResponse(*role)
	return &resp, nil
}

func (s *service) CreateRole(ctx context.Context, req *RoleRequest) (*RoleResponse, error) {
	role, err := s.roleRepo.Create(ctx, &models.Role{Name: req.Name})
	if err != nil {
		return nil, errors.NewInternalServerError("failed to create role")
	}
	resp := toRoleResponse(*role)
	return &resp, nil
}

func (s *service) UpdateRole(ctx context.Context, id string, req *RoleRequest) (*RoleResponse, error) {
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
	resp := toRoleResponse(*updated)
	return &resp, nil
}

func (s *service) DeleteRole(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return errors.NewBadRequestError("invalid role id")
	}
	if _, err := s.roleRepo.FindById(ctx, uid); err != nil {
		return errors.NewNotFoundError("role not found")
	}
	return s.roleRepo.Delete(ctx, uid)
}

func (s *service) SyncRolePermissions(ctx context.Context, id string, req *SyncPermissionsRequest) (*RoleResponse, error) {
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
	resp := toRoleResponse(*role)
	return &resp, nil
}

func (s *service) ListPermissions(ctx context.Context, q *pagination.QueryParamDatatable) (*pagination.PaginationResponse[PermissionResponse], error) {
	result, err := s.permissionRepo.FindAll(ctx, q)
	if err != nil {
		return nil, errors.NewInternalServerError("failed to list permissions")
	}
	items := make([]PermissionResponse, 0, len(result.Data))
	for _, p := range result.Data {
		items = append(items, toPermissionResponse(p))
	}
	return &pagination.PaginationResponse[PermissionResponse]{
		Total:       result.Total,
		PerPage:     result.PerPage,
		CurrentPage: result.CurrentPage,
		LastPage:    result.LastPage,
		Data:        items,
	}, nil
}

func (s *service) GetPermission(ctx context.Context, id string) (*PermissionResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.NewBadRequestError("invalid permission id")
	}
	perm, err := s.permissionRepo.FindById(ctx, uid)
	if err != nil {
		return nil, errors.NewNotFoundError("permission not found")
	}
	resp := toPermissionResponse(*perm)
	return &resp, nil
}

func (s *service) CreatePermission(ctx context.Context, req *PermissionRequest) (*PermissionResponse, error) {
	perm, err := s.permissionRepo.Create(ctx, &models.Permission{Name: req.Name, Modules: req.Modules})
	if err != nil {
		return nil, errors.NewInternalServerError("failed to create permission")
	}
	resp := toPermissionResponse(*perm)
	return &resp, nil
}

func (s *service) UpdatePermission(ctx context.Context, id string, req *PermissionRequest) (*PermissionResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.NewBadRequestError("invalid permission id")
	}
	perm, err := s.permissionRepo.FindById(ctx, uid)
	if err != nil {
		return nil, errors.NewNotFoundError("permission not found")
	}
	perm.Name = req.Name
	perm.Modules = req.Modules
	updated, err := s.permissionRepo.Update(ctx, perm)
	if err != nil {
		return nil, errors.NewInternalServerError("failed to update permission")
	}
	resp := toPermissionResponse(*updated)
	return &resp, nil
}

func (s *service) DeletePermission(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return errors.NewBadRequestError("invalid permission id")
	}
	if _, err := s.permissionRepo.FindById(ctx, uid); err != nil {
		return errors.NewNotFoundError("permission not found")
	}
	return s.permissionRepo.Delete(ctx, uid)
}

func (s *service) SyncUserRoles(ctx context.Context, userID string, req *SyncRolesRequest) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return errors.NewBadRequestError("invalid user id")
	}
	user, err := s.userRepo.FindById(ctx, uid)
	if err != nil {
		return errors.NewNotFoundError("user not found")
	}
	roleIDs := make([]uuid.UUID, 0, len(req.RoleIDs))
	for _, rid := range req.RoleIDs {
		ruid, err := uuid.Parse(rid)
		if err != nil {
			return errors.NewBadRequestError("invalid role id: " + rid)
		}
		roleIDs = append(roleIDs, ruid)
	}
	roles := make([]models.Role, 0, len(roleIDs))
	for _, ruid := range roleIDs {
		role, err := s.roleRepo.FindById(ctx, ruid)
		if err != nil {
			return errors.NewNotFoundError("role not found: " + ruid.String())
		}
		roles = append(roles, *role)
	}
	return s.userRepo.SyncRoles(ctx, user, roles)
}

func toRoleResponse(r models.Role) RoleResponse {
	resp := RoleResponse{
		ID:        r.ID.String(),
		Name:      r.Name,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
	for _, p := range r.Permissions {
		resp.Permissions = append(resp.Permissions, toPermissionResponse(p))
	}
	return resp
}

func toPermissionResponse(p models.Permission) PermissionResponse {
	return PermissionResponse{
		ID:        p.ID.String(),
		Name:      p.Name,
		Modules:   p.Modules,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
