package helpers

import (
	"testing"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/multiformats/go-multiaddr"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/blessnetwork/b7s/host"
)

const (
	loopback = "127.0.0.1"
)

func NewLoopbackHost(t *testing.T, log zerolog.Logger) *host.Host {
	t.Helper()

	host, err := host.New(log, loopback, 0)
	require.NoError(t, err)

	return host
}

func HostGetAddrInfo(t *testing.T, host *host.Host) *peer.AddrInfo {

	addresses := host.Addresses()
	require.NotEmpty(t, addresses)

	maddrs := make([]multiaddr.Multiaddr, len(addresses))
	for i, addr := range addresses {

		maddr, err := multiaddr.NewMultiaddr(addr)
		require.NoError(t, err)

		maddrs[i] = maddr
	}

	info := peer.AddrInfo{
		ID:    host.ID(),
		Addrs: maddrs,
	}

	return &info
}

func HostAddNewPeer(t *testing.T, host *host.Host, newPeer *host.Host) {
	t.Helper()

	info := HostGetAddrInfo(t, newPeer)
	host.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)
}
