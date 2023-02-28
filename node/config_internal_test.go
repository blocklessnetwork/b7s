package node

import (
	"testing"
	"time"

	"github.com/blocklessnetworking/b7s/models/blockless"
	"github.com/blocklessnetworking/b7s/testing/mocks"

	"github.com/stretchr/testify/require"
)

func TestConfig_NodeRole(t *testing.T) {

	const role = blockless.WorkerNode

	cfg := Config{
		Role: blockless.WorkerNode,
	}

	WithRole(role)(&cfg)
	require.Equal(t, role, cfg.Role)
}

func TestConfig_Topic(t *testing.T) {

	const topic = "super-secret-topic"

	cfg := Config{
		Topic: "",
	}

	WithTopic(topic)(&cfg)
	require.Equal(t, topic, cfg.Topic)
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
