package dto

import (
	"time"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/split"
)

type SplitExerciseResponseDTO struct {
	ExerciseID   string  `json:"exercise_id"`
	ExerciseName string  `json:"exercise_name"`
	TargetSets   int     `json:"target_sets"`
	TargetReps   int     `json:"target_reps"`
	TargetWeight float64 `json:"target_weight"`
	Notes        string  `json:"notes"`
}

type SplitDayResponseDTO struct {
	ID        string                     `json:"id"`
	DayOrder  int                        `json:"day_order"`
	Name      string                     `json:"name"`
	Exercises []SplitExerciseResponseDTO `json:"exercises"`
}

type SplitTemplateResponseDTO struct {
	ID          string                `json:"id"`
	UserID      string                `json:"user_id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	CreatedBy   string                `json:"created_by"`
	FocusMuscle string                `json:"focus_muscle"`
	IsActive    bool                  `json:"is_active"`
	Days        []SplitDayResponseDTO `json:"days"`
	CreatedAt   time.Time             `json:"created_at"`
}

func FromDomainSplitTemplate(t split.SplitTemplate) SplitTemplateResponseDTO {
	out := SplitTemplateResponseDTO{
		ID:          t.ID,
		UserID:      t.UserID,
		Name:        t.Name,
		Description: t.Description,
		CreatedBy:   t.CreatedBy,
		FocusMuscle: t.FocusMuscle,
		IsActive:    t.IsActive,
		CreatedAt:   t.CreatedAt,
	}

	for _, day := range t.Days {
		dayOut := SplitDayResponseDTO{
			ID:       day.ID,
			DayOrder: day.DayOrder,
			Name:     day.Name,
		}
		for _, ex := range day.Exercises {
			dayOut.Exercises = append(dayOut.Exercises, SplitExerciseResponseDTO{
				ExerciseID:   ex.ExerciseID,
				ExerciseName: ex.ExerciseName,
				TargetSets:   ex.TargetSets,
				TargetReps:   ex.TargetReps,
				TargetWeight: ex.TargetWeight,
				Notes:        ex.Notes,
			})
		}
		out.Days = append(out.Days, dayOut)
	}

	return out
}

func FromDomainSplitTemplates(items []split.SplitTemplate) []SplitTemplateResponseDTO {
	out := make([]SplitTemplateResponseDTO, 0, len(items))
	for _, item := range items {
		out = append(out, FromDomainSplitTemplate(item))
	}
	return out
}
