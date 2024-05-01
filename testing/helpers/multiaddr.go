package helpers

import (
	"fmt"
	"testing"

	ma "github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/require"
)

// NOTE: Copeid over from go-libp2p/core/test and tweaked

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
