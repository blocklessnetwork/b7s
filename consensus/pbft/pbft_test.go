package pbft

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/blocklessnetworking/b7s/host"
	"github.com/blocklessnetworking/b7s/testing/mocks"
)

const (
	loopback = "127.0.0.1"
)

func newDummyReplica(t *testing.T) *Replica {
	t.Helper()

	var (
		logger    = mocks.NoopLogger
		executor  = mocks.BaselineExecutor(t)
		clusterID = mocks.GenericUUID.String()
		peers     = mocks.GenericPeerIDs[:4]
	)

	host, err := host.New(logger, loopback, 0)
	require.NoError(t, err)

	replica, err := NewReplica(logger, host, executor, peers, clusterID)
	require.NoError(t, err)

	return replica
}
