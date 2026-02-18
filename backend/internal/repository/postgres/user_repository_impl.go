package postgres

import (
	"context"
	"database/sql"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/user"
	domainrepo "S.P.A.R.T.A/backend/internal/domain/repository"
)

type userRepository struct {
	db DBTX
}

func NewUserRepository(db DBTX) domainrepo.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, u *user.User) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO users(id,name,email,password_hash,created_at,updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6)`,
		u.ID, u.Name, u.Email, u.PasswordHash, u.CreatedAt, u.UpdatedAt,
	)
	return err
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*user.User, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id,name,email,password_hash,created_at,updated_at FROM users WHERE id=$1`,
		id,
	)

	var out user.User
	if err := row.Scan(&out.ID, &out.Name, &out.Email, &out.PasswordHash, &out.CreatedAt, &out.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	return &out, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id,name,email,password_hash,created_at,updated_at FROM users WHERE email=$1`,
		email,
	)

	var out user.User
	if err := row.Scan(&out.ID, &out.Name, &out.Email, &out.PasswordHash, &out.CreatedAt, &out.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	return &out, nil
}
