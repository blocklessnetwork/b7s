package node

import (
	"testing"

	"github.com/blocklessnetworking/b7s/models/blockless"
	"github.com/blocklessnetworking/b7s/testing/mocks"
	"github.com/stretchr/testify/require"
)

func TestNode_Sync(t *testing.T) {

	var (
		installed = []string{
			"func1",
			"func2",
			"func3",
		}

		synced []string
	)

	fstore := mocks.BaselineFStore(t)
	fstore.InstalledFunctionsFunc = func() []string {
		return installed
	}
	fstore.SyncFunc = func(cid string) error {
		synced = append(synced, cid)
		return nil
	}

	node := createNode(t, blockless.WorkerNode)
	node.fstore = fstore

	node.syncFunctions()

	// Verify all functions were synced.
	require.Equal(t, installed, synced)
}
