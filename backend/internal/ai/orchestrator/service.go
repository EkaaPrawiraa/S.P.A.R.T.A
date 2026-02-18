package orchestrator

import (
	"context"
	"errors"
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

	resp, err := s.client.Generate(ctx, prompt)
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
	return nil, errors.New("not implemented")
}

func (s *service) SuggestOverload(ctx context.Context, input OverloadInput) (*OverloadOutput, error) {
	return nil, errors.New("not implemented")
}
