package dto

import (
	"time"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/workout"
	"github.com/google/uuid"
)

func ToDomainWorkoutSession(d CreateWorkoutSessionDTO) workout.WorkoutSession {
    sessionDate, _ := time.Parse("2006-01-02", d.SessionDate)

    ws := workout.WorkoutSession{
        ID:          uuid.NewString(),
        UserID:      d.UserID,
        SplitDayID:  d.SplitDayID,
        SessionDate: sessionDate,
        DurationMin: d.DurationMin,
        Notes:       d.Notes,
        CreatedAt:   time.Now(),
    }

    for _, exDTO := range d.Exercises {
        we := workout.WorkoutExercise{
            ID:         uuid.NewString(),
            ExerciseID: exDTO.ExerciseID,
        }

        for _, setDTO := range exDTO.Sets {
            we.Sets = append(we.Sets, workout.WorkoutSet{
                ID:        uuid.NewString(),
                SetOrder:  setDTO.SetOrder,
                Reps:      setDTO.Reps,
                Weight:    setDTO.Weight,
                RPE:       setDTO.RPE,
                SetType:   setDTO.SetType,
                CreatedAt: time.Now(),
            })
        }

        ws.Exercises = append(ws.Exercises, we)
    }

    return ws
}
