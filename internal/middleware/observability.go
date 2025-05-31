package middleware

import (
	"context"
	"strconv"
	"time"

	"github.com/chat-socio/backend/pkg/observability"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/google/uuid"
)

func ObservabilityMiddleware(obs *observability.Observability) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		start := time.Now()

		// Generate correlation ID
		correlationID := uuid.New().String()
		ctx = obs.Logger.WithCorrelationID(ctx, correlationID)

		// Start span
		ctx, span := obs.StartSpan(ctx, "HTTP "+string(c.Method())+" "+string(c.Path()))
		defer span()

		// Set request context
		c.Set("observability_context", ctx)

		// Increment in-flight requests
		obs.Metrics.HTTPRequestsInFlight.Inc()
		defer obs.Metrics.HTTPRequestsInFlight.Dec()

		// Process request
		c.Next(ctx)

		// Record metrics
		duration := time.Since(start)
		method := string(c.Method())
		path := string(c.Path())
		status := strconv.Itoa(c.Response.StatusCode())

		obs.Metrics.RecordHTTPRequest(method, path, status, duration)

		// Log request
		logger := obs.Logger.WithContext(ctx)
		logger.Info("HTTP request completed", map[string]any{
			"method":         method,
			"path":           path,
			"status":         status,
			"duration_ms":    duration.Milliseconds(),
			"correlation_id": correlationID,
		})
	}
}
