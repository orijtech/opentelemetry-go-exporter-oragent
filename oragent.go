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

// Derived from OpenTelemetry exporter/opencensusexporter
// https://github.com/open-telemetry/opentelemetry-collector/tree/master/exporter/opencensusexporter

package oragentexporter

import (
	"context"
	"errors"
	"fmt"
	"os"

	commonpb "github.com/census-instrumentation/opencensus-proto/gen-go/agent/common/v1"
	agentmetricspb "github.com/census-instrumentation/opencensus-proto/gen-go/agent/metrics/v1"
	agenttracepb "github.com/census-instrumentation/opencensus-proto/gen-go/agent/trace/v1"
	resourcepb "github.com/census-instrumentation/opencensus-proto/gen-go/resource/v1"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer/pdata"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"go.opentelemetry.io/collector/translator/internaldata"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const authHeader = "x-oragent-auth"

type tracesClient struct {
	agenttracepb.TraceService_ExportClient
	cancelFunc context.CancelFunc
}

type metricsClient struct {
	agentmetricspb.MetricsService_ExportClient
	cancelFunc context.CancelFunc
}

type orExporter struct {
	cfg              *Config
	traceSvcClient   agenttracepb.TraceServiceClient
	metricsSvcClient agentmetricspb.MetricsServiceClient
	tClientCh        chan *tracesClient
	mClientCh        chan *metricsClient
	conn             *grpc.ClientConn
	metadata         metadata.MD
}

func newOrExporter(ctx context.Context, cfg *Config) (*orExporter, error) {
	if cfg.Endpoint == "" {
		return nil, errors.New("Oragent exporter cfg requires an Endpoint")
	}

	if cfg.NumWorkers <= 0 {
		return nil, errors.New("Oragent exporter cfg requires at least one worker")
	}

	if apiKey := os.Getenv("ORAGENT_API_KEY"); apiKey != "" {
		cfg.APIKey = apiKey
	}
	if cfg.APIKey == "" {
		return nil, errors.New("Oragent exporter cfg requires api key")
	}
	if cfg.Headers == nil {
		cfg.Headers = make(map[string]string)
	}
	cfg.Headers[authHeader] = cfg.APIKey

	dialOpts, err := cfg.GRPCClientSettings.ToDialOptions()
	if err != nil {
		return nil, err
	}

	clientConn, err := grpc.DialContext(ctx, cfg.GRPCClientSettings.Endpoint, dialOpts...)
	if err != nil {
		return nil, err
	}

	return &orExporter{
		cfg:      cfg,
		conn:     clientConn,
		metadata: metadata.New(cfg.GRPCClientSettings.Headers),
	}, nil
}

func (or *orExporter) shutdown(context.Context) error {
	if or.tClientCh != nil {
		for i := 0; i < or.cfg.NumWorkers; i++ {
			<-or.tClientCh
		}
		close(or.tClientCh)
	}
	if or.mClientCh != nil {
		for i := 0; i < or.cfg.NumWorkers; i++ {
			<-or.mClientCh
		}
		close(or.mClientCh)
	}
	return or.conn.Close()
}

func newTraceExporter(ctx context.Context, cfg *Config, logger *zap.Logger) (component.TracesExporter, error) {
	ore, err := newOrExporter(ctx, cfg)
	if err != nil {
		return nil, err
	}
	ore.traceSvcClient = agenttracepb.NewTraceServiceClient(ore.conn)
	ore.tClientCh = make(chan *tracesClient, cfg.NumWorkers)
	for i := 0; i < cfg.NumWorkers; i++ {
		ore.tClientCh <- nil
	}

	return exporterhelper.NewTraceExporter(
		cfg,
		logger,
		ore.pushTraceData,
		exporterhelper.WithShutdown(ore.shutdown))
}

func newMetricsExporter(ctx context.Context, cfg *Config, logger *zap.Logger) (component.MetricsExporter, error) {
	ore, err := newOrExporter(ctx, cfg)
	if err != nil {
		return nil, err
	}
	ore.metricsSvcClient = agentmetricspb.NewMetricsServiceClient(ore.conn)
	ore.mClientCh = make(chan *metricsClient, cfg.NumWorkers)
	for i := 0; i < cfg.NumWorkers; i++ {
		ore.mClientCh <- nil
	}

	return exporterhelper.NewMetricsExporter(
		cfg,
		logger,
		ore.pushMetricsData,
		exporterhelper.WithShutdown(ore.shutdown))
}

func (or *orExporter) pushTraceData(_ context.Context, td pdata.Traces) (int, error) {
	spanCount := td.SpanCount()
	tClient, ok := <-or.tClientCh
	if !ok {
		return spanCount, errors.New("failed to push traces, Oragent exporter was already stopped")
	}
	defer func() {
		or.tClientCh <- tClient
	}()

	if tClient == nil {
		c, err := or.newTracesClient()
		if err != nil {
			or.tClientCh <- nil
			return spanCount, err
		}
		tClient = c
	}

	for _, octd := range internaldata.TraceDataToOC(td) {
		// OC protocol requires a Node for the initial message.
		node := octd.Node
		if node == nil {
			node = &commonpb.Node{}
		}
		resource := octd.Resource
		if resource == nil {
			resource = &resourcepb.Resource{}
		}
		req := &agenttracepb.ExportTraceServiceRequest{
			Spans:    octd.Spans,
			Resource: resource,
			Node:     node,
		}
		if err := tClient.Send(req); err != nil {
			tClient.cancelFunc()
			tClient = nil
			return spanCount, err
		}
	}
	return 0, nil
}

func (or *orExporter) pushMetricsData(_ context.Context, md pdata.Metrics) (int, error) {
	_, mpc := md.MetricAndDataPointCount()
	mClient, ok := <-or.mClientCh
	if !ok {
		return mpc, errors.New("failed to push metrics, Oragent exporter was already stopped")
	}
	defer func() {
		or.mClientCh <- mClient
	}()

	if mClient == nil {
		c, err := or.newMetricsClient()
		if err != nil {
			return mpc, err
		}
		mClient = c
	}

	for _, ocmd := range internaldata.MetricsToOC(md) {
		// OC protocol requires a Node for the initial message.
		node := ocmd.Node
		if node == nil {
			node = &commonpb.Node{}
		}
		resource := ocmd.Resource
		if resource == nil {
			resource = &resourcepb.Resource{}
		}
		req := &agentmetricspb.ExportMetricsServiceRequest{
			Metrics:  ocmd.Metrics,
			Resource: resource,
			Node:     node,
		}
		if err := mClient.Send(req); err != nil {
			mClient.cancelFunc()
			mClient = nil
			return mpc, err
		}
	}
	return 0, nil
}

func (or *orExporter) newTracesClient() (*tracesClient, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	ctx = metadata.NewOutgoingContext(ctx, or.metadata)

	tc, err := or.traceSvcClient.Export(ctx)
	if err != nil {
		cancelFunc()
		return nil, fmt.Errorf("TraceServiceClient: %w", err)
	}
	return &tracesClient{tc, cancelFunc}, nil
}

func (or *orExporter) newMetricsClient() (*metricsClient, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	ctx = metadata.NewOutgoingContext(ctx, or.metadata)

	mc, err := or.metricsSvcClient.Export(ctx)
	if err != nil {
		cancelFunc()
		return nil, fmt.Errorf("MetricsServiceClient: %w", err)
	}
	return &metricsClient{mc, cancelFunc}, nil
}
