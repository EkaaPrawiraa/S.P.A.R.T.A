package usecase

import (
	"context"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/exercise"
)

type ExerciseUsecase interface {
	CreateExercise(ctx context.Context, ex *exercise.Exercise) error
	GetExercise(ctx context.Context, id string) (*exercise.Exercise, error)
	ListExercises(ctx context.Context) ([]exercise.Exercise, error)
	AddExerciseMedia(ctx context.Context, media *exercise.ExerciseMedia) error
}
