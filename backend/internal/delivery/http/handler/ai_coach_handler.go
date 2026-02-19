package handler

import (
	"time"

	"S.P.A.R.T.A/backend/internal/delivery/http/dto"
	"S.P.A.R.T.A/backend/internal/delivery/http/response"
	"S.P.A.R.T.A/backend/internal/domain/aggregate/workout"
	domainuc "S.P.A.R.T.A/backend/internal/domain/usecase"
	"S.P.A.R.T.A/backend/pkg/validator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AICoachHandler struct {
	usecase domainuc.AICoachUsecase
}

func NewAICoachHandler(u domainuc.AICoachUsecase) *AICoachHandler {
	return &AICoachHandler{
		usecase: u,
	}
}

func (h *AICoachHandler) GenerateSplit(c *gin.Context) {
	userIDStr := c.GetString("user_id")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	var req dto.GenerateSplitRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid body")
		return
	}

	if err := validator.ValidateStruct(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.usecase.GenerateSplitTemplate(
		c.Request.Context(),
		userID,
		req.DaysPerWeek,
		req.FocusMuscle,
	)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, dto.FromDomainSplitTemplate(*result))
}

func (h *AICoachHandler) SuggestOverload(c *gin.Context) {
	userIDStr := c.GetString("user_id")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	var req dto.SuggestOverloadRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid body")
		return
	}

	if err := validator.ValidateStruct(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	exerciseID, err := uuid.Parse(req.ExerciseID)
	if err != nil {
		response.BadRequest(c, "invalid exercise id")
		return
	}

	result, err := h.usecase.SuggestProgressiveOverload(
		c.Request.Context(),
		userID,
		exerciseID,
	)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, dto.FromDomainPlannerRecommendation(*result))
}

func (h *AICoachHandler) GetDailyMotivation(c *gin.Context) {
	userIDStr := c.GetString("user_id")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	msg, err := h.usecase.GetDailyMotivation(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, dto.DailyMotivationResponseDTO{
		Date:    time.Now().UTC().Format("2006-01-02"),
		Message: msg,
	})
}

func (h *AICoachHandler) ResetDailyMotivation(c *gin.Context) {
	userIDStr := c.GetString("user_id")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	if err := h.usecase.ResetDailyMotivation(c.Request.Context(), userID); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"reset": true})
}

func (h *AICoachHandler) GenerateWorkoutPlan(c *gin.Context) {
	userIDStr := c.GetString("user_id")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	var req dto.GenerateWorkoutPlanRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid body")
		return
	}

	if err := validator.ValidateStruct(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	splitDayID, err := uuid.Parse(req.SplitDayID)
	if err != nil {
		response.BadRequest(c, "invalid split day id")
		return
	}

	plan, err := h.usecase.GenerateWorkoutPlan(c.Request.Context(), userID, splitDayID, req.Fatigue)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, dto.FromDomainWorkoutPlan(*plan))
}

func (h *AICoachHandler) GetCoachingSuggestions(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	suggestions, err := h.usecase.GetCoachingSuggestions(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, dto.CoachingSuggestionsResponseDTO{
		Date:        time.Now().UTC().Format("2006-01-02"),
		Suggestions: suggestions,
	})
}

func (h *AICoachHandler) ExplainWorkoutPlan(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	var req dto.ExplainWorkoutPlanRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid body")
		return
	}
	if err := validator.ValidateStruct(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	plan := workout.WorkoutPlan{
		UserID: userID.String(),
		Date:   time.Now().UTC(),
	}
	for _, ex := range req.Exercises {
		plan.Exercises = append(plan.Exercises, workout.WorkoutPlanExercise{
			Name:     ex.Name,
			Sets:     ex.Sets,
			RepRange: ex.RepRange,
			Weight:   ex.Weight,
		})
	}

	expl, err := h.usecase.ExplainWorkoutPlan(c.Request.Context(), userID, plan, req.SplitDayName, req.Fatigue)
	if err != nil {
		response.Error(c, err)
		return
	}

	out := dto.WorkoutExplanationResponseDTO{Summary: expl.Summary}
	for _, n := range expl.ExerciseNotes {
		out.ExerciseNotes = append(out.ExerciseNotes, dto.WorkoutExplanationExerciseNoteDTO{Name: n.Name, Note: n.Note})
	}

	response.Success(c, out)
}
