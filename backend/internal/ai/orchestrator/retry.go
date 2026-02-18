package orchestrator

import (
	"context"
	"time"
)

type RetryConfig struct {
	MaxAttempts int
	Delay       time.Duration
}

func WithRetry(ctx context.Context, cfg RetryConfig, fn func(context.Context) error) error {
	if cfg.MaxAttempts <= 0 {
		cfg.MaxAttempts = 1
	}
	if cfg.Delay <= 0 {
		cfg.Delay = 200 * time.Millisecond
	}

	var lastErr error
	for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
		if err := fn(ctx); err != nil {
			lastErr = err
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(cfg.Delay):
				continue
			}
		}
		return nil
	}
	return lastErr
}

