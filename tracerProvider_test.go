/**
 *  Copyright (c) 2021  Comcast Cable Communications Management, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
package candlelight

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"
)

func TestConfigureTracerProvider(t *testing.T) {
	tcs := []struct {
		Description string
		Config      Config
		Err         error
	}{
		{
			Description: "Otlp/gRPC: Valid",
			Config: Config{
				Provider: "otlp/grpc",
				Endpoint: "http://localhost",
			},
		},
		{
			Description: "Otlp/gRPC: Missing Endpoint",
			Config: Config{
				Provider: "otlp/grpc",
			},
			Err: ErrTracerProviderBuildFailed,
		},
		{
			Description: "Otlp/HTTP: Valid",
			Config: Config{
				Provider: "otlp/http",
				Endpoint: "http://localhost",
			},
		},
		{
			Description: "Otlp/HTTP: Missing Endpoint",
			Config: Config{
				Provider: "otlp/http",
			},
			Err: ErrTracerProviderBuildFailed,
		},
		{
			Description: "Jaeger: Missing endpoint",
			Config: Config{
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
					"coolest": func(_ Config) (trace.TracerProvider, error) {
						return trace.NewNoopTracerProvider(), nil
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
