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
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/consumer/pdata"
	"go.uber.org/zap"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/configtls"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/receiver/opencensusreceiver"
	"go.opentelemetry.io/collector/testutil"
)

func TestSendTraces(t *testing.T) {
	sink := new(consumertest.TracesSink)
	rFactory := opencensusreceiver.NewFactory()
	rCfg := rFactory.CreateDefaultConfig().(*opencensusreceiver.Config)
	endpoint := testutil.GetAvailableLocalAddress(t)
	rCfg.GRPCServerSettings.NetAddr.Endpoint = endpoint
	params := component.ReceiverCreateParams{Logger: zap.NewNop()}
	recv, err := rFactory.CreateTracesReceiver(context.Background(), params, rCfg, sink)
	assert.NoError(t, err)
	assert.NoError(t, recv.Start(context.Background(), componenttest.NewNopHost()))
	t.Cleanup(func() {
		assert.NoError(t, recv.Shutdown(context.Background()))
	})

	factory := NewFactory()
	cfg := factory.CreateDefaultConfig().(*Config)
	cfg.GRPCClientSettings = configgrpc.GRPCClientSettings{
		Endpoint: endpoint,
		TLSSetting: configtls.TLSClientSetting{
			Insecure: true,
		},
	}
	cfg.NumWorkers = 2
	cfg.APIKey = "oragent-api-key"
	exp, err := factory.CreateTracesExporter(context.Background(), component.ExporterCreateParams{Logger: zap.NewNop()}, cfg)
	require.NoError(t, err)
	require.NotNil(t, exp)
	host := componenttest.NewNopHost()
	require.NoError(t, exp.Start(context.Background(), host))
	t.Cleanup(func() {
		assert.NoError(t, exp.Shutdown(context.Background()))
	})

	size := 2
	numSpans := 5
	traces := make([]pdata.Traces, size)
	for i := range traces {
		td := pdata.NewTraces()
		resourceSpans := td.ResourceSpans()
		resourceSpans.Resize(1)
		resourceSpans.At(0).InstrumentationLibrarySpans().Resize(1)
		resourceSpans.At(0).InstrumentationLibrarySpans().At(0).Spans().Resize(numSpans)
		for j := 0; j < numSpans; j++ {
			span := resourceSpans.At(0).InstrumentationLibrarySpans().At(0).Spans().At(j)
			span.SetName("oragent")
		}
		traces[i] = td
	}

	for _, td := range traces {
		assert.NoError(t, exp.ConsumeTraces(context.Background(), td))
	}

	testutil.WaitFor(t, func() bool {
		return len(sink.AllTraces()) == size
	})
	gotTraces := sink.AllTraces()
	require.Len(t, gotTraces, size)
	for i := range gotTraces {
		assert.Equal(t, traces[i], gotTraces[i])
		assert.Equal(t, numSpans, gotTraces[i].SpanCount())
	}
}
