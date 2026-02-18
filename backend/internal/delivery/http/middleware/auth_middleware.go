package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	jwtSecret string
}

func NewAuthMiddleware(secret string) *AuthMiddleware {
	return &AuthMiddleware{jwtSecret: secret}
}

func (a *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "missing authorization header",
			})
			return
		}

		// Support:
		// - "Bearer <token>" (standard)
		// - "bearer <token>" (case-insensitive)
		// - accidental "Bearer Bearer <token>" (common when users paste full header into Postman Bearer token field)
		// - raw token value (some tools)
		fields := strings.Fields(authHeader)
		var tokenString string
		switch {
		case len(fields) == 1:
			tokenString = fields[0]
		case len(fields) == 2 && strings.EqualFold(fields[0], "Bearer"):
			tokenString = fields[1]
		case len(fields) == 3 && strings.EqualFold(fields[0], "Bearer") && strings.EqualFold(fields[1], "Bearer"):
			tokenString = fields[2]
		default:
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "invalid authorization header",
			})
			return
		}

		tokenString = strings.TrimSpace(tokenString)
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "missing bearer token",
			})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenSignatureInvalid
			}
			return []byte(a.jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "invalid token",
			})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "invalid token",
			})
			return
		}

		userIDVal, ok := claims["user_id"]
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "invalid token",
			})
			return
		}

		userID, ok := userIDVal.(string)
		if !ok || strings.TrimSpace(userID) == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "invalid token",
			})
			return
		}

		// inject into context
		c.Set("user_id", userID)

		role := "user"
		if roleVal, ok := claims["role"]; ok {
			if roleStr, ok := roleVal.(string); ok {
				trimmed := strings.TrimSpace(roleStr)
				if trimmed != "" {
					role = trimmed
				}
			}
		}
		c.Set("role", role)

		c.Next()
	}
}
