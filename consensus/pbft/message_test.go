package pbft

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/blocklessnetworking/b7s/testing/mocks"
)

func TestMessage_Encode(t *testing.T) {

	t.Run("pre-prepare", func(t *testing.T) {

		orig := PrePrepare{
			View:           1,
			SequenceNumber: 2,
			Digest:         "123456789",
			ReplicaID:      mocks.GenericPeerID,
		}

		encoded, err := json.Marshal(orig)
		require.NoError(t, err)

		var unpacked messageRecord
		err = json.Unmarshal(encoded, &unpacked)
		require.NoError(t, err)

		// Verify that the message type was correctly set.
		require.Equal(t, MessagePrePrepare, unpacked.Type)

		// Verofy that the original data is correctly encoded.
		var msg PrePrepare
		err = json.Unmarshal(unpacked.Data, &msg)
		require.NoError(t, err)

		require.Equal(t, orig, msg)
	})
	t.Run("prepare", func(t *testing.T) {

		orig := Prepare{
			View:           14,
			SequenceNumber: 45,
			Digest:         "abc123def",
			ReplicaID:      mocks.GenericPeerID,
		}

		encoded, err := json.Marshal(orig)
		require.NoError(t, err)

		var unpacked messageRecord
		err = json.Unmarshal(encoded, &unpacked)
		require.NoError(t, err)

		// Verify that the message type was correctly set.
		require.Equal(t, MessagePrepare, unpacked.Type)

		// Verofy that the original data is correctly encoded.
		var msg Prepare
		err = json.Unmarshal(unpacked.Data, &msg)
		require.NoError(t, err)

		require.Equal(t, orig, msg)
	})
	t.Run("commit", func(t *testing.T) {

		orig := Commit{
			View:           23,
			SequenceNumber: 51,
			Digest:         "987xyz",
			ReplicaID:      mocks.GenericPeerID,
		}

		encoded, err := json.Marshal(orig)
		require.NoError(t, err)

		var unpacked messageRecord
		err = json.Unmarshal(encoded, &unpacked)
		require.NoError(t, err)

		// Verify that the message type was correctly set.
		require.Equal(t, MessageCommit, unpacked.Type)

		// Verofy that the original data is correctly encoded.
		var msg Commit
		err = json.Unmarshal(unpacked.Data, &msg)
		require.NoError(t, err)

		require.Equal(t, orig, msg)
	})
	// TODO (pbft): Add request handling
}

func TestMessage_Decode(t *testing.T) {

	t.Run("pre-prepare", func(t *testing.T) {

		orig := PrePrepare{
			View:           1,
			SequenceNumber: 2,
			Digest:         "123456789",
			ReplicaID:      mocks.GenericPeerID,
		}

		encoded, err := json.Marshal(orig)
		require.NoError(t, err)

		unpacked, err := unpackMessage(encoded)
		require.NoError(t, err)

		// Verify that the message type was correctly identified.
		prePrepare, ok := unpacked.(PrePrepare)
		require.True(t, ok)

		// Verify that the data is correctly unpacked.
		require.Equal(t, orig, prePrepare)
	})
	t.Run("prepare", func(t *testing.T) {

		orig := Prepare{
			View:           14,
			SequenceNumber: 45,
			Digest:         "abc123def",
			ReplicaID:      mocks.GenericPeerID,
		}

		encoded, err := json.Marshal(orig)
		require.NoError(t, err)

		unpacked, err := unpackMessage(encoded)
		require.NoError(t, err)

		// Verify that the message type was correctly identified.
		prepare, ok := unpacked.(Prepare)
		require.True(t, ok)

		// Verify that the data is correctly unpacked.
		require.Equal(t, orig, prepare)
	})
	t.Run("commit", func(t *testing.T) {

		orig := Commit{
			View:           23,
			SequenceNumber: 51,
			Digest:         "987xyz",
			ReplicaID:      mocks.GenericPeerID,
		}

		encoded, err := json.Marshal(orig)
		require.NoError(t, err)

		unpacked, err := unpackMessage(encoded)
		require.NoError(t, err)

		// Verify that the message type was correctly identified.
		commit, ok := unpacked.(Commit)
		require.True(t, ok)

		// Verify that the data is correctly unpacked.
		require.Equal(t, orig, commit)
	})
	// TODO (pbft): Add request handling
}
