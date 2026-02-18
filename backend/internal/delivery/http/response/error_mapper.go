package response

import (
	"net/http"

	domainerr "S.P.A.R.T.A/backend/internal/domain/errors"
)

func MapErrorToStatus(err error) int {
	switch err {
	case domainerr.ErrNotFound:
		return http.StatusNotFound
	case domainerr.ErrInvalidInput:
		return http.StatusBadRequest
	case domainerr.ErrUnauthorized:
		return http.StatusUnauthorized
	case domainerr.ErrForbidden:
		return http.StatusForbidden
	case domainerr.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
