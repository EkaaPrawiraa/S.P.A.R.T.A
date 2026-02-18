package postgres

import (
	"context"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/exercise"
	domainrepo "S.P.A.R.T.A/backend/internal/domain/repository"
)

type exerciseRepository struct {
	db DBTX
}

func NewExerciseRepository(db DBTX) domainrepo.ExerciseRepository {
	return &exerciseRepository{db: db}
}

func (r *exerciseRepository) Create(ctx context.Context, ex *exercise.Exercise) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO exercises(id,name,primary_muscle,secondary_muscles,equipment,created_at)
		 VALUES ($1,$2,$3,$4,$5,$6)`,
		ex.ID, ex.Name, ex.PrimaryMuscle, ex.SecondaryMuscles, ex.Equipment, ex.CreatedAt)
	return err
}

func (r *exerciseRepository) GetByID(ctx context.Context, id string) (*exercise.Exercise, error) {
	return nil, nil
}

func (r *exerciseRepository) List(ctx context.Context) ([]exercise.Exercise, error) {
	return nil, nil
}

func (r *exerciseRepository) AddMedia(ctx context.Context, media *exercise.ExerciseMedia) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO exercise_media(id,exercise_id,media_type,media_url,thumbnail_url,created_at)
		 VALUES ($1,$2,$3,$4,$5,$6)`,
		media.ID, media.MediaType, media.MediaURL, media.ThumbnailURL, media.CreatedAt)
	return err
}
