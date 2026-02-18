package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/user"
	domainerr "S.P.A.R.T.A/backend/internal/domain/errors"
	domainrepo "S.P.A.R.T.A/backend/internal/domain/repository"
	domainuc "S.P.A.R.T.A/backend/internal/domain/usecase"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const defaultTokenTTL = 7 * 24 * time.Hour

type authUsecase struct {
	userRepo   domainrepo.UserRepository
	inviteRepo domainrepo.AdminInviteRepository
	jwtSecret  string
}

func NewAuthUsecase(userRepo domainrepo.UserRepository, inviteRepo domainrepo.AdminInviteRepository, jwtSecret string) domainuc.AuthUsecase {
	return &authUsecase{userRepo: userRepo, inviteRepo: inviteRepo, jwtSecret: jwtSecret}
}

func (u *authUsecase) Register(ctx context.Context, name, email, password, inviteToken string) (*domainuc.AuthResult, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(strings.ToLower(email))
	password = strings.TrimSpace(password)
	inviteToken = strings.TrimSpace(inviteToken)

	if name == "" || email == "" || password == "" {
		return nil, domainerr.ErrInvalidInput
	}

	if existing, err := u.userRepo.GetByEmail(ctx, email); err == nil && existing != nil {
		return nil, domainerr.ErrConflict
	} else if err != nil && err != domainerr.ErrNotFound {
		return nil, err
	}

	role := "user"
	count, err := u.userRepo.Count(ctx)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		role = "admin" // one-time bootstrap admin
	}

	userID := uuid.NewString()

	// Invite-based admin: optional invite token during registration.
	var reservedInviteTokenHash string
	if inviteToken != "" {
		if u.inviteRepo == nil {
			return nil, domainerr.ErrInternal
		}
		h := sha256.Sum256([]byte(inviteToken))
		reservedInviteTokenHash = hex.EncodeToString(h[:])
		invite, err := u.inviteRepo.ReserveByTokenHash(ctx, reservedInviteTokenHash, time.Now().UTC())
		if err != nil {
			return nil, err
		}
		if invite != nil && strings.TrimSpace(invite.Role) != "" {
			role = invite.Role
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, domainerr.ErrInternal
	}

	now := time.Now().UTC()
	domainUser := &user.User{
		ID:           userID,
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		Role:         role,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := u.userRepo.Create(ctx, domainUser); err != nil {
		return nil, err
	}

	if reservedInviteTokenHash != "" {
		_ = u.inviteRepo.AttachUsedBy(ctx, reservedInviteTokenHash, domainUser.ID)
	}

	token, err := u.generateToken(domainUser.ID, domainUser.Role)
	if err != nil {
		return nil, err
	}

	return &domainuc.AuthResult{UserID: domainUser.ID, Token: token, Role: domainUser.Role}, nil
}

func (u *authUsecase) Login(ctx context.Context, email, password string) (*domainuc.AuthResult, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	password = strings.TrimSpace(password)
	if email == "" || password == "" {
		return nil, domainerr.ErrInvalidInput
	}

	usr, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if err == domainerr.ErrNotFound {
			return nil, domainerr.ErrUnauthorized
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(usr.PasswordHash), []byte(password)); err != nil {
		return nil, domainerr.ErrUnauthorized
	}

	token, err := u.generateToken(usr.ID, usr.Role)
	if err != nil {
		return nil, err
	}

	return &domainuc.AuthResult{UserID: usr.ID, Token: token, Role: usr.Role}, nil
}

func (u *authUsecase) generateToken(userID string, role string) (string, error) {
	now := time.Now().UTC()
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"iat":     now.Unix(),
		"exp":     now.Add(defaultTokenTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return "", domainerr.ErrInternal
	}
	return signed, nil
}
