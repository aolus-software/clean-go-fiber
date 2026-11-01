package auth

import "github.com/aolus-software/clean-go-fiber/internal/repositories"

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=1"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Name     string `json:"name" validate:"required,min=1,max=255"`
	Password string `json:"password" validate:"required,min=8,strongpassword"`
}

type EmailRequest struct {
	Email string `json:"email" validate:"required,email,max=255"`
}

type TokenRequest struct {
	Token string `json:"token" validate:"required,min=1"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required,min=1"`
	NewPassword string `json:"newPassword" validate:"required,min=8,strongpassword"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required,min=1"`
}

type LoginResponse struct {
	User         *repositories.UserInformation `json:"user"`
	AccessToken  string                        `json:"accessToken"`
	RefreshToken string                        `json:"refreshToken"`
}

type RefreshResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
