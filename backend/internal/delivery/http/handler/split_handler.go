package handler

import (
	"S.P.A.R.T.A/backend/internal/delivery/http/dto"
	"S.P.A.R.T.A/backend/internal/delivery/http/middleware"
	"S.P.A.R.T.A/backend/internal/delivery/http/response"
	domainerr "S.P.A.R.T.A/backend/internal/domain/errors"
	domainuc "S.P.A.R.T.A/backend/internal/domain/usecase"
	"S.P.A.R.T.A/backend/pkg/validator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SplitHandler struct {
	uc domainuc.SplitUsecase
}

func NewSplitHandler(uc domainuc.SplitUsecase) *SplitHandler {
	return &SplitHandler{uc: uc}
}

func (h *SplitHandler) CreateTemplate(c *gin.Context) {
	var req dto.CreateSplitTemplateDTO
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

	domainTpl := dto.ToDomainSplitTemplate(req)

	if err := h.uc.CreateTemplate(c.Request.Context(), &domainTpl); err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, dto.FromDomainSplitTemplate(domainTpl))
}

func (h *SplitHandler) GetUserTemplates(c *gin.Context) {
	userID := c.Param("user_id")
	authedUserID := middleware.GetUserID(c)
	if authedUserID != "" && userID != authedUserID {
		response.Error(c, domainerr.ErrForbidden)
		return
	}

	res, err := h.uc.GetUserTemplates(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, dto.FromDomainSplitTemplates(res))
}

func (h *SplitHandler) GetTemplate(c *gin.Context) {
	id := c.Param("id")

	tpl, err := h.uc.GetTemplate(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}

	authedUserID := middleware.GetUserID(c)
	if authedUserID != "" && tpl.UserID != authedUserID {
		response.Error(c, domainerr.ErrForbidden)
		return
	}

	response.Success(c, dto.FromDomainSplitTemplate(*tpl))
}

func (h *SplitHandler) UpdateTemplate(c *gin.Context) {
	templateID := c.Param("id")
	if _, err := uuid.Parse(templateID); err != nil {
		response.BadRequest(c, "invalid template id")
		return
	}

	userID := middleware.GetUserID(c)
	if userID == "" {
		response.Error(c, domainerr.ErrUnauthorized)
		return
	}

	var req dto.UpdateSplitTemplateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid body")
		return
	}
	if err := validator.ValidateStruct(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	domainTpl := dto.ToDomainSplitTemplateForUpdate(templateID, userID, req)
	if err := h.uc.UpdateTemplate(c.Request.Context(), &domainTpl); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, dto.FromDomainSplitTemplate(domainTpl))
}

func (h *SplitHandler) ActivateTemplate(c *gin.Context) {
	templateID := c.Param("id")
	if _, err := uuid.Parse(templateID); err != nil {
		response.BadRequest(c, "invalid template id")
		return
	}

	userID := middleware.GetUserID(c)
	if userID == "" {
		response.Error(c, domainerr.ErrUnauthorized)
		return
	}

	if err := h.uc.ActivateTemplate(c.Request.Context(), userID, templateID); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"activated": true})
}
