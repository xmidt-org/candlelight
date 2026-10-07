// SPDX-FileCopyrightText: 2026 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

/*
Package candlelight configures OpenTelemetry tracing for XMiDT web services
and provides the middleware and helpers they share.

# Setup

New builds a Tracing value from a Config:

	tracing, err := candlelight.New(candlelight.Config{
		ApplicationName: "talaria",
		Provider:        "otlp/http",
		Endpoint:        "collector.example.com:4318",
		ParentBased:     "honor",
		NoParent:        "always",
		HeaderPrefix:    "X-Xmidt-Trace",
	})

Tracing carries the TracerProvider and the W3C Trace Context propagator, which
services hand to otelmux, otelhttp and their own instrumentation.

The built-in providers are otlp/grpc, otlp/http, stdout and noop, along with
the deprecated jaeger and zipkin exporters.  A provider name not in that list
is looked up in Config.Providers, so a service can supply its own
ProviderConstructor.  With no provider configured, tracing is a no-op;
Tracing.IsNoop reports this.

ParentBased and NoParent select the sampler.  The default, ParentBased "ignore",
disables span creation.  ParentBased "honor" follows the parent span's decision,
and NoParent decides whether a root span is started when there is no parent.

# Middleware

EchoFirstTraceNodeInfo is the middleware that services place after the
otelmux handler.  It reads the incoming trace context from the HeaderPrefix
header or, when told the request is decodable, from the Headers field of the
WRP message the request carries.  The message may be in the request body as
JSON or msgpack, or spread across the X-Xmidt-* headers; both forms are read
by github.com/xmidt-org/wrphttp.  The middleware then puts the trace context
in the request's context for the handlers that follow.

TraceMiddleware and TraceConfig are the earlier form of this and are
deprecated.

# Logging and propagation

ExtractTraceInfo returns the current trace and span IDs from a context, and
AppendTraceInfo and InjectTraceInfoInLogger add them to the key/value pairs
given to a logger.  InjectTraceInfo writes the current span's traceparent and
tracestate into a carrier such as an outbound request's headers.  GenTID
generates a transaction identifier for use when a request arrives without one.
*/
package candlelight
