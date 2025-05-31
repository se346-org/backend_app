package observability

import (
	"context"
	"fmt"
)

type Observability struct {
	Logger  *Logger
	Metrics *Metrics
	Tracer  *Tracer
}

type Config struct {
	ServiceName     string
	JaegerEndpoint  string
	LogLevel        string
	MetricsEnabled  bool
	TracingEnabled  bool
}

func New(config Config) (*Observability, error) {
	obs := &Observability{
		Logger:  NewLogger(),
		Metrics: NewMetrics(),
	}

	if config.TracingEnabled {
		tracer, err := NewTracer(config.ServiceName, config.JaegerEndpoint)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize tracer: %w", err)
		}
		obs.Tracer = tracer
	}

	return obs, nil
}

func (o *Observability) StartSpan(ctx context.Context, name string) (context.Context, func()) {
	if o.Tracer == nil {
		return ctx, func() {}
	}

	ctx, span := o.Tracer.StartSpan(ctx, name)
	return ctx, func() {
		span.End()
	}
} 