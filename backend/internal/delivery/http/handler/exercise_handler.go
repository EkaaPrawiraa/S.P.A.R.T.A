package handler

import (
	"time"

	"S.P.A.R.T.A/backend/internal/delivery/http/dto"
	"S.P.A.R.T.A/backend/internal/delivery/http/response"
	domainuc "S.P.A.R.T.A/backend/internal/domain/usecase"
	"S.P.A.R.T.A/backend/pkg/validator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/exercise"
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
		response.Error(c, err)
		return
	}

	response.Success(c, dto.FromDomainExercises(res))
}

func (h *ExerciseHandler) GetExercise(c *gin.Context) {
	id := c.Param("id")

	res, err := h.uc.GetExercise(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, dto.FromDomainExercise(*res))
}

func (h *ExerciseHandler) CreateExercise(c *gin.Context) {
	var req dto.CreateExerciseRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid body")
		return
	}
	if err := validator.ValidateStruct(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	domainEx := &exercise.Exercise{
		ID:               uuid.NewString(),
		Name:             req.Name,
		PrimaryMuscle:    req.PrimaryMuscle,
		SecondaryMuscles: req.SecondaryMuscles,
		Equipment:        req.Equipment,
		CreatedAt:        time.Now().UTC(),
	}

	if err := h.uc.CreateExercise(c.Request.Context(), domainEx); err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, dto.FromDomainExercise(*domainEx))
}

func (h *ExerciseHandler) AddExerciseMedia(c *gin.Context) {
	exerciseID := c.Param("id")
	if _, err := uuid.Parse(exerciseID); err != nil {
		response.BadRequest(c, "invalid exercise id")
		return
	}

	var req dto.AddExerciseMediaRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid body")
		return
	}
	if err := validator.ValidateStruct(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	media := &exercise.ExerciseMedia{
		ID:           uuid.NewString(),
		ExerciseID:   exerciseID,
		MediaType:    req.MediaType,
		MediaURL:     req.MediaURL,
		ThumbnailURL: req.ThumbnailURL,
		CreatedAt:    time.Now().UTC(),
	}

	if err := h.uc.AddExerciseMedia(c.Request.Context(), media); err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, gin.H{"id": media.ID})
}
