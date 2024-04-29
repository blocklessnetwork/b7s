package fstore

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/testing/mocks"
)

func TestFstore_CheckFunctionFiles(t *testing.T) {

	workdir, err := os.MkdirTemp("", "b7s-function-get-")
	require.NoError(t, err)

	defer os.RemoveAll(workdir)

	store := mocks.BaselineStore(t)
	fh := New(mocks.NoopLogger, store, workdir)

	var (
		archiveName      = "archive.tar.gz"
		functionFileName = "function-file"
	)

	rec := blockless.FunctionRecord{
		Archive: archiveName,
		Files:   functionFileName,
	}

	t.Run("archive and function not found", func(t *testing.T) {

		haveArchive, haveFiles, err := fh.checkFunctionFiles(rec)
		require.NoError(t, err)

		require.False(t, haveArchive)
		require.False(t, haveFiles)
	})
	t.Run("archive and function found", func(t *testing.T) {

		archivePath := filepath.Join(workdir, archiveName)
		_, err := os.Create(archivePath)
		require.NoError(t, err)
		defer os.Remove(archivePath)

		functionPath := filepath.Join(workdir, functionFileName)
		_, err = os.Create(functionPath)
		require.NoError(t, err)
		defer os.Remove(functionPath)

		haveArchive, haveFiles, err := fh.checkFunctionFiles(rec)
		require.NoError(t, err)
		require.True(t, haveArchive)
		require.True(t, haveFiles)
	})
	t.Run("archive found, function files missing", func(t *testing.T) {

		archivePath := filepath.Join(workdir, archiveName)
		_, err := os.Create(archivePath)
		require.NoError(t, err)
		defer os.Remove(archivePath)

		haveArchive, haveFiles, err := fh.checkFunctionFiles(rec)
		require.NoError(t, err)
		require.True(t, haveArchive)
		require.False(t, haveFiles)
	})
	t.Run("archive missing, function files found", func(t *testing.T) {

		functionPath := filepath.Join(workdir, functionFileName)
		_, err = os.Create(functionPath)
		require.NoError(t, err)
		defer os.Remove(functionPath)

		haveArchive, haveFiles, err := fh.checkFunctionFiles(rec)
		require.NoError(t, err)
		require.False(t, haveArchive)
		require.True(t, haveFiles)
	})
}
