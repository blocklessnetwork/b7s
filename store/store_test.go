package store_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/libp2p/go-libp2p/core/peer"
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

	peer := mocks.GenericPeer
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

func TestStore_RetrievePeers(t *testing.T) {
	db := helpers.InMemoryDB(t)
	defer db.Close()
	store := store.New(db, codec.NewJSONCodec())

	count := 10
	// Generate peers
	peers := make(map[peer.ID]blockless.Peer)
	for i := 0; i < count; i++ {

		id := helpers.RandPeerIDFatal(t)
		addrs := helpers.GenerateTestAddrs(t, 1)

		p := blockless.Peer{
			ID:        id,
			MultiAddr: addrs[0].String(),
			AddrInfo: peer.AddrInfo{
				ID:    id,
				Addrs: addrs,
			},
		}

		peers[id] = p
	}

	// Save peers.
	for _, peer := range peers {
		err := store.SavePeer(peer)
		require.NoError(t, err)
	}

	retrieved, err := store.RetrievePeers()
	require.NoError(t, err)
	require.Len(t, retrieved, count)

	// Verify peers.
	for _, peer := range retrieved {
		require.Equal(t, peer, peers[peer.ID])
	}

}

func TestStore_SaveAndRetrieveFunction(t *testing.T) {
	db := helpers.InMemoryDB(t)
	defer db.Close()

	function := mocks.GenericFunctionRecord
	store := store.New(db, codec.NewJSONCodec())

	t.Run("save function", func(t *testing.T) {
		err := store.SaveFunction(function)
		require.NoError(t, err)
	})
	t.Run("retrieve function", func(t *testing.T) {
		retrieved, err := store.RetrieveFunction(function.CID)
		require.NoError(t, err)

		require.Equal(t, function, retrieved)
	})
}

func TestStore_HandlesFailures(t *testing.T) {

	db := helpers.InMemoryDB(t)
	defer db.Close()

	t.Run("retrieving missing peer fails", func(t *testing.T) {
		store := store.New(db, codec.NewJSONCodec())

		_, err := store.RetrievePeer(mocks.GenericPeerID)
		require.Error(t, err)
	})
	t.Run("retrieving missing function fails", func(t *testing.T) {
		store := store.New(db, codec.NewJSONCodec())

		_, err := store.RetrieveFunction(mocks.GenericString)
		require.Error(t, err)
	})
	t.Run("save peer handles marshalling failures", func(t *testing.T) {

		codec := mocks.BaselineCodec(t)
		codec.MarshalFunc = func(any) ([]byte, error) {
			return nil, mocks.GenericError
		}
		store := store.New(db, codec)

		err := store.SavePeer(mocks.GenericPeer)
		require.Error(t, err)
	})
	t.Run("save function handles marshalling failures", func(t *testing.T) {

		codec := mocks.BaselineCodec(t)
		codec.MarshalFunc = func(any) ([]byte, error) {
			return nil, mocks.GenericError
		}
		store := store.New(db, codec)

		err := store.SaveFunction(mocks.GenericFunctionRecord)
		require.Error(t, err)
	})
	t.Run("retrieve peer handles unmarshalling failures", func(t *testing.T) {

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
		err := store.SavePeer(mocks.GenericPeer)
		require.NoError(t, err)

		_, err = store.RetrievePeer(mocks.GenericPeerID)
		require.Error(t, err)
		require.ErrorIs(t, err, unmarshalErr)
	})
	t.Run("retrieve function handles unmarshalling failures", func(t *testing.T) {

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
		err := store.SaveFunction(mocks.GenericFunctionRecord)
		require.NoError(t, err)

		_, err = store.RetrieveFunction(mocks.GenericFunctionRecord.CID)
		require.Error(t, err)
		require.ErrorIs(t, err, unmarshalErr)
	})
}
