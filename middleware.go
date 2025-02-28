// SPDX-FileCopyrightText: 2021 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package candlelight

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/xmidt-org/wrp-go/v3/wrpcontext"
	"github.com/xmidt-org/wrp-go/v3/wrphttp"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const (
	spanIDHeaderName  = "X-Xmidt-Span-ID"
	traceIDHeaderName = "X-Xmidt-Trace-ID"
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

// EchoFirstNodeTraceInfo captures the trace information from a request, writes it
// back in the response headers, and adds it to the request's context
// It can also decode the request and save the resulting WRP object in the context if isDecodable is true
func EchoFirstTraceNodeInfo(tracing Tracing, isDecodable bool) func(http.Handler) http.Handler {
	return func(delegate http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			var ctx context.Context
			var headerPrefix = tracing.headerPrefix
			var propagator = tracing.propagator

			if isDecodable {
				if req, err := wrphttp.DecodeRequest(r, nil); err == nil {
					r = req
				}
			}

			var traceHeaders []string
			ctx = propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))
			if msg, ok := wrpcontext.GetMessage(ctx); ok {
				traceHeaders = msg.Headers
			} else if headers := r.Header.Values(headerPrefix); len(headers) != 0 {
				traceHeaders = headers
			}

			// Iterate through the trace headers (if any), format them, and add them to ctx
			var tmp propagation.TextMapCarrier = propagation.MapCarrier{}
			for _, f := range traceHeaders {
				if f != "" {
					parts := strings.Split(f, ":")
					if len(parts) > 1 {
						// Remove leading space if there's any
						parts[1] = strings.Trim(parts[1], " ")
						tmp.Set(parts[0], parts[1])
					}
				}
			}

			ctx = propagation.TraceContext{}.Extract(ctx, tmp)
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
