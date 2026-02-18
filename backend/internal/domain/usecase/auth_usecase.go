package usecase

import "context"

type AuthResult struct {
	UserID string `json:"user_id"`
	Token  string `json:"token"`
	Role   string `json:"role"`
}

type AuthUsecase interface {
	Register(ctx context.Context, name, email, password, inviteToken string) (*AuthResult, error)
	Login(ctx context.Context, email, password string) (*AuthResult, error)
}
