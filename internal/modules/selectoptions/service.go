package selectoptions

import (
	"context"

	"github.com/aolus-software/clean-go-fiber/internal/libs/cache"
	"github.com/aolus-software/clean-go-fiber/internal/repositories"
)

type Service interface {
	Permissions(ctx context.Context) ([]PermissionGroupOption, error)
	Roles(ctx context.Context) ([]RoleOption, error)
}

type service struct {
	permissionRepo repositories.PermissionRepository
	roleRepo       repositories.RoleRepository
	cache          *cache.Cache
}

func NewService(permissionRepo repositories.PermissionRepository, roleRepo repositories.RoleRepository, c *cache.Cache) Service {
	return &service{permissionRepo: permissionRepo, roleRepo: roleRepo, cache: c}
}

func (s *service) Permissions(ctx context.Context) ([]PermissionGroupOption, error) {
	var groups []PermissionGroupOption
	err := s.cache.Remember(ctx, cache.PermissionSelectOptionsKey, cache.SelectOptionsTTL, &groups, func() (any, error) {
		return s.loadPermissions(ctx)
	})
	return groups, err
}

func (s *service) Roles(ctx context.Context) ([]RoleOption, error) {
	var options []RoleOption
	err := s.cache.Remember(ctx, cache.RoleSelectOptionsKey, cache.SelectOptionsTTL, &options, func() (any, error) {
		return s.loadRoles(ctx)
	})
	return options, err
}

func (s *service) loadPermissions(ctx context.Context) ([]PermissionGroupOption, error) {
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

func (s *service) loadRoles(ctx context.Context) ([]RoleOption, error) {
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
