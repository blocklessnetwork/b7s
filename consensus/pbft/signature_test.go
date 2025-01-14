package pbft

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/blessnetwork/b7s/testing/mocks"
)

func TestSign_PrePrepare(t *testing.T) {

	var (
		samplePrePrepare = PrePrepare{
			View:           14,
			SequenceNumber: 45,
			Digest:         "abcdef123456789",
			Request: Request{
				ID:        mocks.GenericUUID.String(),
				Timestamp: time.Now().UTC(),
				Origin:    mocks.GenericPeerID,
				Execute:   mocks.GenericExecutionRequest,
			},
		}
	)

	t.Run("nominal case", func(t *testing.T) {

		prePrepare := samplePrePrepare
		require.Empty(t, prePrepare.Signature)

		replica := newDummyReplica(t)
		err := replica.sign(&prePrepare)
		require.NoError(t, err)
		require.NotEmpty(t, prePrepare.Signature)

		verifier := newDummyReplica(t)

		err = verifier.verifySignature(&prePrepare, replica.host.ID())
		require.NoError(t, err)
	})
	t.Run("catch tampering", func(t *testing.T) {

		prePrepare := samplePrePrepare

		replica := newDummyReplica(t)
		err := replica.sign(&prePrepare)
		require.NoError(t, err)
		require.NotEmpty(t, prePrepare.Signature)

		verifier := newDummyReplica(t)

		prePrepare.View++

		err = verifier.verifySignature(&prePrepare, replica.host.ID())
		require.Error(t, err)
	})
	t.Run("validating empty signature fails", func(t *testing.T) {

		prePrepare := samplePrePrepare
		replica := newDummyReplica(t)

		err := replica.verifySignature(&prePrepare, replica.host.ID())
		require.Error(t, err)
	})
}

func TestSign_Prepare(t *testing.T) {

	var (
		samplePrepare = Prepare{
			View:           78,
			SequenceNumber: 15,
			Digest:         "987654321abcdef",
		}
	)

	t.Run("nominal case", func(t *testing.T) {

		prepare := samplePrepare
		require.Empty(t, prepare.Signature)

		replica := newDummyReplica(t)
		err := replica.sign(&prepare)
		require.NoError(t, err)
		require.NotEmpty(t, prepare.Signature)

		verifier := newDummyReplica(t)

		err = verifier.verifySignature(&prepare, replica.host.ID())
		require.NoError(t, err)
	})
	t.Run("catch tampering", func(t *testing.T) {

		prepare := samplePrepare

		replica := newDummyReplica(t)
		err := replica.sign(&prepare)
		require.NoError(t, err)
		require.NotEmpty(t, prepare.Signature)

		verifier := newDummyReplica(t)

		prepare.Digest += " "

		err = verifier.verifySignature(&prepare, replica.host.ID())
		require.Error(t, err)
	})
	t.Run("validating empty signature fails", func(t *testing.T) {

		prepare := samplePrepare
		replica := newDummyReplica(t)

		err := replica.verifySignature(&prepare, replica.host.ID())
		require.Error(t, err)
	})
}

func TestSign_Commit(t *testing.T) {

	var (
		sampleCommit = Commit{
			View:           32,
			SequenceNumber: 41,
			Digest:         "123456789pqrs",
		}
	)

	t.Run("nominal case", func(t *testing.T) {

		commit := sampleCommit
		require.Empty(t, commit.Signature)

		replica := newDummyReplica(t)
		err := replica.sign(&commit)
		require.NoError(t, err)
		require.NotEmpty(t, commit.Signature)

		verifier := newDummyReplica(t)

		err = verifier.verifySignature(&commit, replica.host.ID())
		require.NoError(t, err)
	})
	t.Run("catch tampering", func(t *testing.T) {

		commit := sampleCommit

		replica := newDummyReplica(t)
		err := replica.sign(&commit)
		require.NoError(t, err)
		require.NotEmpty(t, commit.Signature)

		verifier := newDummyReplica(t)

		commit.SequenceNumber = 0

		err = verifier.verifySignature(&commit, replica.host.ID())
		require.Error(t, err)
	})
	t.Run("validating empty signature fails", func(t *testing.T) {

		commit := sampleCommit
		replica := newDummyReplica(t)

		err := replica.verifySignature(&commit, replica.host.ID())
		require.Error(t, err)
	})
}
