package users

import "time"

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=255"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Status   string `json:"status" validate:"required,oneof=active inactive"`
}

type UpdateUserRequest struct {
	Name   string `json:"name" validate:"required,min=2,max=255"`
	Email  string `json:"email" validate:"required,email"`
	Status string `json:"status" validate:"required,oneof=active inactive"`
}

type SyncRolesRequest struct {
	RoleIDs []string `json:"role_ids" validate:"required,min=1"`
}

type RoleItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UserResponse struct {
	ID              string     `json:"id"`
	Name            string     `json:"name"`
	Email           string     `json:"email"`
	Status          string     `json:"status"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	Roles           []RoleItem `json:"roles,omitempty"`
}
