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

// Config specifies parameters relevant for otel trace provider.
type Config struct {
	// ApplicationName is the name for this application.
	ApplicationName string `json:"applicationName"`

	// Provider is the name of the trace provider to use.
	Provider string `json:"provider"`

	// Endpoint is the endpoint to which spans need to be submitted.
	Endpoint string `json:"endpoint"`

	// SkipTraceExport works only in case of provider stdout. Set
	// SkipTraceExport = true if you don't want to print the span
	// and tracer information in stdout.
	SkipTraceExport bool `json:"skipTraceExport"`

	// Providers are useful when client wants to add their own custom
	// TracerProvider.
	Providers map[string]ProviderConstructor `json:"-"`
}

// HeaderConfig specifies relevant header parameters.
type HeaderConfig struct {

	// SpanIDHeaderName is used to set custom header in response and logs.
	SpanIDHeaderName string `json:"spanIDHeaderName"`

	// TraceIDHeaderName is used to set custom header in response and logs.
	TraceIDHeaderName string `json:"traceIDHeaderName"`
}

// TraceConfig will be used in TraceMiddleware to use config and TraceProvider
// objects created by ConfigureTracerProvider.
type TraceConfig struct {
	HeaderConfig  HeaderConfig
	TraceProvider trace.TracerProvider
}

// TracingConfig represents a component for tracing. It will be used to unmarshal
// the configuration provided in application.yaml file.
type TracingConfig struct {
	// Provider describe the configuration for this application..
	Provider Config
	// Headers  is the header name  configuration for this application.
	Headers HeaderConfig
}
