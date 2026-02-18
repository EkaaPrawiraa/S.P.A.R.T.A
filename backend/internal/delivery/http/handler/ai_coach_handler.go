package handler

import (
	"net/http"

	domainuc "S.P.A.R.T.A/backend/internal/domain/usecase"
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

type GenerateSplitRequest struct {
	DaysPerWeek int    `json:"days_per_week"`
	FocusMuscle string `json:"focus_muscle"`
}

func (h *AICoachHandler) GenerateSplit(c *gin.Context) {
	userIDStr := c.GetString("user_id")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req GenerateSplitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.usecase.GenerateSplitTemplate(
		c.Request.Context(),
		userID,
		req.DaysPerWeek,
		req.FocusMuscle,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

type OverloadRequest struct {
	ExerciseID string `json:"exercise_id"`
}

func (h *AICoachHandler) SuggestOverload(c *gin.Context) {
	userIDStr := c.GetString("user_id")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req OverloadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exerciseID, err := uuid.Parse(req.ExerciseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid exercise id"})
		return
	}

	result, err := h.usecase.SuggestProgressiveOverload(
		c.Request.Context(),
		userID,
		exerciseID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
