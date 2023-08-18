package fstore

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/blocklessnetwork/b7s/store"
	"github.com/blocklessnetwork/b7s/testing/helpers"
	"github.com/blocklessnetwork/b7s/testing/mocks"
)

func TestFunction_UnpackArchive(t *testing.T) {

	const (
		filename = "testdata/testFunction.tar.gz"
	)

	workdir, err := os.MkdirTemp("", "b7s-function-unpack-")
	require.NoError(t, err)

	defer os.RemoveAll(workdir)

	store := store.New(helpers.InMemoryDB(t))
	fh := New(mocks.NoopLogger, store, workdir)

	err = fh.unpackArchive(filename, workdir)
	require.NoError(t, err)
}

func TestFunction_UnpackArchiveHandlesErrors(t *testing.T) {
	t.Run("handles missing archive", func(t *testing.T) {

		const (
			filename = "testdata/nonExistantFile.tar.gz"
		)

		workdir, err := os.MkdirTemp("", "b7s-function-unpack-")
		require.NoError(t, err)

		defer os.RemoveAll(workdir)

		store := store.New(helpers.InMemoryDB(t))
		fh := New(mocks.NoopLogger, store, workdir)

		err = fh.unpackArchive(filename, workdir)
		require.Error(t, err)
	})
}
