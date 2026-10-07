// SPDX-FileCopyrightText: 2022 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package candlelight

import (
	"bytes"
	"io"
	"maps"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xmidt-org/wrp-go/v5"
	"go.opentelemetry.io/otel/trace"
)

const (
	testTraceID     = "4bf92f3577b34da6a3ce929d0e0e4736"
	testTraceparent = "traceparent: 00-" + testTraceID + "-00f067aa0ba902b7-01"
	testTraceHeader = "X-Trace"
	contentType     = "Content-Type"
	headerA         = "a: 1"
	msgpackType     = "application/msgpack"
	testTracestate  = "tracestate: congo=t61rcWkgMzE,rojo=00f067aa0ba902b7"

	// Header names used when a WRP message is expressed as HTTP headers.
	wrpMessageTypeHeader       = "X-Xmidt-Message-Type"
	wrpLegacyMessageTypeHeader = "X-Midt-Msg-Type"
	wrpHeadersHeader           = "X-Xmidt-Headers"
	wrpLegacyHeadersHeader     = "X-Midt-Headers"
	simpleEvent                = "SimpleEvent"
)

func TestGenTID(t *testing.T) {
	assert := assert.New(t)
	tid := GenTID()
	assert.NotEmpty(tid)
}

func encodeWRP(t *testing.T, f wrp.Format, headers ...string) []byte {
	var buf []byte
	msg := wrp.Message{
		Type:        wrp.SimpleEventMessageType,
		Source:      "mac:112233445566",
		Destination: "event:device-status",
		Headers:     headers,
	}
	require.NoError(t, wrp.NewEncoderBytes(&buf, f).Encode(&msg))
	return buf
}

func TestDecodeWRPHeaders(t *testing.T) {
	tests := []struct {
		description string
		header      http.Header
		body        []byte
		expected    []string
		expectedOK  bool
	}{
		{
			description: "msgpack body",
			header:      http.Header{contentType: {msgpackType}},
			body:        encodeWRP(t, wrp.Msgpack, headerA, "b: 2"),
			expected:    []string{headerA, "b: 2"},
			expectedOK:  true,
		}, {
			description: "json body",
			header:      http.Header{contentType: {"application/json"}},
			body:        encodeWRP(t, wrp.JSON, headerA),
			expected:    []string{headerA},
			expectedOK:  true,
		}, {
			description: "json body without content type",
			body:        encodeWRP(t, wrp.JSON, headerA),
			expected:    []string{headerA},
			expectedOK:  true,
		}, {
			description: "msgpack body without content type",
			body:        encodeWRP(t, wrp.Msgpack, headerA),
			expected:    []string{headerA},
			expectedOK:  true,
		}, {
			description: "message as http headers",
			header: http.Header{
				wrpMessageTypeHeader: {simpleEvent},
				wrpHeadersHeader:     {"a: 1, b: 2", "c: 3"},
			},
			expected:   []string{headerA, "b: 2", "c: 3"},
			expectedOK: true,
		}, {
			// A proxy may fold the header lines into one; a tracestate's own
			// commas must survive that.
			description: "message as http headers folded by a proxy",
			header: http.Header{
				wrpMessageTypeHeader: {simpleEvent},
				wrpHeadersHeader:     {testTraceparent + ", " + testTracestate},
			},
			expected:   []string{testTraceparent, testTracestate},
			expectedOK: true,
		}, {
			description: "message as legacy http headers",
			header: http.Header{
				wrpLegacyMessageTypeHeader: {simpleEvent},
				wrpLegacyHeadersHeader:     {headerA},
			},
			expected:   []string{headerA},
			expectedOK: true,
		}, {
			// Several messages can't share one request span, so the HTTP
			// headers decide instead.
			description: "several messages in one request",
			header:      http.Header{contentType: {"application/jsonl"}},
			body: bytes.Join([][]byte{
				encodeWRP(t, wrp.JSON, headerA),
				encodeWRP(t, wrp.JSON, "b: 2"),
			}, []byte("\n")),
		}, {
			description: "unsupported content type",
			header:      http.Header{contentType: {"text/plain"}},
			body:        []byte("hello"),
		}, {
			description: "undecodable body",
			header:      http.Header{contentType: {msgpackType}},
			body:        []byte("not msgpack"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			assert := assert.New(t)
			r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(tc.body))
			maps.Copy(r.Header, tc.header)

			headers, ok := decodeWRPHeaders(r)
			assert.Equal(tc.expectedOK, ok)
			assert.Equal(tc.expected, headers)

			// The body must remain readable by downstream handlers.
			body, err := io.ReadAll(r.Body)
			assert.NoError(err)
			assert.Equal(len(tc.body), len(body))
		})
	}
}

func TestEchoFirstTraceNodeInfo(t *testing.T) {
	tests := []struct {
		description string
		isDecodable bool
		header      http.Header
		body        []byte
		expectValid bool
	}{
		{
			description: "trace from wrp message headers",
			isDecodable: true,
			header:      http.Header{contentType: {msgpackType}},
			body:        encodeWRP(t, wrp.Msgpack, testTraceparent),
			expectValid: true,
		}, {
			description: "trace from wrp message headers without content type",
			isDecodable: true,
			body:        encodeWRP(t, wrp.Msgpack, testTraceparent),
			expectValid: true,
		}, {
			description: "trace from wrp http headers folded by a proxy",
			isDecodable: true,
			header: http.Header{
				wrpMessageTypeHeader: {simpleEvent},
				wrpHeadersHeader:     {testTraceparent + ", " + testTracestate},
			},
			expectValid: true,
		}, {
			description: "trace from configured http header",
			header:      http.Header{testTraceHeader: {testTraceparent}},
			expectValid: true,
		}, {
			description: "trace from configured http header when body is not wrp",
			isDecodable: true,
			header: http.Header{
				contentType:     {"text/plain"},
				testTraceHeader: {testTraceparent},
			},
			body:        []byte("hello"),
			expectValid: true,
		}, {
			description: "trace from configured http header when the request carries several messages",
			isDecodable: true,
			header: http.Header{
				contentType:     {"application/jsonl"},
				testTraceHeader: {testTraceparent},
			},
			body: bytes.Join([][]byte{
				encodeWRP(t, wrp.JSON, "traceparent: 00-00000000000000000000000000000001-0000000000000001-01"),
				encodeWRP(t, wrp.JSON, "traceparent: 00-00000000000000000000000000000002-0000000000000002-01"),
			}, []byte("\n")),
			expectValid: true,
		}, {
			description: "wrp message is ignored when not decodable",
			header:      http.Header{contentType: {msgpackType}},
			body:        encodeWRP(t, wrp.Msgpack, testTraceparent),
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			assert := assert.New(t)
			tracing := Tracing{headerPrefix: testTraceHeader}

			var sc trace.SpanContext
			var body []byte
			handler := EchoFirstTraceNodeInfo(tracing, tc.isDecodable)(
				http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
					sc = trace.SpanContextFromContext(r.Context())
					body, _ = io.ReadAll(r.Body)
				}),
			)

			r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(tc.body))
			maps.Copy(r.Header, tc.header)
			handler.ServeHTTP(httptest.NewRecorder(), r)

			assert.Equal(tc.expectValid, sc.IsValid())
			if tc.expectValid {
				assert.Equal(testTraceID, sc.TraceID().String())
			}
			assert.Equal(len(tc.body), len(body))
		})
	}
}
