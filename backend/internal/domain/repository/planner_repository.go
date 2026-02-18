package repository

import (
	"context"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/planner"
)

type PlannerRepository interface {
    SaveRecommendation(ctx context.Context, rec *planner.PlannerRecommendation) error
    GetUserRecommendations(ctx context.Context, userID string) ([]planner.PlannerRecommendation, error)
}
