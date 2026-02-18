package handler

import (
	"S.P.A.R.T.A/backend/internal/delivery/http/dto"
	"S.P.A.R.T.A/backend/internal/delivery/http/middleware"
	"S.P.A.R.T.A/backend/internal/delivery/http/response"
	domainerr "S.P.A.R.T.A/backend/internal/domain/errors"
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
	authedUserID := middleware.GetUserID(c)
	if authedUserID != "" && userID != authedUserID {
		response.Error(c, domainerr.ErrForbidden)
		return
	}

	res, err := h.uc.GenerateRecommendation(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, dto.FromDomainPlannerRecommendation(*res))
}

func (h *PlannerHandler) GetUserRecommendations(c *gin.Context) {
	userID := c.Param("user_id")
	authedUserID := middleware.GetUserID(c)
	if authedUserID != "" && userID != authedUserID {
		response.Error(c, domainerr.ErrForbidden)
		return
	}

	res, err := h.uc.GetUserRecommendations(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, dto.FromDomainPlannerRecommendations(res))
}
