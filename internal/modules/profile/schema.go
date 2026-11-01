package profile

type UpdateProfileRequest struct {
	Name  string `json:"name" validate:"required,min=1,max=255"`
	Email string `json:"email" validate:"required,email,max=255"`
}

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" validate:"required,min=1"`
	NewPassword     string `json:"newPassword" validate:"required,min=8,strongpassword"`
}
