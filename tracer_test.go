// SPDX-FileCopyrightText: 2021 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package candlelight

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/propagation"
)

func TestExtractTraceInfo(t *testing.T) {
	assert := assert.New(t)
	traceId, spanId, ok := ExtractTraceInfo(context.TODO())
	assert.Equal(traceId, "00000000000000000000000000000000")
	assert.Equal(spanId, "0000000000000000")
	assert.False(ok)
}

func TestInjectTraceInfo(t *testing.T) {
	headers := http.Header{}
	InjectTraceInfo(context.TODO(), propagation.HeaderCarrier(headers))
	assert.Empty(t, headers)
}
