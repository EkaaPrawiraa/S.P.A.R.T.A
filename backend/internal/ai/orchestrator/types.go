package orchestrator

type SplitInput struct {
	UserID          string
	DaysPerWeek     int
	ExperienceLevel string
	FocusMuscle     string
}

type SplitOutput struct {
	Name string
	Days []SplitDayOutput
}

type SplitDayOutput struct {
	DayName   string
	Focus     []string
	Exercises []SplitExerciseOutput
}

type SplitExerciseOutput struct {
	Name     string
	Sets     int
	RepRange string
	Priority string
}

type WorkoutInput struct {
	UserID     string
	SplitDayID string
	Fatigue    int
	LastVolume int
}

type WorkoutOutput struct {
	Exercises []WorkoutExerciseOutput
}

type WorkoutExerciseOutput struct {
	Name     string
	Sets     int
	RepRange string
	Weight   float64
}

type OverloadInput struct {
	UserID      string
	ExerciseID  string
	LastWeight  float64
	LastReps    int
	Performance string
}

type OverloadOutput struct {
	Action  string
	Message string
}
