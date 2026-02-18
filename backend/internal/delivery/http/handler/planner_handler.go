package handler

import (
	"net/http"

	domainuc "S.P.A.R.T.A/backend/internal/domain/usecase"
	"github.com/gin-gonic/gin"
)

type PlannerHandler struct {
	uc domainuc.PlannerUsecase
}

func NewPlannerHandler(uc domainuc.PlannerUsecase) *PlannerHandler {
	return &PlannerHandler{uc: uc}
}

func (h *PlannerHandler) GenerateRecommendation(c *gin.Context) {
	userID := c.Param("user_id")

	res, err := h.uc.GenerateRecommendation(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *PlannerHandler) GetUserRecommendations(c *gin.Context) {
	userID := c.Param("user_id")

	res, err := h.uc.GetUserRecommendations(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
