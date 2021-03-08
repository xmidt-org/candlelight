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
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestConfigureTracerProvider(t *testing.T) {

	testData := []struct {
		config Config
		err    error
	}{
		{
			config: Config{
				Provider: "jaeger",
			},
			err: errors.New("collectorEndpoint must not be empty"),
		},
		{
			config: Config{
				Provider: "jaeger",
				Endpoint: "http://localhost",
			},
			err: nil,
		},
		{
			config: Config{
				Provider: "Zipkin",
				Endpoint: "http://localhost",
			},
			err: nil,
		},
		{
			config: Config{
				Provider: "Zipkin",
			},
			err: errors.New("collector URL cannot be empty"),
		},
		{
			config: Config{
				Provider: "undefined",
			},
			err: nilProviderErr,
		},
		{
			config: Config{
				Provider: "stdOut",
			},
			err: nil,
		},
		{
			config: Config{
				Provider:        "stdoUt",
				SkipTraceExport: true,
			},
			err: nil,
		},
		{
			config: Config{},
			err:    nilProviderErr,
		},
	}
	for i, record := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var (
				assert = assert.New(t)
				_, err = ConfigureTracerProvider(record.config)
			)
			assert.Equal(record.err, err)

		})
	}

}
