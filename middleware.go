// SPDX-FileCopyrightText: 2021 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package candlelight

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"strings"

	"github.com/xmidt-org/wrp-go/v5"
	"github.com/xmidt-org/wrphttp"

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
// If isDecodable is true, it also decodes the request as a WRP message and uses the
// message headers as the source of trace information.
func EchoFirstTraceNodeInfo(tracing Tracing, isDecodable bool) func(http.Handler) http.Handler {
	return func(delegate http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			var ctx context.Context
			headerPrefix := tracing.headerPrefix
			propagator := tracing.Propagator()

			var traceHeaders []string
			var decoded bool
			if isDecodable {
				traceHeaders, decoded = decodeWRPHeaders(r)
			}

			ctx = propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))
			if !decoded {
				traceHeaders = r.Header.Values(headerPrefix)
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

// decodeWRPHeaders returns the headers of the WRP message carried by r, which is
// either expressed as HTTP headers or encoded in the body as JSON or msgpack.
// The request body is left intact for downstream handlers.
func decodeWRPHeaders(r *http.Request) ([]string, bool) {
	// wrphttp consumes the body, so decode from a copy and put the body back.
	req := *r
	req.Body = http.NoBody
	if r.Body != nil {
		contents, err := io.ReadAll(r.Body)
		_ = r.Body.Close()
		r.Body = io.NopCloser(bytes.NewReader(contents))
		if err != nil {
			return nil, false
		}
		req.Body = io.NopCloser(bytes.NewReader(contents))
	}

	msgs, err := wrphttp.DecodeRequest(&req, wrp.NoStandardValidation())
	if err != nil || len(msgs) == 0 {
		return nil, false
	}

	var msg wrp.Message
	if err := msgs[0].To(&msg, wrp.NoStandardValidation()); err != nil {
		return nil, false
	}

	return msg.Headers, true
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
