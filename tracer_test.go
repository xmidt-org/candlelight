package candlelight

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestExtractTraceInformation(t *testing.T) {
	traceId, spanId, valid := ExtractTraceInformation(context.TODO())
	assert.False(t, valid)
	assert.Equal(t, traceId, "00000000000000000000000000000000")
	assert.Equal(t, spanId, "0000000000000000")
}

func TestInjectTraceInformation(t *testing.T) {
	headers := http.Header{}
	InjectTraceInformation(context.TODO(), headers)
	assert.Empty(t, headers)
}
