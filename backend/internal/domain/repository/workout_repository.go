package repository

import (
	"S.P.A.R.T.A/backend/internal/domain/aggregate/workout"
	"context"
)

type WorkoutRepository interface {
	CreateSession(ctx context.Context, session *workout.WorkoutSession) error
	UpdateSession(ctx context.Context, session *workout.WorkoutSession) error
	GetSessionByID(ctx context.Context, id string) (*workout.WorkoutSession, error)
	GetSessionsByUser(ctx context.Context, userID string) ([]workout.WorkoutSession, error)
}
