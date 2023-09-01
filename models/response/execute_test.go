package response

import (
	"testing"

	"github.com/blocklessnetworking/b7s/host"
	"github.com/blocklessnetworking/b7s/models/blockless"
	"github.com/blocklessnetworking/b7s/models/codes"
	"github.com/blocklessnetworking/b7s/models/execute"
	"github.com/blocklessnetworking/b7s/testing/mocks"
	"github.com/stretchr/testify/require"
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
		host, err := host.New(mocks.NoopLogger, "127.0.0.1", 0)
		require.NoError(t, err)

		err = res.Sign(host.PrivateKey())
		require.NoError(t, err)

		err = res.VerifySignature(host.PublicKey())
		require.NoError(t, err)
	})
	t.Run("empty signature verification fails", func(t *testing.T) {

		res := sampleRes
		res.Signature = ""

		host, err := host.New(mocks.NoopLogger, "127.0.0.1", 0)
		require.NoError(t, err)

		err = res.VerifySignature(host.PublicKey())
		require.Error(t, err)
	})
	t.Run("tampered data signature verification fails", func(t *testing.T) {

		res := sampleRes
		host, err := host.New(mocks.NoopLogger, "127.0.0.1", 0)
		require.NoError(t, err)

		err = res.Sign(host.PrivateKey())
		require.NoError(t, err)

		res.RequestID += " "

		err = res.VerifySignature(host.PublicKey())
		require.Error(t, err)
	})

}
