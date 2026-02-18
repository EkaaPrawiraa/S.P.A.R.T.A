package persistence

import (
	"S.P.A.R.T.A/backend/internal/domain/repository"
	"context"
	"database/sql"
)

type unitOfWork struct {
	db *sql.DB
}

func NewUnitOfWork(db *sql.DB) repository.UnitOfWork {
	return &unitOfWork{db: db}
}

func (u *unitOfWork) Do(ctx context.Context, fn func(r repository.Registry) error) error {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	registry := NewRegistry(tx)

	if err := fn(registry); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
