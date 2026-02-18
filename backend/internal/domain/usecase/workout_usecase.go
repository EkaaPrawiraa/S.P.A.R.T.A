package usecase

import (
	"context"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/workout"
)

type WorkoutUsecase interface {
    CreateWorkoutSession(ctx context.Context, session *workout.WorkoutSession) error
    GetWorkoutSession(ctx context.Context, id string) (*workout.WorkoutSession, error)
    GetUserWorkoutSessions(ctx context.Context, userID string) ([]workout.WorkoutSession, error)
}
