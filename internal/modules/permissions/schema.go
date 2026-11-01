package permissions

import "time"

type CreatePermissionRequest struct {
	Name  []string `json:"name" validate:"required,min=1,dive,required,min=1,max=255"`
	Group string   `json:"group" validate:"required,min=1,max=100"`
}

type UpdatePermissionRequest struct {
	Name  string `json:"name" validate:"required,min=1,max=255"`
	Group string `json:"group" validate:"required,min=1,max=100"`
}

type PermissionList struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Group     string    `json:"group"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
