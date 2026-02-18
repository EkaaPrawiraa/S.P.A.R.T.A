package dto

type CreateSplitTemplateDTO struct {
	UserID      string        `json:"user_id" validate:"required,uuid4"`
	Name        string        `json:"name" validate:"required"`
	Description string        `json:"description"`
	CreatedBy   string        `json:"created_by" validate:"required"`
	FocusMuscle string        `json:"focus_muscle" validate:"required"`
	IsActive    bool          `json:"is_active"`
	Days        []SplitDayDTO `json:"days" validate:"required,dive"`
}

type SplitDayDTO struct {
	DayOrder  int                `json:"day_order" validate:"required,gte=1"`
	Name      string             `json:"name" validate:"required"`
	Exercises []SplitExerciseDTO `json:"exercises" validate:"required,dive"`
}

type SplitExerciseDTO struct {
	ExerciseID   string  `json:"exercise_id" validate:"required,uuid4"`
	TargetSets   int     `json:"target_sets" validate:"required,gte=1"`
	TargetReps   int     `json:"target_reps" validate:"required,gte=1"`
	TargetWeight float64 `json:"target_weight"`
	Notes        string  `json:"notes"`
}

type UpdateSplitTemplateDTO struct {
	Name        string        `json:"name" validate:"required"`
	Description string        `json:"description"`
	FocusMuscle string        `json:"focus_muscle" validate:"required"`
	IsActive    bool          `json:"is_active"`
	Days        []SplitDayDTO `json:"days" validate:"required,dive"`
}
