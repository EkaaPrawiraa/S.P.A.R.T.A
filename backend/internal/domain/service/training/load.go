package training

import (
	"math"
	"time"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/workout"
)

type LoadSummary struct {
	AcuteLoad7d     float64
	ChronicLoad28d  float64
	LastSessionLoad float64
	Sessions7d      int
	Sessions28d     int
	AvgRPE7d        float64
	ACWR            float64
}

func ComputeLoadSummary(sessions []workout.WorkoutSession, now time.Time) LoadSummary {
	sevenDaysAgo := now.AddDate(0, 0, -7)
	twentyEightDaysAgo := now.AddDate(0, 0, -28)

	var acute float64
	var chronic float64
	var lastLoad float64
	var sessions7 int
	var sessions28 int

	var rpeSum float64
	var rpeCount int

	for i, s := range sessions {
		load := sessionLoad(s)
		if i == 0 {
			lastLoad = load
		}

		if s.SessionDate.After(twentyEightDaysAgo) || s.SessionDate.Equal(twentyEightDaysAgo) {
			sessions28++
			chronic += load
		}
		if s.SessionDate.After(sevenDaysAgo) || s.SessionDate.Equal(sevenDaysAgo) {
			sessions7++
			acute += load
			avg, cnt := sessionAvgRPE(s)
			rpeSum += avg * float64(cnt)
			rpeCount += cnt
		}
	}

	avgRPE := 0.0
	if rpeCount > 0 {
		avgRPE = rpeSum / float64(rpeCount)
	}

	acwr := 0.0
	// ACWR proxy: acute / (chronic/4) (28d ~= 4 weeks)
	if chronic > 0 {
		acwr = acute / (chronic / 4.0)
	}

	return LoadSummary{
		AcuteLoad7d:     acute,
		ChronicLoad28d:  chronic,
		LastSessionLoad: lastLoad,
		Sessions7d:      sessions7,
		Sessions28d:     sessions28,
		AvgRPE7d:        avgRPE,
		ACWR:            acwr,
	}
}

// EstimateFatigueScore returns a 0..10 fatigue estimate based on load and perceived effort.
// This is a heuristic starter (Phase 2 foundation), not a clinical model.
func EstimateFatigueScore(sum LoadSummary, userReported int) int {
	score := float64(userReported)

	if sum.Sessions7d == 0 {
		// If user hasn't trained recently but reports fatigue, trust it.
		return clampInt(userReported, 0, 10)
	}

	// Load ratio effects.
	switch {
	case sum.ACWR >= 1.6:
		score += 2
	case sum.ACWR >= 1.3:
		score += 1
	case sum.ACWR > 0 && sum.ACWR <= 0.8:
		score -= 1
	}

	// Perceived effort (RPE).
	if sum.AvgRPE7d >= 8.5 {
		score += 1
	} else if sum.AvgRPE7d > 0 && sum.AvgRPE7d <= 6.5 {
		score -= 1
	}

	return clampInt(int(math.Round(score)), 0, 10)
}

func sessionLoad(s workout.WorkoutSession) float64 {
	var total float64
	for _, ex := range s.Exercises {
		for _, set := range ex.Sets {
			if set.Reps <= 0 {
				continue
			}
			if set.Weight <= 0 {
				continue
			}
			total += float64(set.Reps) * set.Weight
		}
	}
	return total
}

func sessionAvgRPE(s workout.WorkoutSession) (avg float64, count int) {
	var sum float64
	for _, ex := range s.Exercises {
		for _, set := range ex.Sets {
			if set.RPE <= 0 {
				continue
			}
			sum += set.RPE
			count++
		}
	}
	if count == 0 {
		return 0, 0
	}
	return sum / float64(count), count
}

func clampInt(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
