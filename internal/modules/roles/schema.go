package roles

import "time"

type RoleRequest struct {
	Name string `json:"name" validate:"required,min=2,max=255"`
}

type SyncPermissionsRequest struct {
	PermissionIDs []string `json:"permission_ids" validate:"required,min=1"`
}

type PermissionItem struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Modules   string    `json:"modules"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RoleResponse struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	Permissions []PermissionItem `json:"permissions,omitempty"`
}
