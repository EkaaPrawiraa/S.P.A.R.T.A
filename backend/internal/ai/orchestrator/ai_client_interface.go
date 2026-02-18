package orchestrator

import "context"

type AIClient interface {
	Generate(ctx context.Context, prompt string) (string, error)
}
