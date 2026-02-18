package planner

import "time"

type PlannerRecommendation struct {
    ID                 string
    UserID             string
    WorkoutSessionID   *string
    Recommendation     string
    RecommendationType string
    CreatedAt          time.Time
}
