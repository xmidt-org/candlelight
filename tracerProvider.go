// SPDX-FileCopyrightText: 2021 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package candlelight

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

var (
	ErrTracerProviderNotFound    = errors.New("tracerProvider builder could not be found")
	ErrTracerProviderBuildFailed = errors.New("failed building TracerProvider")
	ErrInvalidParentBasedValue   = errors.New("invalid ParentBased value provided in configuration")
	ErrInvalidNoParentValue      = errors.New("invalid No Parent value provided in configuration")

	// ErrTracerProviderRemoved is returned when the configuration names a
	// provider this package used to ship and no longer does.  The error text
	// says what to use instead.
	ErrTracerProviderRemoved = errors.New("tracerProvider was removed")
)

// removedProviders maps the names of providers that were removed to the
// migration advice returned with ErrTracerProviderRemoved.  The names stay
// here so a configuration that still uses one fails with an explanation
// rather than as an unknown provider.  A constructor of the same name in
// Config.Providers takes precedence.
var removedProviders = map[string]string{
	"jaeger": `the "jaeger" provider was removed in candlelight v0.4.0 because OpenTelemetry ` +
		`dropped the Jaeger exporter; Jaeger accepts OTLP natively, so use "otlp/grpc" ` +
		`with the collector's port 4317 or "otlp/http" with port 4318`,
	"zipkin": `the "zipkin" provider was removed in candlelight v0.4.0 because OpenTelemetry ` +
		`deprecated the Zipkin exporter; use "otlp/grpc" or "otlp/http" and send to an ` +
		`OpenTelemetry Collector that exports to Zipkin, or to a Zipkin server that accepts OTLP`,
}

// DefaultTracerProvider is used when no provider is given.
// The Noop tracer provider turns all tracing related operations into
// noops essentially disabling tracing.
const DefaultTracerProvider = "noop"

// ConfigureTracerProvider creates the TracerProvider based on the configuration
// provided. It has built-in support for the otlp/grpc, otlp/http, stdout and
// noop providers.  A different provider can be used if a constructor for it is
// provided in the config.
// If a provider name is not provided, a noop tracerProvider will be returned.
// The jaeger and zipkin providers were removed in v0.4.0; naming one returns
// ErrTracerProviderRemoved with migration advice.
func ConfigureTracerProvider(config Config) (trace.TracerProvider, error) {
	if len(config.Provider) == 0 {
		config.Provider = DefaultTracerProvider
	}
	// Handling camelcase of provider.
	config.Provider = strings.ToLower(config.Provider)
	providerConfig := config.Providers[config.Provider]
	parentBasedTracing := config.ParentBased
	noParentTracing := config.NoParent
	if providerConfig == nil {
		providerConfig = providersConfig[config.Provider]
	}
	if providerConfig == nil {
		if why, removed := removedProviders[config.Provider]; removed {
			return nil, fmt.Errorf("%w: %s", ErrTracerProviderRemoved, why)
		}
		return nil, fmt.Errorf("%w for provider %s", ErrTracerProviderNotFound, config.Provider)
	}

	// If parentBased value is empty, use default value
	if parentBasedTracing == "" {
		// nolint:goconst
		parentBasedTracing = "ignore"
	}

	// If noParent value is empty, use default value
	if noParentTracing == "" {
		// nolint:goconst
		noParentTracing = "never"
	}

	// Setting up trace sampler based on ParentBased and NoParent values in the config
	var sampler sdktrace.Sampler

	switch parentBasedTracing {
	case "ignore":
		sampler = sdktrace.NeverSample()
	case "honor": // nolint:goconst
		switch noParentTracing {
		case "never":
			sampler = sdktrace.ParentBased(sdktrace.NeverSample())

		// nolint:goconst
		case "always":
			sampler = sdktrace.ParentBased(sdktrace.AlwaysSample())

		default:
			return nil, ErrInvalidNoParentValue
		}
	default:
		return nil, ErrInvalidParentBasedValue
	}

	provider, err := providerConfig(config, sampler)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrTracerProviderBuildFailed, err)
	}
	return provider, nil
}

// ProviderConstructor is useful when client wants to add their own custom
// TracerProvider.
type ProviderConstructor func(config Config, sampler sdktrace.Sampler) (trace.TracerProvider, error)

// newResource describes the service that spans come from: its name, the
// deployment environment when one is configured, and the exporter in use.
func newResource(cfg Config) *resource.Resource {
	attrs := []attribute.KeyValue{
		semconv.ServiceNameKey.String(cfg.ApplicationName),
		attribute.String("exporter", cfg.Provider),
	}
	if cfg.Environment != "" {
		attrs = append(attrs, semconv.DeploymentEnvironmentKey.String(cfg.Environment))
	}
	return resource.NewWithAttributes(semconv.SchemaURL, attrs...)
}

// Created pre-defined immutable map of built-in provider's
var providersConfig = map[string]ProviderConstructor{
	// nolint:goconst
	"otlp/grpc": func(cfg Config, smplr sdktrace.Sampler) (trace.TracerProvider, error) {
		// Send traces over gRPC
		if cfg.Endpoint == "" {
			return nil, ErrTracerProviderBuildFailed
		}
		exporter, err := otlptracegrpc.New(context.Background(),

			otlptracegrpc.WithEndpoint(cfg.Endpoint),
			otlptracegrpc.WithInsecure(),
		)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrTracerProviderBuildFailed, err)
		}

		return sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(newResource(cfg)),
			sdktrace.WithSampler(smplr),
		), nil

	},
	// nolint:goconst
	"otlp/http": func(cfg Config, smplr sdktrace.Sampler) (trace.TracerProvider, error) {
		// Send traces over HTTP
		if cfg.Endpoint == "" {
			return nil, ErrTracerProviderBuildFailed
		}
		exporter, err := otlptracehttp.New(context.Background(),

			otlptracehttp.WithEndpoint(cfg.Endpoint),
			otlptracehttp.WithInsecure(),
		)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrTracerProviderBuildFailed, err)
		}

		return sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(newResource(cfg)),
			sdktrace.WithSampler(smplr),
		), nil

	},
	// nolint:goconst
	"stdout": func(cfg Config, smplr sdktrace.Sampler) (trace.TracerProvider, error) {
		var option stdout.Option
		if cfg.SkipTraceExport {
			option = stdout.WithWriter(io.Discard)
		} else {
			option = stdout.WithPrettyPrint()
		}
		exporter, err := stdout.New(option)
		if err != nil {
			return nil, err
		}
		tp := sdktrace.NewTracerProvider(
			sdktrace.WithSyncer(exporter),
			sdktrace.WithResource(newResource(cfg)),
		)
		return tp, nil
	},
	"noop": func(config Config, smplr sdktrace.Sampler) (trace.TracerProvider, error) {
		return noop.NewTracerProvider(), nil
	},
}
