# Oragent Exporter for OpenTelemetry

This exporter supports sending OpenTelemetry data to [Oragent](https://orijtech.com/oragent/)

Supported pipeline types: traces, metrics

### Configuration options

- `api_key` (required): API key uses for authentication.

Other configs can be provided as the same manner with [OpenCensus Exporter](https://github.com/open-telemetry/opentelemetry-collector/blob/master/exporter/opencensusexporter/README.md)

### Example

```yaml
exporters:
  oragent:
    api_key: "<YOUR_API_KEY>"
```
