package orchestrator

import (
	"context"
	"strings"
	"time"
)

type service struct {
	client AIClient
}

func NewOrchestrator(client AIClient) Orchestrator {
	return &service{
		client: client,
	}
}

func (s *service) GenerateSplit(ctx context.Context, input SplitInput) (*SplitOutput, error) {
	prompt := BuildSplitPrompt(input)

	var resp string
	err := WithRetry(ctx, RetryConfig{MaxAttempts: 3, Delay: 250 * time.Millisecond}, func(ctx context.Context) error {
		out, err := s.client.Generate(ctx, prompt)
		if err != nil {
			return err
		}
		resp = out
		return nil
	})
	if err != nil {
		return nil, err
	}

	out, err := ParseSplitResponse(resp)
	if err != nil {
		return nil, err
	}

	if err := ValidateSplit(out); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *service) GenerateWorkout(ctx context.Context, input WorkoutInput) (*WorkoutOutput, error) {
	prompt := BuildWorkoutPrompt(input)

	var resp string
	err := WithRetry(ctx, RetryConfig{MaxAttempts: 3, Delay: 250 * time.Millisecond}, func(ctx context.Context) error {
		out, err := s.client.Generate(ctx, prompt)
		if err != nil {
			return err
		}
		resp = out
		return nil
	})
	if err != nil {
		return nil, err
	}

	out, err := ParseWorkoutResponse(resp)
	if err != nil {
		return nil, err
	}

	if err := ValidateWorkout(out); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *service) SuggestOverload(ctx context.Context, input OverloadInput) (*OverloadOutput, error) {
	prompt := BuildOverloadPrompt(input)

	var resp string
	err := WithRetry(ctx, RetryConfig{MaxAttempts: 3, Delay: 250 * time.Millisecond}, func(ctx context.Context) error {
		out, err := s.client.Generate(ctx, prompt)
		if err != nil {
			return err
		}
		resp = out
		return nil
	})
	if err != nil {
		return nil, err
	}

	out, err := ParseOverloadResponse(resp)
	if err != nil {
		return nil, err
	}

	if err := ValidateOverload(out); err != nil {
		return nil, err
	}

	// Normalize action.
	out.Action = strings.ToLower(strings.TrimSpace(out.Action))
	return out, nil
}

func (s *service) GenerateDailyMotivation(ctx context.Context, input MotivationInput) (*MotivationOutput, error) {
	prompt := BuildMotivationPrompt(input)

	var resp string
	err := WithRetry(ctx, RetryConfig{MaxAttempts: 3, Delay: 250 * time.Millisecond}, func(ctx context.Context) error {
		out, err := s.client.Generate(ctx, prompt)
		if err != nil {
			return err
		}
		resp = out
		return nil
	})
	if err != nil {
		return nil, err
	}

	out, err := ParseMotivationResponse(resp)
	if err != nil {
		return nil, err
	}

	if err := ValidateMotivation(out); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *service) GenerateCoachingSuggestions(ctx context.Context, input CoachingInput) (*CoachingOutput, error) {
	prompt := BuildCoachingPrompt(input)

	var resp string
	err := WithRetry(ctx, RetryConfig{MaxAttempts: 3, Delay: 250 * time.Millisecond}, func(ctx context.Context) error {
		out, err := s.client.Generate(ctx, prompt)
		if err != nil {
			return err
		}
		resp = out
		return nil
	})
	if err != nil {
		return nil, err
	}

	out, err := ParseCoachingResponse(resp)
	if err != nil {
		return nil, err
	}
	if err := ValidateCoaching(out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *service) ExplainWorkoutPlan(ctx context.Context, input ExplainWorkoutPlanInput) (*ExplainWorkoutPlanOutput, error) {
	prompt := BuildExplainWorkoutPlanPrompt(input)

	var resp string
	err := WithRetry(ctx, RetryConfig{MaxAttempts: 3, Delay: 250 * time.Millisecond}, func(ctx context.Context) error {
		out, err := s.client.Generate(ctx, prompt)
		if err != nil {
			return err
		}
		resp = out
		return nil
	})
	if err != nil {
		return nil, err
	}

	out, err := ParseExplainWorkoutPlanResponse(resp)
	if err != nil {
		return nil, err
	}
	if err := ValidateExplainWorkoutPlan(out); err != nil {
		return nil, err
	}
	return out, nil
}
