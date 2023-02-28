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

func TestNode_Notifiee(t *testing.T) {

	const (
		address = "127.0.0.1"
	)

	var (
		logger          = mocks.NoopLogger
		store           = mocks.BaselineStore(t)
		functionHandler = mocks.BaselineFunctionHandler(t)
	)

	server, err := host.New(mocks.NoopLogger, address, 0)
	require.NoError(t, err)

	var (
		storedPeer      bool
		updatedPeerList bool
	)

	peerstore := mocks.BaselinePeerStore(t)
	// Override the peerstore methods so we know if the node correctly handled incoming connection.
	peerstore.StoreFunc = func(peer.ID, multiaddr.Multiaddr, peer.AddrInfo) error {
		storedPeer = true
		return nil
	}
	peerstore.UpdatePeerListFunc = func(peer.ID, multiaddr.Multiaddr, peer.AddrInfo) error {
		updatedPeerList = true
		return nil
	}

	_, err = New(logger, server, store, peerstore, functionHandler, WithRole(blockless.HeadNode))
	require.NoError(t, err)

	serverAddresses := server.Addresses()
	require.NotEmpty(t, serverAddresses)

	serverAddress := serverAddresses[0]

	client, err := host.New(mocks.NoopLogger, address, 0)
	require.NoError(t, err)

	serverInfo := addPeerToPeerStore(t, client, serverAddress)

	err = client.Connect(context.Background(), *serverInfo)
	require.NoError(t, err)

	// Verify that peer store was updated.
	require.True(t, storedPeer)
	require.True(t, updatedPeerList)
}

func addPeerToPeerStore(t *testing.T, host *host.Host, addr string) *peer.AddrInfo {
	t.Helper()

	maddr, err := multiaddr.NewMultiaddr(addr)
	require.NoError(t, err)

	info, err := peer.AddrInfoFromP2pAddr(maddr)
	require.NoError(t, err)

	host.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

	return info
}
