package fstore_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/blocklessnetworking/b7s/fstore"
	"github.com/blocklessnetworking/b7s/testing/mocks"
)

func TestFunction_GetHandlesErrors(t *testing.T) {

	const (
		testCID = "dummy-cid"
	)

	t.Run("handles failure to read manifest from store", func(t *testing.T) {

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
	})
}
