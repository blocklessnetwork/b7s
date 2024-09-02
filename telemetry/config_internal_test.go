package telemetry

import (
	"testing"
	"time"

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
	t.Run("enable GRPC tracing", func(t *testing.T) {

		var cfg Config
		WithGRPCTracing("")
		require.Empty(t, cfg.Trace.GRPC.Endpoint)
		require.False(t, cfg.Trace.GRPC.Enabled)
	})
}
