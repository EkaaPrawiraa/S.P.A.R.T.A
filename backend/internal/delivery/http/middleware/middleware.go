package middleware

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	requestIDHeader = "X-Request-Id"
	requestIDKey    = "request_id"
)

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := strings.TrimSpace(c.GetHeader(requestIDHeader))
		if rid == "" {
			rid = uuid.NewString()
		}
		c.Set(requestIDKey, rid)
		c.Header(requestIDHeader, rid)
		c.Next()
	}
}

func GetRequestID(c *gin.Context) string {
	val, ok := c.Get(requestIDKey)
	if !ok {
		return ""
	}
	s, _ := val.(string)
	return s
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		method := c.Request.Method
		rid := GetRequestID(c)

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		size := c.Writer.Size()
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()

		userID := ""
		if v, ok := c.Get("user_id"); ok {
			if s, ok := v.(string); ok {
				userID = s
			}
		}

		level := slog.LevelInfo
		if status >= http.StatusInternalServerError {
			level = slog.LevelError
		} else if status >= http.StatusBadRequest {
			level = slog.LevelWarn
		}

		slog.Log(c.Request.Context(), level, "http_request",
			slog.String("request_id", rid),
			slog.String("method", method),
			slog.String("path", path),
			slog.String("query", query),
			slog.Int("status", status),
			slog.Int("bytes", size),
			slog.String("client_ip", clientIP),
			slog.String("user_agent", userAgent),
			slog.String("user_id", userID),
			slog.Int64("latency_ms", latency.Milliseconds()),
		)
	}
}
