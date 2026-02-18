package handler

import (
	"S.P.A.R.T.A/backend/internal/delivery/http/dto"
	"S.P.A.R.T.A/backend/internal/delivery/http/middleware"
	"S.P.A.R.T.A/backend/internal/delivery/http/response"
	domainerr "S.P.A.R.T.A/backend/internal/domain/errors"
	domainuc "S.P.A.R.T.A/backend/internal/domain/usecase"
	"S.P.A.R.T.A/backend/pkg/validator"

	"github.com/gin-gonic/gin"
)

type WorkoutHandler struct {
	workoutUC domainuc.WorkoutUsecase
}

func NewWorkoutHandler(workoutUC domainuc.WorkoutUsecase) *WorkoutHandler {
	return &WorkoutHandler{workoutUC: workoutUC}
}

func (h *WorkoutHandler) CreateWorkoutSession(c *gin.Context) {
	var req dto.CreateWorkoutSessionDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid body")
		return
	}

	if err := validator.ValidateStruct(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	authedUserID := middleware.GetUserID(c)
	if authedUserID != "" && req.UserID != authedUserID {
		response.Error(c, domainerr.ErrForbidden)
		return
	}

	domainSession := dto.ToDomainWorkoutSession(req)

	if err := h.workoutUC.CreateWorkoutSession(c.Request.Context(), &domainSession); err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, dto.FromDomainWorkoutSession(domainSession))
}

func (h *WorkoutHandler) GetWorkoutSession(c *gin.Context) {
	id := c.Param("id")

	result, err := h.workoutUC.GetWorkoutSession(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, dto.FromDomainWorkoutSession(*result))
}

func (h *WorkoutHandler) GetUserWorkoutSessions(c *gin.Context) {
	userID := c.Param("user_id")
	authedUserID := middleware.GetUserID(c)
	if authedUserID != "" && userID != authedUserID {
		response.Error(c, domainerr.ErrForbidden)
		return
	}

	result, err := h.workoutUC.GetUserWorkoutSessions(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, dto.FromDomainWorkoutSessions(result))
}
