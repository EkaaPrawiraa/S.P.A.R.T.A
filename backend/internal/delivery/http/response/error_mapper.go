package response

import (
	"errors"
	"net/http"

	domainerr "S.P.A.R.T.A/backend/internal/domain/errors"
)

func MapErrorToStatus(err error) int {
	switch {
	case errors.Is(err, domainerr.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, domainerr.ErrInvalidInput):
		return http.StatusBadRequest
	case errors.Is(err, domainerr.ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, domainerr.ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, domainerr.ErrConflict):
		return http.StatusConflict
	case errors.Is(err, domainerr.ErrAIUnavailable):
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}
