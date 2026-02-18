package repository

import (
	"context"
	"time"
)

type AdminInvite struct {
	ID        string
	TokenHash string
	Role      string
	ExpiresAt time.Time
	UsedAt    *time.Time
	UsedBy    *string
	CreatedBy string
	CreatedAt time.Time
}

type AdminInviteRepository interface {
	Create(ctx context.Context, invite *AdminInvite) error
	ReserveByTokenHash(ctx context.Context, tokenHash string, reservedAt time.Time) (*AdminInvite, error)
	AttachUsedBy(ctx context.Context, tokenHash string, usedBy string) error
}
