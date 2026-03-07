package permissions

import "time"

type PermissionRequest struct {
	Name    string `json:"name" validate:"required,min=2,max=255"`
	Modules string `json:"modules" validate:"required,min=2,max=255"`
}

type PermissionResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Modules   string    `json:"modules"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
