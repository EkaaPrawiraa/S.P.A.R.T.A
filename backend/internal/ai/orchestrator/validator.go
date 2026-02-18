package orchestrator

import (
	"fmt"
	"strings"

	domainerr "S.P.A.R.T.A/backend/internal/domain/errors"
)

func ValidateSplit(out *SplitOutput) error {
	if out.Name == "" {
		return fmt.Errorf("%w: invalid split name", domainerr.ErrInvalidInput)
	}

	if len(out.Days) == 0 {
		return fmt.Errorf("%w: split must contain days", domainerr.ErrInvalidInput)
	}

	for _, d := range out.Days {
		if len(d.Exercises) == 0 {
			return fmt.Errorf("%w: each day must contain exercises", domainerr.ErrInvalidInput)
		}
	}

	return nil
}

func ValidateWorkout(out *WorkoutOutput) error {
	if out == nil {
		return fmt.Errorf("%w: nil workout output", domainerr.ErrInternal)
	}
	if len(out.Exercises) == 0 {
		return fmt.Errorf("%w: workout must contain exercises", domainerr.ErrInvalidInput)
	}
	for _, ex := range out.Exercises {
		if ex.Name == "" || ex.Sets <= 0 || ex.RepRange == "" {
			return fmt.Errorf("%w: invalid workout exercise", domainerr.ErrInvalidInput)
		}
	}
	return nil
}

func ValidateOverload(out *OverloadOutput) error {
	if out == nil {
		return fmt.Errorf("%w: nil overload output", domainerr.ErrInternal)
	}
	if out.Action == "" || out.Message == "" {
		return fmt.Errorf("%w: invalid overload output", domainerr.ErrInvalidInput)
	}
	return nil
}

func ValidateMotivation(out *MotivationOutput) error {
	if out == nil {
		return fmt.Errorf("%w: nil motivation output", domainerr.ErrInternal)
	}
	if strings.TrimSpace(out.Message) == "" {
		return fmt.Errorf("%w: empty motivation message", domainerr.ErrInvalidInput)
	}
	return nil
}

func ValidateCoaching(out *CoachingOutput) error {
	if out == nil {
		return fmt.Errorf("%w: nil coaching output", domainerr.ErrInternal)
	}
	if len(out.Suggestions) == 0 {
		return fmt.Errorf("%w: empty suggestions", domainerr.ErrInvalidInput)
	}
	for _, s := range out.Suggestions {
		if strings.TrimSpace(s) == "" {
			return fmt.Errorf("%w: empty suggestion", domainerr.ErrInvalidInput)
		}
	}
	return nil
}

func ValidateExplainWorkoutPlan(out *ExplainWorkoutPlanOutput) error {
	if out == nil {
		return fmt.Errorf("%w: nil explain output", domainerr.ErrInternal)
	}
	if strings.TrimSpace(out.Summary) == "" {
		return fmt.Errorf("%w: empty summary", domainerr.ErrInvalidInput)
	}
	return nil
}
