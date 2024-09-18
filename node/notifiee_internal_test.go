package node

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/blocklessnetwork/b7s/host"
	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/testing/mocks"
)

func TestNode_Notifiee(t *testing.T) {

	var (
		logger          = mocks.NoopLogger
		functionHandler = mocks.BaselineFStore(t)
	)

	server, err := host.New(mocks.NoopLogger, loopback, 0)
	require.NoError(t, err)

	var (
		storedPeer bool
	)

	store := mocks.BaselineStore(t)
	// Override the peerstore methods so we know if the node correctly handled incoming connection.
	store.SavePeerFunc = func(context.Context, blockless.Peer) error {
		storedPeer = true
		return nil
	}

	node, err := New(logger, server, store, functionHandler, WithRole(blockless.HeadNode))
	require.NoError(t, err)

	client, err := host.New(mocks.NoopLogger, loopback, 0)
	require.NoError(t, err)

	hostAddNewPeer(t, client, node.host)

	serverInfo := hostGetAddrInfo(t, server)
	err = client.Connect(context.Background(), *serverInfo)
	require.NoError(t, err)

	// Verify that peer store was updated.
	require.True(t, storedPeer)
}
