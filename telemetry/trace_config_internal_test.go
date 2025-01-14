package telemetry

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/blessnetwork/b7s/models/blockless"
)

func TestTraceConfig_ID(t *testing.T) {

	const id = "dummy-id"

	cfg := TraceConfig{
		ID: "super-legit-id",
	}

	WithID(id)(&cfg)
	require.Equal(t, id, cfg.ID)
}

func TestTraceConfig_NodeRole(t *testing.T) {

	const role = blockless.WorkerNode

	cfg := TraceConfig{
		Role: blockless.HeadNode,
	}

	WithNodeRole(role)(&cfg)
	require.Equal(t, role, cfg.Role)
}

func TestTraceConfig_BatchTraceTimeout(t *testing.T) {

	const timeout = time.Minute

	var cfg TraceConfig
	cfg.ExporterBatchTimeout = time.Second

	WithBatchTraceTimeout(timeout)(&cfg)
	require.Equal(t, timeout, cfg.ExporterBatchTimeout)
}

func TestTraceConfig_TracingGRPC(t *testing.T) {

	t.Run("enable GRPC tracing", func(t *testing.T) {

		const endpoint = "localhost:1234"

		var cfg TraceConfig
		WithGRPCTracing(endpoint)(&cfg)
		require.Equal(t, endpoint, cfg.GRPC.Endpoint)
		require.True(t, cfg.GRPC.Enabled)
	})
	t.Run("disable GRPC tracing", func(t *testing.T) {

		var cfg TraceConfig
		cfg.GRPC.Endpoint = "localhost:9876"
		WithGRPCTracing("")(&cfg)
		require.Empty(t, cfg.GRPC.Endpoint)
		require.False(t, cfg.GRPC.Enabled)
	})
}

func TestTraceConfig_TracingHTTP(t *testing.T) {

	t.Run("enable HTTP tracing", func(t *testing.T) {

		const endpoint = "localhost:1234"

		var cfg TraceConfig
		WithHTTPTracing(endpoint)(&cfg)
		require.Equal(t, endpoint, cfg.HTTP.Endpoint)
		require.True(t, cfg.HTTP.Enabled)
	})
	t.Run("disable HTTP tracing", func(t *testing.T) {

		var cfg TraceConfig
		cfg.HTTP.Endpoint = "localhost:9876"
		WithHTTPTracing("")(&cfg)
		require.Empty(t, cfg.HTTP.Endpoint)
		require.False(t, cfg.HTTP.Enabled)
	})
}
