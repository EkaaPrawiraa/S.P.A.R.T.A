package orchestrator

type SplitInput struct {
	UserID          string
	DaysPerWeek     int
	ExperienceLevel string
	FocusMuscle     string
}

type SplitOutput struct {
	Name string           `json:"name"`
	Days []SplitDayOutput `json:"days"`
}

type SplitDayOutput struct {
	DayName   string              `json:"day_name"`
	Focus     []string            `json:"focus"`
	Exercises []SplitExerciseOutput `json:"exercises"`
}

type SplitExerciseOutput struct {
	Name     string `json:"name"`
	Sets     int    `json:"sets"`
	RepRange string `json:"rep_range"`
	Priority string `json:"priority"`
}

type WorkoutInput struct {
	UserID           string
	SplitDayID       string
	SplitDayName     string
	PlannedExercises []string
	Fatigue          int
	AcuteLoad7d      float64
	ChronicLoad28d   float64
	ACWR             float64
	FatigueEstimated int
	LastVolume       int
}

type WorkoutOutput struct {
	Exercises []WorkoutExerciseOutput `json:"exercises"`
}

type WorkoutExerciseOutput struct {
	Name     string  `json:"name"`
	Sets     int     `json:"sets"`
	RepRange string  `json:"rep_range"`
	Weight   float64 `json:"weight"`
}

type OverloadInput struct {
	UserID      string
	ExerciseID  string
	LastWeight  float64
	LastReps    int
	Performance string
}

type OverloadOutput struct {
	Action  string `json:"action"`
	Message string `json:"message"`
}

type MotivationInput struct {
	UserID              string
	Date                string
	WorkoutsLast7Days   int
	LastWorkoutDate     string
	LastWorkoutDuration int
	LastWorkoutNotes    string
}

type MotivationOutput struct {
	Message string `json:"message"`
}

type CoachingInput struct {
	UserID            string
	Date              string
	AcuteLoad7d       float64
	ChronicLoad28d    float64
	ACWR              float64
	RecentWorkouts    string
	RecentNutrition   string
	RecentPlannerRecs string
}

type CoachingOutput struct {
	Suggestions []string `json:"suggestions"`
}

type ExplainWorkoutPlanInput struct {
	UserID       string
	SplitDayName string
	Fatigue      int
	Exercises    []ExplainWorkoutExercise
}

type ExplainWorkoutExercise struct {
	Name     string
	Sets     int
	RepRange string
	Weight   float64
}

type ExplainWorkoutPlanOutput struct {
	Summary       string
	ExerciseNotes []ExplainExerciseNote `json:"exercise_notes"`
}

type ExplainExerciseNote struct {
	Name string
	Note string
}
