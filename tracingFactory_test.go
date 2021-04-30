package candlelight

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/propagation"
)

func TestNew(t *testing.T) {
	tcs := []struct {
		Description  string
		ShouldFail   bool
		ExpectIsNoop bool
		Config       Config
	}{
		{
			Description:  "Default means disabled tracing",
			ShouldFail:   false,
			ExpectIsNoop: true,
			Config:       Config{},
		},
		{
			Description: "Provider not found",
			ShouldFail:  true,
			Config: Config{
				Provider: "somethingElse",
			},
		},
		{
			Description:  "Stdout provider",
			ShouldFail:   false,
			ExpectIsNoop: false,
			Config: Config{
				Provider: "stdout",
			},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.Description, func(t *testing.T) {
			assert := assert.New(t)
			tracing, err := New(tc.Config)
			if tc.ShouldFail {
				assert.NotNil(err)
			} else {
				assert.Nil(err)
				assert.Equal(tc.ExpectIsNoop, tracing.IsNoop())
			}
		})
	}
}

func TestTracing(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	tracerProvider, err := ConfigureTracerProvider(Config{Provider: "stdout"})
	require.Nil(err)
	propagator := propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{})

	tracing := Tracing{
		tracerProvider: tracerProvider,
		propagator:     propagator,
	}

	assert.False(tracing.IsNoop())
	assert.NotNil(tracing.TracerProvider())
	assert.Equal(tracerProvider, tracing.TracerProvider())
	assert.Equal(propagator, tracing.Propagator())
}

func TestTracingDefault(t *testing.T) {
	assert := assert.New(t)
	var tracing Tracing
	assert.True(tracing.IsNoop())
	assert.NotNil(tracing.TracerProvider())
	assert.NotNil(tracing.Propagator())
}
