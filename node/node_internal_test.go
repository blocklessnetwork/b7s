package node

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/require"

	"github.com/blocklessnetwork/b7s/host"
	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/testing/mocks"
)

const (
	loopback = "127.0.0.1"

	// How long can the client wait for a published message before giving up.
	publishTimeout = 10 * time.Second

	// It seems like a delay is needed so that the hosts exchange information about the fact
	// that they are subscribed to the same topic. If that does not happen, node might publish
	// a message too soon and the client might miss it. It will then wait for a published message in vain.
	// This is the pause we make after subscribing to the topic and before publishing a message.
	// In reality as little as 250ms is enough, but lets allow a longer time for when
	// tests are executed in parallel or on weaker machines.
	subscriptionDiseminationPause = 2 * time.Second
)

func TestNode_New(t *testing.T) {

	var (
		logger          = mocks.NoopLogger
		peerstore       = mocks.BaselinePeerStore(t)
		functionHandler = mocks.BaselineFStore(t)
		executor        = mocks.BaselineExecutor(t)
	)

	host, err := host.New(logger, loopback, 0)
	require.NoError(t, err)

	t.Run("create a head node", func(t *testing.T) {
		t.Parallel()

		node, err := New(logger, host, peerstore, functionHandler, WithRole(blockless.HeadNode))
		require.NoError(t, err)
		require.NotNil(t, node)

		// Creating a head node with executor fails.
		_, err = New(logger, host, peerstore, functionHandler, WithRole(blockless.HeadNode), WithExecutor(executor))
		require.Error(t, err)
	})
	t.Run("create a worker node", func(t *testing.T) {
		t.Parallel()

		node, err := New(logger, host, peerstore, functionHandler, WithRole(blockless.WorkerNode), WithExecutor(executor), WithWorkspace(t.TempDir()))
		require.NoError(t, err)
		require.NotNil(t, node)

		// Creating a worker node without executor fails.
		_, err = New(logger, host, peerstore, functionHandler, WithRole(blockless.WorkerNode))
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

	var (
		logger          = mocks.NoopLogger
		peerstore       = mocks.BaselinePeerStore(t)
		functionHandler = mocks.BaselineFStore(t)
	)

	host, err := host.New(logger, loopback, 0)
	require.NoError(t, err)

	opts := []Option{
		WithRole(role),
	}

	if role == blockless.WorkerNode {
		executor := mocks.BaselineExecutor(t)
		opts = append(opts, WithExecutor(executor))
		opts = append(opts, WithWorkspace(t.TempDir()))
	}

	node, err := New(logger, host, peerstore, functionHandler, opts...)
	require.NoError(t, err)

	return node
}

func hostAddNewPeer(t *testing.T, host *host.Host, newPeer *host.Host) {
	t.Helper()

	info := hostGetAddrInfo(t, newPeer)
	host.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)
}

func hostGetAddrInfo(t *testing.T, host *host.Host) *peer.AddrInfo {

	addresses := host.Addresses()
	require.NotEmpty(t, addresses)

	addr := addresses[0]

	maddr, err := multiaddr.NewMultiaddr(addr)
	require.NoError(t, err)

	info, err := peer.AddrInfoFromP2pAddr(maddr)
	require.NoError(t, err)

	return info
}

func getStreamPayload(t *testing.T, stream network.Stream, output any) {
	t.Helper()

	buf := bufio.NewReader(stream)
	payload, err := buf.ReadBytes('\n')
	require.ErrorIs(t, err, io.EOF)

	err = json.Unmarshal(payload, output)
	require.NoError(t, err)
}

func serialize(t *testing.T, message any) []byte {
	payload, err := json.Marshal(message)
	require.NoError(t, err)

	return payload
}
