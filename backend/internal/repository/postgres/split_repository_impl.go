package postgres

import (
	"context"
	"database/sql"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/split"
	domainerr "S.P.A.R.T.A/backend/internal/domain/errors"
	domainrepo "S.P.A.R.T.A/backend/internal/domain/repository"
	"github.com/google/uuid"
)

type splitRepository struct {
	db DBTX
}

func NewSplitRepository(db DBTX) domainrepo.SplitRepository {
	return &splitRepository{db: db}
}

func (r *splitRepository) CreateTemplate(ctx context.Context, tpl *split.SplitTemplate) error {
	if tpl == nil || tpl.ID == "" {
		return domainerr.ErrInvalidInput
	}

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO split_templates(id,user_id,name,description,created_by,focus_muscle,is_active,created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		tpl.ID, tpl.UserID, tpl.Name, tpl.Description, tpl.CreatedBy, tpl.FocusMuscle, tpl.IsActive, tpl.CreatedAt)
	if err != nil {
		return domainerr.ErrInternal
	}

	for _, day := range tpl.Days {
		_, err := r.db.ExecContext(ctx,
			`INSERT INTO split_days(id,split_template_id,day_order,name)
			 VALUES ($1,$2,$3,$4)`,
			day.ID, tpl.ID, day.DayOrder, day.Name,
		)
		if err != nil {
			return domainerr.ErrInternal
		}

		for _, ex := range day.Exercises {
			var exerciseID any
			if ex.ExerciseID != "" {
				exerciseID = ex.ExerciseID
			}

			_, err := r.db.ExecContext(ctx,
				`INSERT INTO split_day_exercises(id,split_day_id,exercise_id,target_sets,target_reps,target_weight,notes)
				 VALUES ($1,$2,$3,$4,$5,$6,$7)`,
				uuid.NewString(),
				day.ID,
				exerciseID,
				ex.TargetSets,
				ex.TargetReps,
				ex.TargetWeight,
				ex.Notes,
			)
			if err != nil {
				return domainerr.ErrInternal
			}
		}
	}

	return nil
}

func (r *splitRepository) UpdateTemplate(ctx context.Context, tpl *split.SplitTemplate) error {
	if tpl == nil || tpl.ID == "" || tpl.UserID == "" {
		return domainerr.ErrInvalidInput
	}

	res, err := r.db.ExecContext(ctx,
		`UPDATE split_templates
		 SET name=$3, description=$4, focus_muscle=$5, is_active=$6
		 WHERE id=$1 AND user_id=$2`,
		tpl.ID,
		tpl.UserID,
		tpl.Name,
		tpl.Description,
		tpl.FocusMuscle,
		tpl.IsActive,
	)
	if err != nil {
		return domainerr.ErrInternal
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return domainerr.ErrNotFound
	}

	// Rewrite days/exercises for the template.
	_, err = r.db.ExecContext(ctx,
		`DELETE FROM split_day_exercises
		 WHERE split_day_id IN (SELECT id FROM split_days WHERE split_template_id=$1)`,
		tpl.ID,
	)
	if err != nil {
		return domainerr.ErrInternal
	}
	_, err = r.db.ExecContext(ctx,
		`DELETE FROM split_days WHERE split_template_id=$1`,
		tpl.ID,
	)
	if err != nil {
		return domainerr.ErrInternal
	}

	for _, day := range tpl.Days {
		_, err := r.db.ExecContext(ctx,
			`INSERT INTO split_days(id,split_template_id,day_order,name)
			 VALUES ($1,$2,$3,$4)`,
			day.ID, tpl.ID, day.DayOrder, day.Name,
		)
		if err != nil {
			return domainerr.ErrInternal
		}

		for _, ex := range day.Exercises {
			var exerciseID any
			if ex.ExerciseID != "" {
				exerciseID = ex.ExerciseID
			}

			_, err := r.db.ExecContext(ctx,
				`INSERT INTO split_day_exercises(id,split_day_id,exercise_id,target_sets,target_reps,target_weight,notes)
				 VALUES ($1,$2,$3,$4,$5,$6,$7)`,
				uuid.NewString(),
				day.ID,
				exerciseID,
				ex.TargetSets,
				ex.TargetReps,
				ex.TargetWeight,
				ex.Notes,
			)
			if err != nil {
				return domainerr.ErrInternal
			}
		}
	}

	return nil
}

func (r *splitRepository) ActivateTemplate(ctx context.Context, userID string, templateID string) error {
	if userID == "" || templateID == "" {
		return domainerr.ErrInvalidInput
	}

	// Deactivate all templates for the user.
	if _, err := r.db.ExecContext(ctx,
		`UPDATE split_templates SET is_active=false WHERE user_id=$1`,
		userID,
	); err != nil {
		return domainerr.ErrInternal
	}

	res, err := r.db.ExecContext(ctx,
		`UPDATE split_templates SET is_active=true WHERE id=$1 AND user_id=$2`,
		templateID,
		userID,
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

func (r *splitRepository) GetTemplateByID(ctx context.Context, id string) (*split.SplitTemplate, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id,user_id,name,description,created_by,focus_muscle,is_active,created_at
		 FROM split_templates
		 WHERE id=$1`,
		id,
	)

	var out split.SplitTemplate
	if err := row.Scan(
		&out.ID,
		&out.UserID,
		&out.Name,
		&out.Description,
		&out.CreatedBy,
		&out.FocusMuscle,
		&out.IsActive,
		&out.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, domainerr.ErrNotFound
		}
		return nil, domainerr.ErrInternal
	}

	days, err := r.getDays(ctx, out.ID)
	if err != nil {
		return nil, err
	}
	out.Days = days
	return &out, nil
}

func (r *splitRepository) GetUserTemplates(ctx context.Context, userID string) ([]split.SplitTemplate, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id,user_id,name,description,created_by,focus_muscle,is_active,created_at
		 FROM split_templates
		 WHERE user_id=$1
		 ORDER BY created_at DESC
		 LIMIT 20`,
		userID,
	)
	if err != nil {
		return nil, domainerr.ErrInternal
	}
	defer rows.Close()

	items := make([]split.SplitTemplate, 0)
	for rows.Next() {
		var tpl split.SplitTemplate
		if err := rows.Scan(
			&tpl.ID,
			&tpl.UserID,
			&tpl.Name,
			&tpl.Description,
			&tpl.CreatedBy,
			&tpl.FocusMuscle,
			&tpl.IsActive,
			&tpl.CreatedAt,
		); err != nil {
			return nil, domainerr.ErrInternal
		}

		days, err := r.getDays(ctx, tpl.ID)
		if err != nil {
			return nil, err
		}
		tpl.Days = days
		items = append(items, tpl)
	}
	if err := rows.Err(); err != nil {
		return nil, domainerr.ErrInternal
	}
	return items, nil
}

func (r *splitRepository) GetSplitDayByID(ctx context.Context, id string) (*split.SplitDay, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id,day_order,name
		 FROM split_days
		 WHERE id=$1`,
		id,
	)

	var d split.SplitDay
	if err := row.Scan(&d.ID, &d.DayOrder, &d.Name); err != nil {
		if err == sql.ErrNoRows {
			return nil, domainerr.ErrNotFound
		}
		return nil, domainerr.ErrInternal
	}

	ex, err := r.getDayExercises(ctx, d.ID)
	if err != nil {
		return nil, err
	}
	d.Exercises = ex
	return &d, nil
}

func (r *splitRepository) getDays(ctx context.Context, templateID string) ([]split.SplitDay, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id,day_order,name
		 FROM split_days
		 WHERE split_template_id=$1
		 ORDER BY day_order ASC`,
		templateID,
	)
	if err != nil {
		return nil, domainerr.ErrInternal
	}
	defer rows.Close()

	items := make([]split.SplitDay, 0)
	for rows.Next() {
		var d split.SplitDay
		if err := rows.Scan(&d.ID, &d.DayOrder, &d.Name); err != nil {
			return nil, domainerr.ErrInternal
		}
		ex, err := r.getDayExercises(ctx, d.ID)
		if err != nil {
			return nil, err
		}
		d.Exercises = ex
		items = append(items, d)
	}
	if err := rows.Err(); err != nil {
		return nil, domainerr.ErrInternal
	}
	return items, nil
}

func (r *splitRepository) getDayExercises(ctx context.Context, dayID string) ([]split.SplitExercise, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT COALESCE(exercise_id::text,''), target_sets, target_reps, COALESCE(target_weight,0), COALESCE(notes,'')
		 FROM split_day_exercises
		 WHERE split_day_id=$1`,
		dayID,
	)
	if err != nil {
		return nil, domainerr.ErrInternal
	}
	defer rows.Close()

	items := make([]split.SplitExercise, 0)
	for rows.Next() {
		var ex split.SplitExercise
		if err := rows.Scan(&ex.ExerciseID, &ex.TargetSets, &ex.TargetReps, &ex.TargetWeight, &ex.Notes); err != nil {
			return nil, domainerr.ErrInternal
		}
		items = append(items, ex)
	}
	if err := rows.Err(); err != nil {
		return nil, domainerr.ErrInternal
	}
	return items, nil
}
