package candlelight

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/propagation"
)

func TestExtractTraceInformation(t *testing.T) {
	traceId, spanId := ExtractTraceInformation(context.TODO())
	assert.Equal(t, traceId, "00000000000000000000000000000000")
	assert.Equal(t, spanId, "0000000000000000")
}

func TestInjectTraceInformation(t *testing.T) {
	headers := http.Header{}
	InjectTraceInformation(context.TODO(), propagation.HeaderCarrier(headers))
	assert.Empty(t, headers)
}
