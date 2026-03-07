package profile

import "time"

type UpdateProfileRequest struct {
	Name string `json:"name" validate:"required,min=2,max=255"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type ProfileResponse struct {
	ID              string     `json:"id"`
	Name            string     `json:"name"`
	Email           string     `json:"email"`
	Status          string     `json:"status"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	Roles           []RoleItem `json:"roles"`
}

type RoleItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
