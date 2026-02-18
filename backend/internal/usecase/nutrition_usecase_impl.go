package usecase

import (
	"context"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/nutrition"
	domainrepo "S.P.A.R.T.A/backend/internal/domain/repository"
	domainuc "S.P.A.R.T.A/backend/internal/domain/usecase"
)

type nutritionUsecase struct {
	repo domainrepo.NutritionRepository
}

func NewNutritionUsecase(repo domainrepo.NutritionRepository) domainuc.NutritionUsecase {
	return &nutritionUsecase{repo: repo}
}

func (u *nutritionUsecase) SaveDaily(ctx context.Context, n *nutrition.DailyNutrition) error {
	return u.repo.SaveDaily(ctx, n)
}

func (u *nutritionUsecase) GetByDate(ctx context.Context, userID string, date string) (*nutrition.DailyNutrition, error) {
	return u.repo.GetByDate(ctx, userID, date)
}
