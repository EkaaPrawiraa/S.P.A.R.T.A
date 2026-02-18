package middleware

import (
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

type MetricsSnapshot struct {
	UptimeSeconds int64  `json:"uptime_seconds"`
	RequestsTotal uint64 `json:"requests_total"`
	Responses2xx  uint64 `json:"responses_2xx"`
	Responses4xx  uint64 `json:"responses_4xx"`
	Responses5xx  uint64 `json:"responses_5xx"`
}

type metrics struct {
	startedAt time.Time
	reqTotal  atomic.Uint64
	resp2xx   atomic.Uint64
	resp4xx   atomic.Uint64
	resp5xx   atomic.Uint64
}

var globalMetrics = &metrics{startedAt: time.Now()}

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		globalMetrics.reqTotal.Add(1)
		status := c.Writer.Status()
		switch {
		case status >= 200 && status < 300:
			globalMetrics.resp2xx.Add(1)
		case status >= 400 && status < 500:
			globalMetrics.resp4xx.Add(1)
		case status >= 500:
			globalMetrics.resp5xx.Add(1)
		}
	}
}

func GetMetricsSnapshot() MetricsSnapshot {
	uptime := time.Since(globalMetrics.startedAt)
	return MetricsSnapshot{
		UptimeSeconds: int64(uptime.Seconds()),
		RequestsTotal: globalMetrics.reqTotal.Load(),
		Responses2xx:  globalMetrics.resp2xx.Load(),
		Responses4xx:  globalMetrics.resp4xx.Load(),
		Responses5xx:  globalMetrics.resp5xx.Load(),
	}
}

func MetricsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "data": GetMetricsSnapshot()})
	}
}
