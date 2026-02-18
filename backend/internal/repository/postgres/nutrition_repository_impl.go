package postgres

import (
	"context"
	"database/sql"

	domainerr "S.P.A.R.T.A/backend/internal/domain/errors"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/nutrition"
	domainrepo "S.P.A.R.T.A/backend/internal/domain/repository"
)

type nutritionRepository struct {
	db DBTX
}

func NewNutritionRepository(db DBTX) domainrepo.NutritionRepository {
	return &nutritionRepository{db: db}
}

func (r *nutritionRepository) SaveDaily(ctx context.Context, n *nutrition.DailyNutrition) error {
	if n == nil || n.ID == "" || n.UserID == "" {
		return domainerr.ErrInvalidInput
	}

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO daily_nutritions(id,user_id,date,protein_grams,calories,notes)
		 VALUES ($1,$2,$3,$4,$5,$6)
		 ON CONFLICT (user_id, date)
		 DO UPDATE SET protein_grams=EXCLUDED.protein_grams, calories=EXCLUDED.calories, notes=EXCLUDED.notes`,
		n.ID, n.UserID, n.Date, n.ProteinGrams, n.Calories, n.Notes,
	)
	if err != nil {
		return domainerr.ErrInternal
	}
	return nil
}

func (r *nutritionRepository) GetByDate(ctx context.Context, userID string, date string) (*nutrition.DailyNutrition, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id,user_id,date,protein_grams,calories,COALESCE(notes,'')
		 FROM daily_nutritions
		 WHERE user_id=$1 AND date=$2`,
		userID,
		date,
	)

	var out nutrition.DailyNutrition
	if err := row.Scan(&out.ID, &out.UserID, &out.Date, &out.ProteinGrams, &out.Calories, &out.Notes); err != nil {
		if err == sql.ErrNoRows {
			return nil, domainerr.ErrNotFound
		}
		return nil, domainerr.ErrInternal
	}
	return &out, nil
}
