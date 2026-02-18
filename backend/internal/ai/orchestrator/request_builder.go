package orchestrator

import (
	"fmt"
	"strings"

	"S.P.A.R.T.A/backend/internal/ai/prompts"
)

func BuildSplitPrompt(input SplitInput) string {
	return fmt.Sprintf(`
You are an elite strength coach.

Generate a structured workout split in STRICT JSON format.

User experience: %s
Training days per week: %d
Primary focus muscle: %s

Return JSON schema:
%s
`, input.ExperienceLevel, input.DaysPerWeek, input.FocusMuscle, prompts.SplitPromptTemplate)
}

func BuildWorkoutPrompt(input WorkoutInput) string {
	planned := ""
	if len(input.PlannedExercises) > 0 {
		planned = "- " + strings.Join(input.PlannedExercises, "\n- ")
	}
	return fmt.Sprintf(`
You are an elite strength coach.

Generate today's workout in STRICT JSON format.

User ID: %s
Split day ID: %s
Split day name: %s
Planned exercises (if provided):
%s
Fatigue (0-10): %d
Estimated fatigue (0-10): %d
Acute load 7d: %.0f
Chronic load 28d: %.0f
ACWR: %.2f
Last volume (arbitrary units): %d

Return JSON schema:
%s
`, input.UserID, input.SplitDayID, input.SplitDayName, planned, input.Fatigue, input.FatigueEstimated, input.AcuteLoad7d, input.ChronicLoad28d, input.ACWR, input.LastVolume, prompts.WorkoutPromptTemplate)
}

func BuildOverloadPrompt(input OverloadInput) string {
	return fmt.Sprintf(`
You are an elite strength coach.

Suggest progressive overload in STRICT JSON.

User ID: %s
Exercise ID: %s
Last weight: %.2f
Last reps: %d
Performance notes: %s

Return JSON schema:
%s
`, input.UserID, input.ExerciseID, input.LastWeight, input.LastReps, input.Performance, prompts.OverloadPromptTemplate)
}

func BuildMotivationPrompt(input MotivationInput) string {
	return fmt.Sprintf(`
You are S.P.A.R.T.A — a supportive, no-BS gym motivation coach.
Generate a short daily motivation message based on recent activity.

Date: %s
Workouts last 7 days: %d
Last workout date: %s
Last workout duration minutes: %d
Last workout notes: %s

Rules:
- Keep it concise (max ~10-30 words).
- Be encouraging but practical.
- If workouts last 7 days is 0, focus on getting started today.
- DO NOT use hypen '—' in the quotes.

Return STRICT JSON:
{
  "message": "..."
}
`, input.Date, input.WorkoutsLast7Days, input.LastWorkoutDate, input.LastWorkoutDuration, input.LastWorkoutNotes)
}

func BuildCoachingPrompt(input CoachingInput) string {
	return fmt.Sprintf(`
You are S.P.A.R.T.A — a supportive, no-BS gym coach.

Generate coaching suggestions in STRICT JSON.

Date: %s
User ID: %s

Acute load 7d: %.0f
Chronic load 28d: %.0f
ACWR: %.2f

Recent workouts:
%s

Recent nutrition:
%s

Recent recommendations:
%s

Return JSON schema:
%s
`, input.Date, input.UserID, input.AcuteLoad7d, input.ChronicLoad28d, input.ACWR, input.RecentWorkouts, input.RecentNutrition, input.RecentPlannerRecs, prompts.CoachingPromptTemplate)
}

func BuildExplainWorkoutPlanPrompt(input ExplainWorkoutPlanInput) string {
	lines := make([]string, 0, len(input.Exercises))
	for _, ex := range input.Exercises {
		lines = append(lines, fmt.Sprintf("- %s | %dx%s | weight %.2f", ex.Name, ex.Sets, ex.RepRange, ex.Weight))
	}
	plan := strings.Join(lines, "\n")

	return fmt.Sprintf(`
You are an elite strength coach.

Explain the workout plan in a practical way.

User ID: %s
Split day name: %s
Fatigue (0-10): %d

Workout plan:
%s

Rules:
- Keep it concise.
- Explain intent, what to focus on, and key technique cues.
- Do not invent exercises not in the plan.

Return JSON schema:
%s
`, input.UserID, input.SplitDayName, input.Fatigue, plan, prompts.ExplainWorkoutPromptTemplate)
}
