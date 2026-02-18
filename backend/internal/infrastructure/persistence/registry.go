package persistence

import (
	"database/sql"

	"S.P.A.R.T.A/backend/internal/domain/repository"
)

type registry struct {
	tx *sql.Tx
}

func NewRegistry(tx *sql.Tx) repository.Registry {
	return &registry{tx: tx}
}

func (r *registry) User() repository.UserRepository {
	return NewUserRepository(r.tx)
}

func (r *registry) Workout() repository.WorkoutRepository {
	return NewWorkoutRepository(r.tx)
}

func (r *registry) Split() repository.SplitRepository {
	return NewSplitRepository(r.tx)
}

func (r *registry) Exercise() repository.ExerciseRepository {
	return NewExerciseRepository(r.tx)
}
