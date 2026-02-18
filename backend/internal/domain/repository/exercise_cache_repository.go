package repository

import (
	"context"
	"time"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/exercise"
)

type ExerciseCacheRepository interface {
	GetExerciseList(ctx context.Context) ([]exercise.Exercise, bool, error)
	SetExerciseList(ctx context.Context, items []exercise.Exercise, ttl time.Duration) error
	DeleteExerciseList(ctx context.Context) error
}
