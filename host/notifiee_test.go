package host_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/blessnetwork/b7s/host"
	"github.com/blessnetwork/b7s/models/bls"
	"github.com/blessnetwork/b7s/testing/helpers"
	"github.com/blessnetwork/b7s/testing/mocks"
)

const (
	loopback = "127.0.0.1"
)

func TestNotifiee(t *testing.T) {

	var (
		storedPeer bool
	)

	store := mocks.BaselineStore(t)
	// Override the peerstore methods so we know if the node correctly handled incoming connection.
	store.SavePeerFunc = func(context.Context, bls.Peer) error {
		storedPeer = true
		return nil
	}

	server, err := host.New(mocks.NoopLogger, loopback, 0)
	require.NoError(t, err)

	notifiee := host.NewNotifee(mocks.NoopLogger, store)
	server.Network().Notify(notifiee)

	client, err := host.New(mocks.NoopLogger, loopback, 0)
	require.NoError(t, err)

	helpers.HostAddNewPeer(t, client, server)

	serverInfo := helpers.HostGetAddrInfo(t, server)
	err = client.Connect(context.Background(), *serverInfo)
	require.NoError(t, err)

	// Verify that peer store was updated.
	require.True(t, storedPeer)
}
