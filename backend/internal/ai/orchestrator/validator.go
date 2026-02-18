package orchestrator

import "errors"

func ValidateSplit(out *SplitOutput) error {
	if out.Name == "" {
		return errors.New("invalid split name")
	}

	if len(out.Days) == 0 {
		return errors.New("split must contain days")
	}

	for _, d := range out.Days {
		if len(d.Exercises) == 0 {
			return errors.New("each day must contain exercises")
		}
	}

	return nil
}
