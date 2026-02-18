package postgres

import (
	"context"
	"database/sql"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/planner"
	domainrepo "S.P.A.R.T.A/backend/internal/domain/repository"
)

type plannerRepository struct {
	db *sql.DB
}

func NewPlannerRepository(db *sql.DB) domainrepo.PlannerRepository {
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
	return nil, nil
}
