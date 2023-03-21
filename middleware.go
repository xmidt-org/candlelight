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
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/xmidt-org/webpa-common/xhttp"
	"github.com/xmidt-org/wrp-go/v3"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const (
	spanIDHeaderName  = "X-Midt-Span-ID"
	traceIDHeaderName = "X-Midt-Trace-ID"
	SpanIDLogKeyName  = "span-id"
	TraceIdLogKeyName = "trace-id"
	// HeaderWPATIDKeyName is the header key for the WebPA transaction UUID
	HeaderWPATIDKeyName = "X-WebPA-Transaction-Id"
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
		sc := trace.SpanContextFromContext(ctx)
		tracer := traceConfig.TraceProvider.Tracer(r.URL.Path)
		ctx, span := tracer.Start(ctx, r.URL.Path)
		defer span.End()
		if !sc.IsValid() {
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
			sc := trace.SpanContextFromContext(ctx)

			// If there is no traceparent information in HTTP headers, check for them elsewhere
			if !sc.IsRemote() {

				// Check for trace info in Xmidt headers
				var traceHeaders []string
				xmidtHeaders := r.Header.Values("X-midt-Headers")
				if len(xmidtHeaders) != 0 {
					traceHeaders = xmidtHeaders
				} else {

					// Check for trace info in the request body
					reader := r.Body
					bodyData, err := io.ReadAll(reader)
					if err != nil {
						log.Printf("failed to read body: %s", err)
					}

					r.Body, r.GetBody = xhttp.NewRewindBytes(bodyData)
					var msg wrp.Message

					if err := json.Unmarshal(bodyData, &msg); err != nil {
						log.Printf("failed to unmarshal payload: %s", err)
					}

					if r.Header.Get("Content-Type") == "application/msgpack" {
						// Request body is msgpack
						err = wrp.NewDecoderBytes(bodyData, wrp.Msgpack).Decode(&msg)
						if err != nil {
							log.Printf("failed to decode msgpack: %s", err)
						}
						traceHeaders = msg.Headers
					} else if msg.Headers != nil {
						// Request body is JSON
						traceHeaders = msg.Headers
					}
				}

				// Add trace headers to TextMapCarrier
				var tmp propagation.TextMapCarrier = propagation.MapCarrier{}
				for _, f := range traceHeaders {
					if f != "" {
						parts := strings.Split(f, ":")
						// Remove leading space if there's any
						parts[1] = strings.Trim(parts[1], " ")
						tmp.Set(parts[0], parts[1])
					}
				}

				// Propagate context with traceheaders
				ctx = propagation.TraceContext{}.Extract(ctx, tmp)
				sc = trace.SpanContextFromContext(ctx)
			}

			if sc.IsValid() {
				w.Header().Set("X-Midt-Span-ID", sc.SpanID().String())
				w.Header().Set("X-Midt-Trace-ID", sc.TraceID().String())
			}
			delegate.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GenTID generates a 16-byte long string
// it returns "N/A" in the extreme case the random string could not be generated
func GenTID() (tid string) {
	buf := make([]byte, 16)
	tid = "N/A"
	if _, err := rand.Read(buf); err == nil {
		tid = base64.RawURLEncoding.EncodeToString(buf)
	}
	return
}
