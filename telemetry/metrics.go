package telemetry

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// TODO: Maybe just use the prometheus default registry?
	globalRegistry *prometheus.Registry
	metricsOnce    sync.Once
)

func initPrometheusRegistry() error {

	registry := prometheus.NewRegistry()
	colls := []prometheus.Collector{
		collectors.NewGoCollector(),
	}

	po := collectors.ProcessCollectorOpts{}
	procCollector := collectors.NewProcessCollector(po)

	colls = append(colls, procCollector)

	for _, col := range colls {
		err := registry.Register(col)
		if err != nil && !errors.As(err, &prometheus.AlreadyRegisteredError{}) {
			return fmt.Errorf("could not register collector: %w", err)
		}
	}

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
