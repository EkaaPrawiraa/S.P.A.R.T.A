package dto

type UpsertDailyNutritionDTO struct {
	UserID       string `json:"user_id" validate:"required,uuid4"`
	Date         string `json:"date" validate:"required"`
	ProteinGrams int    `json:"protein_grams" validate:"required,gte=0"`
	Calories     int    `json:"calories" validate:"gte=0"`
	Notes        string `json:"notes"`
}
