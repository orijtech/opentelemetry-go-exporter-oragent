package main

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenterror"
	"go.opentelemetry.io/collector/service/defaultcomponents"

	"github.com/orijtech/opentelemetry-go-exporter-oragent/otelcollector"
)

func components() (component.Factories, error) {
	var errs []error
	factories, err := defaultcomponents.Components()
	if err != nil {
		return component.Factories{}, err
	}

	exporters := make([]component.ExporterFactory, 0, len(factories.Exporters)+1)
	exporters = append(exporters, otelcollector.NewFactory())
	for _, exporter := range factories.Exporters {
		exporters = append(exporters, exporter)
	}

	factories.Exporters, err = component.MakeExporterFactoryMap(exporters...)
	if err != nil {
		errs = append(errs, err)
	}

	return factories, componenterror.CombineErrors(errs)
}
