package handler

import (
	"net/http"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/workout"
	domainuc "S.P.A.R.T.A/backend/internal/domain/usecase"
	"github.com/gin-gonic/gin"
)

type WorkoutHandler struct {
    workoutUC domainuc.WorkoutUsecase
}

func NewWorkoutHandler(workoutUC domainuc.WorkoutUsecase) *WorkoutHandler {
    return &WorkoutHandler{workoutUC: workoutUC}
}

func (h *WorkoutHandler) CreateWorkoutSession(c *gin.Context) {
    var req workout.WorkoutSession

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.workoutUC.CreateWorkoutSession(c.Request.Context(), &req); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, req)
}

func (h *WorkoutHandler) GetWorkoutSession(c *gin.Context) {
    id := c.Param("id")

    result, err := h.workoutUC.GetWorkoutSession(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, result)
}

func (h *WorkoutHandler) GetUserWorkoutSessions(c *gin.Context) {
    userID := c.Param("user_id")

    result, err := h.workoutUC.GetUserWorkoutSessions(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, result)
}
