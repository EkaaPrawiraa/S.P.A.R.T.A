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

type NutritionHandler struct {
	uc domainuc.NutritionUsecase
}

func NewNutritionHandler(uc domainuc.NutritionUsecase) *NutritionHandler {
	return &NutritionHandler{uc: uc}
}

func (h *NutritionHandler) UpsertDailyNutrition(c *gin.Context) {
	var req dto.UpsertDailyNutritionDTO

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

	domainNutrition, err := dto.ToDomainDailyNutrition(req)
	if err != nil {
		response.BadRequest(c, "invalid date format (expected YYYY-MM-DD)")
		return
	}

	if err := h.uc.SaveDaily(c.Request.Context(), &domainNutrition); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, dto.FromDomainDailyNutrition(domainNutrition))
}

func (h *NutritionHandler) GetDailyNutrition(c *gin.Context) {
	userID := c.Param("user_id")
	authedUserID := middleware.GetUserID(c)
	if authedUserID != "" && userID != authedUserID {
		response.Error(c, domainerr.ErrForbidden)
		return
	}
	date := c.Query("date")

	res, err := h.uc.GetByDate(c.Request.Context(), userID, date)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, dto.FromDomainDailyNutrition(*res))
}
