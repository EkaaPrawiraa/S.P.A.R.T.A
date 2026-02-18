package nutrition

import "time"

type DailyNutrition struct {
	ID           string
	UserID       string
	Date         time.Time
	ProteinGrams int
	Calories     int
	Notes        string
}
