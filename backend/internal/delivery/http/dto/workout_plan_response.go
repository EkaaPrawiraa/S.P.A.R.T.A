package dto

import "S.P.A.R.T.A/backend/internal/domain/aggregate/workout"

type WorkoutPlanExerciseResponseDTO struct {
	Name     string  `json:"name"`
	Sets     int     `json:"sets"`
	RepRange string  `json:"rep_range"`
	Weight   float64 `json:"weight"`
}

type WorkoutPlanResponseDTO struct {
	UserID     string                           `json:"user_id"`
	SplitDayID string                           `json:"split_day_id"`
	Date       string                           `json:"date"`
	Exercises  []WorkoutPlanExerciseResponseDTO `json:"exercises"`
}

func FromDomainWorkoutPlan(p workout.WorkoutPlan) WorkoutPlanResponseDTO {
	out := WorkoutPlanResponseDTO{
		UserID:     p.UserID,
		SplitDayID: p.SplitDayID,
		Date:       p.Date.UTC().Format("2006-01-02"),
	}

	for _, ex := range p.Exercises {
		out.Exercises = append(out.Exercises, WorkoutPlanExerciseResponseDTO{
			Name:     ex.Name,
			Sets:     ex.Sets,
			RepRange: ex.RepRange,
			Weight:   ex.Weight,
		})
	}

	return out
}
