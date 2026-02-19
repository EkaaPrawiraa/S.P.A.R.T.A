package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	domainerr "S.P.A.R.T.A/backend/internal/domain/errors"
	domainrepo "S.P.A.R.T.A/backend/internal/domain/repository"

	"github.com/redis/go-redis/v9"
)

type motivationRepository struct {
	client *redis.Client
}

func NewMotivationRepository(client *redis.Client) domainrepo.MotivationRepository {
	return &motivationRepository{client: client}
}

func (r *motivationRepository) GetDailyMotivation(ctx context.Context, userID string, date time.Time) (string, bool, error) {
	key := motivationKey(userID, date)

	val, err := r.client.Get(ctx, key).Result()
	if err == nil {
		return val, true, nil
	}
	if errors.Is(err, redis.Nil) {
		return "", false, nil
	}
	return "", false, domainerr.ErrInternal
}

func (r *motivationRepository) SetDailyMotivation(ctx context.Context, userID string, date time.Time, message string, ttl time.Duration) error {
	key := motivationKey(userID, date)

	if err := r.client.Set(ctx, key, message, ttl).Err(); err != nil {
		return domainerr.ErrInternal
	}
	return nil
}

func (r *motivationRepository) DeleteDailyMotivation(ctx context.Context, userID string, date time.Time) error {
	key := motivationKey(userID, date)
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return domainerr.ErrInternal
	}
	return nil
}

func motivationKey(userID string, date time.Time) string {
	return fmt.Sprintf("motivation:%s:%s", userID, date.UTC().Format("2006-01-02"))
}
