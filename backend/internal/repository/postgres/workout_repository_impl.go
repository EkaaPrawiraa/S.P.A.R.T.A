package postgres

import (
	"context"
	"fmt"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/workout"
	domainrepo "S.P.A.R.T.A/backend/internal/domain/repository"
)

type workoutRepository struct {
	db DBTX
}

func NewWorkoutRepository(db DBTX) domainrepo.WorkoutRepository {
	return &workoutRepository{db: db}
}

func (r *workoutRepository) CreateSession(
	ctx context.Context,
	session *workout.WorkoutSession,
) error {

	sessionQuery := `
		INSERT INTO workout_sessions
		(id, user_id, split_day_id, session_date, duration_minutes, notes, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
	`

	_, err := r.db.ExecContext(
		ctx,
		sessionQuery,
		session.ID,
		session.UserID,
		session.SplitDayID,
		session.SessionDate,
		session.DurationMin,
		session.Notes,
		session.CreatedAt,
	)
	if err != nil {
		return err
	}

	for _, ex := range session.Exercises {

		exQuery := `
			INSERT INTO workout_exercises
			(id, workout_session_id, exercise_id)
			VALUES ($1,$2,$3)
		`

		_, err = r.db.ExecContext(
			ctx,
			exQuery,
			ex.ID,
			session.ID,
			ex.ExerciseID,
		)
		if err != nil {
			return err
		}

		for _, set := range ex.Sets {

			setQuery := `
				INSERT INTO workout_sets
				(id, workout_exercise_id, set_order, reps, weight, rpe, set_type, created_at)
				VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
			`

			_, err = r.db.ExecContext(
				ctx,
				setQuery,
				set.ID,
				ex.ID,
				set.SetOrder,
				set.Reps,
				set.Weight,
				set.RPE,
				set.SetType,
				set.CreatedAt,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *workoutRepository) UpdateSession(
	ctx context.Context,
	session *workout.WorkoutSession,
) error {
	return fmt.Errorf("not implemented")
}

func (r *workoutRepository) GetSessionByID(
	ctx context.Context,
	id string,
) (*workout.WorkoutSession, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *workoutRepository) GetSessionsByUser(
	ctx context.Context,
	userID string,
) ([]workout.WorkoutSession, error) {
	return nil, fmt.Errorf("not implemented")
}
