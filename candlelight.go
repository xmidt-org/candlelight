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
	"fmt"
	"go.opentelemetry.io/otel/exporters/stdout"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/exporters/trace/zipkin"
	"go.opentelemetry.io/otel/label"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"strings"
)

var (
	nilProviderErr = fmt.Errorf("No provider is configured")
)

// This function is responsible for creating the traceProvider
// Will be creating the traceProvider based on the config.Provider
func ConfigureTracerProvider(config Config) (trace.TracerProvider, error) {
	if len(config.Provider) == 0 {
		return nil, nilProviderErr
	}
	// Handling camelcase of provider.
	config.Provider = strings.ToLower(config.Provider)
	providerConfig := providersConfig[config.Provider]
	if providerConfig == nil {
		providerConfig = config.Providers[config.Provider]
	}
	if providerConfig == nil {
		return nil, nilProviderErr
	}
	provider, err := providerConfig(config)
	if err != nil {
		return nil, err
	}
	return provider, nil
}

type ProviderConstructor func(config Config) (trace.TracerProvider, error)

// Created pre-defined immutable map of built-in provider's
var providersConfig = map[string]ProviderConstructor{
	"jaeger": func(cfg Config) (trace.TracerProvider, error) {
		traceProvider, flushFn, err := jaeger.NewExportPipeline(
			jaeger.WithCollectorEndpoint(cfg.Endpoint),
			jaeger.WithProcess(jaeger.Process{
				ServiceName: cfg.ApplicationName,
				Tags: []label.KeyValue{
					label.String("exporter", cfg.Provider),
				},
			}),
			jaeger.WithSDK(&sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		)
		if err != nil {
			return nil, err
		}
		defer flushFn()
		return traceProvider, nil
	},
	"zipkin": func(cfg Config) (trace.TracerProvider, error) {
		traceProvider, err := zipkin.NewExportPipeline(cfg.Endpoint,
			cfg.ApplicationName,
			zipkin.WithSDK(&sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		)
		return traceProvider, err
	},
	"stdout": func(cfg Config) (trace.TracerProvider, error) {
		var option stdout.Option
		if cfg.SkipTraceExport {
			option = stdout.WithoutTraceExport()
		} else {
			option = stdout.WithPrettyPrint()
		}
		otExporter, err := stdout.NewExporter(option)
		if err != nil {
			return nil, err
		}
		traceProvider := sdktrace.NewTracerProvider(sdktrace.WithSyncer(otExporter))
		return traceProvider, nil
	},
}
