package workout

import "time"

type WorkoutSession struct {
	ID          string
	UserID      string
	SplitDayID  *string
	SessionDate time.Time
	DurationMin int
	Notes       string
	Exercises   []WorkoutExercise
	CreatedAt   time.Time
}

type WorkoutExercise struct {
	ID         string
	ExerciseID string
	Sets       []WorkoutSet
}

type WorkoutSet struct {
	ID        string
	SetOrder  int
	Reps      int
	Weight    float64
	RPE       float64
	SetType   string
	CreatedAt time.Time
}
