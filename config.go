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
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

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

// TraceConfig will be used in TraceMiddleware to use config and TraceProvider
// objects created by ConfigureTracerProvider.
// (Deprecated). Consider using Tracing
type TraceConfig struct {
	TraceProvider trace.TracerProvider
}

// Tracing contains the core dependencies for any component that will be
// part of tracing by either starting spans or propagating trace context across API boundaries.
type Tracing struct {
	TracerProvider trace.TracerProvider
	Propagator     propagation.TextMapPropagator
}
