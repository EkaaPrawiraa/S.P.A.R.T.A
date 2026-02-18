package usecase

import (
	"context"
	"time"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/planner"
	domainrepo "S.P.A.R.T.A/backend/internal/domain/repository"
	domainuc "S.P.A.R.T.A/backend/internal/domain/usecase"
	"github.com/google/uuid"
)

type plannerUsecase struct {
	repo domainrepo.PlannerRepository
}

func NewPlannerUsecase(repo domainrepo.PlannerRepository) domainuc.PlannerUsecase {
	return &plannerUsecase{repo: repo}
}

func (u *plannerUsecase) GenerateRecommendation(ctx context.Context, userID string) (*planner.PlannerRecommendation, error) {
	// Placeholder implementation to keep the vertical slice compiling.
	// Later this should call the AI orchestrator and/or deterministic rules engine.
	rec := &planner.PlannerRecommendation{
		ID:                 uuid.NewString(),
		UserID:             userID,
		Recommendation:     "planner recommendation not implemented yet",
		RecommendationType: "planner",
		CreatedAt:          time.Now(),
	}

	if err := u.repo.SaveRecommendation(ctx, rec); err != nil {
		return nil, err
	}

	return rec, nil
}

func (u *plannerUsecase) SaveRecommendation(ctx context.Context, rec *planner.PlannerRecommendation) error {
	return u.repo.SaveRecommendation(ctx, rec)
}

func (u *plannerUsecase) GetUserRecommendations(ctx context.Context, userID string) ([]planner.PlannerRecommendation, error) {
	return u.repo.GetUserRecommendations(ctx, userID)
}
