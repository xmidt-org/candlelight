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

	// ParentBased and NoParent dictate if and when new spans should be created.
	// ParentBased = "ignore" (default), tracing is effectively turned off and the "NoParent" value is ignored
	// ParentBased = "honor", the sampling decision is made by the parent of the span
	ParentBased string `json:"parentBased"`

	// NoParent decides if a root span should be initiated in the case where there is no existing parent
	// This value is ignored if ParentBased = "ignore"
	// NoParent = "never" (default), root spans are not initiated
	// NoParent = "always", roots spans are initiated
	NoParent string `json:"noParent"`
}

// TraceConfig will be used in TraceMiddleware to use config and TraceProvider
// objects created by ConfigureTracerProvider.
// (Deprecated). Consider using Tracing instead.
type TraceConfig struct {
	TraceProvider trace.TracerProvider
}
