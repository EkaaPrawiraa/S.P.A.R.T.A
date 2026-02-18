package postgres

import (
	"context"
	"database/sql"
	"time"

	domainerr "S.P.A.R.T.A/backend/internal/domain/errors"
	domainrepo "S.P.A.R.T.A/backend/internal/domain/repository"
)

type adminInviteRepository struct {
	db DBTX
}

func NewAdminInviteRepository(db DBTX) domainrepo.AdminInviteRepository {
	return &adminInviteRepository{db: db}
}

func (r *adminInviteRepository) Create(ctx context.Context, invite *domainrepo.AdminInvite) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO admin_invites(id, token_hash, role, expires_at, used_at, used_by, created_by, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		invite.ID, invite.TokenHash, invite.Role, invite.ExpiresAt, invite.UsedAt, invite.UsedBy, invite.CreatedBy, invite.CreatedAt,
	)
	if err != nil {
		return domainerr.ErrInternal
	}
	return nil
}

func (r *adminInviteRepository) ReserveByTokenHash(ctx context.Context, tokenHash string, reservedAt time.Time) (*domainrepo.AdminInvite, error) {
	row := r.db.QueryRowContext(ctx,
		`UPDATE admin_invites
		 SET used_at=$3
		 WHERE token_hash=$1 AND used_at IS NULL AND expires_at > $2
		 RETURNING id, token_hash, role, expires_at, used_at, used_by, created_by, created_at`,
		tokenHash, reservedAt, reservedAt,
	)

	var out domainrepo.AdminInvite
	var usedAtOut sql.NullTime
	var usedByOut sql.NullString

	if err := row.Scan(
		&out.ID,
		&out.TokenHash,
		&out.Role,
		&out.ExpiresAt,
		&usedAtOut,
		&usedByOut,
		&out.CreatedBy,
		&out.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, domainerr.ErrInvalidInput
		}
		return nil, domainerr.ErrInternal
	}

	if usedAtOut.Valid {
		t := usedAtOut.Time
		out.UsedAt = &t
	}
	if usedByOut.Valid {
		s := usedByOut.String
		out.UsedBy = &s
	}

	return &out, nil
}

func (r *adminInviteRepository) AttachUsedBy(ctx context.Context, tokenHash string, usedBy string) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE admin_invites SET used_by=$2 WHERE token_hash=$1`,
		tokenHash, usedBy,
	)
	if err != nil {
		return domainerr.ErrInternal
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return domainerr.ErrInternal
	}
	if rows == 0 {
		return domainerr.ErrInvalidInput
	}
	return nil
}
