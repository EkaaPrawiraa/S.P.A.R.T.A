package workout

import "time"

// WorkoutPlan is an AI-generated plan/suggestion (not a logged session).
// It is safe to return via transport DTOs without coupling the domain to delivery.
type WorkoutPlan struct {
	UserID     string
	SplitDayID string
	Date       time.Time
	Exercises  []WorkoutPlanExercise
}

type WorkoutPlanExercise struct {
	Name     string
	Sets     int
	RepRange string
	Weight   float64
}
