package repository

import (
	"context"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/nutrition"
)

type NutritionRepository interface {
	SaveDaily(ctx context.Context, n *nutrition.DailyNutrition) error
	GetByDate(ctx context.Context, userID string, date string) (*nutrition.DailyNutrition, error)
}
