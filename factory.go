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

package oragentexporter

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"go.opentelemetry.io/collector/exporter/opencensusexporter"
)

const (
	// The value of "type" key in configuration.
	typeStr                = "oragent"
	defaultWriteBufferSize = 512 * 1024
)

// NewFactory creates a factory for OTLP exporter.
func NewFactory() component.ExporterFactory {
	return exporterhelper.NewFactory(
		typeStr,
		createDefaultConfig,
		exporterhelper.WithTraces(createTraceExporter),
		exporterhelper.WithMetrics(createMetricsExporter))
}

func createDefaultConfig() configmodels.Exporter {
	return &Config{
		Config: opencensusexporter.Config{
			ExporterSettings: configmodels.ExporterSettings{
				TypeVal: typeStr,
				NameVal: typeStr,
			},
			GRPCClientSettings: configgrpc.GRPCClientSettings{
				WriteBufferSize: defaultWriteBufferSize,
			},
			NumWorkers: 2,
		},
	}
}

func createTraceExporter(ctx context.Context, params component.ExporterCreateParams, config configmodels.Exporter) (component.TracesExporter, error) {
	return newTraceExporter(ctx, config.(*Config), params.Logger)
}

func createMetricsExporter(ctx context.Context, params component.ExporterCreateParams, config configmodels.Exporter) (component.MetricsExporter, error) {
	return newMetricsExporter(ctx, config.(*Config), params.Logger)
}
