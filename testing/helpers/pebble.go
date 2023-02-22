package helpers

import (
	"testing"

	"github.com/cockroachdb/pebble"
	"github.com/cockroachdb/pebble/vfs"
	"github.com/stretchr/testify/require"
)

func InMemoryDB(t *testing.T) *pebble.DB {
	t.Helper()

	opts := pebble.Options{
		FS: vfs.NewMem(),
	}
	db, err := pebble.Open("", &opts)
	require.NoError(t, err)

	return db
}
