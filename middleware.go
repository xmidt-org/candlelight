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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

const (
	SpanIdHeader  = "X-B3-SpanId"
	TraceIdHeader = "X-B3-TraceId"
)

/**
1. Acts as interceptor should be  the first point of interactions for all request's.
2. Will  be responsible for starting a new span with existing traceId if present in request header as traceparent else it will  generate new trace id.
	  example of traceparent will be version[2]-traceId[32]-spanId[16]-traceFlags[2] is mandatory for continuing existing traces and tracestate is optional.
3. Will be writing the traceId and spanId in response headers for easier debugging in case of any incident.
4. will be adding the newly created span in the request context so that we can use it other places.
*/

func TraceMiddleware(delegate http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var span trace.Span
		ctx := r.Context()
		prop := propagation.TraceContext{}
		ctx = prop.Extract(ctx, r.Header)
		tracer := otel.GetTracerProvider().Tracer(r.URL.Path)
		ctx, span = tracer.Start(ctx, r.URL.Path)
		defer span.End()
		w.Header().Set(SpanIdHeader, span.SpanContext().SpanID.String())
		w.Header().Set(TraceIdHeader, span.SpanContext().TraceID.String())
		delegate.ServeHTTP(w, r.WithContext(ctx))
	})
}
