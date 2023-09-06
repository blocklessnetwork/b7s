package response

import (
	"testing"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/stretchr/testify/require"

	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/models/codes"
	"github.com/blocklessnetwork/b7s/models/execute"
	"github.com/blocklessnetwork/b7s/testing/mocks"
)

func TestExecute_Signing(t *testing.T) {

	sampleRes := Execute{
		Type:      blockless.MessageExecuteResponse,
		RequestID: mocks.GenericUUID.String(),
		From:      mocks.GenericPeerID,
		Code:      codes.OK,
		Results: execute.ResultMap{
			mocks.GenericPeerID: mocks.GenericExecutionResult,
		},
		Cluster: execute.Cluster{
			Peers: mocks.GenericPeerIDs[:4],
		},
	}

	t.Run("nominal case", func(t *testing.T) {

		res := sampleRes
		priv, pub := newKey(t)

		err := res.Sign(priv)
		require.NoError(t, err)

		err = res.VerifySignature(pub)
		require.NoError(t, err)
	})
	t.Run("empty signature verification fails", func(t *testing.T) {

		res := sampleRes
		res.Signature = ""

		_, pub := newKey(t)

		err := res.VerifySignature(pub)
		require.Error(t, err)
	})
	t.Run("tampered data signature verification fails", func(t *testing.T) {

		res := sampleRes
		priv, pub := newKey(t)

		err := res.Sign(priv)
		require.NoError(t, err)

		res.RequestID += " "

		err = res.VerifySignature(pub)
		require.Error(t, err)
	})
}

func newKey(t *testing.T) (crypto.PrivKey, crypto.PubKey) {
	t.Helper()
	priv, pub, err := crypto.GenerateKeyPair(crypto.Ed25519, 0)
	require.NoError(t, err)

	return priv, pub
}
