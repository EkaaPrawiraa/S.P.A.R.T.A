package middleware

import "github.com/gin-gonic/gin"

func GetUserID(c *gin.Context) string {
	val, exists := c.Get("user_id")
	if !exists {
		return ""
	}
	s, ok := val.(string)
	if !ok {
		return ""
	}
	return s
}
