package fstore_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/blocklessnetwork/b7s/fstore"
	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/testing/mocks"
)

func TestFunction_RetrieveHandlesErrors(t *testing.T) {

	const (
		testCID = "dummy-cid"
	)

	workdir, err := os.MkdirTemp("", "b7s-function-get-")
	require.NoError(t, err)

	defer os.RemoveAll(workdir)

	store := mocks.BaselineStore(t)
	store.RetrieveFunctionFunc = func(context.Context, string) (blockless.FunctionRecord, error) {
		return blockless.FunctionRecord{}, mocks.GenericError
	}

	fh := fstore.New(mocks.NoopLogger, store, workdir)

	_, err = fh.Get(context.Background(), testCID)
	require.Error(t, err)
}
