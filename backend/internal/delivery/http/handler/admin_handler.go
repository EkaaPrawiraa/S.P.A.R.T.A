package handler

import (
	"time"

	"S.P.A.R.T.A/backend/internal/delivery/http/dto"
	"S.P.A.R.T.A/backend/internal/delivery/http/response"
	domainuc "S.P.A.R.T.A/backend/internal/domain/usecase"
	"S.P.A.R.T.A/backend/pkg/validator"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	uc domainuc.AdminUsecase
}

func NewAdminHandler(uc domainuc.AdminUsecase) *AdminHandler {
	return &AdminHandler{uc: uc}
}

func (h *AdminHandler) CreateInvite(c *gin.Context) {
	var req dto.CreateAdminInviteRequestDTO
	// Body is optional; default expiry if missing.
	_ = c.ShouldBindJSON(&req)
	if req.ExpiresInHours != 0 {
		if err := validator.ValidateStruct(&req); err != nil {
			response.BadRequest(c, err.Error())
			return
		}
	}

	createdBy := c.GetString("user_id")
	expiresIn := 24 * time.Hour
	if req.ExpiresInHours > 0 {
		expiresIn = time.Duration(req.ExpiresInHours) * time.Hour
	}

	res, err := h.uc.CreateAdminInvite(c.Request.Context(), createdBy, expiresIn)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, res)
}
