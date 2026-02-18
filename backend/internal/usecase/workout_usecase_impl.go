package usecase

import (
	"context"
	"errors"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/workout"
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
		return errors.New("user id required")
	}

	if len(session.Exercises) == 0 {
		return errors.New("workout must contain exercises")
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
