package candlelight

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"strconv"
	"testing"
)

func TestExtractTraceInformation(t *testing.T) {
	traceId, spanId := ExtractTraceInformation(context.TODO())
	assert.Equal(t, traceId, "00000000000000000000000000000000")
	assert.Equal(t, spanId, "0000000000000000")
}

func TestInjectTraceInformation(t *testing.T) {
	headers := http.Header{}
	InjectTraceInformation(context.TODO(), headers)
	assert.Empty(t, headers)
}

func TestExtractSpanIDAndTraceIDHeaderName(t *testing.T) {
	testData := []struct {
		config            Config
		spanIDHeaderName  string
		traceIDHeaderName string
	}{
		{Config{}, DefaultSpanIDHeaderName, DefaultTraceIDHeaderName},
		{
			config: Config{
				SpanIDHeaderName:  "SpanID",
				TraceIDHeaderName: "TraceId",
			},
			spanIDHeaderName:  "SpanID",
			traceIDHeaderName: "TraceId",
		},
	}

	for i, record := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var (
				assert                              = assert.New(t)
				spanIDHeaderName, traceIDHeaderName = ExtractSpanIDAndTraceIDHeaderName(record.config)
			)
			assert.Equal(record.spanIDHeaderName, spanIDHeaderName)
			assert.Equal(record.traceIDHeaderName, traceIDHeaderName)
		})
	}
}

func TestNewTraceConfig(t *testing.T) {
	testData := []struct {
		config        Config
		traceProvider trace.TracerProvider
		traceConfig   TraceConfig
	}{
		{Config{}, nil, NewTraceConfig(Config{}, nil)},
	}
	for i, record := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var (
				assert = assert.New(t)
				output = NewTraceConfig(record.config, record.traceProvider)
			)
			assert.Equal(record.traceConfig, output)
		})
	}

}
