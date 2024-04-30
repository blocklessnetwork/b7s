package fstore_test

import (
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
	store.RetrieveFunctionFunc = func(string) (blockless.FunctionRecord, error) {
		return blockless.FunctionRecord{}, mocks.GenericError
	}

	fh := fstore.New(mocks.NoopLogger, store, workdir)

	_, err = fh.Get(testCID)
	require.Error(t, err)
}

/*
func TestFunction_RetrieveFunctions(t *testing.T) {

	// TODO: Implement this and also handle errors.

	installed := []string{
		"func1",
		"func2",
		"func3",
	}

	workdir, err := os.MkdirTemp("", "b7s-function-get-")
	require.NoError(t, err)

	defer os.RemoveAll(workdir)

	store := mocks.BaselineStore(t)
	store.KeysFunc = func() []string {
		return installed
	}

	fh := fstore.New(mocks.NoopLogger, store, workdir)

	list, err := fh.InstalledFunctions()
	require.NoError(t, err)
	require.Equal(t, installed, list)
}

*/
