package usecase

import (
	"context"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/planner"
	"S.P.A.R.T.A/backend/internal/domain/aggregate/split"
	"github.com/google/uuid"
)

type AICoachUsecase interface {
	GenerateSplitTemplate(ctx context.Context, userID uuid.UUID, daysPerWeek int, focusMuscle string) (*split.SplitTemplate, error)
	SuggestProgressiveOverload(ctx context.Context, userID uuid.UUID, exerciseID uuid.UUID) (*planner.PlannerRecommendation, error)
}
