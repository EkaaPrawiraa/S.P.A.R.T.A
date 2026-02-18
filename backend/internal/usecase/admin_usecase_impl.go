package usecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"time"

	domainerr "S.P.A.R.T.A/backend/internal/domain/errors"
	domainrepo "S.P.A.R.T.A/backend/internal/domain/repository"
	domainuc "S.P.A.R.T.A/backend/internal/domain/usecase"
	"github.com/google/uuid"
)

type adminUsecase struct {
	inviteRepo domainrepo.AdminInviteRepository
}

func NewAdminUsecase(inviteRepo domainrepo.AdminInviteRepository) domainuc.AdminUsecase {
	return &adminUsecase{inviteRepo: inviteRepo}
}

func (u *adminUsecase) CreateAdminInvite(ctx context.Context, createdBy string, expiresIn time.Duration) (*domainuc.AdminInviteResult, error) {
	if u.inviteRepo == nil {
		return nil, domainerr.ErrInternal
	}
	if expiresIn <= 0 {
		expiresIn = 24 * time.Hour
	}
	if expiresIn > 14*24*time.Hour {
		expiresIn = 14 * 24 * time.Hour
	}

	raw, err := newRandomToken(32)
	if err != nil {
		return nil, domainerr.ErrInternal
	}

	h := sha256.Sum256([]byte(raw))
	hashHex := hex.EncodeToString(h[:])

	now := time.Now().UTC()
	invite := &domainrepo.AdminInvite{
		ID:        uuid.NewString(),
		TokenHash: hashHex,
		Role:      "admin",
		ExpiresAt: now.Add(expiresIn),
		CreatedBy: createdBy,
		CreatedAt: now,
	}

	if err := u.inviteRepo.Create(ctx, invite); err != nil {
		return nil, err
	}

	return &domainuc.AdminInviteResult{InviteToken: raw, ExpiresAt: invite.ExpiresAt}, nil
}

func newRandomToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
