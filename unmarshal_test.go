package candlelight

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshal(t *testing.T) {
	v := viper.New()
	v.Set("tracing.provider", "stdout")

	vNotFound := viper.New()
	vNotFound.Set("tracing.provider", "somethingElse")
	tcs := []struct {
		Description          string
		ShouldFail           bool
		ExpectTracingEnabled bool
		V                    *viper.Viper
	}{
		{
			Description:          "Default means disabled tracing",
			ShouldFail:           false,
			ExpectTracingEnabled: false,
			V:                    viper.New(),
		},
		{
			Description: "Provider not found",
			ShouldFail:  true,
			V:           vNotFound,
		},
		{
			Description:          "Stdout provider",
			ShouldFail:           false,
			ExpectTracingEnabled: true,
			V:                    v,
		},
	}
	u := Unmarshal{
		AppName: "testing",
		Key:     "tracing",
	}
	for _, tc := range tcs {
		t.Run(tc.Description, func(t *testing.T) {
			assert := assert.New(t)
			tracing, err := u.New(tc.V)
			if tc.ShouldFail {
				assert.NotNil(err)
				assert.Nil(tracing)
			} else {
				assert.Nil(err)
				assert.NotNil(tracing)
				assert.Equal(tc.ExpectTracingEnabled, tracing.Enabled)
			}
		})
	}
}
