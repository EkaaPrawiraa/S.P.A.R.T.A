package usecase

import (
	"context"
	"time"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/exercise"
	domainrepo "S.P.A.R.T.A/backend/internal/domain/repository"
	domainuc "S.P.A.R.T.A/backend/internal/domain/usecase"
)

type exerciseUsecase struct {
	repo      domainrepo.ExerciseRepository
	cacheRepo domainrepo.ExerciseCacheRepository
}

func NewExerciseUsecase(repo domainrepo.ExerciseRepository, cacheRepo domainrepo.ExerciseCacheRepository) domainuc.ExerciseUsecase {
	return &exerciseUsecase{repo: repo, cacheRepo: cacheRepo}
}

func (u *exerciseUsecase) CreateExercise(ctx context.Context, ex *exercise.Exercise) error {
	if err := u.repo.Create(ctx, ex); err != nil {
		return err
	}
	if u.cacheRepo != nil {
		_ = u.cacheRepo.DeleteExerciseList(ctx)
	}
	return nil
}

func (u *exerciseUsecase) GetExercise(ctx context.Context, id string) (*exercise.Exercise, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *exerciseUsecase) ListExercises(ctx context.Context) ([]exercise.Exercise, error) {
	if u.cacheRepo != nil {
		if cached, ok, err := u.cacheRepo.GetExerciseList(ctx); err == nil && ok {
			return cached, nil
		}
		// If Redis is down/unreachable, treat it as a cache miss.
	}

	items, err := u.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	if u.cacheRepo != nil {
		// Cache is best-effort; never fail the request due to Redis.
		_ = u.cacheRepo.SetExerciseList(ctx, items, 5*time.Minute)
	}

	return items, nil
}

func (u *exerciseUsecase) AddExerciseMedia(ctx context.Context, media *exercise.ExerciseMedia) error {
	return u.repo.AddMedia(ctx, media)
}
