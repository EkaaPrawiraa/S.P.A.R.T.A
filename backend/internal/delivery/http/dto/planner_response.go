package dto

import (
	"time"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/planner"
)

type PlannerRecommendationResponseDTO struct {
	ID                 string    `json:"id"`
	UserID             string    `json:"user_id"`
	WorkoutSessionID   *string   `json:"workout_session_id,omitempty"`
	Recommendation     string    `json:"recommendation"`
	RecommendationType string    `json:"recommendation_type"`
	CreatedAt          time.Time `json:"created_at"`
}

func FromDomainPlannerRecommendation(p planner.PlannerRecommendation) PlannerRecommendationResponseDTO {
	return PlannerRecommendationResponseDTO{
		ID:                 p.ID,
		UserID:             p.UserID,
		WorkoutSessionID:   p.WorkoutSessionID,
		Recommendation:     p.Recommendation,
		RecommendationType: p.RecommendationType,
		CreatedAt:          p.CreatedAt,
	}
}

func FromDomainPlannerRecommendations(items []planner.PlannerRecommendation) []PlannerRecommendationResponseDTO {
	out := make([]PlannerRecommendationResponseDTO, 0, len(items))
	for _, item := range items {
		out = append(out, FromDomainPlannerRecommendation(item))
	}
	return out
}
