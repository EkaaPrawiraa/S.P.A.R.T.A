package redis

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/exercise"
	domainerr "S.P.A.R.T.A/backend/internal/domain/errors"
	domainrepo "S.P.A.R.T.A/backend/internal/domain/repository"

	"github.com/redis/go-redis/v9"
)

type exerciseCacheRepository struct {
	client *redis.Client
}

func NewExerciseCacheRepository(client *redis.Client) domainrepo.ExerciseCacheRepository {
	return &exerciseCacheRepository{client: client}
}

func (r *exerciseCacheRepository) GetExerciseList(ctx context.Context) ([]exercise.Exercise, bool, error) {
	val, err := r.client.Get(ctx, exerciseListKey()).Result()
	if err == nil {
		var items []exercise.Exercise
		if err := json.Unmarshal([]byte(val), &items); err != nil {
			_ = r.client.Del(ctx, exerciseListKey()).Err()
			return nil, false, nil
		}
		return items, true, nil
	}
	if errors.Is(err, redis.Nil) {
		return nil, false, nil
	}
	return nil, false, domainerr.ErrInternal
}

func (r *exerciseCacheRepository) SetExerciseList(ctx context.Context, items []exercise.Exercise, ttl time.Duration) error {
	b, err := json.Marshal(items)
	if err != nil {
		return domainerr.ErrInternal
	}
	if err := r.client.Set(ctx, exerciseListKey(), string(b), ttl).Err(); err != nil {
		return domainerr.ErrInternal
	}
	return nil
}

func (r *exerciseCacheRepository) DeleteExerciseList(ctx context.Context) error {
	if err := r.client.Del(ctx, exerciseListKey()).Err(); err != nil {
		return domainerr.ErrInternal
	}
	return nil
}

func exerciseListKey() string {
	return "exercises:list"
}
