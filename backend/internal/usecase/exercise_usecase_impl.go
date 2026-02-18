package usecase

import (
	"context"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/exercise"
	domainrepo "S.P.A.R.T.A/backend/internal/domain/repository"
	domainuc "S.P.A.R.T.A/backend/internal/domain/usecase"
)

type exerciseUsecase struct {
	repo domainrepo.ExerciseRepository
}

func NewExerciseUsecase(repo domainrepo.ExerciseRepository) domainuc.ExerciseUsecase {
	return &exerciseUsecase{repo: repo}
}

func (u *exerciseUsecase) CreateExercise(ctx context.Context, ex *exercise.Exercise) error {
	return u.repo.Create(ctx, ex)
}

func (u *exerciseUsecase) GetExercise(ctx context.Context, id string) (*exercise.Exercise, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *exerciseUsecase) ListExercises(ctx context.Context) ([]exercise.Exercise, error) {
	return u.repo.List(ctx)
}

func (u *exerciseUsecase) AddExerciseMedia(ctx context.Context, media *exercise.ExerciseMedia) error {
	return u.repo.AddMedia(ctx, media)
}
