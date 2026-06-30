// SPDX-FileCopyrightText: 2021 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package candlelight

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
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
