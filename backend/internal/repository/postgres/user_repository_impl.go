package postgres

import (
	"context"
	"database/sql"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/user"
	domainerr "S.P.A.R.T.A/backend/internal/domain/errors"
	domainrepo "S.P.A.R.T.A/backend/internal/domain/repository"
	"github.com/lib/pq"
)

type userRepository struct {
	db DBTX
}

func NewUserRepository(db DBTX) domainrepo.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, u *user.User) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO users(id,name,email,password_hash,role,created_at,updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		u.ID, u.Name, u.Email, u.PasswordHash, u.Role, u.CreatedAt, u.UpdatedAt,
	)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			// 23505 = unique_violation (e.g. users.email)
			if string(pqErr.Code) == "23505" {
				return domainerr.ErrConflict
			}
		}
		return domainerr.ErrInternal
	}
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*user.User, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id,name,email,password_hash,role,created_at,updated_at FROM users WHERE id=$1`,
		id,
	)

	var out user.User
	if err := row.Scan(&out.ID, &out.Name, &out.Email, &out.PasswordHash, &out.Role, &out.CreatedAt, &out.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, domainerr.ErrNotFound
		}
		return nil, domainerr.ErrInternal
	}
	return &out, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id,name,email,password_hash,role,created_at,updated_at FROM users WHERE email=$1`,
		email,
	)

	var out user.User
	if err := row.Scan(&out.ID, &out.Name, &out.Email, &out.PasswordHash, &out.Role, &out.CreatedAt, &out.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, domainerr.ErrNotFound
		}
		return nil, domainerr.ErrInternal
	}
	return &out, nil
}

func (r *userRepository) Count(ctx context.Context) (int, error) {
	row := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users`)
	var n int
	if err := row.Scan(&n); err != nil {
		return 0, domainerr.ErrInternal
	}
	return n, nil
}
