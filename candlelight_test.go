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
	"bytes"
	"github.com/spf13/viper"
	"github.com/xmidt-org/webpa-common/logging"
	"go.opentelemetry.io/otel"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigureTraceProviderStdOut(t *testing.T) {
	var stdoutConfig = []byte(`type: stdout
skipTraceExport: true
`)
	var stdoutViper = viper.New()
	stdoutViper.SetConfigType("yaml")
	stdoutViper.ReadConfig(bytes.NewBuffer(stdoutConfig))
	ConfigureTracerProvider(stdoutViper, logging.DefaultLogger(), "stdOutTestCase")
	assert.NotNil(t, otel.GetTracerProvider())
}

func TestConfigureTraceProviderJaegar(t *testing.T) {
	var jaegarConfig = []byte(`type: jaegar
endpoint: http://localhost
`)
	var viper = viper.New()
	viper.SetConfigType("yaml")
	viper.ReadConfig(bytes.NewBuffer(jaegarConfig))
	ConfigureTracerProvider(viper, logging.DefaultLogger(), "jaegarTestCase")
	assert.NotNil(t, otel.GetTracerProvider())
}

func TestConfigureTraceProviderZipkin(t *testing.T) {
	var zipkingConfig = []byte(`type: zipkin
endpoint: http://127.0.0.1/
`)
	var viper = viper.New()
	viper.SetConfigType("yaml")
	viper.ReadConfig(bytes.NewBuffer(zipkingConfig))
	ConfigureTracerProvider(viper, logging.DefaultLogger(), "ZipkinTestCase")
	assert.NotNil(t, otel.GetTracerProvider())
}


func TestConfigureTracerProviderWhenViperIsNil(t *testing.T) {
	err := ConfigureTracerProvider(nil, logging.DefaultLogger(), "NilViperTestCase")
	assert.NotNil(t, err)
}

func TestConfigureTraceProviderJaegarWhenEndpointIsNil(t *testing.T) {
	var jaegarConfig = []byte(`type: jaegar
`)
	var viper = viper.New()
	viper.SetConfigType("yaml")
	viper.ReadConfig(bytes.NewBuffer(jaegarConfig))
	err := ConfigureTracerProvider(viper, logging.DefaultLogger(), "jaegarTestCase")
	assert.NotNil(t, err)
}

func TestConfigureTraceProviderZipkinWhenEndpointIsNil(t *testing.T) {
	var zipkingConfig = []byte(`type: zipkin
`)
	var viper = viper.New()
	viper.SetConfigType("yaml")
	viper.ReadConfig(bytes.NewBuffer(zipkingConfig))
	err := ConfigureTracerProvider(viper, logging.DefaultLogger(), "ZipkinTestCase")
	assert.NotNil(t, err)
}

func TestConfigureTraceProviderStdOutWithoutSkipTraceExport(t *testing.T) {
	var stdoutConfig = []byte(`type: stdout
`)
	var stdoutViper = viper.New()
	stdoutViper.SetConfigType("yaml")
	stdoutViper.ReadConfig(bytes.NewBuffer(stdoutConfig))
	ConfigureTracerProvider(stdoutViper, logging.DefaultLogger(), "stdOutTestCase")
	assert.NotNil(t, otel.GetTracerProvider())
}

