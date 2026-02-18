package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	domainerr "S.P.A.R.T.A/backend/internal/domain/errors"
)

type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(200, APIResponse{
		Status: "success",
		Data:   data,
	})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(201, APIResponse{
		Status: "success",
		Data:   data,
	})
}

func BadRequest(c *gin.Context, message string) {
	c.JSON(400, APIResponse{
		Status:  "error",
		Message: message,
	})
}

func Error(c *gin.Context, err error) {
	status := MapErrorToStatus(err)
	message := ""
	if err != nil {
		message = err.Error()
	}

	if status == http.StatusInternalServerError {
		message = domainerr.ErrInternal.Error()
	}

	c.JSON(status, APIResponse{
		Status:  "error",
		Message: message,
	})
}
