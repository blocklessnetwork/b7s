package telemetry

import (
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/armon/go-metrics"
	mp "github.com/armon/go-metrics/prometheus"
	"github.com/blocklessnetwork/b7s/consensus/pbft"
	"github.com/blocklessnetwork/b7s/consensus/raft"
	"github.com/blocklessnetwork/b7s/executor"
	"github.com/blocklessnetwork/b7s/fstore"
	"github.com/blocklessnetwork/b7s/host"
	"github.com/blocklessnetwork/b7s/node"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func initPrometheusRegistry() error {

	var (
		opts = mp.PrometheusOpts{
			Registerer:         prometheus.DefaultRegisterer,
			CounterDefinitions: counters(),
			SummaryDefinitions: summaries(),
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

func counters() []mp.CounterDefinition {

	counters := slices.Concat(
		node.Counters,
		host.Counters,
		fstore.Counters,
		executor.Counters,
	)
	prefixed := make([]mp.CounterDefinition, len(counters))

	for i := 0; i < len(counters); i++ {
		c := counters[i]
		c.Name = append([]string{metricPrefix}, c.Name...)
		prefixed[i] = c
	}

	return prefixed
}

func summaries() []mp.SummaryDefinition {

	summaries := slices.Concat(
		executor.Summaries,
		fstore.Summaries,
		pbft.Summaries,
		raft.Summaries,
	)
	prefixed := make([]mp.SummaryDefinition, len(summaries))

	for i := 0; i < len(summaries); i++ {
		s := summaries[i]
		s.Name = append([]string{metricPrefix}, s.Name...)
		prefixed[i] = s
	}

	return prefixed
}
