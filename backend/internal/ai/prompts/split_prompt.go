package prompts

const SplitPromptTemplate = `You are an elite strength coach.

Generate a structured split plan in STRICT JSON format:

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
