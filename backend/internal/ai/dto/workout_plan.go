package dto

type WorkoutPlan struct {
	DayName   string          `json:"day_name"`
	Exercises []WorkoutExercise `json:"exercises"`
}

type WorkoutExercise struct {
	Name     string  `json:"name"`
	Sets     int     `json:"sets"`
	RepRange string  `json:"rep_range"`
	Weight   float64 `json:"weight"`
}

