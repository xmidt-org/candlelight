/**
 *  Copyright (c) 2021  Comcast Cable Communications Management, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
package candlelight

import (
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/exporters/trace/zipkin"
	"go.opentelemetry.io/otel/label"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

const (
	zipkinName                   = "zipkin"
	jaegarName                   = "jaegar"
	traceProviderType            = "type"
	traceProviderEndpoint        = "endpoint"
	traceProviderSkipTraceExport = "skipTraceExport"
)

/***
 1. will be responsible for creating the traceprovider and setting it back to  opentelemetry
		viper will be having all the fields like  type,endpoint and skipTraceExport
 2. supported traceProviders are zipkin,jaegar and stdout
 3. set skipTraceExport = true if you don't want to print the span and tracer information in stdout
*/
func ConfigureTracerProvider(v *viper.Viper, logger log.Logger, applicationName string) error {
	// added  this condition if traceProvider  is missing in properties file then v will be coming as nil
	if v == nil {
		return  errors.New("viper instance can't be nil")

	}
	var traceProviderName = v.GetString(traceProviderType)

	switch traceProviderName {

	case zipkinName:
		err := zipkin.InstallNewPipeline(
			v.GetString(traceProviderEndpoint),
			applicationName,
			zipkin.WithSDK(&sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		)
		if err != nil {
			logger.Log("message", "failed to create zipkin pipeline", "err", err)
			return err
		}
	case jaegarName:
		flush, err := jaeger.InstallNewPipeline(
			jaeger.WithCollectorEndpoint(v.GetString(traceProviderEndpoint)),
			jaeger.WithProcess(jaeger.Process{
				ServiceName: applicationName,
				Tags: []label.KeyValue{
					label.String("exporter", jaegarName),
				},
			}),
			jaeger.WithSDK(&sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		)
		if err != nil {
			logger.Log("message", "failed to create jaegar pipeline", "err", err)
			return err
		}
		defer flush()
	default:
		var skipTraceExport = v.GetBool(traceProviderSkipTraceExport)
		var option stdout.Option
		if skipTraceExport {
			option = stdout.WithoutTraceExport()
		} else {
			option = stdout.WithPrettyPrint()
		}
		otExporter, err := stdout.NewExporter(option)
		if err != nil {
			logger.Log("message", "failed to create stdout exporter", "err", err)
			return
		}
		traceProvider := sdktrace.NewTracerProvider(sdktrace.WithSyncer(otExporter))
		otel.SetTracerProvider(traceProvider)
	}
	return  nil
}
