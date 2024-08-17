package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	metric "github.com/gothinkster/golang-gin-realworld-example-app/pkg/metric"
)

// Prometheus metrics middleware
func (mv *MiddlewareManager) MetricsMiddleware(p *metric.Prometheus) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		var status = c.Writer.Status()
		p.ObserveResponseTime(status, c.Request.Method, c.Request.URL.Path, time.Since(start).Seconds())
	}
}
