package users

import "time"

type CreateUserRequest struct {
	Name     string   `json:"name" validate:"required,min=1,max=255"`
	Email    string   `json:"email" validate:"required,email,max=255"`
	Password string   `json:"password" validate:"required,min=8,strongpassword"`
	Status   string   `json:"status" validate:"required,oneof=active inactive suspended blocked"`
	Remark   *string  `json:"remark" validate:"omitempty,max=255"`
	RoleIDs  []string `json:"roleIds" validate:"required,dive,uuid"`
}

type UpdateUserRequest struct {
	Name    string   `json:"name" validate:"required,min=1,max=255"`
	Email   string   `json:"email" validate:"required,email,max=255"`
	Status  string   `json:"status" validate:"required,oneof=active inactive suspended blocked"`
	Remark  *string  `json:"remark" validate:"omitempty,max=255"`
	RoleIDs []string `json:"roleIds" validate:"required,dive,uuid"`
}

type ResetPasswordRequest struct {
	NewPassword string `json:"newPassword" validate:"required,min=8,strongpassword"`
}

type RoleOption struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UserList struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Status    string    `json:"status"`
	Remark    *string   `json:"remark"`
	Roles     []string  `json:"roles"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UserDetail struct {
	ID              string       `json:"id"`
	Name            string       `json:"name"`
	Email           string       `json:"email"`
	Status          string       `json:"status"`
	Remark          *string      `json:"remark"`
	EmailVerifiedAt *time.Time   `json:"emailVerifiedAt"`
	Roles           []RoleOption `json:"roles"`
	CreatedAt       time.Time    `json:"createdAt"`
	UpdatedAt       time.Time    `json:"updatedAt"`
}
