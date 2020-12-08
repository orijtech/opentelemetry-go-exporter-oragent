// Copyright 2020 Orijtech, Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package otelcollector

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/config/configtest"
	"go.opentelemetry.io/collector/exporter/opencensusexporter"
)

func TestLoadConfig(t *testing.T) {
	factories, err := componenttest.ExampleComponents()
	assert.NoError(t, err)

	factory := NewFactory()
	factories.Exporters[typeStr] = factory
	cfg, err := configtest.LoadConfigFile(t, path.Join(".", "testdata", "config.yaml"), factories)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, cfg.Exporters["oragent"], factory.CreateDefaultConfig())
	assert.Equal(t,
		cfg.Exporters["oragent/2"],
		&Config{
			Config: opencensusexporter.Config{
				ExporterSettings: configmodels.ExporterSettings{
					NameVal: "oragent/2",
					TypeVal: "oragent",
				},
				GRPCClientSettings: configgrpc.GRPCClientSettings{
					Endpoint:        "oragent.orijtech.com:443",
					WriteBufferSize: defaultWriteBufferSize,
				},
				NumWorkers: 2,
			},
			APIKey: "foo",
		},
	)
}
