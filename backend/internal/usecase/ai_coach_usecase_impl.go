package usecase

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	"S.P.A.R.T.A/backend/internal/ai/orchestrator"
	"S.P.A.R.T.A/backend/internal/domain/aggregate/planner"
	"S.P.A.R.T.A/backend/internal/domain/aggregate/split"
	"S.P.A.R.T.A/backend/internal/domain/aggregate/workout"
	domainerr "S.P.A.R.T.A/backend/internal/domain/errors"
	domainrepo "S.P.A.R.T.A/backend/internal/domain/repository"
	"S.P.A.R.T.A/backend/internal/domain/service/training"
	domainuc "S.P.A.R.T.A/backend/internal/domain/usecase"
	"github.com/google/uuid"
)

type aiCoachUsecase struct {
	orchestrator        orchestrator.Orchestrator
	splitRepository     domainrepo.SplitRepository
	exerciseRepository  domainrepo.ExerciseRepository
	plannerRepository   domainrepo.PlannerRepository
	workoutRepository   domainrepo.WorkoutRepository
	nutritionRepository domainrepo.NutritionRepository
	motivationRepo      domainrepo.MotivationRepository
}

func NewAICoachUsecase(
	orchestrator orchestrator.Orchestrator,
	splitRepository domainrepo.SplitRepository,
	exerciseRepository domainrepo.ExerciseRepository,
	plannerRepository domainrepo.PlannerRepository,
	workoutRepository domainrepo.WorkoutRepository,
	nutritionRepository domainrepo.NutritionRepository,
	motivationRepository domainrepo.MotivationRepository,
) domainuc.AICoachUsecase {
	return &aiCoachUsecase{
		orchestrator:        orchestrator,
		splitRepository:     splitRepository,
		exerciseRepository:  exerciseRepository,
		plannerRepository:   plannerRepository,
		workoutRepository:   workoutRepository,
		nutritionRepository: nutritionRepository,
		motivationRepo:      motivationRepository,
	}
}

func (u *aiCoachUsecase) GenerateSplitTemplate(
	ctx context.Context,
	userID uuid.UUID,
	daysPerWeek int,
	focusMuscle string,
) (*split.SplitTemplate, error) {

	aiResult, err := u.orchestrator.GenerateSplit(ctx, orchestrator.SplitInput{
		UserID:      userID.String(),
		DaysPerWeek: daysPerWeek,
		FocusMuscle: focusMuscle,
	})
	if err != nil {
		return nil, err
	}

	template := &split.SplitTemplate{
		ID:          uuid.NewString(),
		UserID:      userID.String(),
		Name:        aiResult.Name,
		Description: "AI Generated Split",
		CreatedBy:   "ai",
		FocusMuscle: focusMuscle,
		IsActive:    false,
		CreatedAt:   time.Now(),
	}

	for i, d := range aiResult.Days {
		day := split.SplitDay{
			ID:       uuid.NewString(),
			DayOrder: i + 1,
			Name:     strings.TrimSpace(d.DayName),
		}
		if day.Name == "" {
			day.Name = fmt.Sprintf("Day %d", day.DayOrder)
		}

		for _, ex := range d.Exercises {
			name := strings.TrimSpace(ex.Name)
			if name == "" {
				continue
			}
			day.Exercises = append(day.Exercises, split.SplitExercise{
				TargetSets: ex.Sets,
				TargetReps: parseRepRange(ex.RepRange),
				Notes:      name,
			})
		}

		template.Days = append(template.Days, day)
	}

	// Best-effort: map AI exercise names to real exercise IDs in our library.
	// This makes the template immediately usable by downstream endpoints.
	u.enrichSplitTemplateExerciseIDs(ctx, template)

	err = u.splitRepository.CreateTemplate(ctx, template)
	if err != nil {
		return nil, err
	}

	return template, nil
}

func (u *aiCoachUsecase) enrichSplitTemplateExerciseIDs(ctx context.Context, tpl *split.SplitTemplate) {
	if tpl == nil || u.exerciseRepository == nil {
		return
	}

	exs, err := u.exerciseRepository.List(ctx)
	if err != nil {
		return
	}

	index := make(map[string]string, len(exs))
	for _, ex := range exs {
		key := normalizeExerciseName(ex.Name)
		if key == "" {
			continue
		}
		if _, exists := index[key]; !exists {
			index[key] = ex.ID
		}
	}

	for di := range tpl.Days {
		for ei := range tpl.Days[di].Exercises {
			ex := &tpl.Days[di].Exercises[ei]
			if strings.TrimSpace(ex.ExerciseID) != "" {
				continue
			}
			name := strings.TrimSpace(ex.Notes)
			if name == "" {
				continue
			}
			if id, ok := index[normalizeExerciseName(name)]; ok {
				ex.ExerciseID = id
			}
		}
	}
}

func normalizeExerciseName(name string) string {
	name = strings.TrimSpace(strings.ToLower(name))
	if name == "" {
		return ""
	}

	var b strings.Builder
	b.Grow(len(name))
	lastSpace := false
	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			b.WriteRune(r)
			lastSpace = false
			continue
		}
		if !lastSpace {
			b.WriteByte(' ')
			lastSpace = true
		}
	}

	out := strings.TrimSpace(b.String())
	return out
}

func (u *aiCoachUsecase) SuggestProgressiveOverload(
	ctx context.Context,
	userID uuid.UUID,
	exerciseID uuid.UUID,
) (*planner.PlannerRecommendation, error) {
	lastWeight, lastReps, perfNotes := u.findLastExercisePerformance(ctx, userID.String(), exerciseID.String())

	result, err := u.orchestrator.SuggestOverload(ctx, orchestrator.OverloadInput{
		UserID:      userID.String(),
		ExerciseID:  exerciseID.String(),
		LastWeight:  lastWeight,
		LastReps:    lastReps,
		Performance: perfNotes,
	})
	if err != nil {
		return nil, err
	}

	rec := &planner.PlannerRecommendation{
		ID:                 uuid.NewString(),
		UserID:             userID.String(),
		Recommendation:     result.Message,
		RecommendationType: "progressive_overload",
		CreatedAt:          time.Now(),
	}

	err = u.plannerRepository.SaveRecommendation(ctx, rec)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

func (u *aiCoachUsecase) GetDailyMotivation(ctx context.Context, userID uuid.UUID) (string, error) {
	now := time.Now().UTC()
	dateStr := now.Format("2006-01-02")

	if u.motivationRepo != nil {
		msg, found, err := u.motivationRepo.GetDailyMotivation(ctx, userID.String(), now)
		if err != nil {
			return "", err
		}
		if found {
			return msg, nil
		}
	}

	sessions, err := u.workoutRepository.GetSessionsByUser(ctx, userID.String())
	if err != nil {
		return "", err
	}

	last7Days := 0
	sevenDaysAgo := now.AddDate(0, 0, -7)
	for _, s := range sessions {
		if s.SessionDate.After(sevenDaysAgo) || s.SessionDate.Equal(sevenDaysAgo) {
			last7Days++
		}
	}

	lastWorkoutDate := ""
	lastWorkoutDuration := 0
	lastWorkoutNotes := ""
	if len(sessions) > 0 {
		lastWorkoutDate = sessions[0].SessionDate.Format("2006-01-02")
		lastWorkoutDuration = sessions[0].DurationMin
		lastWorkoutNotes = strings.TrimSpace(sessions[0].Notes)
	}

	aiOut, err := u.orchestrator.GenerateDailyMotivation(ctx, orchestrator.MotivationInput{
		UserID:              userID.String(),
		Date:                dateStr,
		WorkoutsLast7Days:   last7Days,
		LastWorkoutDate:     lastWorkoutDate,
		LastWorkoutDuration: lastWorkoutDuration,
		LastWorkoutNotes:    lastWorkoutNotes,
	})
	if err != nil {
		return "", err
	}

	msg := strings.TrimSpace(aiOut.Message)
	if msg == "" {
		return "", domainerr.ErrInternal
	}

	if u.motivationRepo != nil {
		next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)
		ttl := next.Sub(now)
		if ttl < 10*time.Second {
			ttl = 24 * time.Hour
		}
		_ = u.motivationRepo.SetDailyMotivation(ctx, userID.String(), now, msg, ttl)
	}

	return msg, nil
}

func (u *aiCoachUsecase) ResetDailyMotivation(ctx context.Context, userID uuid.UUID) error {
	if u.motivationRepo == nil {
		return nil
	}
	now := time.Now().UTC()
	return u.motivationRepo.DeleteDailyMotivation(ctx, userID.String(), now)
}

func (u *aiCoachUsecase) GenerateWorkoutPlan(
	ctx context.Context,
	userID uuid.UUID,
	splitDayID uuid.UUID,
	fatigue int,
) (*workout.WorkoutPlan, error) {
	if fatigue < 0 || fatigue > 10 {
		return nil, domainerr.ErrInvalidInput
	}

	sessions, err := u.workoutRepository.GetSessionsByUser(ctx, userID.String())
	if err != nil {
		return nil, err
	}

	lastVolume := estimateLastVolume(sessions)
	now := time.Now().UTC()
	loadSum := training.ComputeLoadSummary(sessions, now)
	fatigueEstimated := training.EstimateFatigueScore(loadSum, fatigue)

	splitDay, err := u.splitRepository.GetSplitDayByID(ctx, splitDayID.String())
	if err != nil {
		return nil, err
	}

	plannedExercises := make([]string, 0)
	for _, ex := range splitDay.Exercises {
		if ex.ExerciseID != "" && u.exerciseRepository != nil {
			dbEx, err := u.exerciseRepository.GetByID(ctx, ex.ExerciseID)
			if err != nil {
				if errors.Is(err, domainerr.ErrNotFound) {
					continue
				}
				return nil, err
			}
			name := strings.TrimSpace(dbEx.Name)
			if name != "" {
				plannedExercises = append(plannedExercises, name)
			}
			continue
		}

		name := strings.TrimSpace(ex.Notes)
		if name != "" {
			plannedExercises = append(plannedExercises, name)
		}
	}

	aiOut, err := u.orchestrator.GenerateWorkout(ctx, orchestrator.WorkoutInput{
		UserID:           userID.String(),
		SplitDayID:       splitDayID.String(),
		SplitDayName:     splitDay.Name,
		PlannedExercises: plannedExercises,
		Fatigue:          fatigue,
		AcuteLoad7d:      loadSum.AcuteLoad7d,
		ChronicLoad28d:   loadSum.ChronicLoad28d,
		ACWR:             loadSum.ACWR,
		FatigueEstimated: fatigueEstimated,
		LastVolume:       lastVolume,
	})
	if err != nil {
		return nil, err
	}

	plan := &workout.WorkoutPlan{
		UserID:     userID.String(),
		SplitDayID: splitDayID.String(),
		Date:       time.Now().UTC(),
	}

	for _, ex := range aiOut.Exercises {
		plan.Exercises = append(plan.Exercises, workout.WorkoutPlanExercise{
			Name:     ex.Name,
			Sets:     ex.Sets,
			RepRange: ex.RepRange,
			Weight:   ex.Weight,
		})
	}

	return plan, nil
}

func parseRepRange(repRange string) int {
	repRange = strings.TrimSpace(repRange)
	if repRange == "" {
		return 10
	}

	// Formats we accept:
	// - "8-12" -> avg
	// - "8–12" (en dash)
	// - "10" -> exact
	// - "8 to 12" -> best-effort
	clean := strings.ToLower(repRange)
	clean = strings.ReplaceAll(clean, "–", "-")
	clean = strings.ReplaceAll(clean, "to", "-")
	clean = strings.ReplaceAll(clean, " ", "")

	if strings.Contains(clean, "-") {
		parts := strings.Split(clean, "-")
		if len(parts) >= 2 {
			lo, errLo := strconv.Atoi(parts[0])
			hi, errHi := strconv.Atoi(parts[1])
			if errLo == nil && errHi == nil {
				if lo <= 0 || hi <= 0 {
					return 10
				}
				if lo > hi {
					lo, hi = hi, lo
				}
				return (lo + hi) / 2
			}
		}
	}

	if n, err := strconv.Atoi(clean); err == nil {
		if n <= 0 {
			return 10
		}
		return n
	}

	return 10
}

func estimateLastVolume(sessions []workout.WorkoutSession) int {
	if len(sessions) == 0 {
		return 0
	}
	// Simple heuristic: total reps in the most recent session.
	vol := 0
	for _, ex := range sessions[0].Exercises {
		for _, s := range ex.Sets {
			vol += s.Reps
		}
	}
	return vol
}

func (u *aiCoachUsecase) findLastExercisePerformance(ctx context.Context, userID, exerciseID string) (lastWeight float64, lastReps int, notes string) {
	sessions, err := u.workoutRepository.GetSessionsByUser(ctx, userID)
	if err != nil {
		return 0, 0, ""
	}

	for _, sess := range sessions {
		for _, ex := range sess.Exercises {
			if ex.ExerciseID != exerciseID {
				continue
			}
			// pick the last set (highest order) as proxy for top set.
			bestOrder := -1
			for _, set := range ex.Sets {
				if set.SetOrder > bestOrder {
					bestOrder = set.SetOrder
					lastWeight = set.Weight
					lastReps = set.Reps
				}
			}
			notes = strings.TrimSpace(sess.Notes)
			return lastWeight, lastReps, notes
		}
	}

	return 0, 0, ""
}

func (u *aiCoachUsecase) GetCoachingSuggestions(ctx context.Context, userID uuid.UUID) ([]string, error) {
	now := time.Now().UTC()
	dateStr := now.Format("2006-01-02")

	sessions, err := u.workoutRepository.GetSessionsByUser(ctx, userID.String())
	if err != nil {
		return nil, err
	}
	loadSum := training.ComputeLoadSummary(sessions, now)

	workoutsSummary := ""
	if len(sessions) == 0 {
		workoutsSummary = "(no workouts logged)"
	} else {
		last := sessions[0]
		workoutsSummary = strings.TrimSpace(last.SessionDate.Format("2006-01-02") + " | " + last.Notes)
	}

	nutritionSummary := ""
	if u.nutritionRepository != nil {
		n, err := u.nutritionRepository.GetByDate(ctx, userID.String(), dateStr)
		if err != nil {
			if !errors.Is(err, domainerr.ErrNotFound) {
				return nil, err
			}
		} else if n != nil {
			nutritionSummary = "protein=" + strconv.Itoa(n.ProteinGrams) + "g calories=" + strconv.Itoa(n.Calories) + " notes=" + strings.TrimSpace(n.Notes)
		}
	}
	if strings.TrimSpace(nutritionSummary) == "" {
		nutritionSummary = "(no nutrition log found for today)"
	}

	recs, err := u.plannerRepository.GetUserRecommendations(ctx, userID.String())
	if err != nil {
		return nil, err
	}
	recsSummary := ""
	if len(recs) > 0 {
		recsSummary = recs[0].RecommendationType + ": " + strings.TrimSpace(recs[0].Recommendation)
	}
	if strings.TrimSpace(recsSummary) == "" {
		recsSummary = "(no recommendations yet)"
	}

	out, err := u.orchestrator.GenerateCoachingSuggestions(ctx, orchestrator.CoachingInput{
		UserID:            userID.String(),
		Date:              dateStr,
		AcuteLoad7d:       loadSum.AcuteLoad7d,
		ChronicLoad28d:    loadSum.ChronicLoad28d,
		ACWR:              loadSum.ACWR,
		RecentWorkouts:    workoutsSummary,
		RecentNutrition:   nutritionSummary,
		RecentPlannerRecs: recsSummary,
	})
	if err != nil {
		return nil, err
	}

	return out.Suggestions, nil
}

func (u *aiCoachUsecase) ExplainWorkoutPlan(
	ctx context.Context,
	userID uuid.UUID,
	plan workout.WorkoutPlan,
	splitDayName string,
	fatigue int,
) (*workout.WorkoutExplanation, error) {
	ex := make([]orchestrator.ExplainWorkoutExercise, 0, len(plan.Exercises))
	for _, e := range plan.Exercises {
		ex = append(ex, orchestrator.ExplainWorkoutExercise{
			Name:     e.Name,
			Sets:     e.Sets,
			RepRange: e.RepRange,
			Weight:   e.Weight,
		})
	}

	out, err := u.orchestrator.ExplainWorkoutPlan(ctx, orchestrator.ExplainWorkoutPlanInput{
		UserID:       userID.String(),
		SplitDayName: splitDayName,
		Fatigue:      fatigue,
		Exercises:    ex,
	})
	if err != nil {
		return nil, err
	}

	resp := &workout.WorkoutExplanation{Summary: strings.TrimSpace(out.Summary)}
	for _, n := range out.ExerciseNotes {
		resp.ExerciseNotes = append(resp.ExerciseNotes, workout.WorkoutExerciseNote{
			Name: strings.TrimSpace(n.Name),
			Note: strings.TrimSpace(n.Note),
		})
	}
	return resp, nil
}
