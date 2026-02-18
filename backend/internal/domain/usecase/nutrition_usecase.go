package usecase

import (
	"context"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/nutrition"
)

type NutritionUsecase interface {
	SaveDaily(ctx context.Context, n *nutrition.DailyNutrition) error
	GetByDate(ctx context.Context, userID string, date string) (*nutrition.DailyNutrition, error)
}
