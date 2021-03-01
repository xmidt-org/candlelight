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

import "go.opentelemetry.io/otel/trace"

// Config object for otel tracing
type Config struct {
	// ApplicationName is the name for this application
	ApplicationName string `json:"applicationName"`

	// Provider is the name of the trace provider to use
	Provider string `json:"provider"`

	// Endpoint is the endpoint to which spans needs to be submitted.
	Endpoint string `json:"endpoint"`

	// SkipTraceExport works only in case of provider stdout
	// set SkipTraceExport = true if you don't want to print the span and tracer information in stdout
	SkipTraceExport bool `json:"skipTraceExport"`

	// In case of any custom provider which client want to use its own.
	Providers map[string]ProviderConstructor `json:"-"`

	// In case User wants to use his own  headers in response and logs.
	SpanIDHeaderName  string `json:"spanIDHeaderName"`
	TraceIDHeaderName string `json:"traceIDHeaderName"`
}

// TraceConfig  which will have Config and trace provider created by ConfigureTraceProvider.
type TraceConfig struct {
	config        Config
	traceProvider trace.TracerProvider
}

// Constructor for TraceConfig object
func NewTraceConfig(config Config, traceProvider trace.TracerProvider) TraceConfig {
	return TraceConfig{config, traceProvider}
}
