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
	"context"
	"net/http"

	"github.com/xmidt-org/webpa-common/logging/logginghttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// InjectTraceInformationInLogger adds the traceID and spanID to
// key value pairs that can be provided to a logger.
func InjectTraceInformationInLogger() logginghttp.LoggerFunc {
	return func(kv []interface{}, request *http.Request) []interface{} {
		traceID, spanID := ExtractTraceInformation(request.Context())
		return append(kv, SpanIDLogKeyName, spanID, TraceIdLogKeyName, traceID)
	}
}

// ExtractTraceInformation will be extracting the traceID and spanID. If span
// is not started in middleware, then it will be returning noopSpan which will
// have traceID[32 digits] and spanID[16 digits] having all 0's i.e.
// 00000000000000000000000000000000 and 0000000000000000 respectively.
func ExtractTraceInformation(ctx context.Context) (string, string) {
	span := trace.SpanFromContext(ctx)
	traceID := span.SpanContext().TraceID().String()
	spanID := span.SpanContext().SpanID().String()
	return traceID, spanID
}

// InjectTraceInformation will be injecting traceParent and tracestate as
// headers in carrier from span which is available in context.
func InjectTraceInformation(ctx context.Context, carrier propagation.TextMapCarrier) {
	prop := propagation.TraceContext{}
	prop.Inject(ctx, carrier)
}
