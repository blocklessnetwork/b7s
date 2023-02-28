package node

import (
	"context"
	"testing"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/require"

	"github.com/blocklessnetworking/b7s/host"
	"github.com/blocklessnetworking/b7s/models/blockless"
	"github.com/blocklessnetworking/b7s/testing/mocks"
)

func TestNode_New(t *testing.T) {

	const (
		address = "127.0.0.1"
		port    = 0
	)

	var (
		logger          = mocks.NoopLogger
		store           = mocks.BaselineStore(t)
		peerstore       = mocks.BaselinePeerStore(t)
		functionHandler = mocks.BaselineFunctionHandler(t)
		executor        = mocks.BaselineExecutor(t)
	)

	host, err := host.New(logger, address, port)
	require.NoError(t, err)

	t.Run("create a head node", func(t *testing.T) {
		t.Parallel()

		node, err := New(logger, host, store, peerstore, functionHandler, WithRole(blockless.HeadNode))
		require.NoError(t, err)
		require.NotNil(t, node)

		// Creating a head node with executor fails.
		_, err = New(logger, host, store, peerstore, functionHandler, WithRole(blockless.HeadNode), WithExecutor(executor))
		require.Error(t, err)
	})
	t.Run("create a worker node", func(t *testing.T) {
		t.Parallel()

		node, err := New(logger, host, store, peerstore, functionHandler, WithRole(blockless.WorkerNode), WithExecutor(executor))
		require.NoError(t, err)
		require.NotNil(t, node)

		// Creating a worker node without executor fails.
		_, err = New(logger, host, store, peerstore, functionHandler, WithRole(blockless.WorkerNode))
		require.Error(t, err)
	})
}

func TestNode_MessageHandler(t *testing.T) {
	t.Run("unsupported messages should fail", func(t *testing.T) {
		t.Parallel()

		const (
			msgType = "jibberish"
		)

		node := createNode(t, blockless.HeadNode)

		handlerFunc := node.getHandler(msgType)

		err := handlerFunc(context.Background(), mocks.GenericPeerID, []byte{})
		require.Error(t, err)

		require.ErrorIs(t, err, ErrUnsupportedMessage)
	})
}

func createNode(t *testing.T, role blockless.NodeRole) *Node {
	t.Helper()

	const (
		address = "127.0.0.1"
		port    = 0
	)

	var (
		logger          = mocks.NoopLogger
		store           = mocks.BaselineStore(t)
		peerstore       = mocks.BaselinePeerStore(t)
		functionHandler = mocks.BaselineFunctionHandler(t)
	)

	host, err := host.New(logger, address, port)
	require.NoError(t, err)

	opts := []Option{
		WithRole(role),
	}

	if role == blockless.WorkerNode {
		executor := mocks.BaselineExecutor(t)
		opts = append(opts, WithExecutor(executor))
	}

	node, err := New(logger, host, store, peerstore, functionHandler, opts...)
	require.NoError(t, err)

	return node
}

func getAddrInfo(t *testing.T, addr string) *peer.AddrInfo {
	t.Helper()

	maddr, err := multiaddr.NewMultiaddr(addr)
	require.NoError(t, err)

	info, err := peer.AddrInfoFromP2pAddr(maddr)
	require.NoError(t, err)

	return info
}

func addPeerToPeerStore(t *testing.T, host *host.Host, addr string) *peer.AddrInfo {
	t.Helper()

	info := getAddrInfo(t, addr)
	host.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

	return info
}
