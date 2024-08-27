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

func initPrometheusRegistry(cfg MetricsConfig) error {

	var (
		opts = mp.PrometheusOpts{
			Registerer:         prometheus.DefaultRegisterer,
			CounterDefinitions: counters(cfg.Counters),
			SummaryDefinitions: summaries(cfg.Summaries),
			GaugeDefinitions:   gauges(cfg.Gauges),
		}
	)

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

func counters(counters []mp.CounterDefinition) []mp.CounterDefinition {

	prefixed := make([]mp.CounterDefinition, len(counters))

	for i := 0; i < len(counters); i++ {
		c := counters[i]
		c.Name = append([]string{metricPrefix}, c.Name...)
		prefixed[i] = c
	}

	return prefixed
}

func summaries(summaries []mp.SummaryDefinition) []mp.SummaryDefinition {

	prefixed := make([]mp.SummaryDefinition, len(summaries))

	for i := 0; i < len(summaries); i++ {
		s := summaries[i]
		s.Name = append([]string{metricPrefix}, s.Name...)
		prefixed[i] = s
	}

	return prefixed
}

func gauges(gauges []mp.GaugeDefinition) []mp.GaugeDefinition {

	// Right now we have a single gauge - node info.
	// gauges := node.Gauges
	prefixed := make([]mp.GaugeDefinition, len(gauges))

	for i := 0; i < len(gauges); i++ {
		g := gauges[i]
		g.Name = append([]string{metricPrefix}, g.Name...)
		prefixed[i] = g
	}

	return prefixed
}
