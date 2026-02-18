package usecase

import (
	"context"
	"time"

	"S.P.A.R.T.A/backend/internal/ai/orchestrator"
	"S.P.A.R.T.A/backend/internal/domain/aggregate/planner"
	"S.P.A.R.T.A/backend/internal/domain/aggregate/split"
	domainrepo "S.P.A.R.T.A/backend/internal/domain/repository"
	domainuc "S.P.A.R.T.A/backend/internal/domain/usecase"
	"github.com/google/uuid"
)

type aiCoachUsecase struct {
	orchestrator      orchestrator.Orchestrator
	splitRepository   domainrepo.SplitRepository
	plannerRepository domainrepo.PlannerRepository
}

func NewAICoachUsecase(
	orchestrator orchestrator.Orchestrator,
	splitRepository domainrepo.SplitRepository,
	plannerRepository domainrepo.PlannerRepository,
) domainuc.AICoachUsecase {
	return &aiCoachUsecase{
		orchestrator:      orchestrator,
		splitRepository:   splitRepository,
		plannerRepository: plannerRepository,
	}
}

func (u *aiCoachUsecase) GenerateSplitTemplate(
	ctx context.Context,
	userID uuid.UUID,
	daysPerWeek int,
	focusMuscle string,
) (*split.SplitTemplate, error) {

	aiResult, err := u.orchestrator.GenerateSplit(ctx, orchestrator.SplitInput{
		UserID:      userID.String(),
		DaysPerWeek: daysPerWeek,
		FocusMuscle: focusMuscle,
	})
	if err != nil {
		return nil, err
	}

	template := &split.SplitTemplate{
		ID:          uuid.NewString(),
		UserID:      userID.String(),
		Name:        aiResult.Name,
		Description: "AI Generated Split",
		CreatedBy:   "ai",
		FocusMuscle: focusMuscle,
		IsActive:    false,
		CreatedAt:   time.Now(),
	}

	for i, d := range aiResult.Days {
		day := split.SplitDay{
			ID:       uuid.NewString(),
			DayOrder:        i + 1,
			Name:            d.DayName,
		}

		for _, ex := range d.Exercises {
			day.Exercises = append(day.Exercises, split.SplitExercise{
				TargetSets:  ex.Sets,
				TargetReps:  parseRepRange(ex.RepRange),
				Notes:       ex.Name,
			})
		}

		template.Days = append(template.Days, day)
	}

	err = u.splitRepository.CreateTemplate(ctx, template)
	if err != nil {
		return nil, err
	}

	return template, nil
}

func (u *aiCoachUsecase) SuggestProgressiveOverload(
	ctx context.Context,
	userID uuid.UUID,
	exerciseID uuid.UUID,
) (*planner.PlannerRecommendation, error) {

	result, err := u.orchestrator.SuggestOverload(ctx, orchestrator.OverloadInput{
		UserID:     userID.String(),
		ExerciseID: exerciseID.String(),
	})
	if err != nil {
		return nil, err
	}

	rec := &planner.PlannerRecommendation{
		ID:                 uuid.NewString(),
		UserID:             userID.String(),
		Recommendation:     result.Message,
		RecommendationType: "progressive_overload",
		CreatedAt:          time.Now(),
	}

	err = u.plannerRepository.SaveRecommendation(ctx, rec)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

func parseRepRange(repRange string) int {
	// simple parser example (8-12 -> 10)
	return 10
}
