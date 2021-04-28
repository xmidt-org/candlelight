package candlelight

import (
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/propagation"
)

// Unmarshal helps load tracing components from configuration.
type Unmarshal struct {
	// AppName is the name of the application to be traced.
	AppName string

	// Key is the viper configuration key containing the tracing options.
	Key string
}

func (u Unmarshal) New(v *viper.Viper) (*Tracing, error) {
	var tracing = Tracing{
		Propagator: propagation.TraceContext{},
	}
	var traceConfig Config
	err := v.UnmarshalKey(u.Key, &traceConfig)
	if err != nil {
		return nil, err
	}
	traceConfig.ApplicationName = u.AppName
	tracerProvider, err := ConfigureTracerProvider(traceConfig)
	if err != nil {
		return nil, err
	}
	if len(traceConfig.Provider) != 0 && traceConfig.Provider != DefaultTracerProvider {
		tracing.Enabled = true
	}
	tracing.TracerProvider = tracerProvider
	return &tracing, nil
}
