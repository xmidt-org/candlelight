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
	"go.opentelemetry.io/otel/exporters/jaeger" // nolint:staticcheck
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/exporters/zipkin"
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
)

// DefaultTracerProvider is used when no provider is given.
// The Noop tracer provider turns all tracing related operations into
// noops essentially disabling tracing.
const DefaultTracerProvider = "noop"

// ConfigureTracerProvider creates the TracerProvider based on the configuration
// provided. It has built-in support for jaeger, zipkin, stdout and noop providers.
// A different provider can be used if a constructor for it is provided in the
// config.
// If a provider name is not provided, a noop tracerProvider will be returned.
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
		return nil, fmt.Errorf("%w for provider %s", ErrTracerProviderNotFound, config.Provider)
	}

	// If parentBased value is empty, use default value
	if parentBasedTracing == "" {
		parentBasedTracing = "ignore"
	}

	// If noParent value is empty, use default value
	if noParentTracing == "" {
		noParentTracing = "never"
	}

	// Setting up trace sampler based on ParentBased and NoParent values in the config
	var sampler sdktrace.Sampler

	if parentBasedTracing == "ignore" {
		sampler = sdktrace.NeverSample()
	} else if parentBasedTracing == "honor" {
		if noParentTracing == "never" {
			sampler = sdktrace.ParentBased(sdktrace.NeverSample())

		} else if noParentTracing == "always" {
			sampler = sdktrace.ParentBased(sdktrace.AlwaysSample())

		} else {
			return nil, ErrInvalidNoParentValue
		}
	} else {
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

// Created pre-defined immutable map of built-in provider's
var providersConfig = map[string]ProviderConstructor{
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
			sdktrace.WithResource(
				resource.NewWithAttributes(
					semconv.SchemaURL,
					semconv.ServiceNameKey.String(cfg.ApplicationName),
				),
			),
			sdktrace.WithSampler(smplr),
		), nil

	},
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
			sdktrace.WithResource(
				resource.NewWithAttributes(
					semconv.SchemaURL,
					semconv.ServiceNameKey.String(cfg.ApplicationName),
				),
			),
			sdktrace.WithSampler(smplr),
		), nil

	},
	"jaeger": func(cfg Config, smplr sdktrace.Sampler) (trace.TracerProvider, error) {
		if cfg.Endpoint == "" {
			return nil, ErrTracerProviderBuildFailed
		}

		exporter, err := jaeger.New(
			jaeger.WithCollectorEndpoint(
				jaeger.WithEndpoint(cfg.Endpoint)))
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrTracerProviderBuildFailed, err)
		}

		tp := sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(exporter),
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithResource(
				resource.NewWithAttributes(
					semconv.SchemaURL,
					semconv.ServiceNameKey.String(cfg.ApplicationName),
					attribute.String("exporter", cfg.Provider),
				)),
		)
		return tp, nil
	},
	"zipkin": func(cfg Config, smplr sdktrace.Sampler) (trace.TracerProvider, error) {
		if cfg.Endpoint == "" {
			return nil, ErrTracerProviderBuildFailed
		}

		exporter, err := zipkin.New(cfg.Endpoint)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrTracerProviderBuildFailed, err)
		}

		tp := sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(exporter),
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithResource(
				resource.NewWithAttributes(
					semconv.SchemaURL,
					semconv.ServiceNameKey.String(cfg.ApplicationName),
					attribute.String("exporter", cfg.Provider),
				)),
		)
		return tp, nil
	},
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
		tp := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))
		return tp, nil
	},
	"noop": func(config Config, smplr sdktrace.Sampler) (trace.TracerProvider, error) {
		return noop.NewTracerProvider(), nil
	},
}
