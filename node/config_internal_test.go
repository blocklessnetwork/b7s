package node

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/testing/mocks"
)

func TestConfig_NodeRole(t *testing.T) {

	const role = blockless.WorkerNode

	cfg := Config{
		Role: blockless.HeadNode,
	}

	WithRole(role)(&cfg)
	require.Equal(t, role, cfg.Role)
}

func TestConfig_Topic(t *testing.T) {

	topics := []string{"super-secret-topic"}

	cfg := Config{
		Topics: []string{},
	}

	WithTopics(topics)(&cfg)
	require.Equal(t, topics, cfg.Topics)
}

func TestConfig_Executor(t *testing.T) {

	executor := mocks.BaselineExecutor(t)

	cfg := Config{
		Execute: nil,
	}

	WithExecutor(executor)(&cfg)

	require.Equal(t, executor, cfg.Execute)
}

func TestConfig_HealthInterval(t *testing.T) {

	const interval = 30 * time.Second

	cfg := Config{
		HealthInterval: 0,
	}

	WithHealthInterval(interval)(&cfg)

	require.Equal(t, interval, cfg.HealthInterval)
}

func TestConfig_RollCallTimeout(t *testing.T) {

	const timeout = 10 * time.Second

	cfg := Config{
		RollCallTimeout: 0,
	}

	WithRollCallTimeout(timeout)(&cfg)

	require.Equal(t, timeout, cfg.RollCallTimeout)
}

func TestConfig_ExecutionTimeout(t *testing.T) {

	const timeout = 10 * time.Second

	cfg := Config{
		ExecutionTimeout: 0,
	}

	WithExecutionTimeout(timeout)(&cfg)

	require.Equal(t, timeout, cfg.ExecutionTimeout)
}

func TestConfig_Concurrency(t *testing.T) {

	const concurrency = uint(10)

	cfg := Config{
		Concurrency: 0,
	}

	WithConcurrency(concurrency)(&cfg)

	require.Equal(t, concurrency, cfg.Concurrency)
}
