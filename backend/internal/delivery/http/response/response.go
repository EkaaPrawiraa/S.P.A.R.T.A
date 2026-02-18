package response

import "github.com/gin-gonic/gin"

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

func Error(c *gin.Context, err error) {
	status := MapErrorToStatus(err)

	c.JSON(status, APIResponse{
		Status:  "error",
		Message: err.Error(),
	})
}
