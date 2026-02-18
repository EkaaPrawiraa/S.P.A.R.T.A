package dto

import (
	"time"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/split"
	"github.com/google/uuid"
)

func ToDomainSplitTemplate(d CreateSplitTemplateDTO) split.SplitTemplate {
	tpl := split.SplitTemplate{
		ID:          uuid.NewString(),
		UserID:      d.UserID,
		Name:        d.Name,
		Description: d.Description,
		CreatedBy:   d.CreatedBy,
		FocusMuscle: d.FocusMuscle,
		IsActive:    d.IsActive,
		CreatedAt:   time.Now(),
	}

	for _, dayDTO := range d.Days {
		day := split.SplitDay{
			ID:       uuid.NewString(),
			DayOrder: dayDTO.DayOrder,
			Name:     dayDTO.Name,
		}

		for _, exDTO := range dayDTO.Exercises {
			day.Exercises = append(day.Exercises, split.SplitExercise{
				ExerciseID:   exDTO.ExerciseID,
				TargetSets:   exDTO.TargetSets,
				TargetReps:   exDTO.TargetReps,
				TargetWeight: exDTO.TargetWeight,
				Notes:        exDTO.Notes,
			})
		}

		tpl.Days = append(tpl.Days, day)
	}

	return tpl
}

func ToDomainSplitTemplateForUpdate(templateID string, userID string, d UpdateSplitTemplateDTO) split.SplitTemplate {
	tpl := split.SplitTemplate{
		ID:          templateID,
		UserID:      userID,
		Name:        d.Name,
		Description: d.Description,
		CreatedBy:   "user",
		FocusMuscle: d.FocusMuscle,
		IsActive:    d.IsActive,
		CreatedAt:   time.Now(),
	}

	for _, dayDTO := range d.Days {
		day := split.SplitDay{
			ID:       uuid.NewString(),
			DayOrder: dayDTO.DayOrder,
			Name:     dayDTO.Name,
		}

		for _, exDTO := range dayDTO.Exercises {
			day.Exercises = append(day.Exercises, split.SplitExercise{
				ExerciseID:   exDTO.ExerciseID,
				TargetSets:   exDTO.TargetSets,
				TargetReps:   exDTO.TargetReps,
				TargetWeight: exDTO.TargetWeight,
				Notes:        exDTO.Notes,
			})
		}

		tpl.Days = append(tpl.Days, day)
	}

	return tpl
}
