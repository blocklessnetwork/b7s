package fstore_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/blocklessnetwork/b7s/fstore"
	"github.com/blocklessnetwork/b7s/testing/mocks"
)

func TestFunction_GetHandlesErrors(t *testing.T) {

	const (
		testCID = "dummy-cid"
	)

	workdir, err := os.MkdirTemp("", "b7s-function-get-")
	require.NoError(t, err)

	defer os.RemoveAll(workdir)

	store := mocks.BaselineStore(t)
	store.GetRecordFunc = func(string, interface{}) error {
		return mocks.GenericError
	}

	fh := fstore.New(mocks.NoopLogger, store, workdir)

	_, err = fh.Get(testCID)
	require.Error(t, err)
}

func TestFunction_InstalledFunctions(t *testing.T) {

	installed := []string{
		"func1",
		"func2",
		"func3",
	}

	workdir, err := os.MkdirTemp("", "b7s-function-get-")
	require.NoError(t, err)

	defer os.RemoveAll(workdir)

	store := mocks.BaselineStore(t)
	store.KeysFunc = func() ([]string, error) {
		return installed, nil
	}

	fh := fstore.New(mocks.NoopLogger, store, workdir)

	list, err := fh.InstalledFunctions()
	require.NoError(t, err)
	require.Equal(t, installed, list)
}
