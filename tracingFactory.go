package candlelight

import (
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// New creates a structure with components that apps can use to initialize OpenTelemetry
// tracing instrumentation code.
func New(config Config) (Tracing, error) {
	var tracing = Tracing{
		propagator: propagation.TraceContext{},
	}
	tracerProvider, err := ConfigureTracerProvider(config)
	if err != nil {
		return Tracing{}, err
	}
	tracing.tracerProvider = tracerProvider
	return tracing, nil
}

// Tracing contains the core dependencies to make tracing possible across an
// application.
type Tracing struct {
	tracerProvider trace.TracerProvider
	propagator     propagation.TextMapPropagator
}

// IsNoop returns true if the tracer provider component is a noop. False otherwise.
func (t Tracing) IsNoop() bool {
	return t.TracerProvider() == trace.NewNoopTracerProvider()
}

// TracerProvider returns the tracer provider component. By default, the noop
// tracer provider is returned.
func (t Tracing) TracerProvider() trace.TracerProvider {
	if t.tracerProvider == nil {
		return trace.NewNoopTracerProvider()
	}
	return t.tracerProvider
}

// Propagator returns the component that helps propagate trace context across
// API boundaries. By default, a W3C Trace Context format propagator is returned.
func (t Tracing) Propagator() propagation.TextMapPropagator {
	if t.propagator == nil {
		return propagation.TraceContext{}
	}
	return t.propagator
}
