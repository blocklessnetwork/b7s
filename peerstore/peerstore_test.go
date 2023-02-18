package peerstore_test

import (
	"testing"

	"github.com/cockroachdb/pebble"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/require"

	"github.com/blocklessnetworking/b7s/models/blockless"
	"github.com/blocklessnetworking/b7s/peerstore"
	"github.com/blocklessnetworking/b7s/store"
	"github.com/blocklessnetworking/b7s/testing/helpers"
	"github.com/blocklessnetworking/b7s/testing/mocks"
)

func Test_PeerStore(t *testing.T) {
	t.Run("empty peer store", func(t *testing.T) {
		t.Parallel()

		peerstore, _, db := setupPeerStore(t)
		defer db.Close()

		peers, err := peerstore.Peers()
		require.NoError(t, err)
		require.Empty(t, peers)
	})
	t.Run("store peer", func(t *testing.T) {
		t.Parallel()

		peerstore, store, db := setupPeerStore(t)
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

		// Verify peer is written to the underlying store.
		var peer blockless.Peer
		err = store.GetRecord(peerID.String(), &peer)
		require.NoError(t, err)

		require.Equal(t, peerID, peer.ID)
		require.Equal(t, addr.String(), peer.MultiAddr)
		require.Equal(t, info, peer.AddrInfo)
	})
	t.Run("update peer list", func(t *testing.T) {
		t.Parallel()

		peerstore, _, db := setupPeerStore(t)
		defer db.Close()

		var (
			peerID = mocks.GenericPeerID
			addr   = genericMultiAddr(t)
			info   = peer.AddrInfo{
				ID:    peerID,
				Addrs: []multiaddr.Multiaddr{addr},
			}
		)

		err := peerstore.UpdatePeerList(peerID, addr, info)
		require.NoError(t, err)

		peers, err := peerstore.Peers()
		require.NoError(t, err)
		require.Len(t, peers, 1)

		peer := peers[0]
		require.Equal(t, peerID, peer.ID)
		require.Equal(t, addr.String(), peer.MultiAddr)
		require.Equal(t, info, peer.AddrInfo)
	})
	t.Run("adding known peer to peer list", func(t *testing.T) {
		t.Parallel()

		peerstore, _, db := setupPeerStore(t)
		defer db.Close()

		var (
			peerID = mocks.GenericPeerID
			addr   = genericMultiAddr(t)
			info   = peer.AddrInfo{
				ID:    peerID,
				Addrs: []multiaddr.Multiaddr{addr},
			}
		)

		err := peerstore.UpdatePeerList(peerID, addr, info)
		require.NoError(t, err)

		// Add the same peer again - we should still only have one peer in the list.
		err = peerstore.UpdatePeerList(peerID, addr, info)
		require.NoError(t, err)

		peers, err := peerstore.Peers()
		require.NoError(t, err)
		require.Len(t, peers, 1)
	})
}

func Test_PeerStore_Store(t *testing.T) {

	t.Run("handles failure to store peer", func(t *testing.T) {

		store := mocks.BaselineStore(t)
		store.GetRecordFunc = func(string, interface{}) error {
			// We first check if the peer exists - make sure it doesn't.
			return blockless.ErrNotFound
		}
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
	t.Run("handles failure to get existing peer", func(t *testing.T) {

		store := mocks.BaselineStore(t)
		store.GetRecordFunc = func(string, interface{}) error {
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
	t.Run("handles noop on existing peer", func(t *testing.T) {

		store := mocks.BaselineStore(t)
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
		require.NoError(t, err)
	})
	t.Run("handles failure to get existing peer list", func(t *testing.T) {
		store := mocks.BaselineStore(t)
		store.GetRecordFunc = func(string, interface{}) error {
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

		err := peerstore.UpdatePeerList(peerID, addr, info)
		require.ErrorIs(t, err, mocks.GenericError)
	})
	t.Run("handles failure to update peer list", func(t *testing.T) {
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

		err := peerstore.UpdatePeerList(peerID, addr, info)
		require.ErrorIs(t, err, mocks.GenericError)
	})
	t.Run("handles peer list retrieval error", func(t *testing.T) {
		store := mocks.BaselineStore(t)
		store.GetRecordFunc = func(string, interface{}) error {
			return mocks.GenericError
		}

		peerstore := peerstore.New(store)

		_, err := peerstore.Peers()
		require.ErrorIs(t, err, mocks.GenericError)
	})
}
func setupPeerStore(t *testing.T) (*peerstore.PeerStore, *store.Store, *pebble.DB) {
	t.Helper()

	db := helpers.InMemoryDB(t)
	store := store.New(db)
	ps := peerstore.New(store)

	return ps, store, db
}

func genericMultiAddr(t *testing.T) multiaddr.Multiaddr {
	t.Helper()

	addr, err := multiaddr.NewMultiaddr(mocks.GenericAddress)
	require.NoError(t, err)

	return addr
}
