package store_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/require"

	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/store"
	"github.com/blocklessnetwork/b7s/store/codec"
	"github.com/blocklessnetwork/b7s/testing/helpers"
	"github.com/blocklessnetwork/b7s/testing/mocks"
)

func TestStore_SaveAndRetrievePeer(t *testing.T) {
	db := helpers.InMemoryDB(t)
	defer db.Close()

	peer := createGenericPeer(t)
	store := store.New(db, codec.NewJSONCodec())

	t.Run("save peer", func(t *testing.T) {
		err := store.SavePeer(peer)
		require.NoError(t, err)
	})
	t.Run("retrieve peer", func(t *testing.T) {
		retrieved, err := store.RetrievePeer(mocks.GenericPeerID)
		require.NoError(t, err)

		require.Equal(t, peer, retrieved)
	})
}

func TestStore_PeerFunctionsHandleFailures(t *testing.T) {

	db := helpers.InMemoryDB(t)
	defer db.Close()

	t.Run("retrieving missing peer fails", func(t *testing.T) {
		store := store.New(db, codec.NewJSONCodec())

		_, err := store.RetrievePeer(mocks.GenericPeerID)
		require.Error(t, err)
	})
	t.Run("save handles marshalling failures", func(t *testing.T) {

		codec := mocks.BaselineCodec(t)
		codec.MarshalFunc = func(any) ([]byte, error) {
			return nil, mocks.GenericError
		}
		store := store.New(db, codec)

		peer := createGenericPeer(t)
		err := store.SavePeer(peer)
		require.Error(t, err)
	})
	t.Run("retrieve handles unmarshalling failures", func(t *testing.T) {

		unmarshalErr := errors.New("unmarshalling error")
		codec := mocks.BaselineCodec(t)
		codec.MarshalFunc = func(obj any) ([]byte, error) {
			return json.Marshal(obj)
		}
		codec.UnmarshalFunc = func([]byte, any) error {
			return unmarshalErr
		}
		store := store.New(db, codec)

		// First, save the peer so we don't end up with a "not found" error.
		err := store.SavePeer(createGenericPeer(t))
		require.NoError(t, err)

		_, err = store.RetrievePeer(mocks.GenericPeerID)
		require.Error(t, err)
		require.ErrorIs(t, err, unmarshalErr)
	})
}

func createGenericPeer(t *testing.T) blockless.Peer {
	t.Helper()

	return blockless.Peer{
		ID:        mocks.GenericPeerID,
		MultiAddr: mocks.GenericMultiaddress.String(),
		AddrInfo: peer.AddrInfo{
			ID: mocks.GenericPeerID,
			Addrs: []multiaddr.Multiaddr{
				mocks.GenericMultiaddress,
			},
		},
	}
}
