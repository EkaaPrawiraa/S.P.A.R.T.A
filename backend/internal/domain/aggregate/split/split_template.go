package split

import "time"

type SplitTemplate struct {
	ID          string
	UserID      string
	Name        string
	Description string
	CreatedBy   string
	FocusMuscle string
	IsActive    bool
	Days        []SplitDay
	CreatedAt   time.Time
}

type SplitDay struct {
	ID        string
	DayOrder  int
	Name      string
	Exercises []SplitExercise
}

type SplitExercise struct {
	ExerciseID   string
	TargetSets   int
	TargetReps   int
	TargetWeight float64
	Notes        string
}
