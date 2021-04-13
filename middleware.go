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
	"net/http"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const (
	spanIDHeaderName  = "X-Midt-Span-ID"
	traceIDHeaderName = "X-Midt-Trace-ID"
	SpanIDLogKeyName  = "span-id"
	TraceIdLogKeyName = "trace-id"
)

// TraceMiddleware acts as interceptor that is the first point of interaction
// for all requests. It will be responsible for starting a new span with existing
// traceId if present in the request header as traceparent. Otherwise it will
// generate new trace id. Example of traceparent will be
// version[2]-traceId[32]-spanId[16]-traceFlags[2]. It is mandatory for continuing
// existing traces while tracestate is optional.
// Deprecated. Please consider using EchoFirstTraceNodeInfo.
func (traceConfig *TraceConfig) TraceMiddleware(delegate http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prop := propagation.TraceContext{}
		ctx := prop.Extract(r.Context(), propagation.HeaderCarrier(r.Header))
		rsc := trace.RemoteSpanContextFromContext(ctx)
		tracer := traceConfig.TraceProvider.Tracer(r.URL.Path)
		ctx, span := tracer.Start(ctx, r.URL.Path)
		defer span.End()
		if !rsc.IsValid() {
			w.Header().Set(spanIDHeaderName, span.SpanContext().SpanID().String())
			w.Header().Set(traceIDHeaderName, span.SpanContext().TraceID().String())
		}
		delegate.ServeHTTP(w, r.WithContext(ctx))
	})
}

// EchoFirstNodeTraceInfo captures the trace information from a request and writes it
// back in the response headers if the request is the first one in the trace path.
func EchoFirstTraceNodeInfo(propagator propagation.TextMapPropagator) func(http.Handler) http.Handler {
	return func(delegate http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))
			rsc := trace.RemoteSpanContextFromContext(ctx)
			sc := trace.SpanContextFromContext(ctx)
			if sc.IsValid() && !rsc.IsValid() {
				w.Header().Set("X-Midt-Span-ID", sc.SpanID().String())
				w.Header().Set("X-Midt-Trace-ID", sc.TraceID().String())
			}
			delegate.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
