package usecase

import (
	"context"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/planner"
	domainrepo "S.P.A.R.T.A/backend/internal/domain/repository"
	domainuc "S.P.A.R.T.A/backend/internal/domain/usecase"
)

type plannerUsecase struct {
	repo domainrepo.PlannerRepository
}

func NewPlannerUsecase(repo domainrepo.PlannerRepository) domainuc.PlannerUsecase {
	return &plannerUsecase{repo: repo}
}

func (u *plannerUsecase) SaveRecommendation(ctx context.Context, rec *planner.PlannerRecommendation) error {
	return u.repo.SaveRecommendation(ctx, rec)
}

func (u *plannerUsecase) GetUserRecommendations(ctx context.Context, userID string) ([]planner.PlannerRecommendation, error) {
	return u.repo.GetUserRecommendations(ctx, userID)
}
