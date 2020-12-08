# Oragent OpenTelemetry Collector Exporter

## Config

Edit config with your api key.

## Running

```shell
$ go build
$ ./otel-test -config=./config.yaml
2020-12-08T22:27:16.188+0700	INFO	service/service.go:409	Starting Oragent OpenTelemetry Collector distribution...	{"Version": "1.0.0", "GitHash": "", "NumCPU": 8}
2020-12-08T22:27:16.188+0700	INFO	service/service.go:253	Setting up own telemetry...
2020-12-08T22:27:16.207+0700	INFO	service/telemetry.go:101	Serving Prometheus metrics	{"address": "localhost:8888", "level": 0, "service.instance.id": "5387b218-ee6c-46a9-a1d8-3fac22375640"}
2020-12-08T22:27:16.208+0700	INFO	service/service.go:290	Loading configuration...
2020-12-08T22:27:16.209+0700	INFO	service/service.go:301	Applying configuration...
2020-12-08T22:27:16.209+0700	INFO	service/service.go:322	Starting extensions...
2020-12-08T22:27:16.210+0700	INFO	builder/exporters_builder.go:306	Exporter is enabled.	{"component_kind": "exporter", "exporter": "oragent"}
2020-12-08T22:27:16.212+0700	INFO	builder/exporters_builder.go:306	Exporter is enabled.	{"component_kind": "exporter", "exporter": "logging"}
2020-12-08T22:27:16.212+0700	INFO	service/service.go:337	Starting exporters...
2020-12-08T22:27:16.212+0700	INFO	builder/exporters_builder.go:92	Exporter is starting...	{"component_kind": "exporter", "component_type": "oragent", "component_name": "oragent"}
2020-12-08T22:27:16.212+0700	INFO	builder/exporters_builder.go:97	Exporter started.	{"component_kind": "exporter", "component_type": "oragent", "component_name": "oragent"}
2020-12-08T22:27:16.212+0700	INFO	builder/exporters_builder.go:92	Exporter is starting...	{"component_kind": "exporter", "component_type": "logging", "component_name": "logging"}
2020-12-08T22:27:16.212+0700	INFO	builder/exporters_builder.go:97	Exporter started.	{"component_kind": "exporter", "component_type": "logging", "component_name": "logging"}
2020-12-08T22:27:16.212+0700	INFO	builder/pipelines_builder.go:207	Pipeline is enabled.	{"pipeline_name": "traces", "pipeline_datatype": "traces"}
2020-12-08T22:27:16.212+0700	INFO	builder/pipelines_builder.go:207	Pipeline is enabled.	{"pipeline_name": "metrics", "pipeline_datatype": "metrics"}
2020-12-08T22:27:16.212+0700	INFO	service/service.go:350	Starting processors...
2020-12-08T22:27:16.212+0700	INFO	builder/pipelines_builder.go:51	Pipeline is starting...	{"pipeline_name": "traces", "pipeline_datatype": "traces"}
2020-12-08T22:27:16.212+0700	INFO	builder/pipelines_builder.go:61	Pipeline is started.	{"pipeline_name": "traces", "pipeline_datatype": "traces"}
2020-12-08T22:27:16.212+0700	INFO	builder/pipelines_builder.go:51	Pipeline is starting...	{"pipeline_name": "metrics", "pipeline_datatype": "metrics"}
2020-12-08T22:27:16.212+0700	INFO	builder/pipelines_builder.go:61	Pipeline is started.	{"pipeline_name": "metrics", "pipeline_datatype": "metrics"}
2020-12-08T22:27:16.212+0700	INFO	builder/receivers_builder.go:235	Receiver is enabled.	{"component_kind": "receiver", "component_type": "opencensus", "component_name": "opencensus", "datatype": "traces"}
2020-12-08T22:27:16.212+0700	INFO	builder/receivers_builder.go:235	Receiver is enabled.	{"component_kind": "receiver", "component_type": "opencensus", "component_name": "opencensus", "datatype": "metrics"}
2020-12-08T22:27:16.212+0700	INFO	service/service.go:362	Starting receivers...
2020-12-08T22:27:16.212+0700	INFO	builder/receivers_builder.go:70	Receiver is starting...	{"component_kind": "receiver", "component_type": "opencensus", "component_name": "opencensus"}
2020-12-08T22:27:16.213+0700	INFO	builder/receivers_builder.go:75	Receiver started.	{"component_kind": "receiver", "component_type": "opencensus", "component_name": "opencensus"}
2020-12-08T22:27:16.213+0700	INFO	service/service.go:265	Everything is ready. Begin running and processing data.
```

Now the collector is up, you can feed data using any OpenCensus exporter, then the collector will push data to Oragent.
