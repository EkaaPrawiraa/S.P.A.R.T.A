package dto

type WorkoutSetDTO struct {
	SetOrder int     `json:"set_order" validate:"required"`
	Reps     int     `json:"reps" validate:"required"`
	Weight   float64 `json:"weight" validate:"required"`
	RPE      float64 `json:"rpe" validate:"required"`
	SetType  string  `json:"set_type" validate:"required"`
}

type WorkoutExerciseDTO struct {
	ExerciseID string          `json:"exercise_id" validate:"required,uuid4"`
	Sets       []WorkoutSetDTO `json:"sets" validate:"required,dive"`
}

type CreateWorkoutSessionDTO struct {
	UserID      string               `json:"user_id" validate:"required,uuid4"`
	SplitDayID  *string              `json:"split_day_id"`
	SessionDate string               `json:"session_date" validate:"required"`
	DurationMin int                  `json:"duration_minutes" validate:"required"`
	Notes       string               `json:"notes"`
	Exercises   []WorkoutExerciseDTO `json:"exercises" validate:"required,dive"`
}
