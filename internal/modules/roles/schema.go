package roles

import "time"

type RoleRequest struct {
	Name          string   `json:"name" validate:"required,min=1,max=100"`
	PermissionIDs []string `json:"permissionIds" validate:"required,dive,uuid"`
}

type RoleList struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type PermissionAssignment struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	IsAssigned bool   `json:"isAssigned"`
}

type PermissionGroup struct {
	Group string                 `json:"group"`
	Names []PermissionAssignment `json:"names"`
}

type RoleDetail struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
	Permissions []PermissionGroup `json:"permissions"`
}
