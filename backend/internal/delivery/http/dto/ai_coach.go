package dto

type GenerateSplitRequestDTO struct {
	DaysPerWeek int    `json:"days_per_week" validate:"required,gte=1,lte=7"`
	FocusMuscle string `json:"focus_muscle" validate:"required"`
}

type SuggestOverloadRequestDTO struct {
	ExerciseID string `json:"exercise_id" validate:"required,uuid4"`
}

type GenerateWorkoutPlanRequestDTO struct {
	SplitDayID string `json:"split_day_id" validate:"required,uuid4"`
	Fatigue    int    `json:"fatigue" validate:"required,gte=0,lte=10"`
}

type ExplainWorkoutExerciseDTO struct {
	Name     string  `json:"name" validate:"required"`
	Sets     int     `json:"sets" validate:"required,gte=1"`
	RepRange string  `json:"rep_range" validate:"required"`
	Weight   float64 `json:"weight" validate:"gte=0"`
}

type ExplainWorkoutPlanRequestDTO struct {
	SplitDayName string                      `json:"split_day_name"`
	Fatigue      int                         `json:"fatigue" validate:"required,gte=0,lte=10"`
	Exercises    []ExplainWorkoutExerciseDTO `json:"exercises" validate:"required,dive"`
}

type CoachingSuggestionsResponseDTO struct {
	Date        string   `json:"date"`
	Suggestions []string `json:"suggestions"`
}

type WorkoutExplanationExerciseNoteDTO struct {
	Name string `json:"name"`
	Note string `json:"note"`
}

type WorkoutExplanationResponseDTO struct {
	Summary       string                              `json:"summary"`
	ExerciseNotes []WorkoutExplanationExerciseNoteDTO `json:"exercise_notes"`
}
