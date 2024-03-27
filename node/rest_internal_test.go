package node

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/testing/mocks"
)

func TestNode_RestExecuteNotSupportedOnWorker(t *testing.T) {
	node := createNode(t, blockless.WorkerNode)
	_, _, _, _, err := node.ExecuteFunction(context.Background(), mocks.GenericExecutionRequest, "")
	require.Error(t, err)
}

func TestNode_InstallMessageFromCID(t *testing.T) {

	const (
		cid                 = "dummy-cid"
		expectedManifestURL = "https://dummy-cid.ipfs.w3s.link/manifest.json"
	)

	req := createInstallMessageFromCID(cid)

	require.Equal(t, cid, req.CID)
	require.Equal(t, expectedManifestURL, req.ManifestURL)
}
