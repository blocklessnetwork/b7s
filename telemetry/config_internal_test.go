package telemetry

import (
	"testing"
	"time"

	"github.com/armon/go-metrics/prometheus"
	"github.com/stretchr/testify/require"

	"github.com/blocklessnetwork/b7s/models/blockless"
)

func TestConfig_NodeRole(t *testing.T) {

	const role = blockless.WorkerNode

	cfg := Config{
		Role: blockless.HeadNode,
	}

	WithNodeRole(role)(&cfg)
	require.Equal(t, role, cfg.Role)
}

func TestConfig_BatchTraceTimeout(t *testing.T) {

	const timeout = time.Minute

	var cfg Config
	cfg.Trace.ExporterBatchTimeout = time.Second

	WithBatchTraceTimeout(timeout)(&cfg)
	require.Equal(t, timeout, cfg.Trace.ExporterBatchTimeout)
}

func TestConfig_ID(t *testing.T) {

	const id = "dummy-id"

	cfg := Config{
		ID: "super-legit-id",
	}

	WithID(id)(&cfg)
	require.Equal(t, id, cfg.ID)
}

func TestConfig_TracingGRPC(t *testing.T) {

	t.Run("enable GRPC tracing", func(t *testing.T) {

		const endpoint = "localhost:1234"

		var cfg Config
		WithGRPCTracing(endpoint)(&cfg)
		require.Equal(t, endpoint, cfg.Trace.GRPC.Endpoint)
		require.True(t, cfg.Trace.GRPC.Enabled)
	})
	t.Run("disable GRPC tracing", func(t *testing.T) {

		var cfg Config
		cfg.Trace.GRPC.Endpoint = "localhost:9876"
		WithGRPCTracing("")(&cfg)
		require.Empty(t, cfg.Trace.GRPC.Endpoint)
		require.False(t, cfg.Trace.GRPC.Enabled)
	})
}

func TestConfig_TracingHTTP(t *testing.T) {

	t.Run("enable HTTP tracing", func(t *testing.T) {

		const endpoint = "localhost:1234"

		var cfg Config
		WithHTTPTracing(endpoint)(&cfg)
		require.Equal(t, endpoint, cfg.Trace.HTTP.Endpoint)
		require.True(t, cfg.Trace.HTTP.Enabled)
	})
	t.Run("disable HTTP tracing", func(t *testing.T) {

		var cfg Config
		cfg.Trace.HTTP.Endpoint = "localhost:9876"
		WithHTTPTracing("")(&cfg)
		require.Empty(t, cfg.Trace.HTTP.Endpoint)
		require.False(t, cfg.Trace.HTTP.Enabled)
	})
}

func TestConfig_MetricCounters(t *testing.T) {

	var counters = []prometheus.CounterDefinition{
		{
			Name: []string{"random", "counter", "value", "1"},
			Help: "Dummy counter description",
		},
		{
			Name: []string{"generic", "counter", "value", "2"},
			Help: "Dummy counter description",
		},
	}

	var cfg Config
	WithCounters(counters)(&cfg)
	require.Equal(t, counters, cfg.Metrics.Counters)
}

func TestConfig_MetricSummaries(t *testing.T) {

	var summary = []prometheus.SummaryDefinition{
		{
			Name: []string{"random", "summary", "value", "1"},
			Help: "Dummy summary description #1",
		},
		{
			Name: []string{"generic", "summary", "value", "2"},
			Help: "Dummy summary description #2",
		},
	}

	var cfg Config
	WithSummaries(summary)(&cfg)
	require.Equal(t, summary, cfg.Metrics.Summaries)
}

func TestConfig_MetricGauges(t *testing.T) {

	var gauges = []prometheus.GaugeDefinition{
		{
			Name: []string{"random", "gauges", "value", "1"},
			Help: "Dummy gauges description",
		},
		{
			Name: []string{"generic", "gauges", "value", "2"},
			Help: "Dummy gauges description",
		},
	}

	var cfg Config
	WithGauges(gauges)(&cfg)
	require.Equal(t, gauges, cfg.Metrics.Gauges)
}
