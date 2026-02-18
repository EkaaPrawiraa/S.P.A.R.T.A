package usecase

import (
	"context"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/workout"
	domainerr "S.P.A.R.T.A/backend/internal/domain/errors"
	domainrepo "S.P.A.R.T.A/backend/internal/domain/repository"
	domainuc "S.P.A.R.T.A/backend/internal/domain/usecase"
)

type workoutUsecase struct {
	workoutRepo domainrepo.WorkoutRepository
}

func NewWorkoutUsecase(
	workoutRepo domainrepo.WorkoutRepository,
) domainuc.WorkoutUsecase {
	return &workoutUsecase{
		workoutRepo: workoutRepo,
	}
}

func (u *workoutUsecase) CreateWorkoutSession(
	ctx context.Context,
	session *workout.WorkoutSession,
) error {

	if session.UserID == "" {
		return domainerr.ErrInvalidInput
	}

	if len(session.Exercises) == 0 {
		return domainerr.ErrInvalidInput
	}

	return u.workoutRepo.CreateSession(ctx, session)
}

func (u *workoutUsecase) GetWorkoutSession(
	ctx context.Context,
	id string,
) (*workout.WorkoutSession, error) {
	return u.workoutRepo.GetSessionByID(ctx, id)
}

func (u *workoutUsecase) GetUserWorkoutSessions(
	ctx context.Context,
	userID string,
) ([]workout.WorkoutSession, error) {
	return u.workoutRepo.GetSessionsByUser(ctx, userID)
}
