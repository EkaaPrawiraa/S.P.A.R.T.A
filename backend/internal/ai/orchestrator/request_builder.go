package orchestrator

import "fmt"

func BuildSplitPrompt(input SplitInput) string {
	return fmt.Sprintf(`
You are an elite strength coach.

Generate a structured workout split in STRICT JSON format.

User experience: %s
Training days per week: %d
Primary focus muscle: %s

Return JSON:
{
  "name": "...",
  "days": [
    {
      "day_name": "...",
      "focus": [],
      "exercises": [
        {
          "name": "...",
          "sets": 4,
          "rep_range": "8-12",
          "priority": "primary"
        }
      ]
    }
  ]
}
`, input.ExperienceLevel, input.DaysPerWeek, input.FocusMuscle)
}
