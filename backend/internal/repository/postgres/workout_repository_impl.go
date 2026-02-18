package postgres

import (
	"context"
	"database/sql"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/workout"
	domainerr "S.P.A.R.T.A/backend/internal/domain/errors"
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
	if session == nil || session.ID == "" {
		return domainerr.ErrInvalidInput
	}

	var splitDayID any
	if session.SplitDayID != nil {
		splitDayID = *session.SplitDayID
	}

	res, err := r.db.ExecContext(ctx,
		`UPDATE workout_sessions
		 SET split_day_id=$2, session_date=$3, duration_minutes=$4, notes=$5
		 WHERE id=$1`,
		session.ID, splitDayID, session.SessionDate, session.DurationMin, session.Notes,
	)
	if err != nil {
		return domainerr.ErrInternal
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return domainerr.ErrNotFound
	}
	return nil
}

func (r *workoutRepository) GetSessionByID(
	ctx context.Context,
	id string,
) (*workout.WorkoutSession, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id,user_id,split_day_id,session_date,duration_minutes,notes,created_at
		 FROM workout_sessions
		 WHERE id=$1`,
		id,
	)

	var out workout.WorkoutSession
	var splitDay sql.NullString
	if err := row.Scan(
		&out.ID,
		&out.UserID,
		&splitDay,
		&out.SessionDate,
		&out.DurationMin,
		&out.Notes,
		&out.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, domainerr.ErrNotFound
		}
		return nil, domainerr.ErrInternal
	}
	if splitDay.Valid {
		s := splitDay.String
		out.SplitDayID = &s
	}

	exRows, err := r.db.QueryContext(ctx,
		`SELECT id, exercise_id
		 FROM workout_exercises
		 WHERE workout_session_id=$1`,
		out.ID,
	)
	if err != nil {
		return nil, domainerr.ErrInternal
	}
	defer exRows.Close()

	for exRows.Next() {
		var ex workout.WorkoutExercise
		if err := exRows.Scan(&ex.ID, &ex.ExerciseID); err != nil {
			return nil, domainerr.ErrInternal
		}

		setRows, err := r.db.QueryContext(ctx,
			`SELECT id,set_order,reps,weight,rpe,set_type,created_at
			 FROM workout_sets
			 WHERE workout_exercise_id=$1
			 ORDER BY set_order ASC`,
			ex.ID,
		)
		if err != nil {
			return nil, domainerr.ErrInternal
		}

		for setRows.Next() {
			var s workout.WorkoutSet
			if err := setRows.Scan(
				&s.ID,
				&s.SetOrder,
				&s.Reps,
				&s.Weight,
				&s.RPE,
				&s.SetType,
				&s.CreatedAt,
			); err != nil {
				_ = setRows.Close()
				return nil, domainerr.ErrInternal
			}
			ex.Sets = append(ex.Sets, s)
		}
		_ = setRows.Close()
		out.Exercises = append(out.Exercises, ex)
	}
	if err := exRows.Err(); err != nil {
		return nil, domainerr.ErrInternal
	}

	return &out, nil
}

func (r *workoutRepository) GetSessionsByUser(
	ctx context.Context,
	userID string,
) ([]workout.WorkoutSession, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id
		 FROM workout_sessions
		 WHERE user_id=$1
		 ORDER BY session_date DESC, created_at DESC
		 LIMIT 30`,
		userID,
	)
	if err != nil {
		return nil, domainerr.ErrInternal
	}
	defer rows.Close()

	ids := make([]string, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, domainerr.ErrInternal
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, domainerr.ErrInternal
	}

	items := make([]workout.WorkoutSession, 0, len(ids))
	for _, id := range ids {
		s, err := r.GetSessionByID(ctx, id)
		if err != nil {
			return nil, err
		}
		items = append(items, *s)
	}

	return items, nil
}
