package orchestrator

import "encoding/json"

func ParseSplitResponse(resp string) (*SplitOutput, error) {
	var out SplitOutput
	err := json.Unmarshal([]byte(resp), &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
