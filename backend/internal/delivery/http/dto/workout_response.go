package dto

import (
	"time"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/workout"
)

type WorkoutSetResponseDTO struct {
	ID        string    `json:"id"`
	SetOrder  int       `json:"set_order"`
	Reps      int       `json:"reps"`
	Weight    float64   `json:"weight"`
	RPE       float64   `json:"rpe"`
	SetType   string    `json:"set_type"`
	CreatedAt time.Time `json:"created_at"`
}

type WorkoutExerciseResponseDTO struct {
	ID         string                  `json:"id"`
	ExerciseID string                  `json:"exercise_id"`
	Sets       []WorkoutSetResponseDTO `json:"sets"`
}

type WorkoutSessionResponseDTO struct {
	ID              string                       `json:"id"`
	UserID          string                       `json:"user_id"`
	SplitDayID      *string                      `json:"split_day_id,omitempty"`
	SessionDate     string                       `json:"session_date"`
	DurationMinutes int                          `json:"duration_minutes"`
	Notes           string                       `json:"notes"`
	Exercises       []WorkoutExerciseResponseDTO `json:"exercises"`
	CreatedAt       time.Time                    `json:"created_at"`
}

func FromDomainWorkoutSession(s workout.WorkoutSession) WorkoutSessionResponseDTO {
	out := WorkoutSessionResponseDTO{
		ID:              s.ID,
		UserID:          s.UserID,
		SplitDayID:      s.SplitDayID,
		SessionDate:     s.SessionDate.Format("2006-01-02"),
		DurationMinutes: s.DurationMin,
		Notes:           s.Notes,
		CreatedAt:       s.CreatedAt,
	}

	for _, ex := range s.Exercises {
		exOut := WorkoutExerciseResponseDTO{
			ID:         ex.ID,
			ExerciseID: ex.ExerciseID,
		}
		for _, set := range ex.Sets {
			exOut.Sets = append(exOut.Sets, WorkoutSetResponseDTO{
				ID:        set.ID,
				SetOrder:  set.SetOrder,
				Reps:      set.Reps,
				Weight:    set.Weight,
				RPE:       set.RPE,
				SetType:   set.SetType,
				CreatedAt: set.CreatedAt,
			})
		}
		out.Exercises = append(out.Exercises, exOut)
	}

	return out
}

func FromDomainWorkoutSessions(items []workout.WorkoutSession) []WorkoutSessionResponseDTO {
	out := make([]WorkoutSessionResponseDTO, 0, len(items))
	for _, item := range items {
		out = append(out, FromDomainWorkoutSession(item))
	}
	return out
}
