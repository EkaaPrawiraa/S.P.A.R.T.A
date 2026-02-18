package usecase

import (
	"context"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/planner"
)

type PlannerUsecase interface {
	GenerateRecommendation(ctx context.Context, userID string) (*planner.PlannerRecommendation, error)
	SaveRecommendation(ctx context.Context, rec *planner.PlannerRecommendation) error
	GetUserRecommendations(ctx context.Context, userID string) ([]planner.PlannerRecommendation, error)
}
