package prompts

const ExplainWorkoutPromptTemplate = `You are an elite strength coach.

Explain the workout plan in practical terms.

Return STRICT JSON:

{
  "summary": "...",
  "exercise_notes": [
    {
      "name": "...",
      "note": "..."
    }
  ]
}
`
