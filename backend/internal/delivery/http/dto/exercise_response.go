package dto

import (
	"time"

	"S.P.A.R.T.A/backend/internal/domain/aggregate/exercise"
)

type ExerciseMediaResponseDTO struct {
	ID           string    `json:"id"`
	MediaType    string    `json:"media_type"`
	MediaURL     string    `json:"media_url"`
	ThumbnailURL *string   `json:"thumbnail_url,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type ExerciseResponseDTO struct {
	ID               string                     `json:"id"`
	Name             string                     `json:"name"`
	PrimaryMuscle    string                     `json:"primary_muscle"`
	SecondaryMuscles []string                   `json:"secondary_muscles"`
	Equipment        string                     `json:"equipment"`
	CreatedAt        time.Time                  `json:"created_at"`
	Media            []ExerciseMediaResponseDTO `json:"media"`
}

type CreateExerciseRequestDTO struct {
	Name             string   `json:"name" validate:"required"`
	PrimaryMuscle    string   `json:"primary_muscle" validate:"required"`
	SecondaryMuscles []string `json:"secondary_muscles"`
	Equipment        string   `json:"equipment" validate:"required"`
}

type AddExerciseMediaRequestDTO struct {
	MediaType    string  `json:"media_type" validate:"required,oneof=image video"`
	MediaURL     string  `json:"media_url" validate:"required,url"`
	ThumbnailURL *string `json:"thumbnail_url" validate:"omitempty,url"`
}

func FromDomainExercise(ex exercise.Exercise) ExerciseResponseDTO {
	out := ExerciseResponseDTO{
		ID:               ex.ID,
		Name:             ex.Name,
		PrimaryMuscle:    ex.PrimaryMuscle,
		SecondaryMuscles: ex.SecondaryMuscles,
		Equipment:        ex.Equipment,
		CreatedAt:        ex.CreatedAt,
		Media:            make([]ExerciseMediaResponseDTO, 0, len(ex.Media)),
	}

	for _, media := range ex.Media {
		out.Media = append(out.Media, ExerciseMediaResponseDTO{
			ID:           media.ID,
			MediaType:    media.MediaType,
			MediaURL:     media.MediaURL,
			ThumbnailURL: media.ThumbnailURL,
			CreatedAt:    media.CreatedAt,
		})
	}

	return out
}

func FromDomainExercises(items []exercise.Exercise) []ExerciseResponseDTO {
	out := make([]ExerciseResponseDTO, 0, len(items))
	for _, item := range items {
		out = append(out, FromDomainExercise(item))
	}
	return out
}
