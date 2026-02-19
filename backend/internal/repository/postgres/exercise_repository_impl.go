package postgres

import (
	"context"
	"database/sql"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/exercise"
	domainerr "S.P.A.R.T.A/backend/internal/domain/errors"
	domainrepo "S.P.A.R.T.A/backend/internal/domain/repository"

	"github.com/lib/pq"
)

type exerciseRepository struct {
	db DBTX
}

func NewExerciseRepository(db DBTX) domainrepo.ExerciseRepository {
	return &exerciseRepository{db: db}
}

func (r *exerciseRepository) Create(ctx context.Context, ex *exercise.Exercise) error {
	if ex == nil || ex.ID == "" || ex.Name == "" {
		return domainerr.ErrInvalidInput
	}

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO exercises(id,name,primary_muscle,secondary_muscles,equipment,created_at)
		 VALUES ($1,$2,$3,$4,$5,$6)`,
		ex.ID, ex.Name, ex.PrimaryMuscle, pq.Array(ex.SecondaryMuscles), ex.Equipment, ex.CreatedAt)
	if err != nil {
		return domainerr.ErrInternal
	}
	return nil
}

func (r *exerciseRepository) GetByID(ctx context.Context, id string) (*exercise.Exercise, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id,name,primary_muscle,secondary_muscles,equipment,created_at
		 FROM exercises
		 WHERE id=$1`,
		id,
	)

	var out exercise.Exercise
	var secondary pq.StringArray
	if err := row.Scan(&out.ID, &out.Name, &out.PrimaryMuscle, &secondary, &out.Equipment, &out.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, domainerr.ErrNotFound
		}
		return nil, domainerr.ErrInternal
	}
	out.SecondaryMuscles = []string(secondary)

	mediaRows, err := r.db.QueryContext(ctx,
		`SELECT id, media_type, media_url, thumbnail_url, created_at
		 FROM exercise_media
		 WHERE exercise_id=$1
		 ORDER BY created_at DESC`,
		out.ID,
	)
	if err != nil {
		return nil, domainerr.ErrInternal
	}
	defer mediaRows.Close()

	for mediaRows.Next() {
		var m exercise.ExerciseMedia
		var thumb sql.NullString
		if err := mediaRows.Scan(&m.ID, &m.MediaType, &m.MediaURL, &thumb, &m.CreatedAt); err != nil {
			return nil, domainerr.ErrInternal
		}
		if thumb.Valid {
			s := thumb.String
			m.ThumbnailURL = &s
		}
		out.Media = append(out.Media, m)
	}
	if err := mediaRows.Err(); err != nil {
		return nil, domainerr.ErrInternal
	}

	return &out, nil
}

func (r *exerciseRepository) List(ctx context.Context) ([]exercise.Exercise, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id,name,primary_muscle,secondary_muscles,equipment,created_at
		 FROM exercises
		 ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, domainerr.ErrInternal
	}
	defer rows.Close()

	items := make([]exercise.Exercise, 0)
	exerciseIDs := make([]string, 0)
	for rows.Next() {
		var ex exercise.Exercise
		var secondary pq.StringArray
		if err := rows.Scan(&ex.ID, &ex.Name, &ex.PrimaryMuscle, &secondary, &ex.Equipment, &ex.CreatedAt); err != nil {
			return nil, domainerr.ErrInternal
		}
		ex.SecondaryMuscles = []string(secondary)
		items = append(items, ex)
		exerciseIDs = append(exerciseIDs, ex.ID)
	}
	if err := rows.Err(); err != nil {
		return nil, domainerr.ErrInternal
	}

	if len(items) == 0 {
		return items, nil
	}

	mediaRows, err := r.db.QueryContext(ctx,
		`SELECT id, exercise_id, media_type, media_url, thumbnail_url, created_at
		 FROM exercise_media
		 WHERE exercise_id = ANY($1)
		 ORDER BY exercise_id, created_at DESC`,
		pq.Array(exerciseIDs),
	)
	if err != nil {
		return nil, domainerr.ErrInternal
	}
	defer mediaRows.Close()

	mediaByExerciseID := make(map[string][]exercise.ExerciseMedia)
	for mediaRows.Next() {
		var m exercise.ExerciseMedia
		var thumb sql.NullString
		if err := mediaRows.Scan(&m.ID, &m.ExerciseID, &m.MediaType, &m.MediaURL, &thumb, &m.CreatedAt); err != nil {
			return nil, domainerr.ErrInternal
		}
		if thumb.Valid {
			s := thumb.String
			m.ThumbnailURL = &s
		}
		mediaByExerciseID[m.ExerciseID] = append(mediaByExerciseID[m.ExerciseID], m)
	}
	if err := mediaRows.Err(); err != nil {
		return nil, domainerr.ErrInternal
	}

	for i := range items {
		if media, ok := mediaByExerciseID[items[i].ID]; ok {
			items[i].Media = media
		}
	}
	return items, nil
}

func (r *exerciseRepository) AddMedia(ctx context.Context, media *exercise.ExerciseMedia) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO exercise_media(id,exercise_id,media_type,media_url,thumbnail_url,created_at)
		 VALUES ($1,$2,$3,$4,$5,$6)`,
		media.ID, media.ExerciseID, media.MediaType, media.MediaURL, media.ThumbnailURL, media.CreatedAt)
	if err != nil {
		return domainerr.ErrInternal
	}
	return nil
}
