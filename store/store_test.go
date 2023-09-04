package store_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/store"
	"github.com/blocklessnetwork/b7s/testing/helpers"
)

func Test_Store(t *testing.T) {
	t.Run("setting value", func(t *testing.T) {
		t.Parallel()

		db := helpers.InMemoryDB(t)
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

		db := helpers.InMemoryDB(t)
		defer db.Close()
		store := store.New(db)

		read, err := store.Get("missing-key")
		require.Equal(t, "", read)
		require.ErrorIs(t, err, blockless.ErrNotFound)
	})
	t.Run("overwriting value", func(t *testing.T) {
		t.Parallel()

		db := helpers.InMemoryDB(t)
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

		db := helpers.InMemoryDB(t)
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

		db := helpers.InMemoryDB(t)
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

		db := helpers.InMemoryDB(t)
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

		db := helpers.InMemoryDB(t)
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
	t.Run("listing keys", func(t *testing.T) {
		t.Parallel()

		db := helpers.InMemoryDB(t)
		defer db.Close()
		store := store.New(db)

		readKeys, err := store.Keys()
		require.NoError(t, err)
		require.Empty(t, readKeys)

		keys := []string{
			"key1",
			"key2",
			"key3",
			"key4",
		}

		for _, key := range keys {
			err := store.SetRecord(key, struct{}{})
			require.NoError(t, err)
		}

		readKeys, err = store.Keys()
		require.NoError(t, err)
		require.Equal(t, keys, readKeys)
	})
	t.Run("deleting key", func(t *testing.T) {
		t.Parallel()

		db := helpers.InMemoryDB(t)
		defer db.Close()
		store := store.New(db)

		const (
			key   = "some-key"
			value = "some-value"
		)

		err := store.Set(key, value)
		require.NoError(t, err)

		// Deleting valid key works.
		err = store.Delete(key)
		require.NoError(t, err)

		// Value is no longer found.
		_, err = store.Get(key)
		require.ErrorIs(t, err, blockless.ErrNotFound)
	})
}
