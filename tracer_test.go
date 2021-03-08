package candlelight

import (
	"context"
	"github.com/stretchr/testify/assert"
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
		headerConfig      HeaderConfig
		spanIDHeaderName  string
		traceIDHeaderName string
	}{
		{HeaderConfig{}, DefaultSpanIDHeaderName, DefaultTraceIDHeaderName},
		{
			headerConfig: HeaderConfig{
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
				spanIDHeaderName, traceIDHeaderName = ExtractSpanIDAndTraceIDHeaderName(record.headerConfig)
			)
			assert.Equal(record.spanIDHeaderName, spanIDHeaderName)
			assert.Equal(record.traceIDHeaderName, traceIDHeaderName)
		})
	}
}
