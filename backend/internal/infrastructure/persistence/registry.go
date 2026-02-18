package persistence

import (
	"database/sql"

	"S.P.A.R.T.A/backend/internal/domain/repository"
	postgresRepo "S.P.A.R.T.A/backend/internal/repository/postgres"
)

type registry struct {
	tx *sql.Tx
}

func NewRegistry(tx *sql.Tx) repository.Registry {
	return &registry{tx: tx}
}

func (r *registry) User() repository.UserRepository {
	return postgresRepo.NewUserRepository(r.tx)
}

func (r *registry) Workout() repository.WorkoutRepository {
	return postgresRepo.NewWorkoutRepository(r.tx)
}

func (r *registry) Split() repository.SplitRepository {
	return postgresRepo.NewSplitRepository(r.tx)
}

func (r *registry) Exercise() repository.ExerciseRepository {
	return postgresRepo.NewExerciseRepository(r.tx)
}
