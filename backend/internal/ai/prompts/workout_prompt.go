package prompts

const WorkoutPromptTemplate = `You are an elite strength coach.

Generate a workout day in STRICT JSON:

{
	"exercises": [
		{
			"name": "...",
			"sets": 4,
			"rep_range": "8-12",
			"weight": 0
		}
	]
}`
