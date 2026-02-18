package orchestrator

import "context"

type Orchestrator interface {
	GenerateSplit(ctx context.Context, input SplitInput) (*SplitOutput, error)
	GenerateWorkout(ctx context.Context, input WorkoutInput) (*WorkoutOutput, error)
	SuggestOverload(ctx context.Context, input OverloadInput) (*OverloadOutput, error)
}
