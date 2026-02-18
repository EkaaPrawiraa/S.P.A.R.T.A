package dto

import (
	"S.P.A.R.T.A/backend/internal/domain/aggregate/nutrition"
)

type DailyNutritionResponseDTO struct {
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	Date         string `json:"date"`
	ProteinGrams int    `json:"protein_grams"`
	Calories     int    `json:"calories"`
	Notes        string `json:"notes"`
}

func FromDomainDailyNutrition(n nutrition.DailyNutrition) DailyNutritionResponseDTO {
	return DailyNutritionResponseDTO{
		ID:           n.ID,
		UserID:       n.UserID,
		Date:         n.Date.Format("2006-01-02"),
		ProteinGrams: n.ProteinGrams,
		Calories:     n.Calories,
		Notes:        n.Notes,
	}
}
