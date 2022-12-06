// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kafkareceiver

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/configtls"
	"go.opentelemetry.io/collector/confmap/confmaptest"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/kafkaexporter"
)

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	cm, err := confmaptest.LoadConf(filepath.Join("testdata", "config.yaml"))
	require.NoError(t, err)

	tests := []struct {
		id          component.ID
		expected    component.Config
		expectedErr error
	}{
		{
			id: component.NewIDWithName(typeStr, ""),
			expected: &Config{
				ReceiverSettings: config.NewReceiverSettings(component.NewID(typeStr)),
				Topic:            "spans",
				Encoding:         "otlp_proto",
				Brokers:          []string{"foo:123", "bar:456"},
				ClientID:         "otel-collector",
				GroupID:          "otel-collector",
				Authentication: kafkaexporter.Authentication{
					TLS: &configtls.TLSClientSetting{
						TLSSetting: configtls.TLSSetting{
							CAFile:   "ca.pem",
							CertFile: "cert.pem",
							KeyFile:  "key.pem",
						},
					},
				},
				Metadata: kafkaexporter.Metadata{
					Full: true,
					Retry: kafkaexporter.MetadataRetry{
						Max:     10,
						Backoff: time.Second * 5,
					},
				},
				AutoCommit: AutoCommit{
					Enable:   true,
					Interval: 1 * time.Second,
				},
			},
		},
		{

			id: component.NewIDWithName(typeStr, "logs"),
			expected: &Config{
				ReceiverSettings: config.NewReceiverSettings(component.NewID(typeStr)),
				Topic:            "logs",
				Encoding:         "direct",
				Brokers:          []string{"coffee:123", "foobar:456"},
				ClientID:         "otel-collector",
				GroupID:          "otel-collector",
				Authentication: kafkaexporter.Authentication{
					TLS: &configtls.TLSClientSetting{
						TLSSetting: configtls.TLSSetting{
							CAFile:   "ca.pem",
							CertFile: "cert.pem",
							KeyFile:  "key.pem",
						},
					},
				},
				Metadata: kafkaexporter.Metadata{
					Full: true,
					Retry: kafkaexporter.MetadataRetry{
						Max:     10,
						Backoff: time.Second * 5,
					},
				},
				AutoCommit: AutoCommit{
					Enable:   true,
					Interval: 1 * time.Second,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.id.String(), func(t *testing.T) {
			factory := NewFactory()
			cfg := factory.CreateDefaultConfig()

			sub, err := cm.Sub(tt.id.String())
			require.NoError(t, err)
			require.NoError(t, component.UnmarshalConfig(sub, cfg))

			assert.NoError(t, component.ValidateConfig(cfg))
			assert.Equal(t, tt.expected, cfg)
		})
	}
}