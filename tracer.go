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
	"github.com/xmidt-org/webpa-common/logging/logginghttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

/**
Will be injecting the trace id and span id  in the logger.
*/
func InjectTraceInformationInLogger() logginghttp.LoggerFunc {
	return func(kv []interface{}, request *http.Request) []interface{} {
		traceId, spanId := ExtractTraceInformation(request.Context())
		return append(kv, SpanIdHeader, spanId, TraceIdHeader, traceId)
	}
}

/**
	Will be extracting the traceId and spanId
	if  span is not started in middleware then it will be returning noopSpan
	which will result traceId[32 digits] and spanId[16 digits] as 0
	i.e. 00000000000000000000000000000000 and 0000000000000000
*/
func ExtractTraceInformation(ctx context.Context) (string, string) {
	span := trace.SpanFromContext(ctx)
	traceId := span.SpanContext().TraceID.String()
	spanId := span.SpanContext().SpanID.String()
	return traceId, spanId
}

/**
	Will be injecting  traceParent and tracestate headers in carrier
	from span which is available in context.
*/
func InjectTraceInformation(ctx context.Context, carrier propagation.TextMapCarrier) {
	prop := propagation.TraceContext{}
	prop.Inject(ctx, carrier)

}
