package repository

import "context"

type UnitOfWork interface {
	Do(ctx context.Context, fn func(r Registry) error) error
}

type Registry interface {
	User() UserRepository
	Workout() WorkoutRepository
	Split() SplitRepository
	Exercise() ExerciseRepository
}
