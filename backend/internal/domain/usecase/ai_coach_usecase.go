package usecase

import (
	"context"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/planner"
	"S.P.A.R.T.A/backend/internal/domain/aggregate/split"
	"S.P.A.R.T.A/backend/internal/domain/aggregate/workout"
	"github.com/google/uuid"
)

type AICoachUsecase interface {
	GenerateSplitTemplate(ctx context.Context, userID uuid.UUID, daysPerWeek int, focusMuscle string) (*split.SplitTemplate, error)
	SuggestProgressiveOverload(ctx context.Context, userID uuid.UUID, exerciseID uuid.UUID) (*planner.PlannerRecommendation, error)
	GetDailyMotivation(ctx context.Context, userID uuid.UUID) (string, error)
	ResetDailyMotivation(ctx context.Context, userID uuid.UUID) error
	GenerateWorkoutPlan(ctx context.Context, userID uuid.UUID, splitDayID uuid.UUID, fatigue int) (*workout.WorkoutPlan, error)
	GetCoachingSuggestions(ctx context.Context, userID uuid.UUID) ([]string, error)
	ExplainWorkoutPlan(ctx context.Context, userID uuid.UUID, plan workout.WorkoutPlan, splitDayName string, fatigue int) (*workout.WorkoutExplanation, error)
}
