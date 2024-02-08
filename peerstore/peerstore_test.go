package peerstore_test

import (
	"testing"

	"github.com/cockroachdb/pebble"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/require"

	"github.com/blocklessnetwork/b7s/peerstore"
	"github.com/blocklessnetwork/b7s/store"
	"github.com/blocklessnetwork/b7s/testing/helpers"
	"github.com/blocklessnetwork/b7s/testing/mocks"
)

func Test_PeerStore(t *testing.T) {
	t.Run("empty peer store", func(t *testing.T) {
		t.Parallel()

		peerstore, db := setupPeerStore(t)
		defer db.Close()

		peers, err := peerstore.Peers()
		require.NoError(t, err)
		require.Empty(t, peers)
	})
	t.Run("store/get/delete peer", func(t *testing.T) {
		t.Parallel()

		peerstore, db := setupPeerStore(t)
		defer db.Close()

		var (
			peerID = mocks.GenericPeerID
			addr   = genericMultiAddr(t)
			info   = peer.AddrInfo{
				ID:    peerID,
				Addrs: []multiaddr.Multiaddr{addr},
			}
		)

		// Verify peerstore is empty.
		peers, err := peerstore.Peers()
		require.NoError(t, err)
		require.Len(t, peers, 0)

		err = peerstore.Store(mocks.GenericPeerID, addr, info)
		require.NoError(t, err)

		// Verify peer is written to the peerstore.
		read, err := peerstore.Get(peerID)
		require.NoError(t, err)

		require.Equal(t, addr.String(), read.MultiAddr)
		require.Equal(t, info, read.AddrInfo)

		// Verify peer list has one peer.
		peers, err = peerstore.Peers()
		require.NoError(t, err)
		require.Len(t, peers, 1)

		err = peerstore.Remove(peerID)
		require.NoError(t, err)

		// Verify peer cannot be retrieved anymore.
		_, err = peerstore.Get(peerID)
		require.Error(t, err)

		// Verify peer list is empty now.
		peers, err = peerstore.Peers()
		require.NoError(t, err)
		require.Len(t, peers, 0)
	})
	t.Run("adding known peer", func(t *testing.T) {
		t.Parallel()

		peerstore, db := setupPeerStore(t)
		defer db.Close()

		var (
			peerID = mocks.GenericPeerID
			addr   = genericMultiAddr(t)
			info   = peer.AddrInfo{
				ID:    peerID,
				Addrs: []multiaddr.Multiaddr{addr},
			}
		)

		err := peerstore.Store(mocks.GenericPeerID, addr, info)
		require.NoError(t, err)

		// Add the same peer again - we should still only have one peer in the list.
		err = peerstore.Store(mocks.GenericPeerID, addr, info)
		require.NoError(t, err)

		peers, err := peerstore.Peers()
		require.NoError(t, err)
		require.Len(t, peers, 1)
	})
}

func Test_PeerStore_Store(t *testing.T) {

	t.Run("handles failure to store peer", func(t *testing.T) {

		store := mocks.BaselineStore(t)
		store.SetRecordFunc = func(string, interface{}) error {
			return mocks.GenericError
		}

		peerstore := peerstore.New(store)

		var (
			peerID = mocks.GenericPeerID
			addr   = genericMultiAddr(t)
			info   = peer.AddrInfo{
				ID:    peerID,
				Addrs: []multiaddr.Multiaddr{addr},
			}
		)

		err := peerstore.Store(peerID, addr, info)
		require.ErrorIs(t, err, mocks.GenericError)
	})
	t.Run("handles failure to get peer", func(t *testing.T) {

		store := mocks.BaselineStore(t)
		store.GetRecordFunc = func(string, interface{}) error {
			return mocks.GenericError
		}
		peerstore := peerstore.New(store)

		var (
			peerID = mocks.GenericPeerID
		)

		_, err := peerstore.Get(peerID)
		require.ErrorIs(t, err, mocks.GenericError)
	})
	t.Run("handles peer list retrieval error", func(t *testing.T) {
		store := mocks.BaselineStore(t)
		store.KeysFunc = func() []string {
			return []string{"dummy-key"}
		}
		store.GetRecordFunc = func(string, interface{}) error {
			return mocks.GenericError
		}

		peerstore := peerstore.New(store)

		_, err := peerstore.Peers()
		require.ErrorIs(t, err, mocks.GenericError)
	})
	t.Run("handles peer removal error", func(t *testing.T) {
		t.Parallel()

		store := mocks.BaselineStore(t)
		store.DeleteFunc = func(string) error {
			return mocks.GenericError
		}

		var (
			peerID = mocks.GenericPeerID
		)

		peerstore := peerstore.New(store)

		err := peerstore.Remove(peerID)
		require.ErrorIs(t, err, mocks.GenericError)
	})
}
func setupPeerStore(t *testing.T) (*peerstore.PeerStore, *pebble.DB) {
	t.Helper()

	db := helpers.InMemoryDB(t)
	store := store.New(db)
	ps := peerstore.New(store)

	return ps, db
}

func genericMultiAddr(t *testing.T) multiaddr.Multiaddr {
	t.Helper()

	addr, err := multiaddr.NewMultiaddr(mocks.GenericAddress)
	require.NoError(t, err)

	return addr
}
