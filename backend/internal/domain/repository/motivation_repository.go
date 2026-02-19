package repository

import (
	"context"
	"time"
)

type MotivationRepository interface {
	GetDailyMotivation(ctx context.Context, userID string, date time.Time) (message string, found bool, err error)
	SetDailyMotivation(ctx context.Context, userID string, date time.Time, message string, ttl time.Duration) error
	DeleteDailyMotivation(ctx context.Context, userID string, date time.Time) error
}
