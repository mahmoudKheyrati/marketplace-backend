package middleware

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/mahmoudKheyrati/marketplace-backend/pkg"
	"github.com/mahmoudKheyrati/marketplace-backend/pkg/metric"
	"github.com/opentracing/opentracing-go"
	"time"
)

type MetricsMiddleware struct {
	metrics *metric.Metrics
}

func NewMetricsMiddleware(metrics *metric.Metrics) *MetricsMiddleware {
	return &MetricsMiddleware{metrics: metrics}
}
func (m *MetricsMiddleware) MetricsMiddleware(c *fiber.Ctx) error {
	// metrics
	method := string(c.Request().URI().Path())
	m.metrics.MethodCount.WithLabelValues("safir-api", method, "0").Add(1)
	startTime := time.Now()

	ctx := context.Background()
	span, ctx := opentracing.StartSpanFromContext(ctx, "safir-api "+method)
	defer span.Finish()

	err := c.Next()
	if err != nil {
		pkg.Logger().Error(err)
	}

	latency := float64(time.Since(startTime))
	if c.Response().StatusCode() >= 200 && c.Response().StatusCode() < 300 {
		m.metrics.MethodSuccessCount.WithLabelValues("safir-api", method, "0").Add(1)
		m.metrics.MethodDurations.WithLabelValues("safir-api", method, "0").Observe(latency)
		m.metrics.MethodDurationsHistogram.WithLabelValues("safir-api", method, "0").Observe(latency / float64(time.Second))
	} else if c.Response().StatusCode() >= 400 && c.Response().StatusCode() < 500 {
		m.metrics.MethodUserErrorCount.WithLabelValues("safir-api", method, "0", "user_error").Add(1)
		m.metrics.MethodErrorDurations.WithLabelValues("safir-api", method, "0").Observe(latency)
		m.metrics.MethodErrorDurationsHistogram.WithLabelValues("safir-api", method, "0", "").Observe(latency / float64(time.Second))
	} else {
		m.metrics.MethodFailCount.WithLabelValues("safir-api", method, "0").Add(1)
		m.metrics.MethodErrorDurations.WithLabelValues("safir-api", method, "0").Observe(latency)
		m.metrics.MethodErrorDurationsHistogram.WithLabelValues("safir-api", method, "0", "").Observe(latency / float64(time.Second))
	}
	return nil
}
