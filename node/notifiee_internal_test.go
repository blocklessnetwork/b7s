package node

import (
	"context"
	"testing"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/require"

	"github.com/blocklessnetworking/b7s/host"
	"github.com/blocklessnetworking/b7s/models/blockless"
	"github.com/blocklessnetworking/b7s/testing/mocks"
)

func TestNode_Notifiee(t *testing.T) {

	var (
		logger          = mocks.NoopLogger
		functionHandler = mocks.BaselineFunctionHandler(t)
	)

	server, err := host.New(mocks.NoopLogger, loopback, 0)
	require.NoError(t, err)

	var (
		storedPeer bool
	)

	peerstore := mocks.BaselinePeerStore(t)
	// Override the peerstore methods so we know if the node correctly handled incoming connection.
	peerstore.StoreFunc = func(peer.ID, multiaddr.Multiaddr, peer.AddrInfo) error {
		storedPeer = true
		return nil
	}

	node, err := New(logger, server, peerstore, functionHandler, WithRole(blockless.HeadNode))
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
