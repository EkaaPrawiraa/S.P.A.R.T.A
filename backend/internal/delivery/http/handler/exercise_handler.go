package handler

import (
	"net/http"

	domainuc "S.P.A.R.T.A/backend/internal/domain/usecase"
	"github.com/gin-gonic/gin"
)

type ExerciseHandler struct {
	uc domainuc.ExerciseUsecase
}

func NewExerciseHandler(uc domainuc.ExerciseUsecase) *ExerciseHandler {
	return &ExerciseHandler{uc: uc}
}

func (h *ExerciseHandler) ListExercises(c *gin.Context) {
	res, err := h.uc.ListExercises(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *ExerciseHandler) GetExercise(c *gin.Context) {
	id := c.Param("id")

	res, err := h.uc.GetExercise(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
