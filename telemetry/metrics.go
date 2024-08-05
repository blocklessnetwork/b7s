package telemetry

import (
	"cmp"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	globalRegistry *prometheus.Registry
	metricsOnce    sync.Once
)

func initPrometheusRegistry() error {

	registry := prometheus.NewRegistry()

	po := collectors.ProcessCollectorOpts{}
	procCollector := collectors.NewProcessCollector(po)

	colls := []prometheus.Collector{
		collectors.NewGoCollector(), // Add Go metrics.
		procCollector,               // Add process metrics.
	}

	colls = append(colls, procCollector)

	for _, col := range colls {
		err := registry.Register(col)
		if err != nil && !errors.As(err, &prometheus.AlreadyRegisteredError{}) {
			return fmt.Errorf("could not register collector: %w", err)
		}
	}

	globalRegistry = registry

	return nil
}

func GetMetricsHTTPHandler() http.Handler {

	metricsOnce.Do(func() {

		// TODO: Handle error.
		err := initPrometheusRegistry()
		_ = err
	})

	opts := promhttp.HandlerOpts{
		Registry: globalRegistry,
	}

	return promhttp.HandlerFor(globalRegistry, opts)
}

// TODO: think again, whether we want/need this done manually or just work with the default/global registry?
// Upside here is we manually add in what we want.

func PrometheusRegisterer() prometheus.Registerer {
	return cmp.Or(prometheus.Registerer(globalRegistry), prometheus.DefaultRegisterer)
}

func PrometheusGatherer() prometheus.Gatherer {
	return cmp.Or(prometheus.Gatherer(globalRegistry), prometheus.DefaultGatherer)
}
