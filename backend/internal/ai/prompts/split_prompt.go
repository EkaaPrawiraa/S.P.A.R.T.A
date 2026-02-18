package prompts

const SplitPromptTemplate = `You are an elite strength coach.

Generate a structured split plan in STRICT JSON format:

Rules:
- day_name MUST be non-empty (example: "Push", "Pull", "Legs", "Upper", "Lower", "Full Body")
- rep_range MUST be a string like "6-8" or "8-12"

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
}`
