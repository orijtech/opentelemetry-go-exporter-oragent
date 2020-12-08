# Oragent Exporter for OpenTelemetry Collector

This exporter supports sending OpenTelemetry data to [Oragent](https://orijtech.com/oragent/)

Supported pipeline types: traces, metrics

### Configuration options

- `api_key` (required): API key uses for authentication.

`api_key` can be set both from configuration file or from environment variable `ORAGENT_API_KEY`. If both are specified,
the value from configuration file will be used.

Other configs can be provided as the same manner with [OpenCensus Exporter](https://github.com/open-telemetry/opentelemetry-collector/blob/master/exporter/opencensusexporter/README.md)

### Example

```yaml
exporters:
  oragent:
    api_key: "<YOUR_API_KEY>"
```
