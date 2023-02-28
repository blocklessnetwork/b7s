package waitmap

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWaitMap(t *testing.T) {
	t.Run("setting and getting a value works", func(t *testing.T) {
		t.Parallel()

		const (
			key   = "dummy-key"
			value = "dummy-value"
		)

		wm := New()

		wm.Set(key, value)
		require.Len(t, wm.m, 1)
		require.Equal(t, value, wm.m[key][0])

		retrieved, ok := wm.Get(key)
		require.True(t, ok)
		require.Equal(t, value, retrieved)
	})
	t.Run("getting a value not yet set works", func(t *testing.T) {
		t.Parallel()

		const (
			key = "dummy-key"
		)

		wm := New()

		_, ok := wm.Get(key)
		require.False(t, ok)
	})
	t.Run("waiting on a value works", func(t *testing.T) {
		t.Parallel()

		const (
			key   = "dummy-key"
			value = "dummy-value"
		)

		var (
			retrieved string
		)

		wm := New()

		var wg sync.WaitGroup
		wg.Add(1)
		// Spin up a goroutine that will asynchronously wait for a value to be set.
		go func() {
			defer wg.Done()
			waited := wm.Wait(key)

			retrieved = waited.(string)
		}()

		time.Sleep(200 * time.Millisecond)

		// Confirm that there is no value set yet.
		_, ok := wm.Get(key)
		require.False(t, ok)

		wm.Set(key, value)

		// Make sure to wait for the goroutine to complete.
		wg.Wait()

		require.Equal(t, value, retrieved)
	})
	t.Run("wait returns immediately if the value exists", func(*testing.T) {
		t.Parallel()

		const (
			key   = "dummy-key"
			value = "dummy-value"
		)

		wm := New()

		wm.Set(key, value)

		retrieved := wm.Wait(key)
		require.Equal(t, value, retrieved)
	})
	t.Run("multiple waiters get notified", func(t *testing.T) {
		t.Parallel()

		const (
			key   = "dummy-key"
			value = "dummy-value"
		)

		wm := New()

		var wg sync.WaitGroup
		wg.Add(3)

		fn := func() {
			defer wg.Done()
			waited := wm.Wait(key)

			require.Equal(t, value, waited.(string))
		}

		// Spin up three goroutines - they should all get the same result.
		go fn()
		go fn()
		go fn()

		wm.Set(key, value)

		wg.Wait()
	})
}
