package usecase

import (
	"context"
	"strings"
	"time"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/planner"
	domainerr "S.P.A.R.T.A/backend/internal/domain/errors"
	domainrepo "S.P.A.R.T.A/backend/internal/domain/repository"
	domainuc "S.P.A.R.T.A/backend/internal/domain/usecase"
	"github.com/google/uuid"
)

type plannerUsecase struct {
	repo    domainrepo.PlannerRepository
	aiCoach domainuc.AICoachUsecase
}

func NewPlannerUsecase(repo domainrepo.PlannerRepository, aiCoach domainuc.AICoachUsecase) domainuc.PlannerUsecase {
	return &plannerUsecase{repo: repo, aiCoach: aiCoach}
}

func (u *plannerUsecase) GenerateRecommendation(ctx context.Context, userID string) (*planner.PlannerRecommendation, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, domainerr.ErrInvalidInput
	}
	if u.aiCoach == nil {
		return nil, domainerr.ErrInternal
	}

	suggestions, err := u.aiCoach.GetCoachingSuggestions(ctx, uid)
	if err != nil {
		return nil, err
	}

	lines := make([]string, 0, len(suggestions))
	for _, s := range suggestions {
		item := strings.TrimSpace(s)
		if item == "" {
			continue
		}
		lines = append(lines, "- "+item)
	}

	recText := strings.TrimSpace(strings.Join(lines, "\n"))
	if recText == "" {
		recText = "(no coaching suggestions generated)"
	}

	rec := &planner.PlannerRecommendation{
		ID:                 uuid.NewString(),
		UserID:             userID,
		Recommendation:     recText,
		RecommendationType: "ai_coaching",
		CreatedAt:          time.Now().UTC(),
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
