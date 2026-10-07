// SPDX-FileCopyrightText: 2021 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package candlelight

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

func TestConfigureTracerProvider(t *testing.T) {
	tcs := []struct {
		Description string
		Config      Config
		Err         error
	}{
		{
			// nolint:goconst
			Description: "Otlp/gRPC: Valid",
			Config: Config{
				// nolint:goconst
				Provider: "otlp/grpc",
				// nolint:goconst
				Endpoint: "http://localhost",
				// nolint:goconst
				ParentBased: "ignore",
				// nolint:goconst
				NoParent: "never",
			},
		},
		{
			// nolint:goconst
			Description: "Otlp/gRPC: Valid",
			Config: Config{
				// nolint:goconst
				Provider: "otlp/grpc",
				Endpoint: "http://localhost",
				// nolint:goconst
				ParentBased: "honor",
				// nolint:goconst
				NoParent: "never",
			},
		},
		{
			// nolint:goconst
			Description: "Otlp/gRPC: Valid",
			Config: Config{
				// nolint:goconst
				Provider: "otlp/grpc",
				Endpoint: "http://localhost",
				// nolint:goconst
				ParentBased: "honor",
				// nolint:goconst
				NoParent: "always",
			},
		},
		{
			// nolint:goconst
			Description: "Otlp/gRPC: Valid",
			Config: Config{
				// nolint:goconst
				Provider:    "otlp/grpc",
				Endpoint:    "http://localhost",
				ParentBased: "ignore",
				// nolint:goconst
				NoParent: "always",
			},
		},
		{
			Description: "Otlp/gRPC: Missing Endpoint",
			Config: Config{
				// nolint:goconst
				Provider:    "otlp/grpc",
				ParentBased: "ignore",
				// nolint:goconst
				NoParent: "never",
			},
			Err: ErrTracerProviderBuildFailed,
		},
		{
			Description: "Valid Missing ParentBased Value",
			Config: Config{
				Provider: "otlp/grpc",
				Endpoint: "http://localhost",
			},
		},
		{
			Description: "Valid Missing NoParent Value",
			Config: Config{
				// nolint:goconst
				Provider: "otlp/grpc",
				Endpoint: "http://localhost",
				// nolint:goconst
				ParentBased: "honor",
			},
		},
		{
			Description: "Invalid ParentBased Value",
			Config: Config{
				// nolint:goconst
				Provider:    "otlp/grpc",
				Endpoint:    "http://localhost",
				ParentBased: "dishonor",
			},
			Err: ErrInvalidParentBasedValue,
		},
		{
			Description: "Invalid No Parent Value",
			Config: Config{
				// nolint:goconst
				Provider: "otlp/grpc",
				Endpoint: "http://localhost",
				// nolint:goconst
				ParentBased: "honor",
				NoParent:    "sometimes",
			},
			Err: ErrInvalidNoParentValue,
		},
		{
			Description: "Otlp/HTTP: Valid",
			Config: Config{
				// nolint:goconst
				Provider: "otlp/http",
				Endpoint: "http://localhost",
			},
		},
		{
			Description: "Otlp/HTTP: Missing Endpoint",
			Config: Config{
				// nolint:goconst
				Provider: "otlp/http",
			},
			Err: ErrTracerProviderBuildFailed,
		},
		{
			Description: "Jaeger: Missing endpoint",
			Config: Config{
				// nolint:goconst
				Provider: "jaeger",
			},
			Err: ErrTracerProviderBuildFailed,
		},
		{
			Description: "Zipkin: Missing endpoint",
			Config: Config{
				Provider: "Zipkin",
			},
			Err: ErrTracerProviderBuildFailed,
		},
		{
			Description: "Jaeger: Valid",
			Config: Config{
				Provider: "jaeger",
				Endpoint: "http://localhost",
			},
		},
		{
			Description: "Zipkin: Valid",
			Config: Config{
				Provider: "Zipkin",
				Endpoint: "http://localhost",
			},
		},
		{
			Description: "Unknown Provider",
			Config: Config{
				Provider: "undefined",
			},
			Err: ErrTracerProviderNotFound,
		},
		{
			Description: "Stdout: Valid",
			Config: Config{
				Provider: "stdOut",
			},
		},
		{
			Description: "Stdout: Valid skip export",
			Config: Config{
				Provider:        "stdoUt",
				SkipTraceExport: true,
			},
		},
		{
			Description: "Default",
			Config:      Config{},
		},
		{
			Description: "NoOp: Valid",
			Config: Config{
				Provider: "noop",
			},
		},
		{
			Description: "Custom provider",
			Config: Config{
				Provider: "coolest",
				Providers: map[string]ProviderConstructor{
					"coolest": func(_ Config, _ sdktrace.Sampler) (trace.TracerProvider, error) {
						return noop.NewTracerProvider(), nil
					},
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.Description, func(t *testing.T) {
			var (
				assert  = assert.New(t)
				tp, err = ConfigureTracerProvider(tc.Config)
			)
			if tc.Err == nil {
				assert.NotNil(tp)
			}
			assert.True(errors.Is(err, tc.Err))
		})
	}
}

const (
	testApp      = "talaria"
	testExporter = "otlp/http"
)

func TestNewResource(t *testing.T) {
	tests := []struct {
		description string
		config      Config
		expected    []attribute.KeyValue
	}{
		{
			description: "with an environment",
			config:      Config{ApplicationName: testApp, Provider: testExporter, Environment: "prod"},
			expected: []attribute.KeyValue{
				attribute.String("deployment.environment", "prod"),
				attribute.String("exporter", testExporter),
				attribute.String("service.name", testApp),
			},
		}, {
			description: "without an environment",
			config:      Config{ApplicationName: testApp, Provider: testExporter},
			expected: []attribute.KeyValue{
				attribute.String("exporter", testExporter),
				attribute.String("service.name", testApp),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			// Attributes come back sorted by key.
			assert.Equal(t, tc.expected, newResource(tc.config).Attributes())
		})
	}
}

// The environment reaches the spans a configured provider produces.
func TestEnvironmentOnSpans(t *testing.T) {
	tp, err := ConfigureTracerProvider(Config{
		ApplicationName: testApp,
		Provider:        "stdout", // nolint:goconst
		SkipTraceExport: true,
		Environment:     "staging",
		ParentBased:     "honor",
		NoParent:        "always",
	})
	require.NoError(t, err)

	_, span := tp.Tracer("test").Start(context.Background(), "op")
	span.End()

	ro, ok := span.(sdktrace.ReadOnlySpan)
	require.True(t, ok)
	attrs := ro.Resource().Set()
	env, ok := attrs.Value("deployment.environment")
	require.True(t, ok)
	assert.Equal(t, "staging", env.AsString())
	name, ok := attrs.Value("service.name")
	require.True(t, ok)
	assert.Equal(t, testApp, name.AsString())
}
