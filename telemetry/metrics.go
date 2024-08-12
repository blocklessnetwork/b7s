package telemetry

import (
	"fmt"
	"net/http"
	"time"

	"github.com/armon/go-metrics"
	mp "github.com/armon/go-metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func initPrometheusRegistry() error {

	// registry := prometheus.NewRegistry()

	// po := collectors.ProcessCollectorOpts{}
	// procCollector := collectors.NewProcessCollector(po)

	// colls := []prometheus.Collector{
	// 	collectors.NewGoCollector(), // Add Go metrics.
	// 	procCollector,               // Add process metrics.
	// }

	// colls = append(colls, procCollector)

	// for _, col := range colls {
	// 	err := registry.Register(col)
	// 	if err != nil && !errors.As(err, &prometheus.AlreadyRegisteredError{}) {
	// 		return fmt.Errorf("could not register collector: %w", err)
	// 	}
	// }

	// globalRegistry = registry

	// TODO: Check - Do we want to predeclare some metrics?
	opts := mp.PrometheusOpts{
		Registerer: prometheus.DefaultRegisterer,
		Name:       "b7s",
	}

	sink, err := mp.NewPrometheusSinkFrom(opts)
	if err != nil {
		return fmt.Errorf("could not create prometheus sink: %w", err)
	}

	mcfg := &metrics.Config{
		ServiceName:          "b7s",
		FilterDefault:        true,
		EnableRuntimeMetrics: true,
		TimerGranularity:     time.Millisecond,
	}
	_, err = metrics.NewGlobal(mcfg, sink)
	if err != nil {
		return fmt.Errorf("could not initialize metrics: %w", err)
	}

	return nil
}

func GetMetricsHTTPHandler() http.Handler {

	opts := promhttp.HandlerOpts{
		Registry: prometheus.DefaultRegisterer,
	}

	return promhttp.HandlerFor(prometheus.DefaultGatherer, opts)
}

// TODO: think again, whether we want/need this done manually or just work with the default/global registry?
// Upside here is we manually add in what we want.

// func PrometheusRegisterer() prometheus.Registerer {
// 	return cmp.Or(prometheus.Registerer(globalRegistry), prometheus.DefaultRegisterer)
// }
//
// func PrometheusGatherer() prometheus.Gatherer {
// 	return cmp.Or(prometheus.Gatherer(globalRegistry), prometheus.DefaultGatherer)
// }
