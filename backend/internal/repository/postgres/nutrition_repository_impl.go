package postgres

import (
	"context"

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
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO daily_nutritions(id,user_id,date,protein_grams,calories,notes)
		 VALUES ($1,$2,$3,$4,$5,$6)`,
		n.ID, n.UserID, n.Date, n.ProteinGrams, n.Calories, n.Notes)
	return err
}

func (r *nutritionRepository) GetByDate(ctx context.Context, userID string, date string) (*nutrition.DailyNutrition, error) {
	return nil, nil
}
