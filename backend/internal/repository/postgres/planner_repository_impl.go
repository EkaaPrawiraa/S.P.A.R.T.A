package postgres

import (
	"context"
	"database/sql"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/planner"
	domainerr "S.P.A.R.T.A/backend/internal/domain/errors"
	domainrepo "S.P.A.R.T.A/backend/internal/domain/repository"
)

type plannerRepository struct {
	db DBTX
}

func NewPlannerRepository(db DBTX) domainrepo.PlannerRepository {
	return &plannerRepository{db: db}
}

func (r *plannerRepository) SaveRecommendation(ctx context.Context, rec *planner.PlannerRecommendation) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO planner_recommendations(id,user_id,workout_session_id,recommendation,recommendation_type,created_at)
		 VALUES ($1,$2,$3,$4,$5,$6)`,
		rec.ID, rec.UserID, rec.WorkoutSessionID, rec.Recommendation, rec.RecommendationType, rec.CreatedAt)
	return err
}

func (r *plannerRepository) GetUserRecommendations(ctx context.Context, userID string) ([]planner.PlannerRecommendation, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id,user_id,workout_session_id,recommendation,recommendation_type,created_at
		 FROM planner_recommendations
		 WHERE user_id=$1
		 ORDER BY created_at DESC
		 LIMIT 50`,
		userID,
	)
	if err != nil {
		return nil, domainerr.ErrInternal
	}
	defer rows.Close()

	items := make([]planner.PlannerRecommendation, 0)
	for rows.Next() {
		var rec planner.PlannerRecommendation
		var workoutSessionID sql.NullString
		if err := rows.Scan(
			&rec.ID,
			&rec.UserID,
			&workoutSessionID,
			&rec.Recommendation,
			&rec.RecommendationType,
			&rec.CreatedAt,
		); err != nil {
			return nil, domainerr.ErrInternal
		}
		if workoutSessionID.Valid {
			s := workoutSessionID.String
			rec.WorkoutSessionID = &s
		}
		items = append(items, rec)
	}
	if err := rows.Err(); err != nil {
		return nil, domainerr.ErrInternal
	}
	return items, nil
}
