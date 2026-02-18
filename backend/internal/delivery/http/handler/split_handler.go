package handler

import (
	"net/http"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/split"
	domainuc "S.P.A.R.T.A/backend/internal/domain/usecase"
	"github.com/gin-gonic/gin"
)

type SplitHandler struct {
	uc domainuc.SplitUsecase
}

func NewSplitHandler(uc domainuc.SplitUsecase) *SplitHandler {
	return &SplitHandler{uc: uc}
}

func (h *SplitHandler) CreateTemplate(c *gin.Context) {
	var req split.SplitTemplate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.uc.CreateTemplate(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, req)
}

func (h *SplitHandler) GetUserTemplates(c *gin.Context) {
	userID := c.Param("user_id")

	res, err := h.uc.GetUserTemplates(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
