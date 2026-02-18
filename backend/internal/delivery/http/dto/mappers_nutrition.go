package dto

import (
	"time"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/nutrition"
	"github.com/google/uuid"
)

func ToDomainDailyNutrition(d UpsertDailyNutritionDTO) (nutrition.DailyNutrition, error) {
	parsedDate, err := time.Parse("2006-01-02", d.Date)
	if err != nil {
		return nutrition.DailyNutrition{}, err
	}

	return nutrition.DailyNutrition{
		ID:           uuid.NewString(),
		UserID:       d.UserID,
		Date:         parsedDate,
		ProteinGrams: d.ProteinGrams,
		Calories:     d.Calories,
		Notes:        d.Notes,
	}, nil
}
