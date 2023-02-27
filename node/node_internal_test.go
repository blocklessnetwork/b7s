package node

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/blocklessnetworking/b7s/host"
	"github.com/blocklessnetworking/b7s/models/blockless"
	"github.com/blocklessnetworking/b7s/testing/mocks"
)

func TestNode_New(t *testing.T) {

	const (
		address = "127.0.0.1"
		port    = 9000
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
			address = "127.0.0.1"
			port    = 9000
			msgType = "jibberish"
		)

		var (
			logger          = mocks.NoopLogger
			store           = mocks.BaselineStore(t)
			peerstore       = mocks.BaselinePeerStore(t)
			functionHandler = mocks.BaselineFunctionHandler(t)
		)

		host, err := host.New(logger, address, port)
		require.NoError(t, err)

		node, err := New(logger, host, store, peerstore, functionHandler, WithRole(blockless.HeadNode))
		require.NoError(t, err)

		handlerFunc := node.getHandler(msgType)

		err = handlerFunc(context.Background(), mocks.GenericPeerID, []byte{})
		require.Error(t, err)

		require.ErrorIs(t, err, ErrUnsupportedMessage)
	})
}
