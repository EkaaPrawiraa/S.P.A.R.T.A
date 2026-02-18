package repository

import (
	"context"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/user"
)

type UserRepository interface {
	Create(ctx context.Context, u *user.User) error
	GetByID(ctx context.Context, id string) (*user.User, error)
	GetByEmail(ctx context.Context, email string) (*user.User, error)
}
