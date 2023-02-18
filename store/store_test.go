package store_test

import (
	"testing"

	"github.com/cockroachdb/pebble"
	"github.com/cockroachdb/pebble/vfs"
	"github.com/stretchr/testify/require"

	"github.com/blocklessnetworking/b7s/models/blockless"
	"github.com/blocklessnetworking/b7s/store"
)

func Test_Get(t *testing.T) {
	t.Run("setting value", func(t *testing.T) {
		t.Parallel()

		db := setupDB(t)
		defer db.Close()
		store := store.New(db)

		const (
			key   = "some-key"
			value = "some-value"
		)

		err := store.Set(key, value)
		require.NoError(t, err)

		read, err := store.Get(key)
		require.NoError(t, err)

		require.Equal(t, value, read)
	})
	t.Run("missing value correctly reported", func(t *testing.T) {
		t.Parallel()

		db := setupDB(t)
		defer db.Close()
		store := store.New(db)

		read, err := store.Get("missing-key")
		require.Equal(t, "", read)
		require.ErrorIs(t, err, blockless.ErrNotFound)
	})
	t.Run("overwriting value", func(t *testing.T) {
		t.Parallel()

		db := setupDB(t)
		defer db.Close()
		store := store.New(db)

		const (
			key     = "some-key"
			valueV1 = "some-value-v1"
			valueV2 = "some-value-v2"
		)

		// Set value V1.
		err := store.Set(key, valueV1)
		require.NoError(t, err)

		read, err := store.Get(key)
		require.NoError(t, err)

		require.Equal(t, valueV1, read)

		// Set value V2.
		err = store.Set(key, valueV2)
		require.NoError(t, err)

		read, err = store.Get(key)
		require.NoError(t, err)

		require.Equal(t, read, valueV2)
	})
	t.Run("setting record", func(t *testing.T) {
		t.Parallel()

		db := setupDB(t)
		defer db.Close()
		store := store.New(db)

		const (
			key = "some-key"
		)

		type person struct {
			Name string
			Age  uint
		}

		var value = person{
			Name: "John",
			Age:  30,
		}

		err := store.SetRecord(key, value)
		require.NoError(t, err)

		var read person
		err = store.GetRecord(key, &read)
		require.NoError(t, err)

		require.Equal(t, value, read)
	})
	t.Run("handling missing record", func(t *testing.T) {
		t.Parallel()

		db := setupDB(t)
		defer db.Close()
		store := store.New(db)

		const (
			key = "some-key"
		)

		type person struct {
			Name string
			Age  uint
		}

		var read person
		err := store.GetRecord(key, &read)
		require.Equal(t, person{}, read)
		require.ErrorIs(t, err, blockless.ErrNotFound)
	})
	t.Run("overwriting record", func(t *testing.T) {
		t.Parallel()

		db := setupDB(t)
		defer db.Close()
		store := store.New(db)

		const (
			key = "some-key"
		)

		type person struct {
			Name string
			Age  uint
		}

		var value = person{
			Name: "John",
			Age:  30,
		}

		err := store.SetRecord(key, value)
		require.NoError(t, err)

		var read person
		err = store.GetRecord(key, &read)
		require.NoError(t, err)

		require.Equal(t, value, read)

		// Change record values.
		valueV2 := person{
			Name: "Paul",
			Age:  20,
		}

		err = store.SetRecord(key, valueV2)
		require.NoError(t, err)

		err = store.GetRecord(key, &read)
		require.NoError(t, err)
		require.Equal(t, valueV2, read)
	})
	t.Run("handle invalid output type", func(t *testing.T) {
		t.Parallel()

		db := setupDB(t)
		defer db.Close()
		store := store.New(db)

		const (
			key = "some-key"
		)

		type person struct {
			Name string
			Age  uint
		}

		var value = person{
			Name: "John",
			Age:  30,
		}

		err := store.SetRecord(key, value)
		require.NoError(t, err)

		type invalidModel struct {
			Name float64
			Age  string
		}

		var rec invalidModel
		err = store.GetRecord(key, &rec)
		require.Error(t, err)
	})

}

// Setup a new in-memory pebble database.
func setupDB(t *testing.T) *pebble.DB {
	t.Helper()

	opts := pebble.Options{
		FS: vfs.NewMem(),
	}
	db, err := pebble.Open("", &opts)
	require.NoError(t, err)

	return db
}
