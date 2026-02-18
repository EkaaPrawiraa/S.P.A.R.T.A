package handler

import (
	"net/http"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/nutrition"
	domainuc "S.P.A.R.T.A/backend/internal/domain/usecase"
	"github.com/gin-gonic/gin"
)

type NutritionHandler struct {
	uc domainuc.NutritionUsecase
}

func NewNutritionHandler(uc domainuc.NutritionUsecase) *NutritionHandler {
	return &NutritionHandler{uc: uc}
}

func (h *NutritionHandler) UpsertDailyNutrition(c *gin.Context) {
	var req nutrition.DailyNutrition

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.uc.UpsertDailyNutrition(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, req)
}

func (h *NutritionHandler) GetDailyNutrition(c *gin.Context) {
	userID := c.Param("user_id")
	date := c.Query("date")

	res, err := h.uc.GetDailyNutrition(c.Request.Context(), userID, date)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
