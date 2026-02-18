package repository

import (
	"context"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/exercise"
)

type ExerciseRepository interface {
	Create(ctx context.Context, ex *exercise.Exercise) error
	GetByID(ctx context.Context, id string) (*exercise.Exercise, error)
	List(ctx context.Context) ([]exercise.Exercise, error)
	AddMedia(ctx context.Context, media *exercise.ExerciseMedia) error
}
