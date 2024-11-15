package syncmap_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/blocklessnetwork/b7s/node/internal/syncmap"
)

func TestSyncMap(t *testing.T) {

	t.Run("setting and getting a value works", func(t *testing.T) {
		t.Parallel()

		const (
			key   = "key"
			value = "value"
		)

		m := syncmap.New[string, string]()

		m.Set(key, value)
		read, ok := m.Get(key)
		require.True(t, ok)
		require.Equal(t, value, read)
	})
	t.Run("getting nonexistant value works", func(t *testing.T) {
		t.Parallel()

		m := syncmap.New[string, string]()

		val, ok := m.Get("whatever")
		require.False(t, ok)
		require.Zero(t, val)
	})
	t.Run("deleting a value works", func(t *testing.T) {
		t.Parallel()

		const (
			key   = "key"
			value = "value"
		)

		m := syncmap.New[string, string]()

		m.Set(key, value)
		_, ok := m.Get(key)
		require.True(t, ok)

		m.Delete(key)
		_, ok = m.Get(key)
		require.False(t, ok)
	})
}
