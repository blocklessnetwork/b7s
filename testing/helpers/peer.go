package helpers

import (
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
	mh "github.com/multiformats/go-multihash"

	"github.com/stretchr/testify/require"
)

// NOTE: Inspiration by go-libp2p/core/test

func RandPeerID(t *testing.T) peer.ID {
	t.Helper()

	buf := make([]byte, 16)
	rand.Read(buf)
	h, err := mh.Sum(buf, mh.SHA2_256, -1)
	require.NoError(t, err)

	return peer.ID(h)
}

func GenerateTestAddrs(t *testing.T, n int) []ma.Multiaddr {
	t.Helper()

	out := make([]ma.Multiaddr, n)
	for i := 0; i < n; i++ {
		a, err := ma.NewMultiaddr(fmt.Sprintf("/ip4/1.2.3.4/tcp/%d", i))
		require.NoError(t, err)

		out[i] = a
	}
	return out
}
