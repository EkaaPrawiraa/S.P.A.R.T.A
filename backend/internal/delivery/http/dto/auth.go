package dto

type RegisterRequestDTO struct {
	Name        string `json:"name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
	InviteToken string `json:"invite_token"`
}

type LoginRequestDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type CreateAdminInviteRequestDTO struct {
	ExpiresInHours int `json:"expires_in_hours" validate:"omitempty,min=1,max=336"`
}
