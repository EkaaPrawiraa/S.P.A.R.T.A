package usecase

import (
	"context"
	"time"
)

type AdminInviteResult struct {
	InviteToken string    `json:"invite_token"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type AdminUsecase interface {
	CreateAdminInvite(ctx context.Context, createdBy string, expiresIn time.Duration) (*AdminInviteResult, error)
}
