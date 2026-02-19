package orchestrator

import (
	"encoding/json"
	"fmt"
	"strings"

	domainerr "S.P.A.R.T.A/backend/internal/domain/errors"
)

func ParseSplitResponse(resp string) (*SplitOutput, error) {
	var out SplitOutput
	clean := extractJSONObject(resp)
	err := json.Unmarshal([]byte(clean), &out)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid AI JSON", domainerr.ErrAIUnavailable)
	}
	return &out, nil
}

func ParseWorkoutResponse(resp string) (*WorkoutOutput, error) {
	var out WorkoutOutput
	clean := extractJSONObject(resp)
	err := json.Unmarshal([]byte(clean), &out)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid AI JSON", domainerr.ErrAIUnavailable)
	}
	return &out, nil
}

func ParseOverloadResponse(resp string) (*OverloadOutput, error) {
	var out OverloadOutput
	clean := extractJSONObject(resp)
	err := json.Unmarshal([]byte(clean), &out)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid AI JSON", domainerr.ErrAIUnavailable)
	}
	return &out, nil
}

func ParseMotivationResponse(resp string) (*MotivationOutput, error) {
	var out MotivationOutput
	clean := extractJSONObject(resp)
	err := json.Unmarshal([]byte(clean), &out)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid AI JSON", domainerr.ErrAIUnavailable)
	}
	return &out, nil
}

func ParseCoachingResponse(resp string) (*CoachingOutput, error) {
	var out CoachingOutput
	clean := extractJSONObject(resp)
	err := json.Unmarshal([]byte(clean), &out)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid AI JSON", domainerr.ErrAIUnavailable)
	}
	return &out, nil
}

func ParseExplainWorkoutPlanResponse(resp string) (*ExplainWorkoutPlanOutput, error) {
	var out ExplainWorkoutPlanOutput
	clean := extractJSONObject(resp)
	err := json.Unmarshal([]byte(clean), &out)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid AI JSON", domainerr.ErrAIUnavailable)
	}
	return &out, nil
}

func extractJSONObject(s string) string {
	// Best-effort: strip code fences and take the first {...} block.
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "```json")
	s = strings.TrimPrefix(s, "```")
	s = strings.TrimSuffix(s, "```")
	s = strings.TrimSpace(s)

	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start >= 0 && end > start {
		return strings.TrimSpace(s[start : end+1])
	}
	return s
}
