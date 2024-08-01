package telemetry

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func initializeMetrics(cfg MetricsConfig) error {

	registry := prometheus.NewRegistry()

	colls := []prometheus.Collector{
		collectors.NewGoCollector(),
	}

	po := collectors.ProcessCollectorOpts{}
	procCollector := collectors.NewProcessCollector(po)

	colls = append(colls, procCollector)

	for _, col := range colls {
		err := registry.Register(col)
		if err != nil {
			return fmt.Errorf("could not register collector: %w", err)
		}
	}

	return nil

}

func GetHTTPHandler() http.Handler {

	var registry *prometheus.Registry

	opts := promhttp.HandlerOpts{
		Registry: registry,
	}

	return promhttp.HandlerFor(registry, opts)
}
