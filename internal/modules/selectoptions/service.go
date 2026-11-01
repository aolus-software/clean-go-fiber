package selectoptions

import (
	"context"

	"github.com/aolus-software/clean-go-fiber/internal/repositories"
)

type Service interface {
	Permissions(ctx context.Context) ([]PermissionGroupOption, error)
	Roles(ctx context.Context) ([]RoleOption, error)
}

type service struct {
	permissionRepo repositories.PermissionRepository
	roleRepo       repositories.RoleRepository
}

func NewService(permissionRepo repositories.PermissionRepository, roleRepo repositories.RoleRepository) Service {
	return &service{permissionRepo: permissionRepo, roleRepo: roleRepo}
}

func (s *service) Permissions(ctx context.Context) ([]PermissionGroupOption, error) {
	perms, err := s.permissionRepo.FindAllOrdered(ctx)
	if err != nil {
		return nil, err
	}
	groups := make([]PermissionGroupOption, 0)
	for _, p := range perms {
		if len(groups) == 0 || groups[len(groups)-1].Group != p.Group {
			groups = append(groups, PermissionGroupOption{Group: p.Group, Permissions: []PermissionOption{}})
		}
		last := &groups[len(groups)-1]
		last.Permissions = append(last.Permissions, PermissionOption{ID: p.ID.String(), Name: p.Name, Group: p.Group})
	}
	return groups, nil
}

func (s *service) Roles(ctx context.Context) ([]RoleOption, error) {
	roles, err := s.roleRepo.SelectOptions(ctx)
	if err != nil {
		return nil, err
	}
	options := make([]RoleOption, 0, len(roles))
	for _, r := range roles {
		options = append(options, RoleOption{ID: r.ID.String(), Name: r.Name})
	}
	return options, nil
}
